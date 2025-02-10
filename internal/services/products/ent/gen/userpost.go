// Code generated by ent, DO NOT EDIT.

package gen

import (
	"fmt"
	"products/ent/gen/post"
	"products/ent/gen/user"
	"products/ent/gen/userpost"
	"products/ent/schema"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// UserPost is the model entity for the UserPost schema.
type UserPost struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID int `json:"user_id,omitempty"`
	// PostID holds the value of the "post_id" field.
	PostID int `json:"post_id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Role holds the value of the "role" field.
	Role schema.Role `json:"role,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserPostQuery when eager-loading is set.
	Edges        UserPostEdges `json:"edges"`
	selectValues sql.SelectValues
}

// UserPostEdges holds the relations/edges for other nodes in the graph.
type UserPostEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// Post holds the value of the post edge.
	Post *Post `json:"post,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
	// totalCount holds the count of the edges above.
	totalCount [2]map[string]int
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserPostEdges) UserOrErr() (*User, error) {
	if e.User != nil {
		return e.User, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "user"}
}

// PostOrErr returns the Post value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserPostEdges) PostOrErr() (*Post, error) {
	if e.Post != nil {
		return e.Post, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: post.Label}
	}
	return nil, &NotLoadedError{edge: "post"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserPost) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case userpost.FieldID, userpost.FieldUserID, userpost.FieldPostID:
			values[i] = new(sql.NullInt64)
		case userpost.FieldRole:
			values[i] = new(sql.NullString)
		case userpost.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserPost fields.
func (up *UserPost) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case userpost.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			up.ID = int(value.Int64)
		case userpost.FieldUserID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				up.UserID = int(value.Int64)
			}
		case userpost.FieldPostID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field post_id", values[i])
			} else if value.Valid {
				up.PostID = int(value.Int64)
			}
		case userpost.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				up.CreatedAt = value.Time
			}
		case userpost.FieldRole:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field role", values[i])
			} else if value.Valid {
				up.Role = schema.Role(value.String)
			}
		default:
			up.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the UserPost.
// This includes values selected through modifiers, order, etc.
func (up *UserPost) Value(name string) (ent.Value, error) {
	return up.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the UserPost entity.
func (up *UserPost) QueryUser() *UserQuery {
	return NewUserPostClient(up.config).QueryUser(up)
}

// QueryPost queries the "post" edge of the UserPost entity.
func (up *UserPost) QueryPost() *PostQuery {
	return NewUserPostClient(up.config).QueryPost(up)
}

// Update returns a builder for updating this UserPost.
// Note that you need to call UserPost.Unwrap() before calling this method if this UserPost
// was returned from a transaction, and the transaction was committed or rolled back.
func (up *UserPost) Update() *UserPostUpdateOne {
	return NewUserPostClient(up.config).UpdateOne(up)
}

// Unwrap unwraps the UserPost entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (up *UserPost) Unwrap() *UserPost {
	_tx, ok := up.config.driver.(*txDriver)
	if !ok {
		panic("gen: UserPost is not a transactional entity")
	}
	up.config.driver = _tx.drv
	return up
}

// String implements the fmt.Stringer.
func (up *UserPost) String() string {
	var builder strings.Builder
	builder.WriteString("UserPost(")
	builder.WriteString(fmt.Sprintf("id=%v, ", up.ID))
	builder.WriteString("user_id=")
	builder.WriteString(fmt.Sprintf("%v", up.UserID))
	builder.WriteString(", ")
	builder.WriteString("post_id=")
	builder.WriteString(fmt.Sprintf("%v", up.PostID))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(up.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("role=")
	builder.WriteString(fmt.Sprintf("%v", up.Role))
	builder.WriteByte(')')
	return builder.String()
}

// UserPosts is a parsable slice of UserPost.
type UserPosts []*UserPost
