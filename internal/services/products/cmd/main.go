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
	"go.uber.org/zap"
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
			logger.InitLogger[zap.Field],
			conf.InitConfig,
			http.NewContext,
			validator.New,
			db.NewOtelDBConnectionPool,
			httpServer.NewEchoServer,
			inits.NewEntClient,
			gql.NewGQLServer,
			otel.InitOpentelemetry,
		),
		fx.Invoke(server.RunServers),
		fx.Invoke(inits.InitMediator),
		fx.Invoke(inits.ConfigEndpoints),
		fx.Invoke(inits.ConfigSwagger),
		fx.Invoke(inits.ConfigMiddlewares),
	).Run()
}
