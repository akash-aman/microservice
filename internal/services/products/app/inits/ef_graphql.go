package inits

import (
	"errors"
	"net/http"
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
)

func InitGraphQLServer(client *gen.Client, log logger.ILogger) {
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

	http.Handle("/", playground.Handler("Todo", "/query"))
	http.Handle("/query", srv)

	log.Infof("Starting GraphQL Server on port :3001")

	if err := http.ListenAndServe(":30001", nil); !errors.Is(err, http.ErrServerClosed) {
		log.Error("Error starting GraphQL server", err)
	}
}
