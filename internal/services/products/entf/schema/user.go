// schema/user.go
package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

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

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
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

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("authored_posts", Post.Type).
			Through("user_posts", UserPost.Type),
		edge.To("comments", Comment.Type),
		edge.To("likes", Post.Type).
			Through("user_likes", UserLike.Type),
	}
}

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
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

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("authors", User.Type).
			Ref("authored_posts").
			Through("user_posts", UserPost.Type),
		edge.To("comments", Comment.Type),
		edge.From("liked_by", User.Type).
			Ref("likes").
			Through("user_likes", UserLike.Type),
	}
}

// UserPost holds the schema definition for the UserPost entity.
type UserPost struct {
	ent.Schema
}

// Fields of the UserPost.
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

// Edges of the UserPost.
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

// Comment holds the schema definition for the Comment entity.
type Comment struct {
	ent.Schema
}

// Fields of the Comment.
func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable(),
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
	}
}

// Edges of the Comment.
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

// UserLike holds the schema definition for the UserLike entity.
type UserLike struct {
	ent.Schema
}

// Fields of the UserLike.
func (UserLike) Fields() []ent.Field {
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
	}
}

// Edges of the UserLike.
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
