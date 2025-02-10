// schema/schema.go
package schema

import (
    "time"
    "entgo.io/ent"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/edge"
    "entgo.io/contrib/entgql"
    "entgo.io/ent/schema"
)

// First, let's create a file for GraphQL configuration
// File: ent/entc.go
const gqlConfig = `
// +build ignore

package main

import (
    "log"
    "entgo.io/contrib/entgql"
    "entgo.io/ent/entc"
    "entgo.io/ent/entc/gen"
)

func main() {
    ex, err := entgql.NewExtension(
        entgql.WithWhereFilters(true),
        entgql.WithConfigPath("../gqlgen.yml"),
        entgql.WithSchemaPath("../ent.graphql"),
    )
    if err != nil {
        log.Fatalf("creating entgql extension: %v", err)
    }
    opts := []entc.Option{
        entc.Extensions(ex),
    }
    if err := entc.Generate("./schema", &gen.Config{}, opts...); err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}
`

// Define enum types
type Status string

const (
    StatusDraft     Status = "draft"
    StatusPublished Status = "published"
    StatusArchived  Status = "archived"
)

type Role string

const (
    RoleAuthor      Role = "author"
    RoleEditor      Role = "editor"
    RoleContributor Role = "contributor"
)

// User schema
type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id").
            Positive().
            Immutable().
            Annotations(
                entgql.OrderField("ID"),
            ),
        field.String("username").
            Unique().
            NotEmpty().
            Annotations(
                entgql.OrderField("USERNAME"),
            ),
        field.String("email").
            Unique().
            NotEmpty().
            Annotations(
                entgql.OrderField("EMAIL"),
            ),
        field.String("password_hash").
            Sensitive().
            NotEmpty().
            Annotations(
                entgql.Skip(),
            ),
        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entgql.OrderField("CREATED_AT"),
            ),
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now).
            Annotations(
                entgql.OrderField("UPDATED_AT"),
            ),
        field.Bool("is_active").
            Default(true).
            Annotations(
                entgql.OrderField("IS_ACTIVE"),
            ),
    }
}

func (User) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("authored_posts", Post.Type).
            Through("user_posts", UserPost.Type).
            Annotations(entgql.RelayConnection()),
        edge.To("comments", Comment.Type).
            Annotations(entgql.RelayConnection()),
        edge.To("likes", Post.Type).
            Through("user_likes", UserLike.Type).
            Annotations(entgql.RelayConnection()),
    }
}

func (User) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
        entgql.RelayConnection(),
    }
}

// Post schema
type Post struct {
    ent.Schema
}

func (Post) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id").
            Positive().
            Immutable().
            Annotations(
                entgql.OrderField("ID"),
            ),
        field.String("title").
            NotEmpty().
            Annotations(
                entgql.OrderField("TITLE"),
            ),
        field.Text("content").
            NotEmpty(),
        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entgql.OrderField("CREATED_AT"),
            ),
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now).
            Annotations(
                entgql.OrderField("UPDATED_AT"),
            ),
        field.String("status").
            GoType(Status("")).
            Default(string(StatusDraft)).
            Annotations(
                entgql.OrderField("STATUS"),
            ),
    }
}

func (Post) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("authors", User.Type).
            Ref("authored_posts").
            Through("user_posts", UserPost.Type).
            Annotations(entgql.RelayConnection()),
        edge.To("comments", Comment.Type).
            Annotations(entgql.RelayConnection()),
        edge.From("liked_by", User.Type).
            Ref("likes").
            Through("user_likes", UserLike.Type).
            Annotations(entgql.RelayConnection()),
    }
}

func (Post) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
        entgql.RelayConnection(),
    }
}

// UserPost schema
type UserPost struct {
    ent.Schema
}

func (UserPost) Fields() []ent.Field {
    return []ent.Field{
        field.Int("user_id").
            Positive(),
        field.Int("post_id").
            Positive(),
        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entgql.OrderField("CREATED_AT"),
            ),
        field.String("role").
            GoType(Role("")).
            Default(string(RoleAuthor)).
            Annotations(
                entgql.OrderField("ROLE"),
            ),
    }
}

func (UserPost) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).
            Unique().
            Required().
            Field("user_id"),
        edge.To("post", Post.Type).
            Unique().
            Required().
            Field("post_id"),
    }
}

func (UserPost) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
    }
}

// Comment schema
type Comment struct {
    ent.Schema
}

func (Comment) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id").
            Positive().
            Immutable().
            Annotations(
                entgql.OrderField("ID"),
            ),
        field.Text("content").
            NotEmpty().
            Annotations(
                entgql.OrderField("CONTENT"),
            ),
        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entgql.OrderField("CREATED_AT"),
            ),
        field.Time("updated_at").
            Default(time.Now).
            UpdateDefault(time.Now).
            Annotations(
                entgql.OrderField("UPDATED_AT"),
            ),
    }
}

func (Comment) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("author", User.Type).
            Ref("comments").
            Unique().
            Required(),
        edge.From("post", Post.Type).
            Ref("comments").
            Unique().
            Required(),
    }
}

func (Comment) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
        entgql.RelayConnection(),
    }
}

// UserLike schema
type UserLike struct {
    ent.Schema
}

func (UserLike) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id").
            Positive().
            Immutable().
            Annotations(
                entgql.OrderField("ID"),
            ),
        field.Int("user_id").
            Positive(),
        field.Int("post_id").
            Positive(),
        field.Time("created_at").
            Default(time.Now).
            Immutable().
            Annotations(
                entgql.OrderField("CREATED_AT"),
            ),
    }
}

func (UserLike) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("user", User.Type).
            Unique().
            Required().
            Field("user_id"),
        edge.To("post", Post.Type).
            Unique().
            Required().
            Field("post_id"),
    }
}

func (UserLike) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
        entgql.RelayConnection(),
    }
}