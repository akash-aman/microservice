package gql

import (
	"products/cgfx/ent/gen"

	"github.com/99designs/gqlgen/graphql"
)

// Resolver is the resolver root.
type Resolver struct{ Client *gen.Client }

// NewSchema creates a graphql executable schema.
func NewSchema(client *gen.Client) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{Client: client},
		Directives: DirectiveRoot{
			HasPermissions: HasPermission(client),
			HasRole:        HasRole(client),
		},
	})
}
