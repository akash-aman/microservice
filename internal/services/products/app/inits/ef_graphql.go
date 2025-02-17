package inits

/**
 * https://github.com/cmelgarejo/go-gql-server/blob/master/internal/handlers/gql.go
 */
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

	server.Addr = conf.Port
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

	log.Info(ctx, "Starting GraphQL Server", zap.String("port", conf.Port))

	return server.ListenAndServe()
}
