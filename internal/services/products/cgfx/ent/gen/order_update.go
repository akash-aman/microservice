// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"errors"
	"fmt"
	"products/cgfx/ent/gen/order"
	"products/cgfx/ent/gen/predicate"
	"products/cgfx/ent/gen/user"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// OrderUpdate is the builder for updating Order entities.
type OrderUpdate struct {
	config
	hooks     []Hook
	mutation  *OrderMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the OrderUpdate builder.
func (ou *OrderUpdate) Where(ps ...predicate.Order) *OrderUpdate {
	ou.mutation.Where(ps...)
	return ou
}

// SetStatus sets the "status" field.
func (ou *OrderUpdate) SetStatus(o order.Status) *OrderUpdate {
	ou.mutation.SetStatus(o)
	return ou
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (ou *OrderUpdate) SetNillableStatus(o *order.Status) *OrderUpdate {
	if o != nil {
		ou.SetStatus(*o)
	}
	return ou
}

// SetTotal sets the "total" field.
func (ou *OrderUpdate) SetTotal(f float64) *OrderUpdate {
	ou.mutation.ResetTotal()
	ou.mutation.SetTotal(f)
	return ou
}

// SetNillableTotal sets the "total" field if the given value is not nil.
func (ou *OrderUpdate) SetNillableTotal(f *float64) *OrderUpdate {
	if f != nil {
		ou.SetTotal(*f)
	}
	return ou
}

// AddTotal adds f to the "total" field.
func (ou *OrderUpdate) AddTotal(f float64) *OrderUpdate {
	ou.mutation.AddTotal(f)
	return ou
}

// SetUpdatedAt sets the "updated_at" field.
func (ou *OrderUpdate) SetUpdatedAt(t time.Time) *OrderUpdate {
	ou.mutation.SetUpdatedAt(t)
	return ou
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ou *OrderUpdate) SetUserID(id uuid.UUID) *OrderUpdate {
	ou.mutation.SetUserID(id)
	return ou
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ou *OrderUpdate) SetNillableUserID(id *uuid.UUID) *OrderUpdate {
	if id != nil {
		ou = ou.SetUserID(*id)
	}
	return ou
}

// SetUser sets the "user" edge to the User entity.
func (ou *OrderUpdate) SetUser(u *User) *OrderUpdate {
	return ou.SetUserID(u.ID)
}

// Mutation returns the OrderMutation object of the builder.
func (ou *OrderUpdate) Mutation() *OrderMutation {
	return ou.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ou *OrderUpdate) ClearUser() *OrderUpdate {
	ou.mutation.ClearUser()
	return ou
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ou *OrderUpdate) Save(ctx context.Context) (int, error) {
	ou.defaults()
	return withHooks(ctx, ou.sqlSave, ou.mutation, ou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ou *OrderUpdate) SaveX(ctx context.Context) int {
	affected, err := ou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ou *OrderUpdate) Exec(ctx context.Context) error {
	_, err := ou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ou *OrderUpdate) ExecX(ctx context.Context) {
	if err := ou.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ou *OrderUpdate) defaults() {
	if _, ok := ou.mutation.UpdatedAt(); !ok {
		v := order.UpdateDefaultUpdatedAt()
		ou.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ou *OrderUpdate) check() error {
	if v, ok := ou.mutation.Status(); ok {
		if err := order.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`gen: validator failed for field "Order.status": %w`, err)}
		}
	}
	if v, ok := ou.mutation.Total(); ok {
		if err := order.TotalValidator(v); err != nil {
			return &ValidationError{Name: "total", err: fmt.Errorf(`gen: validator failed for field "Order.total": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ou *OrderUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *OrderUpdate {
	ou.modifiers = append(ou.modifiers, modifiers...)
	return ou
}

func (ou *OrderUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ou.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(order.Table, order.Columns, sqlgraph.NewFieldSpec(order.FieldID, field.TypeUUID))
	if ps := ou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ou.mutation.Status(); ok {
		_spec.SetField(order.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := ou.mutation.Total(); ok {
		_spec.SetField(order.FieldTotal, field.TypeFloat64, value)
	}
	if value, ok := ou.mutation.AddedTotal(); ok {
		_spec.AddField(order.FieldTotal, field.TypeFloat64, value)
	}
	if value, ok := ou.mutation.UpdatedAt(); ok {
		_spec.SetField(order.FieldUpdatedAt, field.TypeTime, value)
	}
	if ou.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   order.UserTable,
			Columns: []string{order.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ou.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   order.UserTable,
			Columns: []string{order.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ou.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, ou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{order.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ou.mutation.done = true
	return n, nil
}

// OrderUpdateOne is the builder for updating a single Order entity.
type OrderUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *OrderMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetStatus sets the "status" field.
func (ouo *OrderUpdateOne) SetStatus(o order.Status) *OrderUpdateOne {
	ouo.mutation.SetStatus(o)
	return ouo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (ouo *OrderUpdateOne) SetNillableStatus(o *order.Status) *OrderUpdateOne {
	if o != nil {
		ouo.SetStatus(*o)
	}
	return ouo
}

// SetTotal sets the "total" field.
func (ouo *OrderUpdateOne) SetTotal(f float64) *OrderUpdateOne {
	ouo.mutation.ResetTotal()
	ouo.mutation.SetTotal(f)
	return ouo
}

// SetNillableTotal sets the "total" field if the given value is not nil.
func (ouo *OrderUpdateOne) SetNillableTotal(f *float64) *OrderUpdateOne {
	if f != nil {
		ouo.SetTotal(*f)
	}
	return ouo
}

// AddTotal adds f to the "total" field.
func (ouo *OrderUpdateOne) AddTotal(f float64) *OrderUpdateOne {
	ouo.mutation.AddTotal(f)
	return ouo
}

// SetUpdatedAt sets the "updated_at" field.
func (ouo *OrderUpdateOne) SetUpdatedAt(t time.Time) *OrderUpdateOne {
	ouo.mutation.SetUpdatedAt(t)
	return ouo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (ouo *OrderUpdateOne) SetUserID(id uuid.UUID) *OrderUpdateOne {
	ouo.mutation.SetUserID(id)
	return ouo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (ouo *OrderUpdateOne) SetNillableUserID(id *uuid.UUID) *OrderUpdateOne {
	if id != nil {
		ouo = ouo.SetUserID(*id)
	}
	return ouo
}

// SetUser sets the "user" edge to the User entity.
func (ouo *OrderUpdateOne) SetUser(u *User) *OrderUpdateOne {
	return ouo.SetUserID(u.ID)
}

// Mutation returns the OrderMutation object of the builder.
func (ouo *OrderUpdateOne) Mutation() *OrderMutation {
	return ouo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (ouo *OrderUpdateOne) ClearUser() *OrderUpdateOne {
	ouo.mutation.ClearUser()
	return ouo
}

// Where appends a list predicates to the OrderUpdate builder.
func (ouo *OrderUpdateOne) Where(ps ...predicate.Order) *OrderUpdateOne {
	ouo.mutation.Where(ps...)
	return ouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ouo *OrderUpdateOne) Select(field string, fields ...string) *OrderUpdateOne {
	ouo.fields = append([]string{field}, fields...)
	return ouo
}

// Save executes the query and returns the updated Order entity.
func (ouo *OrderUpdateOne) Save(ctx context.Context) (*Order, error) {
	ouo.defaults()
	return withHooks(ctx, ouo.sqlSave, ouo.mutation, ouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ouo *OrderUpdateOne) SaveX(ctx context.Context) *Order {
	node, err := ouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ouo *OrderUpdateOne) Exec(ctx context.Context) error {
	_, err := ouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ouo *OrderUpdateOne) ExecX(ctx context.Context) {
	if err := ouo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ouo *OrderUpdateOne) defaults() {
	if _, ok := ouo.mutation.UpdatedAt(); !ok {
		v := order.UpdateDefaultUpdatedAt()
		ouo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ouo *OrderUpdateOne) check() error {
	if v, ok := ouo.mutation.Status(); ok {
		if err := order.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`gen: validator failed for field "Order.status": %w`, err)}
		}
	}
	if v, ok := ouo.mutation.Total(); ok {
		if err := order.TotalValidator(v); err != nil {
			return &ValidationError{Name: "total", err: fmt.Errorf(`gen: validator failed for field "Order.total": %w`, err)}
		}
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (ouo *OrderUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *OrderUpdateOne {
	ouo.modifiers = append(ouo.modifiers, modifiers...)
	return ouo
}

func (ouo *OrderUpdateOne) sqlSave(ctx context.Context) (_node *Order, err error) {
	if err := ouo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(order.Table, order.Columns, sqlgraph.NewFieldSpec(order.FieldID, field.TypeUUID))
	id, ok := ouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`gen: missing "Order.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, order.FieldID)
		for _, f := range fields {
			if !order.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("gen: invalid field %q for query", f)}
			}
			if f != order.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ouo.mutation.Status(); ok {
		_spec.SetField(order.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := ouo.mutation.Total(); ok {
		_spec.SetField(order.FieldTotal, field.TypeFloat64, value)
	}
	if value, ok := ouo.mutation.AddedTotal(); ok {
		_spec.AddField(order.FieldTotal, field.TypeFloat64, value)
	}
	if value, ok := ouo.mutation.UpdatedAt(); ok {
		_spec.SetField(order.FieldUpdatedAt, field.TypeTime, value)
	}
	if ouo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   order.UserTable,
			Columns: []string{order.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ouo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   order.UserTable,
			Columns: []string{order.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(ouo.modifiers...)
	_node = &Order{config: ouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{order.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ouo.mutation.done = true
	return _node, nil
}
