package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/config"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xlog"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/pkg/xresp"
)

type (
	PrivateAuthJWTParams struct {
		fx.In

		Cfg   config.Cfg
		Debug *xlog.DebugLogger
	}

	PrivateAuthJWT struct {
		cfg   config.Cfg
		debug xlog.Logger
	}
)

func NewPrivateAuthJWT(p PrivateAuthJWTParams) (*PrivateAuthJWT, error) {
	if p.Debug == nil {
		return nil, errors.New("field 'Debug' with type '*xlog.DebugLogger' is not provided")
	}

	return &PrivateAuthJWT{cfg: p.Cfg, debug: xlog.NewLogger(p.Debug.Logger)}, nil
}

func (a PrivateAuthJWT) Serve(c huma.Context, next func(c huma.Context)) {
	var (
		ctx    = c.Context()
		auth   = c.Header("Authorization")
		code   = http.StatusUnauthorized
		tid, _ = ctx.Value(xlog.XLOG_REQ_TRACE_ID_CTX_KEY).(string)
		resp   = xresp.GeneralResponseError{
			Code:    code,
			Msg:     http.StatusText(code),
			TraceID: tid,
		}
	)

	if !strings.HasPrefix(auth, "Bearer ") {
		c.SetHeader("Content-Type", "application/json")
		json.NewEncoder(c.BodyWriter()).Encode(resp)
		return
	}

	if token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")); token != "abc" {
		c.SetHeader("Content-Type", "application/json")
		json.NewEncoder(c.BodyWriter()).Encode(resp)
		return
	}

	a.debug.Info(ctx, "auth is success")

	next(c)
}
