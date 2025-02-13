package gql

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
)

func HasPermission() func(context.Context, interface{}, graphql.Resolver, []string) (interface{}, error) {
	return func(
		ctx context.Context,
		obj interface{},
		next graphql.Resolver,
		permissions []string,
	) (res interface{}, err error) {
		// you can do your thing here for permissions
		// if permissions are not met, return an error
		if !checkPermissions(ctx, permissions) {
			return nil, errors.New("unauthorized")
		}

		return next(ctx)
	}
}

// checkPermissions is a placeholder function to check permissions
func checkPermissions(ctx context.Context, permissions []string) bool {
	// implement your permission checking logic here
	return false
}
