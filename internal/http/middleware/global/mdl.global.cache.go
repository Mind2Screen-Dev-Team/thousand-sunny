package global

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/constant"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/internal/util"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xfiber"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xsecurity"
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

func (Cache) App(app *fiber.App) {}

func (s Cache) Serve(c *fiber.Ctx) error {
	if next, ok := xfiber.SkipPath(c, constant.FiberSkipablePathFromMiddleware[:]...); ok {
		return next()
	}

	var (
		rawReqBody = c.BodyRaw()
		reqBody    = make([]byte, len(rawReqBody))
	)
	copy(reqBody, rawReqBody)

	var (
		ctx            = c.UserContext()
		acceptEncoding = c.Get("Accept-Encoding")
		reqUrl, _      = url.Parse(string(c.Request().RequestURI()))
		hexHash256     = xsecurity.HexHashSHA256(fmt.Sprintf("%s|%s|%s", c.Method(), reqUrl.String(), string(reqBody)))
		cacheType      = "plain"
		isCompressed   = strings.Contains(acceptEncoding, "gzip") ||
			strings.Contains(acceptEncoding, "deflate") ||
			strings.Contains(acceptEncoding, "br")
	)

	if isCompressed {
		cacheType = acceptEncoding
	}

	var (
		cacheKey = fmt.Sprintf(
			"api_res_cache:%s:method:%s:path:%s:hash_req_data:%s:%s",
			s.cfg.App.Env,
			c.Method(),
			util.GetSnakeCaseKeyURL(reqUrl),
			hexHash256,
			cacheType,
		)
		cacheLogFlagKey = cacheKey + ":hide_res_log"
		cacheHeaderKey  = cacheKey + ":headers"
	)

	ctx, span := xtracer.Start(s.tracer, ctx, "cache http response api")
	defer span.End()

	if cached, err := s.client.Get(ctx, cacheKey).Bytes(); err == nil {
		if val, err := s.client.Get(ctx, cacheLogFlagKey).Result(); err == nil && val == "1" {
			ctx = context.WithValue(ctx, xlog.XLOG_HIDE_RES_FLAG_CTX_KEY, true)
			c.SetUserContext(ctx)
		}

		if val, err := s.client.Get(ctx, cacheHeaderKey).Bytes(); err == nil && len(val) > 0 {
			var m map[string][]string
			if err := json.Unmarshal(val, &m); err == nil {
				for k, s := range m {
					for i, v := range s {
						if i == 0 {
							c.Set(k, v)
						} else {
							c.Append(k, v)
						}
					}
				}
			}
		}

		if isCompressed {
			c.Set("Content-Encoding", acceptEncoding)
		}

		c.Status(fiber.StatusOK)
		return c.Send(cached)
	}

	c.SetUserContext(ctx)

	if err := c.Next(); err != nil {
		return err
	}

	var buf bytes.Buffer
	c.Response().BodyWriteTo(&buf)

	defer buf.Reset()

	// Filter headers to cache
	ignores := []string{"X-Cache", "X-Cache-Exp"}
	headers := make(map[string][]string)
	for key, value := range c.Response().Header.All() {
		if slices.Contains(ignores, string(key)) {
			continue
		}
		headers[string(key)] = append(headers[string(key)], string(value))
	}

	// Cache control
	ttl := 30 * time.Second // default
	xcache := c.Get("X-Cache")
	if expStr := c.Get("X-Cache-Exp"); expStr != "" {
		if exp, err := strconv.ParseInt(expStr, 10, 64); err == nil && exp > 0 {
			ttl = time.Duration(exp) * time.Second
		}
	}

	if xcache == "true" {
		_ = s.client.Set(ctx, cacheKey, buf.Bytes(), ttl).Err()

		// log flag
		if v, ok := ctx.Value(xlog.XLOG_HIDE_RES_FLAG_CTX_KEY).(bool); ok && v {
			_ = s.client.Set(ctx, cacheLogFlagKey, "1", ttl).Err()
		}

		// header cache
		if len(headers) > 0 {
			if hBytes, err := json.Marshal(headers); err == nil {
				_ = s.client.Set(ctx, cacheHeaderKey, hBytes, ttl).Err()
			}
		}
	}

	return nil
}
