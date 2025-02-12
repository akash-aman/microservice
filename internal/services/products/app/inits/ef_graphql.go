package inits

import (
	"net/http"
	"pkg/logger"
	"products/cgfx/ent/gen"
	"products/cgfx/gql"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func InitGraphQLServer(client *gen.Client, log logger.ILogger) {
	srv := handler.NewDefaultServer(gql.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	http.Handle("/",
		playground.Handler("Todo", "/query"),
	)
	http.Handle("/query", srv)

	log.Infof("GraphQL Server listening on 3001")
	server := &http.Server{
		Addr:              ":3001",
		ReadHeaderTimeout: 30 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Errorf("http server terminated %v", err)
	}
}
