// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"fmt"
	"math"
	"products/ent/gen/post"
	"products/ent/gen/predicate"
	"products/ent/gen/user"
	"products/ent/gen/userpost"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// UserPostQuery is the builder for querying UserPost entities.
type UserPostQuery struct {
	config
	ctx        *QueryContext
	order      []userpost.OrderOption
	inters     []Interceptor
	predicates []predicate.UserPost
	withUser   *UserQuery
	withPost   *PostQuery
	modifiers  []func(*sql.Selector)
	loadTotal  []func(context.Context, []*UserPost) error
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the UserPostQuery builder.
func (upq *UserPostQuery) Where(ps ...predicate.UserPost) *UserPostQuery {
	upq.predicates = append(upq.predicates, ps...)
	return upq
}

// Limit the number of records to be returned by this query.
func (upq *UserPostQuery) Limit(limit int) *UserPostQuery {
	upq.ctx.Limit = &limit
	return upq
}

// Offset to start from.
func (upq *UserPostQuery) Offset(offset int) *UserPostQuery {
	upq.ctx.Offset = &offset
	return upq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (upq *UserPostQuery) Unique(unique bool) *UserPostQuery {
	upq.ctx.Unique = &unique
	return upq
}

// Order specifies how the records should be ordered.
func (upq *UserPostQuery) Order(o ...userpost.OrderOption) *UserPostQuery {
	upq.order = append(upq.order, o...)
	return upq
}

// QueryUser chains the current query on the "user" edge.
func (upq *UserPostQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: upq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := upq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := upq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(userpost.Table, userpost.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, userpost.UserTable, userpost.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(upq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryPost chains the current query on the "post" edge.
func (upq *UserPostQuery) QueryPost() *PostQuery {
	query := (&PostClient{config: upq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := upq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := upq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(userpost.Table, userpost.FieldID, selector),
			sqlgraph.To(post.Table, post.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, userpost.PostTable, userpost.PostColumn),
		)
		fromU = sqlgraph.SetNeighbors(upq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first UserPost entity from the query.
// Returns a *NotFoundError when no UserPost was found.
func (upq *UserPostQuery) First(ctx context.Context) (*UserPost, error) {
	nodes, err := upq.Limit(1).All(setContextOp(ctx, upq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{userpost.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (upq *UserPostQuery) FirstX(ctx context.Context) *UserPost {
	node, err := upq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first UserPost ID from the query.
// Returns a *NotFoundError when no UserPost ID was found.
func (upq *UserPostQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = upq.Limit(1).IDs(setContextOp(ctx, upq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{userpost.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (upq *UserPostQuery) FirstIDX(ctx context.Context) int {
	id, err := upq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single UserPost entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one UserPost entity is found.
// Returns a *NotFoundError when no UserPost entities are found.
func (upq *UserPostQuery) Only(ctx context.Context) (*UserPost, error) {
	nodes, err := upq.Limit(2).All(setContextOp(ctx, upq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{userpost.Label}
	default:
		return nil, &NotSingularError{userpost.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (upq *UserPostQuery) OnlyX(ctx context.Context) *UserPost {
	node, err := upq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only UserPost ID in the query.
// Returns a *NotSingularError when more than one UserPost ID is found.
// Returns a *NotFoundError when no entities are found.
func (upq *UserPostQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = upq.Limit(2).IDs(setContextOp(ctx, upq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{userpost.Label}
	default:
		err = &NotSingularError{userpost.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (upq *UserPostQuery) OnlyIDX(ctx context.Context) int {
	id, err := upq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of UserPosts.
func (upq *UserPostQuery) All(ctx context.Context) ([]*UserPost, error) {
	ctx = setContextOp(ctx, upq.ctx, ent.OpQueryAll)
	if err := upq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*UserPost, *UserPostQuery]()
	return withInterceptors[[]*UserPost](ctx, upq, qr, upq.inters)
}

// AllX is like All, but panics if an error occurs.
func (upq *UserPostQuery) AllX(ctx context.Context) []*UserPost {
	nodes, err := upq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of UserPost IDs.
func (upq *UserPostQuery) IDs(ctx context.Context) (ids []int, err error) {
	if upq.ctx.Unique == nil && upq.path != nil {
		upq.Unique(true)
	}
	ctx = setContextOp(ctx, upq.ctx, ent.OpQueryIDs)
	if err = upq.Select(userpost.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (upq *UserPostQuery) IDsX(ctx context.Context) []int {
	ids, err := upq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (upq *UserPostQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, upq.ctx, ent.OpQueryCount)
	if err := upq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, upq, querierCount[*UserPostQuery](), upq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (upq *UserPostQuery) CountX(ctx context.Context) int {
	count, err := upq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (upq *UserPostQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, upq.ctx, ent.OpQueryExist)
	switch _, err := upq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("gen: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (upq *UserPostQuery) ExistX(ctx context.Context) bool {
	exist, err := upq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the UserPostQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (upq *UserPostQuery) Clone() *UserPostQuery {
	if upq == nil {
		return nil
	}
	return &UserPostQuery{
		config:     upq.config,
		ctx:        upq.ctx.Clone(),
		order:      append([]userpost.OrderOption{}, upq.order...),
		inters:     append([]Interceptor{}, upq.inters...),
		predicates: append([]predicate.UserPost{}, upq.predicates...),
		withUser:   upq.withUser.Clone(),
		withPost:   upq.withPost.Clone(),
		// clone intermediate query.
		sql:  upq.sql.Clone(),
		path: upq.path,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (upq *UserPostQuery) WithUser(opts ...func(*UserQuery)) *UserPostQuery {
	query := (&UserClient{config: upq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	upq.withUser = query
	return upq
}

// WithPost tells the query-builder to eager-load the nodes that are connected to
// the "post" edge. The optional arguments are used to configure the query builder of the edge.
func (upq *UserPostQuery) WithPost(opts ...func(*PostQuery)) *UserPostQuery {
	query := (&PostClient{config: upq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	upq.withPost = query
	return upq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		UserID int `json:"user_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.UserPost.Query().
//		GroupBy(userpost.FieldUserID).
//		Aggregate(gen.Count()).
//		Scan(ctx, &v)
func (upq *UserPostQuery) GroupBy(field string, fields ...string) *UserPostGroupBy {
	upq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &UserPostGroupBy{build: upq}
	grbuild.flds = &upq.ctx.Fields
	grbuild.label = userpost.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		UserID int `json:"user_id,omitempty"`
//	}
//
//	client.UserPost.Query().
//		Select(userpost.FieldUserID).
//		Scan(ctx, &v)
func (upq *UserPostQuery) Select(fields ...string) *UserPostSelect {
	upq.ctx.Fields = append(upq.ctx.Fields, fields...)
	sbuild := &UserPostSelect{UserPostQuery: upq}
	sbuild.label = userpost.Label
	sbuild.flds, sbuild.scan = &upq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a UserPostSelect configured with the given aggregations.
func (upq *UserPostQuery) Aggregate(fns ...AggregateFunc) *UserPostSelect {
	return upq.Select().Aggregate(fns...)
}

func (upq *UserPostQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range upq.inters {
		if inter == nil {
			return fmt.Errorf("gen: uninitialized interceptor (forgotten import gen/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, upq); err != nil {
				return err
			}
		}
	}
	for _, f := range upq.ctx.Fields {
		if !userpost.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("gen: invalid field %q for query", f)}
		}
	}
	if upq.path != nil {
		prev, err := upq.path(ctx)
		if err != nil {
			return err
		}
		upq.sql = prev
	}
	return nil
}

func (upq *UserPostQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*UserPost, error) {
	var (
		nodes       = []*UserPost{}
		_spec       = upq.querySpec()
		loadedTypes = [2]bool{
			upq.withUser != nil,
			upq.withPost != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*UserPost).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &UserPost{config: upq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(upq.modifiers) > 0 {
		_spec.Modifiers = upq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, upq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := upq.withUser; query != nil {
		if err := upq.loadUser(ctx, query, nodes, nil,
			func(n *UserPost, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	if query := upq.withPost; query != nil {
		if err := upq.loadPost(ctx, query, nodes, nil,
			func(n *UserPost, e *Post) { n.Edges.Post = e }); err != nil {
			return nil, err
		}
	}
	for i := range upq.loadTotal {
		if err := upq.loadTotal[i](ctx, nodes); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (upq *UserPostQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*UserPost, init func(*UserPost), assign func(*UserPost, *User)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*UserPost)
	for i := range nodes {
		fk := nodes[i].UserID
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
			return fmt.Errorf(`unexpected foreign-key "user_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (upq *UserPostQuery) loadPost(ctx context.Context, query *PostQuery, nodes []*UserPost, init func(*UserPost), assign func(*UserPost, *Post)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*UserPost)
	for i := range nodes {
		fk := nodes[i].PostID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(post.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "post_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (upq *UserPostQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := upq.querySpec()
	if len(upq.modifiers) > 0 {
		_spec.Modifiers = upq.modifiers
	}
	_spec.Node.Columns = upq.ctx.Fields
	if len(upq.ctx.Fields) > 0 {
		_spec.Unique = upq.ctx.Unique != nil && *upq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, upq.driver, _spec)
}

func (upq *UserPostQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(userpost.Table, userpost.Columns, sqlgraph.NewFieldSpec(userpost.FieldID, field.TypeInt))
	_spec.From = upq.sql
	if unique := upq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if upq.path != nil {
		_spec.Unique = true
	}
	if fields := upq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, userpost.FieldID)
		for i := range fields {
			if fields[i] != userpost.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
		if upq.withUser != nil {
			_spec.Node.AddColumnOnce(userpost.FieldUserID)
		}
		if upq.withPost != nil {
			_spec.Node.AddColumnOnce(userpost.FieldPostID)
		}
	}
	if ps := upq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := upq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := upq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := upq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (upq *UserPostQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(upq.driver.Dialect())
	t1 := builder.Table(userpost.Table)
	columns := upq.ctx.Fields
	if len(columns) == 0 {
		columns = userpost.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if upq.sql != nil {
		selector = upq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if upq.ctx.Unique != nil && *upq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range upq.predicates {
		p(selector)
	}
	for _, p := range upq.order {
		p(selector)
	}
	if offset := upq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := upq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// UserPostGroupBy is the group-by builder for UserPost entities.
type UserPostGroupBy struct {
	selector
	build *UserPostQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (upgb *UserPostGroupBy) Aggregate(fns ...AggregateFunc) *UserPostGroupBy {
	upgb.fns = append(upgb.fns, fns...)
	return upgb
}

// Scan applies the selector query and scans the result into the given value.
func (upgb *UserPostGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, upgb.build.ctx, ent.OpQueryGroupBy)
	if err := upgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*UserPostQuery, *UserPostGroupBy](ctx, upgb.build, upgb, upgb.build.inters, v)
}

func (upgb *UserPostGroupBy) sqlScan(ctx context.Context, root *UserPostQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(upgb.fns))
	for _, fn := range upgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*upgb.flds)+len(upgb.fns))
		for _, f := range *upgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*upgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := upgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// UserPostSelect is the builder for selecting fields of UserPost entities.
type UserPostSelect struct {
	*UserPostQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ups *UserPostSelect) Aggregate(fns ...AggregateFunc) *UserPostSelect {
	ups.fns = append(ups.fns, fns...)
	return ups
}

// Scan applies the selector query and scans the result into the given value.
func (ups *UserPostSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ups.ctx, ent.OpQuerySelect)
	if err := ups.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*UserPostQuery, *UserPostSelect](ctx, ups.UserPostQuery, ups, ups.inters, v)
}

func (ups *UserPostSelect) sqlScan(ctx context.Context, root *UserPostQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ups.fns))
	for _, fn := range ups.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ups.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ups.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
