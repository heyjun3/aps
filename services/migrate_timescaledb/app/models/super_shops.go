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

// SuperShop is an object representing the database table.
type SuperShop struct {
	Name   null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	ShopID string      `boil:"shop_id" json:"shop_id" toml:"shop_id" yaml:"shop_id"`

	R *superShopR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L superShopL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SuperShopColumns = struct {
	Name   string
	ShopID string
}{
	Name:   "name",
	ShopID: "shop_id",
}

var SuperShopTableColumns = struct {
	Name   string
	ShopID string
}{
	Name:   "super_shops.name",
	ShopID: "super_shops.shop_id",
}

// Generated where

var SuperShopWhere = struct {
	Name   whereHelpernull_String
	ShopID whereHelperstring
}{
	Name:   whereHelpernull_String{field: "\"super_shops\".\"name\""},
	ShopID: whereHelperstring{field: "\"super_shops\".\"shop_id\""},
}

// SuperShopRels is where relationship names are stored.
var SuperShopRels = struct {
}{}

// superShopR is where relationships are stored.
type superShopR struct {
}

// NewStruct creates a new relationship struct
func (*superShopR) NewStruct() *superShopR {
	return &superShopR{}
}

// superShopL is where Load methods for each relationship are stored.
type superShopL struct{}

var (
	superShopAllColumns            = []string{"name", "shop_id"}
	superShopColumnsWithoutDefault = []string{"shop_id"}
	superShopColumnsWithDefault    = []string{"name"}
	superShopPrimaryKeyColumns     = []string{"shop_id"}
	superShopGeneratedColumns      = []string{}
)

type (
	// SuperShopSlice is an alias for a slice of pointers to SuperShop.
	// This should almost always be used instead of []SuperShop.
	SuperShopSlice []*SuperShop
	// SuperShopHook is the signature for custom SuperShop hook methods
	SuperShopHook func(context.Context, boil.ContextExecutor, *SuperShop) error

	superShopQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	superShopType                 = reflect.TypeOf(&SuperShop{})
	superShopMapping              = queries.MakeStructMapping(superShopType)
	superShopPrimaryKeyMapping, _ = queries.BindMapping(superShopType, superShopMapping, superShopPrimaryKeyColumns)
	superShopInsertCacheMut       sync.RWMutex
	superShopInsertCache          = make(map[string]insertCache)
	superShopUpdateCacheMut       sync.RWMutex
	superShopUpdateCache          = make(map[string]updateCache)
	superShopUpsertCacheMut       sync.RWMutex
	superShopUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var superShopAfterSelectHooks []SuperShopHook

var superShopBeforeInsertHooks []SuperShopHook
var superShopAfterInsertHooks []SuperShopHook

var superShopBeforeUpdateHooks []SuperShopHook
var superShopAfterUpdateHooks []SuperShopHook

var superShopBeforeDeleteHooks []SuperShopHook
var superShopAfterDeleteHooks []SuperShopHook

var superShopBeforeUpsertHooks []SuperShopHook
var superShopAfterUpsertHooks []SuperShopHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SuperShop) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SuperShop) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SuperShop) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SuperShop) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SuperShop) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SuperShop) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SuperShop) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SuperShop) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SuperShop) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range superShopAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSuperShopHook registers your hook function for all future operations.
func AddSuperShopHook(hookPoint boil.HookPoint, superShopHook SuperShopHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		superShopAfterSelectHooks = append(superShopAfterSelectHooks, superShopHook)
	case boil.BeforeInsertHook:
		superShopBeforeInsertHooks = append(superShopBeforeInsertHooks, superShopHook)
	case boil.AfterInsertHook:
		superShopAfterInsertHooks = append(superShopAfterInsertHooks, superShopHook)
	case boil.BeforeUpdateHook:
		superShopBeforeUpdateHooks = append(superShopBeforeUpdateHooks, superShopHook)
	case boil.AfterUpdateHook:
		superShopAfterUpdateHooks = append(superShopAfterUpdateHooks, superShopHook)
	case boil.BeforeDeleteHook:
		superShopBeforeDeleteHooks = append(superShopBeforeDeleteHooks, superShopHook)
	case boil.AfterDeleteHook:
		superShopAfterDeleteHooks = append(superShopAfterDeleteHooks, superShopHook)
	case boil.BeforeUpsertHook:
		superShopBeforeUpsertHooks = append(superShopBeforeUpsertHooks, superShopHook)
	case boil.AfterUpsertHook:
		superShopAfterUpsertHooks = append(superShopAfterUpsertHooks, superShopHook)
	}
}

// One returns a single superShop record from the query.
func (q superShopQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SuperShop, error) {
	o := &SuperShop{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for super_shops")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SuperShop records from the query.
func (q superShopQuery) All(ctx context.Context, exec boil.ContextExecutor) (SuperShopSlice, error) {
	var o []*SuperShop

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SuperShop slice")
	}

	if len(superShopAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SuperShop records in the query.
func (q superShopQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count super_shops rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q superShopQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if super_shops exists")
	}

	return count > 0, nil
}

// SuperShops retrieves all the records using an executor.
func SuperShops(mods ...qm.QueryMod) superShopQuery {
	mods = append(mods, qm.From("\"super_shops\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"super_shops\".*"})
	}

	return superShopQuery{q}
}

// FindSuperShop retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSuperShop(ctx context.Context, exec boil.ContextExecutor, shopID string, selectCols ...string) (*SuperShop, error) {
	superShopObj := &SuperShop{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"super_shops\" where \"shop_id\"=$1", sel,
	)

	q := queries.Raw(query, shopID)

	err := q.Bind(ctx, exec, superShopObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from super_shops")
	}

	if err = superShopObj.doAfterSelectHooks(ctx, exec); err != nil {
		return superShopObj, err
	}

	return superShopObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SuperShop) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no super_shops provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(superShopColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	superShopInsertCacheMut.RLock()
	cache, cached := superShopInsertCache[key]
	superShopInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			superShopAllColumns,
			superShopColumnsWithDefault,
			superShopColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(superShopType, superShopMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(superShopType, superShopMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"super_shops\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"super_shops\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into super_shops")
	}

	if !cached {
		superShopInsertCacheMut.Lock()
		superShopInsertCache[key] = cache
		superShopInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the SuperShop.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SuperShop) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	superShopUpdateCacheMut.RLock()
	cache, cached := superShopUpdateCache[key]
	superShopUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			superShopAllColumns,
			superShopPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update super_shops, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"super_shops\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, superShopPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(superShopType, superShopMapping, append(wl, superShopPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update super_shops row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for super_shops")
	}

	if !cached {
		superShopUpdateCacheMut.Lock()
		superShopUpdateCache[key] = cache
		superShopUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q superShopQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for super_shops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for super_shops")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SuperShopSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), superShopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"super_shops\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, superShopPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in superShop slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all superShop")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SuperShop) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no super_shops provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(superShopColumnsWithDefault, o)

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

	superShopUpsertCacheMut.RLock()
	cache, cached := superShopUpsertCache[key]
	superShopUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			superShopAllColumns,
			superShopColumnsWithDefault,
			superShopColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			superShopAllColumns,
			superShopPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert super_shops, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(superShopPrimaryKeyColumns))
			copy(conflict, superShopPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"super_shops\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(superShopType, superShopMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(superShopType, superShopMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert super_shops")
	}

	if !cached {
		superShopUpsertCacheMut.Lock()
		superShopUpsertCache[key] = cache
		superShopUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single SuperShop record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SuperShop) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no SuperShop provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), superShopPrimaryKeyMapping)
	sql := "DELETE FROM \"super_shops\" WHERE \"shop_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from super_shops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for super_shops")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q superShopQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no superShopQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from super_shops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for super_shops")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SuperShopSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(superShopBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), superShopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"super_shops\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, superShopPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from superShop slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for super_shops")
	}

	if len(superShopAfterDeleteHooks) != 0 {
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
func (o *SuperShop) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSuperShop(ctx, exec, o.ShopID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SuperShopSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SuperShopSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), superShopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"super_shops\".* FROM \"super_shops\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, superShopPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SuperShopSlice")
	}

	*o = slice

	return nil
}

// SuperShopExists checks if the SuperShop row exists.
func SuperShopExists(ctx context.Context, exec boil.ContextExecutor, shopID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"super_shops\" where \"shop_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, shopID)
	}
	row := exec.QueryRowContext(ctx, sql, shopID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if super_shops exists")
	}

	return exists, nil
}

// Exists checks if the SuperShop row exists.
func (o *SuperShop) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SuperShopExists(ctx, exec, o.ShopID)
}
