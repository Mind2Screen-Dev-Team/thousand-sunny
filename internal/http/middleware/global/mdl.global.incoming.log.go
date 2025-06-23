package http_middleware_global

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xasynq"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhttp"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xpanic"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"

	"go.opentelemetry.io/otel/trace"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	min "github.com/tdewolff/minify/v2"
	min_json "github.com/tdewolff/minify/v2/json"
)

// const DefaultMultipartMaxMemory = 100 << 20

func ProvideIncomingLog(cfg config.Cfg, tracer trace.Tracer, iolog *xlog.IOLogger, debuglog *xlog.DebugLogger, async *xasynq.Asynq) IncomingLog {
	return IncomingLog{
		cfg:      cfg,
		tracer:   tracer,
		async:    async,
		iolog:    xlog.NewLogger(iolog.Logger),
		debuglog: xlog.NewLogger(debuglog.Logger),
	}
}

type IncomingLog struct {
	cfg      config.Cfg
	tracer   trace.Tracer
	async    *xasynq.Asynq
	iolog    xlog.Logger
	debuglog xlog.Logger
}

func (IncomingLog) Name() string {
	return "incoming.log"
}

func (IncomingLog) Order() int {
	return 3
}

func (in IncomingLog) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		var (
			req     = c.Request()
			ctx     = req.Context()
			wg      = sync.WaitGroup{}
			reqBody = xhttp.CopyRequestBody(req)
			reqTime = time.Now().Format(time.RFC3339Nano)
			rw      = bytes.Buffer{}
			w       = c.Response().Writer
			ww      = chi_middleware.NewWrapResponseWriter(w, req.ProtoMajor)
		)

		ctx, span := xtracer.Start(in.tracer, ctx, "incoming log")

		// Set Request Body
		if len(reqBody) > 0 {
			c.Set(xlog.XLOG_REQ_BODY_KEY, string(reqBody))
		}

		ww.Tee(&rw)
		c.SetRequest(req.WithContext(ctx))
		c.Response().Writer = ww

		if err := next(c); err != nil {
			return err
		}

		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		var (
			rCopy = xhttp.DeepCopyRequest(req)
		)

		// Add BEFORE defer/spawned goroutine to prevent early Wait() for span.End()
		wg.Add(1)

		defer func() {
			var (
				panicked     bool
				recoverErr   any
				parsedStacks = make([]xpanic.Stack, 0)
			)

			if r := recover(); r != nil && r != http.ErrAbortHandler {
				// assign
				panicked = true
				recoverErr = r

				stacks := debug.Stack()
				parsedStacks = xpanic.ParseStack(bytes.NewReader(stacks))

				// Only try writing error response if not already written
				if !c.Response().Committed {
					c.JSON(http.StatusInternalServerError, map[string]any{
						"code": http.StatusInternalServerError,
						"msg":  http.StatusText(http.StatusInternalServerError),
						"data": nil,
					})
				}
			}

			go func() {
				defer wg.Done()
				in.incomingLogging(ctx, c, rCopy, ww, &rw, reqTime, panicked, recoverErr, parsedStacks)
			}()
		}()

		// Wait for logging and end span in a blocking goroutine
		go func() {
			wg.Wait()
			span.End()
		}()

		return nil
	})
}

func (in IncomingLog) incomingLogging(
	// general
	ctx context.Context,
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
	var skipRes bool
	if v, ok := c.Get("skip_log_response").(bool); ok {
		skipRes = v
	}

	var (
		isGzip = strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "gzip")
		bw     bytes.Buffer
		rw     bytes.Buffer

		min            = min.New()
		reqRawBytes, _ = io.ReadAll(r.Body)
		reqBody        = io.NopCloser(bytes.NewReader(reqRawBytes))
		resMim         = mimetype.Detect(resbody.Bytes())

		parsedBody, contentSize = in.parseRequest(r, reqBody)
		resFilename             = in.getFilenameFromHeader(ww.Header().Get("Content-Disposition"))
		_                       = in.minifyResponse(skipRes, isGzip, min, resMim, &rw, resbody, resFilename)
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
		reqTime, _ = time.Parse(time.RFC3339Nano, reqtime)
		fields     = []any{
			// set to ignore list keys by group key name
			xlog.IGNORE_KEY, "incoming.log",

			// request
			"req_trace_id", c.Get(xlog.XLOG_TRACE_ID_KEY),
			"req_time", reqtime,
			"req_remote_address", r.RemoteAddr,
			"req_path", r.URL.Path,
			"req_header", r.Header,
			"req_proto", r.Proto,
			"req_method", r.Method,
			"req_user_agent", r.UserAgent(),
			"req_body_raw", reqRawBytes,
			"req_body_parsed", parsedBody,
			"req_bytes_in", contentSize,

			// response
			"res_header", ww.Header(),
			"res_status", http.StatusText(ww.Status()),
			"res_status_code", ww.Status(),
			"res_body", rw.String(),
			"res_bytes_out", ww.BytesWritten(),
			"res_latency", time.Since(reqTime).String(),
		}
	)

	if resFilename != "" {
		fields = append(fields, "res_filename", resFilename)
	}

	if panicked {
		fields = append(fields, "panic_recover_err", recorverErr)
		fields = append(fields, "panic_stack", stacks)
	}

	in.iolog.Info(ctx, "incoming request", fields...)

	// # Notify Process Incoming Log
	ioLogCfg, ok := in.cfg.Log.LogType["io"]
	if !ok || !ioLogCfg.Notify.Enabled {
		return
	}

	m := make(map[string]any)
	for i := 0; i < len(fields); i += 2 {
		if isHasKV := i+1 < len(fields); !isHasKV {
			continue
		}

		key, ok := fields[i].(string)
		if !ok {
			continue
		}

		val := fields[i+1]
		if slices.Contains([]string{"err", "error"}, key) {
			key = "err"
			val = fmt.Sprintf("%+v", val)
		}
		m[key] = val
	}

	in.notify(ctx, ioLogCfg, m)
}

func (in *IncomingLog) notify(ctx context.Context, ioLogCfg config.LogType, m map[string]any) {
	var (
		b, _ = json.Marshal(m)
		name = xasynq.BuildWorkerRouteName(in.cfg.App.Env, "notify:incoming:log")
	)

	var (
		retention = time.Duration(ioLogCfg.Notify.Retention)
		task      = asynq.NewTask(name, b, asynq.Queue(name), asynq.Retention(retention*time.Second))
		info, err = in.async.Client.Enqueue(task)
	)
	if err != nil {
		if ioLogCfg.Notify.Debug {
			in.debuglog.Error(ctx, "error send notification incoming log into asynq", "info", info, "error", err)
		}
		return
	}

	if ioLogCfg.Notify.Debug {
		in.debuglog.Debug(ctx, "send notification incoming log into asynq", "info", info)
	}
}

func (IncomingLog) getFilenameFromHeader(d string) string {
	const key = "filename="
	idx := strings.Index(d, key)
	if idx == -1 {
		return ""
	}
	return strings.Trim(d[idx+len(key):], "\"")
}

func (in IncomingLog) minifyResponse(skipRes bool, isGzip bool, min *min.M, mim *mimetype.MIME, w io.Writer, res *bytes.Buffer, filename string) error {
	if skipRes {
		return nil
	}

	if isGzip {
		var (
			enc    = base64.NewEncoder(base64.RawURLEncoding, w)
			_, err = res.WriteTo(enc)
		)

		defer enc.Close() // ensure flush even on early return

		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		return nil
	}

	if mim.Is("application/json") {
		return min_json.Minify(min, w, res, nil)
	}

	if mim.Is("text/plain") && filename == "" {
		_, err := res.WriteTo(w)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}

func (in IncomingLog) parseRequest(req *http.Request, r io.Reader) (any, int64) {
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
