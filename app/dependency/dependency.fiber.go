package dependency

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/infra/http/middleware"
)

func ProvideHumaFiber(a *fiber.App, c huma.Config) huma.API {
	var (
		api = humafiber.New(a, c)
	)
	return api
}

type ProvideFiberFxParam struct {
	fx.In

	Cfg fiber.Config
	Mdl []xhuma.GlobalMiddleware `group:"global:http:middleware"`
}

func ProvideFiber(p ProvideFiberFxParam) *fiber.App {
	var (
		app = fiber.New(p.Cfg)
	)

	sort.Slice(p.Mdl, func(i, j int) bool {
		var (
			a = p.Mdl[i].Name()
			b = p.Mdl[j].Name()
		)

		var (
			aa = middleware.GlobalMiddlewareOrder[a]
			bb = middleware.GlobalMiddlewareOrder[b]
		)

		return aa < bb
	})

	for _, mdl := range p.Mdl {
		mdl.App(app)
		app.Use(mdl.Serve)
	}

	return app
}

func ProvideFiberConfig(c config.Cfg, l *xlog.DebugLogger) fiber.Config {
	var (
		ctx = context.Background()
		log = xlog.NewLogger(l.Logger)
	)

	var (
		svr = c.Server["http"]
		add = svr.Additional
		fbr = fiber.Config{
			ReduceMemoryUsage:     true,
			DisableStartupMessage: true,
		}
	)

	if dur, ok := parseDurrationConfig(ctx, log, add, "idle.timeout"); ok {
		fbr.IdleTimeout = dur
	}
	if dur, ok := parseDurrationConfig(ctx, log, add, "write.timeout"); ok {
		fbr.WriteTimeout = dur
	}
	if dur, ok := parseDurrationConfig(ctx, log, add, "read.timeout"); ok {
		fbr.ReadTimeout = dur
	}

	return fbr
}

func parseDurrationConfig(ctx context.Context, log xlog.Logger, config map[string]string, key string) (time.Duration, bool) {
	v, ok := config[key]
	if !ok {
		return 0, false
	}

	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n < 0 {
		log.Error(ctx, fmt.Sprintf("failed to parse http additional config: '%s'", key), "err", fmt.Sprintf("%+v", err))
		return 0, false
	}

	return time.Duration(n) * time.Second, true
}
