package server

import (
	"context"
	"errors"
	"net/http"
	"pkg/http/server"
	"pkg/logger"
	"products/app/inits"
	"products/cgfx/ent/gen"
	"products/conf"

	// "entgo.io/contrib/entgql"
	// "github.com/99designs/gqlgen/graphql/handler"
	// "github.com/99designs/gqlgen/graphql/handler/debug"
	// "github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"

	"go.uber.org/fx"
)

func RunServers(lc fx.Lifecycle, e *echo.Echo, client *gen.Client, log logger.ILogger, config *conf.Config, ctx context.Context) {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {

			/**
			 * Http Server.
			 */
			go func() {

				log.Infof("Starting echo server on port %v", config.Echo.Port)

				if err := server.RunEchoServer(ctx, e, log, config.Echo); !errors.Is(err, http.ErrServerClosed) {
					log.Errorf("Error starting echo server %s", err)
				}
			}()

			/**
			 * GraphQL Server
			 */
			go func() {
				if err := inits.InitGraphQLServer(ctx, client, log, config.GraphQL); !errors.Is(err, http.ErrServerClosed) {
					log.Errorf("Error starting GraphQL server: %s", err)
				}
			}()

			/**
			 * Migration
			 */
			go func() {
				if err := client.Schema.Create(ctx); err != nil {
					log.Fatalf("failed creating schema resources: %v", err)
				}
			}()

			/**
			 * Service Route.
			 */
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, config.Service.Name)
			})

			return nil
		},

		OnStop: func(stopCtx context.Context) error {

			if err := e.Shutdown(stopCtx); err != nil {
				log.Errorf("error shutting down HTTP server: %v", err)
			}

			log.Info("All servers shut down gracefully")

			log.Sync()

			return nil
		},
	})
}
