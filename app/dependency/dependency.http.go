package dependency

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xhuma"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

func ProvideHTTPServerName(c config.Cfg) config.Server {
	return c.Server["http"]
}

type InvokeHTTPServerParam struct {
	fx.In

	Lifecycle fx.Lifecycle

	Cfg    config.Cfg
	Server config.Server
	Tracer trace.Tracer

	HumaAPI  huma.API
	FiberApp *fiber.App
	Log      *xlog.DebugLogger

	Handlers []xhuma.HandlerRegister `group:"global:http:handler"`
}

func InvokeHTTPServer(p InvokeHTTPServerParam) {
	var (
		cfg    = p.Server
		logger = xlog.NewLogger(p.Log.Logger)
	)

	for _, h := range p.Handlers {
		h.Register(p.HumaAPI)
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info(ctx, "http server started", "address", cfg.Address)
				if err := p.FiberApp.Listen(cfg.Address); err != nil && !errors.Is(err, net.ErrClosed) {
					logger.Info(ctx, "failed to start http server", "err", fmt.Sprintf("%+v", err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer logger.Info(ctx, "http server stopped", "address", cfg.Address)
			if err := p.FiberApp.ShutdownWithContext(ctx); err != nil && !errors.Is(err, net.ErrClosed) {
				return err
			}
			return nil
		},
	})
}
