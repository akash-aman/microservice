package inits

import (
	"context"
	"errors"
	"net/http"
	conf "pkg/gql"
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
)

func InitGraphQLServer(ctx context.Context, client *gen.Client, log logger.ILogger, conf *conf.GraphQLConfig, server *http.Server) error {

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

	server.Addr = conf.Port
	server.Handler = http.DefaultServeMux

	go func() {
		<-ctx.Done()
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Errorf("Error shutting down GraphQL server: %v", err)
		} else {
			log.Infof("GraphQL server shut down gracefully")
		}
	}()

	log.Infof("Starting GraphQL Server on port %s", conf.Port)

	return server.ListenAndServe()
}
