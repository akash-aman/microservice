package gql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"products/cgfx/ent/gen"

	"entgo.io/contrib/entgql"
)

type Resolver struct{}

// Status is the resolver for the status field.
func (r *postResolver) Status(ctx context.Context, obj *gen.Post) (string, error) {
	panic("not implemented")
}

// UserPosts is the resolver for the userPosts field.
func (r *postResolver) UserPosts(ctx context.Context, obj *gen.Post) ([]*gen.UserPost, error) {
	panic("not implemented")
}

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id int) (gen.Noder, error) {
	panic("not implemented")
}

// Nodes is the resolver for the nodes field.
func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]gen.Noder, error) {
	panic("not implemented")
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, orderBy *gen.CommentOrder) (*gen.CommentConnection, error) {
	panic("not implemented")
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, orderBy *gen.PostOrder) (*gen.PostConnection, error) {
	panic("not implemented")
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, orderBy *gen.UserOrder) (*gen.UserConnection, error) {
	panic("not implemented")
}

// UserLikes is the resolver for the userLikes field.
func (r *queryResolver) UserLikes(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, orderBy *gen.UserLikeOrder) (*gen.UserLikeConnection, error) {
	panic("not implemented")
}

// UserPosts is the resolver for the userPosts field.
func (r *queryResolver) UserPosts(ctx context.Context) ([]*gen.UserPost, error) {
	panic("not implemented")
}

// UserPosts is the resolver for the userPosts field.
func (r *userResolver) UserPosts(ctx context.Context, obj *gen.User) ([]*gen.UserPost, error) {
	panic("not implemented")
}

// Role is the resolver for the role field.
func (r *userPostResolver) Role(ctx context.Context, obj *gen.UserPost) (string, error) {
	panic("not implemented")
}

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

// UserPost returns UserPostResolver implementation.
func (r *Resolver) UserPost() UserPostResolver { return &userPostResolver{r} }

type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userPostResolver struct{ *Resolver }
