package bootstrap

import (
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/bootstrap/dependency"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
)

func (Registry) DependencyStartUp() Provider {
	return Provider{
		fx.WithLogger(xlog.SetupFxLogger),
		fx.Invoke(dependency.PingRedisDB),
		fx.Invoke(dependency.PingPGxDB),
		fx.Invoke(dependency.InitEchoServer),
	}
}
