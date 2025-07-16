package user

import (
	"go.uber.org/fx"
)

var (
	ExampleUserModules = fx.Options(
		ExampleUserCreateHandlerModuleFx,
		ExampleUserReadHandlerModuleFx,
		ExampleUserReadAllHandlerModuleFx,
		ExampleUserUpdateHandlerModuleFx,
		ExampleUserDeleteHandlerModuleFx,
	)
)
