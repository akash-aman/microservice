//go:build ignore

/**
 * Ref: https://github.com/ent/contrib/blob/master/entgql/internal/todo/ent/entc.go
 */
package main

import (
	"log"
	"reflect"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

const (
	gqlgenConfigPath  = "gqlgen.yml"
	entitySchemaPath  = "./cgfx/ent/schema"
	targetOutputPath  = "./cgfx/ent/gen"
	graphQLSchemaPath = "./cgfx/gql/ent.graphql"
	outputPackageName = "products/cgfx/ent/gen"
)

func main() {
	rt := reflect.TypeOf(uuid.UUID{})
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath(graphQLSchemaPath),
		entgql.WithConfigPath(gqlgenConfigPath),
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
	if err := entc.Generate(entitySchemaPath, &gen.Config{
		Templates: entgql.AllTemplates,
		Target:    targetOutputPath,
		Package:   outputPackageName,
		IDType: &field.TypeInfo{
			Type:    field.TypeUUID,
			Ident:   rt.String(),
			PkgPath: rt.PkgPath(),
		},

		Features: []gen.Feature{
			gen.FeatureModifier,
		},
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
