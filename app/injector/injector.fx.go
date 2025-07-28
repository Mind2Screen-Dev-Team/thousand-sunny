package injector

import (
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"go.uber.org/fx"
)

var (
	Fx = fx.Options(
		fx.WithLogger(xlog.SetupFxLogger),
	)
)
