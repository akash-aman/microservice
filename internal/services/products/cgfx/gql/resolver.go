package gql

import (
	"products/cgfx/ent/gen"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct{ client *gen.Client }

// NewSchema creates a graphql executable schema.
func NewSchema(client *gen.Client) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{client},
	})
}
