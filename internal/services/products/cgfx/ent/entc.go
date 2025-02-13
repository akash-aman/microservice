//go:build ignore

/**
 * Ref: https://github.com/ent/contrib/blob/master/entgql/internal/todo/ent/entc.go
 */
package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("./cgfx/gql/ent.graphql"),
		entgql.WithConfigPath("gqlgen.yml"),
		entgql.WithWhereInputs(true),
		entgql.WithRelaySpec(true),
		entgql.WithNodeDescriptor(true),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	if err := entc.Generate("./cgfx/ent/schema", &gen.Config{
		Target:  "./cgfx/ent/gen",
		Package: "products/cgfx/ent/gen",
		Features: []gen.Feature{
			gen.FeatureModifier,
		},
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
