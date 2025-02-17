package server

import (
	"context"
	"errors"
	"net/http"
	"pkg/http/server"
	"pkg/logger"
	"pkg/otel/metrics"
	"products/app/inits"
	"products/cgfx/ent/gen"
	"products/conf"

	metricsdk "go.opentelemetry.io/otel/sdk/metric"

	// "entgo.io/contrib/entgql"
	// "github.com/99designs/gqlgen/graphql/handler"
	// "github.com/99designs/gqlgen/graphql/handler/debug"
	// "github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RunServers(lc fx.Lifecycle, e *echo.Echo, client *gen.Client, log logger.Zapper, config *conf.Config, gqlsrv *http.Server, provider *metricsdk.MeterProvider, ctx context.Context) {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {

			/**
			 * Http Server
			 */
			go func() {

				log.Info(ctx, "starting echo server", zap.String("port", config.Echo.Port))

				if err := server.RunEchoServer(ctx, e, log, config.Echo); !errors.Is(err, http.ErrServerClosed) {
					log.Error(ctx, "error starting echo server", zap.Error(err))
				}
			}()

			/**
			 * GraphQL Server
			 */
			go func() {
				if err := inits.InitGraphQLServer(ctx, client, log, config.GraphQL, gqlsrv); !errors.Is(err, http.ErrServerClosed) {
					log.Error(ctx, "Error starting GraphQL server", zap.Error(err))
				}
			}()

			/**
			 * Migration
			 */
			go func() {
				if err := client.Schema.Create(ctx); err != nil {
					log.Fatal(ctx, "failed creating schema resources", zap.Error(err))
				}
			}()

			/**
			 * Service Route
			 */
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, config.Service.Name)
			})

			/**
			 * Setup Metrics Collector
			 */
			meter := provider.Meter(config.Otel.Service)
			metrics.GenerateMetrics(ctx, meter, log)

			return nil
		},

		OnStop: func(stopCtx context.Context) error {

			if err := e.Shutdown(stopCtx); err != nil {
				log.Error(ctx, "error shutting down HTTP server", zap.Error(err))
			} else {
				log.Info(ctx, "HTTP server shut down gracefully")
			}

			if err := gqlsrv.Shutdown(stopCtx); err != nil {
				log.Error(ctx, "error shutting down GraphQL server", zap.Error(err))
			} else {
				log.Info(ctx, "GraphQL server shut down gracefully")
			}

			log.Info(ctx, "All servers shut down gracefully")

			log.Sync()

			return nil
		},
	})
}
