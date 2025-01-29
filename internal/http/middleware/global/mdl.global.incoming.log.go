package http_middleware_global

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhttp"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xpanic"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	min "github.com/tdewolff/minify/v2"
	min_json "github.com/tdewolff/minify/v2/json"
)

// const DefaultMultipartMaxMemory = 100 << 20

func ProvideIncomingLog(cfg config.Cfg, iolog *xlog.IOLogger, debuglog *xlog.DebugLogger, async *xasynq.Asynq) IncomingLog {
	return IncomingLog{
		cfg:      cfg,
		async:    async,
		iolog:    xlog.NewLogger(iolog.Logger),
		debuglog: xlog.NewLogger(debuglog.Logger),
	}
}

type IncomingLog struct {
	cfg      config.Cfg
	async    *xasynq.Asynq
	iolog    xlog.Logger
	debuglog xlog.Logger
}

func (IncomingLog) Name() string {
	return "incoming.log"
}

func (in IncomingLog) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		var (
			req     = c.Request()
			reqBody = xhttp.CopyRequestBody(req)
			reqTime = time.Now().Format(time.RFC3339Nano)
			rw      = bytes.Buffer{}
			w       = c.Response().Writer
			ww      = chi_middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		)

		ww.Tee(&rw)
		c.SetRequest(req)
		c.Response().Writer = ww

		if err := next(c); err != nil {
			return nil
		}

		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		var (
			rCopy = xhttp.DeepCopyRequest(req)
		)

		defer func() {
			var (
				panicked     bool
				recorverErr  any
				parsedStacks = make([]xpanic.Stack, 0)
			)

			if r := recover(); r != nil && r != http.ErrAbortHandler {
				// assign
				panicked = true
				recorverErr = r

				stacks := debug.Stack()
				parsedStacks = xpanic.ParseStack(bytes.NewReader(stacks))

				c.JSON(http.StatusInternalServerError, map[string]any{
					"code": http.StatusInternalServerError,
					"msg":  http.StatusText(http.StatusInternalServerError),
					"data": nil,
				})
			}

			go in._IncomingLogging(c, rCopy, ww, &rw, reqTime, panicked, recorverErr, parsedStacks)
		}()

		return nil
	})
}

func (in IncomingLog) _IncomingLogging(
	// general
	c echo.Context,

	// general params
	r *http.Request,
	ww chi_middleware.WrapResponseWriter,
	resbody *bytes.Buffer,
	reqtime string,

	// panic params
	panicked bool,
	recorverErr any,
	stacks []xpanic.Stack,
) {
	var (
		bw bytes.Buffer
		rw bytes.Buffer

		min            = min.New()
		reqRawBytes, _ = io.ReadAll(r.Body)
		reqBody        = io.NopCloser(bytes.NewReader(reqRawBytes))
		resMim         = mimetype.Detect(resbody.Bytes())

		parsedBody, contentSize = in._ParseRequest(r, reqBody)
		resFilename             = in._GetFilenameFromHeader(ww.Header().Get("Content-Disposition"))
		_                       = in._MinifyResponse(min, resMim, &rw, resbody, resFilename)
	)

	defer func() {
		// remove all multipart from copied request
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}

		// copied body reset
		bw.Reset()

		// copied response reset
		rw.Reset()

		// request reset
		reqBody.Close()

		// response reset
		resbody.Reset()
	}()

	var (
		contentType           = strings.ToLower(r.Header.Get("Content-Type"))
		isReqMultipartEncoded = strings.HasPrefix(contentType, "application/x-www-form-urlencoded")
		isReqMultipartData    = strings.HasPrefix(contentType, "multipart/form-data")
	)

	if (isReqMultipartEncoded || isReqMultipartData) || len(reqRawBytes) <= 0 {
		reqRawBytes = nil
	}

	var (
		fields = []any{
			// request
			"reqTraceId", c.Get(xlog.XLOG_TRACE_ID_KEY),
			"reqTime", reqtime,
			"reqRemoteAddress", r.RemoteAddr,
			"reqPath", r.URL.Path,
			"reqHeader", r.Header,
			"reqProto", r.Proto,
			"reqMethod", r.Method,
			"reqUserAgent", r.UserAgent(),
			"reqBodyRaw", reqRawBytes,
			"reqBodyParsed", parsedBody,
			"reqBytesIn", contentSize,

			// response
			"resHeader", ww.Header(),
			"resStatus", http.StatusText(ww.Status()),
			"resStatusCode", ww.Status(),
			"resBody", rw.String(),
			"resBytesOut", ww.BytesWritten(),
		}
	)

	if resFilename != "" {
		fields = append(fields, "resFileName", resFilename)
	}

	if panicked {
		fields = append(fields, "panicRecoverErr", recorverErr)
		fields = append(fields, "panicStack", stacks)
	}

	in.iolog.Info("incoming request", fields...)

	// # Notify Process Incoming Log
	ioLogCfg, ok := in.cfg.Log["io"]
	if !ok || !ioLogCfg.Notify.Enabled {
		return
	}

	m := make(map[string]any)
	for i := 0; i < len(fields); i++ {
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				key, ok := fields[i].(string)
				if ok {
					if slices.Contains([]string{"err", "error"}, key) {
						m[key] = fmt.Sprintf("%+v", fields[i+1])
						continue
					}
					m[key] = fields[i+1]
				}
			}
		}
	}

	in._Notify(&ioLogCfg, m)
}

func (in *IncomingLog) _Notify(ioLogCfg *config.Log, m map[string]any) {
	var (
		b, _      = json.Marshal(m)
		name      = xasynq.BuildWorkerRouteName(in.cfg.App.Env, "notify:incoming:log")
		task      = asynq.NewTask(name, b, asynq.Queue(name), asynq.Retention(time.Duration(ioLogCfg.Notify.Retention)*time.Hour))
		info, err = in.async.Client.Enqueue(task)
	)
	if err != nil {
		if ioLogCfg.Notify.Debug {
			in.debuglog.Error("error send notification incoming log into asynq", "info", info, "error", err)
		}
		return
	}

	if ioLogCfg.Notify.Debug {
		in.debuglog.Debug("send notification incoming log into asynq", "info", info)
	}
}

func (IncomingLog) _GetFilenameFromHeader(d string) string {
	const key = "filename="
	idx := strings.Index(d, key)
	if idx == -1 {
		return ""
	}
	return strings.Trim(d[idx+len(key):], "\"")
}

func (in IncomingLog) _MinifyResponse(min *min.M, mim *mimetype.MIME, w io.Writer, res *bytes.Buffer, filename string) error {
	if mim.Is("application/json") {
		return min_json.Minify(min, w, res, nil)
	}

	if mim.Is("text/plain") && filename == "" {
		if _, err := res.WriteTo(w); err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		return nil
	}

	return nil
}

func (in IncomingLog) _ParseRequest(req *http.Request, r io.Reader) (any, int64) {
	var (
		reqBody any
		total   int64
	)

	if req.ContentLength > 0 {
		total = req.ContentLength
	}

	var (
		contentType           = strings.ToLower(req.Header.Get("Content-Type"))
		isReqJson             = strings.HasPrefix(contentType, "application/json")
		isReqMultipartEncoded = strings.HasPrefix(contentType, "application/x-www-form-urlencoded")
		isReqMultipartData    = strings.HasPrefix(contentType, "multipart/form-data")
	)

	if isReqJson {
		if err := json.NewDecoder(r).Decode(&reqBody); err != nil {
			return nil, total
		}

		_total, _ := io.Copy(io.Discard, r)
		if _total > 0 {
			total = _total
		}

		return reqBody, total
	}

	if isReqMultipartData || isReqMultipartEncoded {
		if req.MultipartForm == nil {
			return nil, total
		}

		var (
			_total int64
			mr, _  = req.MultipartReader()
		)
		for {
			if mr == nil {
				break
			}

			np, err := mr.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, 0
			}

			size, _ := io.Copy(io.Discard, np)
			_total += size
		}

		if _total > 0 {
			total = _total
		}

		// Extract form fields
		forms := make(map[string][]string)
		for key, values := range req.MultipartForm.Value {
			forms[key] = values
		}

		files := make(map[string]xhttp.FileInfo)
		for key, fileHeaders := range req.MultipartForm.File {
			for idx, fileHeader := range fileHeaders {
				// Read file metadata
				fileInfo := xhttp.FileInfo{
					FileName:    fileHeader.Filename,
					ContentType: fileHeader.Header.Get("Content-Type"),
					Size:        fileHeader.Size,
				}

				files[fmt.Sprintf("%s.%d", key, idx)] = fileInfo
			}
		}

		return map[string]any{"forms": forms, "files": files}, total
	}

	return nil, total
}
