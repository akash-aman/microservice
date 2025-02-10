// Code generated by ent, DO NOT EDIT.

package userpost

import (
	"products/ent/gen/predicate"
	"products/ent/schema"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.UserPost {
	return predicate.UserPost(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.UserPost {
	return predicate.UserPost(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.UserPost {
	return predicate.UserPost(sql.FieldLTE(FieldID, id))
}

// UserID applies equality check predicate on the "user_id" field. It's identical to UserIDEQ.
func UserID(v int) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldUserID, v))
}

// PostID applies equality check predicate on the "post_id" field. It's identical to PostIDEQ.
func PostID(v int) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldPostID, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldCreatedAt, v))
}

// Role applies equality check predicate on the "role" field. It's identical to RoleEQ.
func Role(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldEQ(FieldRole, vc))
}

// UserIDEQ applies the EQ predicate on the "user_id" field.
func UserIDEQ(v int) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldUserID, v))
}

// UserIDNEQ applies the NEQ predicate on the "user_id" field.
func UserIDNEQ(v int) predicate.UserPost {
	return predicate.UserPost(sql.FieldNEQ(FieldUserID, v))
}

// UserIDIn applies the In predicate on the "user_id" field.
func UserIDIn(vs ...int) predicate.UserPost {
	return predicate.UserPost(sql.FieldIn(FieldUserID, vs...))
}

// UserIDNotIn applies the NotIn predicate on the "user_id" field.
func UserIDNotIn(vs ...int) predicate.UserPost {
	return predicate.UserPost(sql.FieldNotIn(FieldUserID, vs...))
}

// PostIDEQ applies the EQ predicate on the "post_id" field.
func PostIDEQ(v int) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldPostID, v))
}

// PostIDNEQ applies the NEQ predicate on the "post_id" field.
func PostIDNEQ(v int) predicate.UserPost {
	return predicate.UserPost(sql.FieldNEQ(FieldPostID, v))
}

// PostIDIn applies the In predicate on the "post_id" field.
func PostIDIn(vs ...int) predicate.UserPost {
	return predicate.UserPost(sql.FieldIn(FieldPostID, vs...))
}

// PostIDNotIn applies the NotIn predicate on the "post_id" field.
func PostIDNotIn(vs ...int) predicate.UserPost {
	return predicate.UserPost(sql.FieldNotIn(FieldPostID, vs...))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.UserPost {
	return predicate.UserPost(sql.FieldLTE(FieldCreatedAt, v))
}

// RoleEQ applies the EQ predicate on the "role" field.
func RoleEQ(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldEQ(FieldRole, vc))
}

// RoleNEQ applies the NEQ predicate on the "role" field.
func RoleNEQ(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldNEQ(FieldRole, vc))
}

// RoleIn applies the In predicate on the "role" field.
func RoleIn(vs ...schema.Role) predicate.UserPost {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.UserPost(sql.FieldIn(FieldRole, v...))
}

// RoleNotIn applies the NotIn predicate on the "role" field.
func RoleNotIn(vs ...schema.Role) predicate.UserPost {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.UserPost(sql.FieldNotIn(FieldRole, v...))
}

// RoleGT applies the GT predicate on the "role" field.
func RoleGT(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldGT(FieldRole, vc))
}

// RoleGTE applies the GTE predicate on the "role" field.
func RoleGTE(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldGTE(FieldRole, vc))
}

// RoleLT applies the LT predicate on the "role" field.
func RoleLT(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldLT(FieldRole, vc))
}

// RoleLTE applies the LTE predicate on the "role" field.
func RoleLTE(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldLTE(FieldRole, vc))
}

// RoleContains applies the Contains predicate on the "role" field.
func RoleContains(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldContains(FieldRole, vc))
}

// RoleHasPrefix applies the HasPrefix predicate on the "role" field.
func RoleHasPrefix(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldHasPrefix(FieldRole, vc))
}

// RoleHasSuffix applies the HasSuffix predicate on the "role" field.
func RoleHasSuffix(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldHasSuffix(FieldRole, vc))
}

// RoleEqualFold applies the EqualFold predicate on the "role" field.
func RoleEqualFold(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldEqualFold(FieldRole, vc))
}

// RoleContainsFold applies the ContainsFold predicate on the "role" field.
func RoleContainsFold(v schema.Role) predicate.UserPost {
	vc := string(v)
	return predicate.UserPost(sql.FieldContainsFold(FieldRole, vc))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.UserPost {
	return predicate.UserPost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.UserPost {
	return predicate.UserPost(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPost applies the HasEdge predicate on the "post" edge.
func HasPost() predicate.UserPost {
	return predicate.UserPost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, PostTable, PostColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPostWith applies the HasEdge predicate on the "post" edge with a given conditions (other predicates).
func HasPostWith(preds ...predicate.Post) predicate.UserPost {
	return predicate.UserPost(func(s *sql.Selector) {
		step := newPostStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.UserPost) predicate.UserPost {
	return predicate.UserPost(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.UserPost) predicate.UserPost {
	return predicate.UserPost(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.UserPost) predicate.UserPost {
	return predicate.UserPost(sql.NotPredicates(p))
}
