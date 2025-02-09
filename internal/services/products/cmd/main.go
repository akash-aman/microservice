package main

import (
	http "pkg/http"
	httpServer "pkg/http/server"
	"pkg/logger"
	"products/conf"
	"products/server"

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
			httpServer.NewEchoServer,
		),
		fx.Invoke(server.RunServers),
	).Run()
}
