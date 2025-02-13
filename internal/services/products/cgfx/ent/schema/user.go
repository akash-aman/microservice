package schema

import (
	"products/cgfx/ent/schema/annotation"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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

// User schema
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Immutable().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.String("firstname").
			MinLen(4).
			MaxLen(20).
			Annotations(
				entgql.OrderField("FIRSTNAME"),
			),
		field.String("lastname").
			MinLen(4).
			MaxLen(20).
			Annotations(
				entgql.OrderField("LASTNAME"),
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
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("orders", Order.Type).
			Annotations(entgql.RelayConnection()),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField().Description("This is the single user"),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entgql.MultiOrder(),
	}
}

type Order struct {
	ent.Schema
}

func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Immutable().
			Annotations(
				entgql.OrderField("ID"),
			),
		field.Enum("status").
			Values("pending", "completed", "cancelled").
			Default("pending").
			Annotations(
				entgql.OrderField("STATUS"),
			),
		field.Float("total").
			Positive().
			Annotations(
				entgql.OrderField("TOTAL"),
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

func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("orders").
			Unique(),
	}
}

func (Order) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField().Description("This is the order item"),
		entgql.Directives(
			annotation.HasPermissions([]string{"ADMIN", "MODERATOR"}),
		),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entgql.MultiOrder(),
	}
}
