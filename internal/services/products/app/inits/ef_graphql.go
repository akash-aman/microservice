package inits

/**
 * https://github.com/cmelgarejo/go-gql-server/blob/master/internal/handlers/gql.go
 */
import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"pkg/discovery"
	conf "pkg/gql"
	"pkg/helper"
	"pkg/logger"
	"products/cgfx/ent/gen"
	"products/cgfx/gql"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

func InitGraphQLServer(ctx context.Context, client *gen.Client, log logger.Zapper, conf *conf.GraphQLConfig, server *http.Server) error {

	if conf == nil {
		return errors.New("graphQL config not loaded properly")
	}

	srv := handler.New(gql.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	srv.AddTransport(transport.SSE{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.Use(extension.Introspection{})

	http.Handle("/graphql", otelhttp.NewHandler(playground.Handler("Products", "/query"), "Http Query Handler Endpoint"))
	http.Handle("/query", otelhttp.NewHandler(srv, "GraphQL Query Handler Endpoint"))
	http.Handle("/health", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}), "Health Check Endpoint"))

	server.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	server.Handler = http.DefaultServeMux

	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Error(ctx, "Error shutting down GraphQL server", zap.Error(err))
		} else {
			log.Info(ctx, "GraphQL server shut down gracefully")
		}
	}()

	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	log.Infof(ctx, "GraphQL Server Listening on %s", addr)

	if err := discovery.RegisterServiceWithConsul(
		ctx,
		"echo-graphql-service",
		fmt.Sprintf("echo-graphql-service-%s", helper.GetMachineID()),
		conf.Host,
		conf.Port,
		discovery.HTTPService,
		log,
	); err != nil {
		log.Errorf(ctx, "Error registering with Consul: %v", err)
	}

	return server.ListenAndServe()
}
