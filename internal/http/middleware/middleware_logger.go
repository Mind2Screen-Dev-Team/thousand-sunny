package http_middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhttp"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xpanic"

	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	min "github.com/tdewolff/minify/v2"
	min_json "github.com/tdewolff/minify/v2/json"
)

func ProvideIncomingLog(iolog *xlog.IOLogger) IncomingLog {
	return IncomingLog{iolog: xlog.NewLogger(iolog.Logger)}
}

type IncomingLog struct {
	iolog xlog.Logger
}

func (in IncomingLog) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	fn := func(c echo.Context) error {
		var (
			r       = c.Request()
			reqtime = time.Now().Format(time.RFC3339Nano)
			rw      = bytes.Buffer{}
		)

		if err := r.ParseMultipartForm(50 << 20); err != nil {
			fmt.Printf("failed parse multipart-form, got error: %+v\n", err)
		}

		var (
			w     = c.Response().Writer
			rCopy = xhttp.DeepCopyRequest(r, true)
			ww    = chi_middleware.NewWrapResponseWriter(w, r.ProtoMajor)
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

			go in._IncomingLogging(c, rCopy, ww, &rw, reqtime, panicked, recorverErr, parsedStacks)
		}()

		ww.Tee(&rw)
		c.SetRequest(r)
		c.Response().Writer = ww

		return next(c)
	}

	return echo.HandlerFunc(fn)
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
		if r.MultipartForm != nil {
			r.MultipartForm.RemoveAll()
		}

		bw.Reset()
		rw.Reset()
		resbody.Reset()
		reqBody.Close()
	}()

	if len(reqRawBytes) <= 0 {
		reqRawBytes = nil
	}

	var (
		fields = []any{
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
			"resHeader", ww.Header(),
			"resStatus", http.StatusText(ww.Status()),
			"resStatusCode", ww.Status(),
			"reqBytesIn", contentSize,
			"resBytesOut", ww.BytesWritten(),
			"resBody", rw.String(),
		}
	)

	if resFilename != "" {
		fields = append(fields, "resFileName", resFilename)
	}

	if panicked {
		fields = append([]any{"panicRecover", recorverErr, "panicStack", stacks}, fields...)
	}

	in.iolog.Info("incoming request", fields...)
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
		_, err := res.WriteTo(w)
		if err != nil && !errors.Is(err, io.EOF) {
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
