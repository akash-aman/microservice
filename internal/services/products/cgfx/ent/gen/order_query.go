// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"fmt"
	"math"
	"products/cgfx/ent/gen/order"
	"products/cgfx/ent/gen/predicate"
	"products/cgfx/ent/gen/user"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// OrderQuery is the builder for querying Order entities.
type OrderQuery struct {
	config
	ctx        *QueryContext
	order      []order.OrderOption
	inters     []Interceptor
	predicates []predicate.Order
	withUser   *UserQuery
	withFKs    bool
	loadTotal  []func(context.Context, []*Order) error
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the OrderQuery builder.
func (oq *OrderQuery) Where(ps ...predicate.Order) *OrderQuery {
	oq.predicates = append(oq.predicates, ps...)
	return oq
}

// Limit the number of records to be returned by this query.
func (oq *OrderQuery) Limit(limit int) *OrderQuery {
	oq.ctx.Limit = &limit
	return oq
}

// Offset to start from.
func (oq *OrderQuery) Offset(offset int) *OrderQuery {
	oq.ctx.Offset = &offset
	return oq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (oq *OrderQuery) Unique(unique bool) *OrderQuery {
	oq.ctx.Unique = &unique
	return oq
}

// Order specifies how the records should be ordered.
func (oq *OrderQuery) Order(o ...order.OrderOption) *OrderQuery {
	oq.order = append(oq.order, o...)
	return oq
}

// QueryUser chains the current query on the "user" edge.
func (oq *OrderQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: oq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := oq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := oq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(order.Table, order.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, order.UserTable, order.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(oq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Order entity from the query.
// Returns a *NotFoundError when no Order was found.
func (oq *OrderQuery) First(ctx context.Context) (*Order, error) {
	nodes, err := oq.Limit(1).All(setContextOp(ctx, oq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{order.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (oq *OrderQuery) FirstX(ctx context.Context) *Order {
	node, err := oq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Order ID from the query.
// Returns a *NotFoundError when no Order ID was found.
func (oq *OrderQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = oq.Limit(1).IDs(setContextOp(ctx, oq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{order.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (oq *OrderQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := oq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Order entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Order entity is found.
// Returns a *NotFoundError when no Order entities are found.
func (oq *OrderQuery) Only(ctx context.Context) (*Order, error) {
	nodes, err := oq.Limit(2).All(setContextOp(ctx, oq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{order.Label}
	default:
		return nil, &NotSingularError{order.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (oq *OrderQuery) OnlyX(ctx context.Context) *Order {
	node, err := oq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Order ID in the query.
// Returns a *NotSingularError when more than one Order ID is found.
// Returns a *NotFoundError when no entities are found.
func (oq *OrderQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = oq.Limit(2).IDs(setContextOp(ctx, oq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{order.Label}
	default:
		err = &NotSingularError{order.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (oq *OrderQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := oq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Orders.
func (oq *OrderQuery) All(ctx context.Context) ([]*Order, error) {
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryAll)
	if err := oq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Order, *OrderQuery]()
	return withInterceptors[[]*Order](ctx, oq, qr, oq.inters)
}

// AllX is like All, but panics if an error occurs.
func (oq *OrderQuery) AllX(ctx context.Context) []*Order {
	nodes, err := oq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Order IDs.
func (oq *OrderQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if oq.ctx.Unique == nil && oq.path != nil {
		oq.Unique(true)
	}
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryIDs)
	if err = oq.Select(order.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (oq *OrderQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := oq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (oq *OrderQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryCount)
	if err := oq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, oq, querierCount[*OrderQuery](), oq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (oq *OrderQuery) CountX(ctx context.Context) int {
	count, err := oq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (oq *OrderQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, oq.ctx, ent.OpQueryExist)
	switch _, err := oq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("gen: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (oq *OrderQuery) ExistX(ctx context.Context) bool {
	exist, err := oq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the OrderQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (oq *OrderQuery) Clone() *OrderQuery {
	if oq == nil {
		return nil
	}
	return &OrderQuery{
		config:     oq.config,
		ctx:        oq.ctx.Clone(),
		order:      append([]order.OrderOption{}, oq.order...),
		inters:     append([]Interceptor{}, oq.inters...),
		predicates: append([]predicate.Order{}, oq.predicates...),
		withUser:   oq.withUser.Clone(),
		// clone intermediate query.
		sql:       oq.sql.Clone(),
		path:      oq.path,
		modifiers: append([]func(*sql.Selector){}, oq.modifiers...),
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (oq *OrderQuery) WithUser(opts ...func(*UserQuery)) *OrderQuery {
	query := (&UserClient{config: oq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	oq.withUser = query
	return oq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Status order.Status `json:"status,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Order.Query().
//		GroupBy(order.FieldStatus).
//		Aggregate(gen.Count()).
//		Scan(ctx, &v)
func (oq *OrderQuery) GroupBy(field string, fields ...string) *OrderGroupBy {
	oq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &OrderGroupBy{build: oq}
	grbuild.flds = &oq.ctx.Fields
	grbuild.label = order.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Status order.Status `json:"status,omitempty"`
//	}
//
//	client.Order.Query().
//		Select(order.FieldStatus).
//		Scan(ctx, &v)
func (oq *OrderQuery) Select(fields ...string) *OrderSelect {
	oq.ctx.Fields = append(oq.ctx.Fields, fields...)
	sbuild := &OrderSelect{OrderQuery: oq}
	sbuild.label = order.Label
	sbuild.flds, sbuild.scan = &oq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a OrderSelect configured with the given aggregations.
func (oq *OrderQuery) Aggregate(fns ...AggregateFunc) *OrderSelect {
	return oq.Select().Aggregate(fns...)
}

func (oq *OrderQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range oq.inters {
		if inter == nil {
			return fmt.Errorf("gen: uninitialized interceptor (forgotten import gen/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, oq); err != nil {
				return err
			}
		}
	}
	for _, f := range oq.ctx.Fields {
		if !order.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("gen: invalid field %q for query", f)}
		}
	}
	if oq.path != nil {
		prev, err := oq.path(ctx)
		if err != nil {
			return err
		}
		oq.sql = prev
	}
	return nil
}

func (oq *OrderQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Order, error) {
	var (
		nodes       = []*Order{}
		withFKs     = oq.withFKs
		_spec       = oq.querySpec()
		loadedTypes = [1]bool{
			oq.withUser != nil,
		}
	)
	if oq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, order.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Order).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Order{config: oq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(oq.modifiers) > 0 {
		_spec.Modifiers = oq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, oq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := oq.withUser; query != nil {
		if err := oq.loadUser(ctx, query, nodes, nil,
			func(n *Order, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	for i := range oq.loadTotal {
		if err := oq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (oq *OrderQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Order, init func(*Order), assign func(*Order, *User)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*Order)
	for i := range nodes {
		if nodes[i].user_orders == nil {
			continue
		}
		fk := *nodes[i].user_orders
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_orders" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (oq *OrderQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := oq.querySpec()
	if len(oq.modifiers) > 0 {
		_spec.Modifiers = oq.modifiers
	}
	_spec.Node.Columns = oq.ctx.Fields
	if len(oq.ctx.Fields) > 0 {
		_spec.Unique = oq.ctx.Unique != nil && *oq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, oq.driver, _spec)
}

func (oq *OrderQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(order.Table, order.Columns, sqlgraph.NewFieldSpec(order.FieldID, field.TypeUUID))
	_spec.From = oq.sql
	if unique := oq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if oq.path != nil {
		_spec.Unique = true
	}
	if fields := oq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, order.FieldID)
		for i := range fields {
			if fields[i] != order.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := oq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := oq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := oq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := oq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (oq *OrderQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(oq.driver.Dialect())
	t1 := builder.Table(order.Table)
	columns := oq.ctx.Fields
	if len(columns) == 0 {
		columns = order.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if oq.sql != nil {
		selector = oq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if oq.ctx.Unique != nil && *oq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range oq.modifiers {
		m(selector)
	}
	for _, p := range oq.predicates {
		p(selector)
	}
	for _, p := range oq.order {
		p(selector)
	}
	if offset := oq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := oq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (oq *OrderQuery) Modify(modifiers ...func(s *sql.Selector)) *OrderSelect {
	oq.modifiers = append(oq.modifiers, modifiers...)
	return oq.Select()
}

// OrderGroupBy is the group-by builder for Order entities.
type OrderGroupBy struct {
	selector
	build *OrderQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ogb *OrderGroupBy) Aggregate(fns ...AggregateFunc) *OrderGroupBy {
	ogb.fns = append(ogb.fns, fns...)
	return ogb
}

// Scan applies the selector query and scans the result into the given value.
func (ogb *OrderGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ogb.build.ctx, ent.OpQueryGroupBy)
	if err := ogb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OrderQuery, *OrderGroupBy](ctx, ogb.build, ogb, ogb.build.inters, v)
}

func (ogb *OrderGroupBy) sqlScan(ctx context.Context, root *OrderQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(ogb.fns))
	for _, fn := range ogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*ogb.flds)+len(ogb.fns))
		for _, f := range *ogb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*ogb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ogb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// OrderSelect is the builder for selecting fields of Order entities.
type OrderSelect struct {
	*OrderQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (os *OrderSelect) Aggregate(fns ...AggregateFunc) *OrderSelect {
	os.fns = append(os.fns, fns...)
	return os
}

// Scan applies the selector query and scans the result into the given value.
func (os *OrderSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, os.ctx, ent.OpQuerySelect)
	if err := os.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*OrderQuery, *OrderSelect](ctx, os.OrderQuery, os, os.inters, v)
}

func (os *OrderSelect) sqlScan(ctx context.Context, root *OrderQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(os.fns))
	for _, fn := range os.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*os.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := os.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (os *OrderSelect) Modify(modifiers ...func(s *sql.Selector)) *OrderSelect {
	os.modifiers = append(os.modifiers, modifiers...)
	return os
}
