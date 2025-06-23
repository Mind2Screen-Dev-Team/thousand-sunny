package http_middleware_global

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	chi_middleware "github.com/go-chi/chi/v5/middleware"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/helper"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhttp"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xtracer"
)

func ProvideCache(cfg config.Cfg, client *redis.Client, tracer trace.Tracer) Cache {
	return Cache{cfg, client, tracer}
}

type Cache struct {
	cfg    config.Cfg
	client *redis.Client
	tracer trace.Tracer
}

func (Cache) Name() string {
	return "cache"
}

func (s Cache) Order() int {
	return 4
}

func (s Cache) Serve(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			req      = c.Request()
			ctx      = req.Context()
			encoding = req.Header.Get("Accept-Encoding")
			reqBody  = xhttp.CopyRequestBody(req)
			hash256  = helper.HashSHA256(req.Method + "|" + req.RequestURI + "|" + string(reqBody))
			isGzip   = strings.Contains(encoding, "gzip")

			cacheType = func() string {
				if isGzip {
					return "gzip"
				}
				return "plain"
			}()

			cacheKey = fmt.Sprintf(
				"api_res_cache:%s:method:%s:path:%s:hash_req_data:%s:%s",
				s.cfg.App.Env,
				req.Method,
				helper.GetSnakeCaseKey(req),
				hash256,
				cacheType,
			)
			cacheLogFlagKey = cacheKey + ":log_flag"

			span trace.Span
		)

		ctx, span = xtracer.Start(s.tracer, ctx, "getter / setter response cache")
		defer span.End()

		// Check Redis for cache
		if cached, err := s.client.Get(ctx, cacheKey).Bytes(); err == nil {
			if val, err := s.client.Get(ctx, cacheLogFlagKey).Result(); err == nil && val == "1" {
				xlog.SkipLogResponse(c)
			}

			if isGzip {
				c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
			}

			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c.Response().WriteHeader(http.StatusOK)
			_, err := c.Response().Write(cached)
			return err
		}

		var (
			buff = &bytes.Buffer{}
			ww   = chi_middleware.NewWrapResponseWriter(c.Response().Writer, req.ProtoMajor)
		)

		ww.Tee(buff)
		c.Response().Writer = ww
		c.Response().After(func() {
			if buff != nil {
				buff.Reset()
			}
		})

		if err := next(c); err != nil {
			return err
		}

		// Only cache if explicitly enabled
		var (
			ttl          = 30 * time.Second
			xcache       = c.Response().Header().Get("X-Cache")
			xcacheExp, _ = strconv.ParseInt(c.Response().Header().Get("X-Cache-Exp"), 10, 64)
		)

		if xcacheExp > 0 {
			ttl = time.Duration(xcacheExp) * time.Second
		}

		if xcache == "true" {
			s.client.Set(ctx, cacheKey, buff.Bytes(), ttl).Err()
			// Check if handler wants to log even from cache
			if v, ok := c.Get("skip_log_response").(bool); ok && v {
				s.client.Set(ctx, cacheLogFlagKey, "1", ttl)
			}
		}

		return nil // already written to client
	}
}
