package user

import (
	"go.uber.org/fx"
)

var (
	ExampleUserHandlerModules = fx.Options(
		ExampleUserCreateHandlerModuleFx,
		ExampleUserReadHandlerModuleFx,
		ExampleUserReadAllHandlerModuleFx,
		ExampleUserUpdateHandlerModuleFx,
		ExampleUserDeleteHandlerModuleFx,
	)
)
