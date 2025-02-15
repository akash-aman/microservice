package main

import (
	"pkg/db"
	"pkg/gql"
	http "pkg/http"
	httpServer "pkg/http/server"
	"pkg/logger"
	"pkg/otel"
	"products/app/inits"
	"products/conf"
	"products/server"

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
			db.NewOtelDBConnectionPool,
			httpServer.NewEchoServer,
			inits.NewEntClient,
			gql.NewGQLServer,
			otel.InitTracer,
		),
		fx.Invoke(server.RunServers),
		fx.Invoke(inits.InitMediator),
		fx.Invoke(inits.ConfigEndpoints),
	).Run()
}
