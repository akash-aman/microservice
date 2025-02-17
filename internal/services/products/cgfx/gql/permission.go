package gql

import (
	"context"
	"errors"

	otel "pkg/otel/tracer"
	"products/cgfx/ent/gen"

	"github.com/99designs/gqlgen/graphql"
)

/**
 * Runs for every entry from the result.
 */
func HasPermission(client *gen.Client) func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []string) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []string) (res interface{}, err error) {
		// Implement your permission logic here using the client and context
		// For example, you can extract user information from the context and check permissions
		tracer, ctx := otel.NewTracer(ctx, "HasPermission Controller")
		defer tracer.End()

		// Your Conditional Logic To Authorize
		if true {
			tracer.RecordError(errors.New("unauthorized"), "Access Denied")
			return nil, errors.New("unauthorized")
		}

		return next(ctx)
	}
}

/**
 * Runs before query.
 */
func HasRole(client *gen.Client) func(ctx context.Context, obj interface{}, next graphql.Resolver, roles string) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles string) (res interface{}, err error) {
		// Implement your permission logic here using the client and context
		// For example, you can extract user information from the context and check permissions
		tracer, ctx := otel.NewTracer(ctx, "HasRole Controller")
		defer tracer.End()

		// Your Conditional Logic To Authorize
		if true {
			tracer.RecordError(errors.New("unauthorized"), "Access Denied")
			return nil, errors.New("unauthorized")
		}

		return next(ctx)
	}
}
