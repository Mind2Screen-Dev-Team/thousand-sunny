package bootstrap

import (
	"go.uber.org/fx"
)

type (
	Registry struct{}
	Provider = []fx.Option
)

func App() Provider {
	var (
		p  Registry
		ps Provider
	)

	// # App Dependency
	ps = append(ps, p.DependencyProvider()...)

	// # App Global Middleware
	ps = append(ps, p.GlobalHTTPMiddlewareProvider()...)

	// # App Start-Up List
	ps = append(ps, p.DependencyStartUp()...)

	return ps
}
