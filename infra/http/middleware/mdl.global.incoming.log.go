package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/constant"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xfiber"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xpanic"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"
)

func ProvideIncomingLog(cfg config.Cfg, tracer trace.Tracer, debugLog *xlog.DebugLogger) IncomingLog {
	return IncomingLog{
		cfg:      cfg,
		tracer:   tracer,
		debugLog: xlog.NewLogger(debugLog.Logger),
	}
}

type IncomingLog struct {
	cfg      config.Cfg
	tracer   trace.Tracer
	debugLog xlog.Logger
}

func (IncomingLog) Name() string {
	return "incoming.log"
}

func (IncomingLog) App(app *fiber.App) {}

func (in IncomingLog) Serve(c *fiber.Ctx) error {
	if next, ok := xfiber.SkipPath(c, constant.FiberSkipablePathFromMiddleware[:]...); ok {
		return next()
	}

	var (
		now    = time.Now()
		ua     = c.Context().UserAgent()
		ae     = c.Get("Accept-Encoding")
		ct     = c.Get("Content-Type")
		tid, _ = c.UserContext().Value(xlog.XLOG_REQ_TRACE_ID_CTX_KEY).(string)

		isCompress                      = in.isCompress(ae)
		isMultipart, isMultipartEncoded = in.isMultipart(ct)

		rawReqBody = c.BodyRaw()
		reqBody    = make([]byte, len(rawReqBody))
		reqSize, _ = io.Copy(io.Discard, strings.NewReader(string(rawReqBody)))
	)

	copy(reqBody, rawReqBody)

	var (
		wg        sync.WaitGroup
		ctx, span = xtracer.Start(in.tracer, c.UserContext(), "incoming log")
	)

	var (
		reqFormBody = xlog.IncomingLogFormData{
			Values: make(map[string][]string),
			Files:  make(map[string]xlog.FileInfo),
		}
		rawForms, err = c.MultipartForm()
	)
	if err == nil {
		maps.Copy(reqFormBody.Values, rawForms.Value)
		reqSize += in.sizeMapSliceOfString(reqFormBody.Values)

		for k, h := range rawForms.File {
			for i, v := range h {
				reqFormBody.Files[fmt.Sprintf("%s.%d", k, i)] = xlog.FileInfo{
					FileName:    v.Filename,
					ContentType: v.Header.Get("Content-Type"),
					Size:        v.Size,
				}
				reqSize += v.Size
			}
		}
	}

	var (
		r   = c.Request()
		ip  = c.IP()
		ips = c.IPs()
		uri = r.URI().FullURI()
		mtd = r.Header.Method()
		prt = c.Protocol()
	)

	var (
		d = xlog.IncomingLogData{
			ReqUA:       ua,
			ReqIP:       ip,
			ReqIPs:      ips,
			ReqURI:      uri,
			ReqMethod:   mtd,
			ReqProtocol: prt,

			ReqSize:     reqSize,
			ReqBody:     reqBody,
			ReqFormBody: reqFormBody,

			ReqHeader: make(map[string][]string),
			ResHeader: make(map[string][]string),
		}
	)

	for key, value := range c.Request().Header.All() {
		if _, ok := d.ReqHeader[string(key)]; !ok {
			d.ReqHeader[string(key)] = make([]string, 0)
		}
		d.ReqHeader[string(key)] = append(d.ReqHeader[string(key)], string(value))
	}

	wg.Add(1)
	defer func() {
		if r := recover(); r != nil {
			d.IsPanic = true
			d.PanicMsg = fmt.Sprintf("%v", r)

			var (
				s  = debug.Stack()
				b  = bytes.NewBuffer(s)
				ss = xpanic.ParseStack(b)
			)

			b.Reset()

			json.NewEncoder(b).Encode(ss)
			d.PanicStack = b.Bytes()

			b.Reset()
		}

		if d.IsPanic {
			var (
				code = http.StatusInternalServerError
				msg  = http.StatusText(code)
				resp = xresp.GeneralResponseError{
					Code: code,
					Msg:  msg,
					Err: &xresp.ErrorModel{
						Title:  "panic error",
						Detail: d.PanicMsg,
					},
					TraceID: tid,
				}
			)
			c.Response().Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			json.NewEncoder(c.Response().BodyWriter()).Encode(resp)
		}

		var (
			res             = c.Response()
			code            = res.StatusCode()
			buffBytesRes, _ = c.Response().BodyUncompressed()
		)

		for key, value := range res.Header.All() {
			if _, ok := d.ResHeader[string(key)]; !ok {
				d.ResHeader[string(key)] = make([]string, 0)
			}
			d.ResHeader[string(key)] = append(d.ResHeader[string(key)], string(value))
		}

		d.ResStatus = code
		d.ResBody = buffBytesRes
		d.ResSize = int64(len(buffBytesRes))

		d.TimeStart = now
		d.TimeEnd = time.Now().Add(time.Since(now))

		d.IsHideRes, _ = ctx.Value(xlog.XLOG_HIDE_RES_FLAG_CTX_KEY).(bool)
		d.ReqTraceID = tid

		d.IsCompressed = isCompress
		d.IsMultipart = isMultipart
		d.IsMultipartEncoded = isMultipartEncoded

		go func() {
			defer wg.Done()
			in.log(ctx, d)
		}()

		// Wait for logging and end span in a blocking goroutine
		go func() {
			wg.Wait()
			span.End()
		}()
	}()

	c.SetUserContext(ctx)

	return c.Next()
}

func (IncomingLog) sizeMapSliceOfString(m map[string][]string) int64 {
	var n int64
	for _, v := range m {
		n += int64(len(v))
	}
	return n
}

func (IncomingLog) isCompress(ae string) bool {
	return strings.Contains(ae, "gzip") || strings.Contains(ae, "deflate") || strings.Contains(ae, "brotli")
}

func (IncomingLog) isMultipart(ct string) (isMultipart, isMultipartEncoded bool) {
	return strings.HasPrefix(ct, "multipart/form-data"), strings.HasPrefix(ct, "application/x-www-form-urlencoded")
}

func (in IncomingLog) log(ctx context.Context, d xlog.IncomingLogData) {
	var (
		fields = []any{
			"reqTraceId", d.ReqTraceID,
			"reqStartTime", d.TimeStart.Format(time.RFC3339Nano),
			"reqEndTime", d.TimeEnd.Format(time.RFC3339Nano),
			"reqIp", d.ReqIP,
			"reqForwardIps", d.ReqIPs,
			"reqUri", d.ReqURI,
			"reqHeader", d.ReqHeader,
			"reqProto", d.ReqProtocol,
			"reqMethod", d.ReqMethod,
			"reqUserAgent", d.ReqUA,
			"reqBytesIn", d.ReqSize,
			"resHeader", d.ResHeader,
			"resStatus", d.ResStatus,
			"resBytesOut", d.ResSize,
			"resLatency", d.TimeEnd.Sub(d.TimeStart).String(),
		}
	)

	if d.IsMultipart || d.IsMultipartEncoded {
		fields = append(fields, "reqFormBody", d.ReqFormBody)
	} else if len(d.ReqBody) > 0 {
		fields = append(fields, "reqRawBody", d.ReqBody)
	}

	var isPrintableRes bool
	for _, v := range d.ResHeader["Content-Type"] {
		if strings.Contains(v, "application/json") || strings.Contains(v, "text/plain") {
			isPrintableRes = true
		}
	}

	if !d.IsHideRes && isPrintableRes {
		fields = append(fields, "resBody", d.ResBody)
	}

	if d.IsPanic {
		fields = append(fields, "panicMsg", d.PanicMsg)
		fields = append(fields, "panicStack", d.PanicStack)
	}

	in.debugLog.Info(ctx, "incoming log request", fields...)
}
