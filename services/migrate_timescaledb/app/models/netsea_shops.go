// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// NetseaShop is an object representing the database table.
type NetseaShop struct {
	Name   null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	ShopID string      `boil:"shop_id" json:"shop_id" toml:"shop_id" yaml:"shop_id"`

	R *netseaShopR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L netseaShopL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var NetseaShopColumns = struct {
	Name   string
	ShopID string
}{
	Name:   "name",
	ShopID: "shop_id",
}

var NetseaShopTableColumns = struct {
	Name   string
	ShopID string
}{
	Name:   "netsea_shops.name",
	ShopID: "netsea_shops.shop_id",
}

// Generated where

var NetseaShopWhere = struct {
	Name   whereHelpernull_String
	ShopID whereHelperstring
}{
	Name:   whereHelpernull_String{field: "\"netsea_shops\".\"name\""},
	ShopID: whereHelperstring{field: "\"netsea_shops\".\"shop_id\""},
}

// NetseaShopRels is where relationship names are stored.
var NetseaShopRels = struct {
}{}

// netseaShopR is where relationships are stored.
type netseaShopR struct {
}

// NewStruct creates a new relationship struct
func (*netseaShopR) NewStruct() *netseaShopR {
	return &netseaShopR{}
}

// netseaShopL is where Load methods for each relationship are stored.
type netseaShopL struct{}

var (
	netseaShopAllColumns            = []string{"name", "shop_id"}
	netseaShopColumnsWithoutDefault = []string{"shop_id"}
	netseaShopColumnsWithDefault    = []string{"name"}
	netseaShopPrimaryKeyColumns     = []string{"shop_id"}
	netseaShopGeneratedColumns      = []string{}
)

type (
	// NetseaShopSlice is an alias for a slice of pointers to NetseaShop.
	// This should almost always be used instead of []NetseaShop.
	NetseaShopSlice []*NetseaShop
	// NetseaShopHook is the signature for custom NetseaShop hook methods
	NetseaShopHook func(context.Context, boil.ContextExecutor, *NetseaShop) error

	netseaShopQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	netseaShopType                 = reflect.TypeOf(&NetseaShop{})
	netseaShopMapping              = queries.MakeStructMapping(netseaShopType)
	netseaShopPrimaryKeyMapping, _ = queries.BindMapping(netseaShopType, netseaShopMapping, netseaShopPrimaryKeyColumns)
	netseaShopInsertCacheMut       sync.RWMutex
	netseaShopInsertCache          = make(map[string]insertCache)
	netseaShopUpdateCacheMut       sync.RWMutex
	netseaShopUpdateCache          = make(map[string]updateCache)
	netseaShopUpsertCacheMut       sync.RWMutex
	netseaShopUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var netseaShopAfterSelectHooks []NetseaShopHook

var netseaShopBeforeInsertHooks []NetseaShopHook
var netseaShopAfterInsertHooks []NetseaShopHook

var netseaShopBeforeUpdateHooks []NetseaShopHook
var netseaShopAfterUpdateHooks []NetseaShopHook

var netseaShopBeforeDeleteHooks []NetseaShopHook
var netseaShopAfterDeleteHooks []NetseaShopHook

var netseaShopBeforeUpsertHooks []NetseaShopHook
var netseaShopAfterUpsertHooks []NetseaShopHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *NetseaShop) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *NetseaShop) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *NetseaShop) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *NetseaShop) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *NetseaShop) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *NetseaShop) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *NetseaShop) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *NetseaShop) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *NetseaShop) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range netseaShopAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddNetseaShopHook registers your hook function for all future operations.
func AddNetseaShopHook(hookPoint boil.HookPoint, netseaShopHook NetseaShopHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		netseaShopAfterSelectHooks = append(netseaShopAfterSelectHooks, netseaShopHook)
	case boil.BeforeInsertHook:
		netseaShopBeforeInsertHooks = append(netseaShopBeforeInsertHooks, netseaShopHook)
	case boil.AfterInsertHook:
		netseaShopAfterInsertHooks = append(netseaShopAfterInsertHooks, netseaShopHook)
	case boil.BeforeUpdateHook:
		netseaShopBeforeUpdateHooks = append(netseaShopBeforeUpdateHooks, netseaShopHook)
	case boil.AfterUpdateHook:
		netseaShopAfterUpdateHooks = append(netseaShopAfterUpdateHooks, netseaShopHook)
	case boil.BeforeDeleteHook:
		netseaShopBeforeDeleteHooks = append(netseaShopBeforeDeleteHooks, netseaShopHook)
	case boil.AfterDeleteHook:
		netseaShopAfterDeleteHooks = append(netseaShopAfterDeleteHooks, netseaShopHook)
	case boil.BeforeUpsertHook:
		netseaShopBeforeUpsertHooks = append(netseaShopBeforeUpsertHooks, netseaShopHook)
	case boil.AfterUpsertHook:
		netseaShopAfterUpsertHooks = append(netseaShopAfterUpsertHooks, netseaShopHook)
	}
}

// One returns a single netseaShop record from the query.
func (q netseaShopQuery) One(ctx context.Context, exec boil.ContextExecutor) (*NetseaShop, error) {
	o := &NetseaShop{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for netsea_shops")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all NetseaShop records from the query.
func (q netseaShopQuery) All(ctx context.Context, exec boil.ContextExecutor) (NetseaShopSlice, error) {
	var o []*NetseaShop

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to NetseaShop slice")
	}

	if len(netseaShopAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all NetseaShop records in the query.
func (q netseaShopQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count netsea_shops rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q netseaShopQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if netsea_shops exists")
	}

	return count > 0, nil
}

// NetseaShops retrieves all the records using an executor.
func NetseaShops(mods ...qm.QueryMod) netseaShopQuery {
	mods = append(mods, qm.From("\"netsea_shops\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"netsea_shops\".*"})
	}

	return netseaShopQuery{q}
}

// FindNetseaShop retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindNetseaShop(ctx context.Context, exec boil.ContextExecutor, shopID string, selectCols ...string) (*NetseaShop, error) {
	netseaShopObj := &NetseaShop{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"netsea_shops\" where \"shop_id\"=$1", sel,
	)

	q := queries.Raw(query, shopID)

	err := q.Bind(ctx, exec, netseaShopObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from netsea_shops")
	}

	if err = netseaShopObj.doAfterSelectHooks(ctx, exec); err != nil {
		return netseaShopObj, err
	}

	return netseaShopObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *NetseaShop) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no netsea_shops provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(netseaShopColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	netseaShopInsertCacheMut.RLock()
	cache, cached := netseaShopInsertCache[key]
	netseaShopInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			netseaShopAllColumns,
			netseaShopColumnsWithDefault,
			netseaShopColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(netseaShopType, netseaShopMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(netseaShopType, netseaShopMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"netsea_shops\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"netsea_shops\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into netsea_shops")
	}

	if !cached {
		netseaShopInsertCacheMut.Lock()
		netseaShopInsertCache[key] = cache
		netseaShopInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the NetseaShop.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *NetseaShop) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	netseaShopUpdateCacheMut.RLock()
	cache, cached := netseaShopUpdateCache[key]
	netseaShopUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			netseaShopAllColumns,
			netseaShopPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update netsea_shops, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"netsea_shops\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, netseaShopPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(netseaShopType, netseaShopMapping, append(wl, netseaShopPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update netsea_shops row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for netsea_shops")
	}

	if !cached {
		netseaShopUpdateCacheMut.Lock()
		netseaShopUpdateCache[key] = cache
		netseaShopUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q netseaShopQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for netsea_shops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for netsea_shops")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o NetseaShopSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), netseaShopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"netsea_shops\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, netseaShopPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in netseaShop slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all netseaShop")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *NetseaShop) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no netsea_shops provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(netseaShopColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	netseaShopUpsertCacheMut.RLock()
	cache, cached := netseaShopUpsertCache[key]
	netseaShopUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			netseaShopAllColumns,
			netseaShopColumnsWithDefault,
			netseaShopColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			netseaShopAllColumns,
			netseaShopPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert netsea_shops, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(netseaShopPrimaryKeyColumns))
			copy(conflict, netseaShopPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"netsea_shops\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(netseaShopType, netseaShopMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(netseaShopType, netseaShopMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert netsea_shops")
	}

	if !cached {
		netseaShopUpsertCacheMut.Lock()
		netseaShopUpsertCache[key] = cache
		netseaShopUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single NetseaShop record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *NetseaShop) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no NetseaShop provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), netseaShopPrimaryKeyMapping)
	sql := "DELETE FROM \"netsea_shops\" WHERE \"shop_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from netsea_shops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for netsea_shops")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q netseaShopQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no netseaShopQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from netsea_shops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for netsea_shops")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o NetseaShopSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(netseaShopBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), netseaShopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"netsea_shops\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, netseaShopPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from netseaShop slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for netsea_shops")
	}

	if len(netseaShopAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *NetseaShop) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindNetseaShop(ctx, exec, o.ShopID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *NetseaShopSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := NetseaShopSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), netseaShopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"netsea_shops\".* FROM \"netsea_shops\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, netseaShopPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in NetseaShopSlice")
	}

	*o = slice

	return nil
}

// NetseaShopExists checks if the NetseaShop row exists.
func NetseaShopExists(ctx context.Context, exec boil.ContextExecutor, shopID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"netsea_shops\" where \"shop_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, shopID)
	}
	row := exec.QueryRowContext(ctx, sql, shopID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if netsea_shops exists")
	}

	return exists, nil
}

// Exists checks if the NetseaShop row exists.
func (o *NetseaShop) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return NetseaShopExists(ctx, exec, o.ShopID)
}
