//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		// Generate a GraphQL schema for the Ent schema
		// and save it as "ent.graphql".
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("./gql/ent.graphql"),
	)
	if err != nil {
		log.Fatalf("failed to create entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	if err := entc.Generate("./ent/schema", &gen.Config{
		Target:  "./ent/gen",
		Package: "products/ent/gen",
	}, opts...); err != nil {
		log.Fatalf("failed to run ent codegen: %v", err)
	}
}
