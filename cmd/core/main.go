package main

import (
	"log"

	"go.uber.org/fx"

	"github.com/Mind2Screen-Dev-Team/thousand-sunny/bootstrap"
	"github.com/Mind2Screen-Dev-Team/thousand-sunny/bootstrap/dependency"
)

func main() {
	var (
		pkgs = bootstrap.App()
		apps = fx.New(pkgs...)
	)

	// Run and Blocking
	apps.Run()

	// custom shutdown capability for logger
	defer func() {
		if dependency.DebugLogger != nil {
			if err := dependency.DebugLogger.LogRotation.Rotate(); err != nil {
				log.Printf("app 'debug-logger' rotation, got err: %+v\n", err)
			}
		}

		if dependency.IOLogger != nil {
			if err := dependency.IOLogger.LogRotation.Rotate(); err != nil {
				log.Printf("app 'io-logger' rotation, got err: %+v\n", err)
			}
		}

		if dependency.TRXLogger != nil {
			if err := dependency.TRXLogger.Rotate(); err != nil {
				log.Printf("app 'trx-logger' rotation, got err: %+v\n", err)
			}

			// clear memory trx logger
			dependency.TRXLogger.Close()
		}
	}()
}
