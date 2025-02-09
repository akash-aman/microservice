package main

import (
	http "pkg/http"
	httpServer "pkg/http/server"
	"pkg/logger"
	"products/conf"
	"products/server"
	"products/app/inits"
	"github.com/go-playground/validator"
	"go.uber.org/fx"
)

/**
 * Main function.
 *
 * This function is the entry point for the application. It is responsible
 * for initializing the application and starting the server.
 */
func main() {
	fx.New(
		fx.Provide(
			logger.InitLogger,
			conf.InitConfig,
			http.NewContext,
			validator.New,
			httpServer.NewEchoServer,
		),
		fx.Invoke(server.RunServers),
		fx.Invoke(inits.InitMediator),
		fx.Invoke(inits.ConfigEndpoints),
	).Run()
}
