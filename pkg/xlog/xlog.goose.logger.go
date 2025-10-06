package xlog

import (
	"context"
	"fmt"

	"github.com/pressly/goose/v3"
)

type GooseLogger struct {
	Log Logger
	Ctx context.Context
}

func NewGooseLogger(log Logger, ctx context.Context) *GooseLogger {
	if ctx == nil {
		ctx = context.Background()
	}
	return &GooseLogger{
		Log: log,
		Ctx: ctx,
	}
}

// Fatalf implements goose.Logger interface
func (g *GooseLogger) Fatalf(format string, v ...any) {
	g.Log.Fatal(g.Ctx, fmt.Sprintf(format, v...), "logType", "goose")
}

// Printf implements goose.Logger interface
func (g *GooseLogger) Printf(format string, v ...any) {
	g.Log.Info(g.Ctx, fmt.Sprintf(format, v...), "logType", "goose")
}

// Ensure GooseLogger implements goose.Logger interface
var _ goose.Logger = (*GooseLogger)(nil)
