package xlog

import (
	"strings"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

func SetupFxLogger(logger *DebugLogger) fxevent.Logger {
	return &FxZeroLogger{DebugLog: logger}
}

// ZeroLogger is an Fx event logger that logs events to Zero.
type FxZeroLogger struct {
	DebugLog *DebugLogger
}

// LogEvent logs the given event to the provided Zap logger.
func (l *FxZeroLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.DebugLog.Logger.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStart hook executing")
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStart hook failed")
			return
		}

		l.DebugLog.Logger.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Str("runtime", e.Runtime.String()).
			Msg("OnStart hook executed")
	case *fxevent.OnStopExecuting:
		l.DebugLog.Logger.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Msg("OnStop hook executing")
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Str("callee", e.FunctionName).
				Str("caller", e.CallerName).
				Msg("OnStop hook failed")
			return
		}

		l.DebugLog.Logger.Debug().
			Str("callee", e.FunctionName).
			Str("caller", e.CallerName).
			Str("runtime", e.Runtime.String()).
			Msg("OnStart hook executed")
	case *fxevent.Supplied:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Str("type", e.TypeName).
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Msg("Debug encountered while applying options")
			return
		}

		l.DebugLog.Logger.Debug().
			Str("type", e.TypeName).
			Strs("stack.trace", e.StackTrace).
			Strs("module.trace", e.ModuleTrace).
			Func(moduleField(e.ModuleName)).
			Msg("supplied")
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.DebugLog.Logger.Debug().
				Str("constructor", e.ConstructorName).
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Str("type", rtype).
				Func(maybeBool("private", e.Private)).
				Msg("provided")
		}

		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Msg("Debug encountered while applying options")
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.DebugLog.Logger.Debug().
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Str("type", rtype).
				Msg("replaced")
		}

		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Msg("Debug encountered while replacing")
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.DebugLog.Logger.Debug().
				Str("decorator", e.DecoratorName).
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Str("type", rtype).
				Msg("decorated")
		}

		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Strs("stack.trace", e.StackTrace).
				Strs("module.trace", e.ModuleTrace).
				Func(moduleField(e.ModuleName)).
				Msg("Debug encountered while applying options")
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Str("name", e.Name).
				Str("kind", e.Kind).
				Func(moduleField(e.ModuleName)).
				Msg("Debug returned")
			return
		}

		l.DebugLog.Logger.Debug().
			Str("name", e.Name).
			Str("kind", e.Kind).
			Func(moduleField(e.ModuleName)).
			Msg("run")
	case *fxevent.Invoking:
		// do not log stack as it will make logs hard to read
		l.DebugLog.Logger.Debug().
			Str("function", e.FunctionName).
			Func(moduleField(e.ModuleName)).
			Msg("invoking")
	case *fxevent.Invoked:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().
				Err(e.Err).
				Str("function", e.FunctionName).
				Str("stack", e.Trace).
				Func(moduleField(e.ModuleName)).
				Msg("invoke failed")
		}
	case *fxevent.Stopping:
		l.DebugLog.Logger.Info().
			Str("signal", strings.ToUpper(e.Signal.String())).
			Msg("received signal")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().Err(e.Err).Msg("stop failed")
		}
	case *fxevent.RollingBack:
		l.DebugLog.Logger.Debug().Err(e.StartErr).Msg("start failed, rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().Err(e.Err).Msg("rollback failed")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().Err(e.Err).Msg("start failed")
			return
		}

		l.DebugLog.Logger.Debug().Msg("started")
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.DebugLog.Logger.Debug().Err(e.Err).Msg("custom logger initialization failed")
			return
		}

		l.DebugLog.Logger.Debug().Str("function", e.ConstructorName).Msg("initialized custom fxevent.Logger")
	}
}

func moduleField(value string) func(*zerolog.Event) {
	return func(e *zerolog.Event) {
		if value != "" {
			e.Str("module", value)
		}
	}
}

func maybeBool(key string, value bool) func(*zerolog.Event) {
	return func(e *zerolog.Event) {
		if value {
			e.Bool(key, value)
		}
	}
}
