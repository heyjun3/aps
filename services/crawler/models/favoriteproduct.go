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

// Favoriteproduct is an object representing the database table.
type Favoriteproduct struct {
	URL  string   `boil:"url" json:"url" toml:"url" yaml:"url"`
	Jan  string   `boil:"jan" json:"jan" toml:"jan" yaml:"jan"`
	Cost null.Int `boil:"cost" json:"cost,omitempty" toml:"cost" yaml:"cost,omitempty"`

	R *favoriteproductR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L favoriteproductL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var FavoriteproductColumns = struct {
	URL  string
	Jan  string
	Cost string
}{
	URL:  "url",
	Jan:  "jan",
	Cost: "cost",
}

var FavoriteproductTableColumns = struct {
	URL  string
	Jan  string
	Cost string
}{
	URL:  "favoriteproduct.url",
	Jan:  "favoriteproduct.jan",
	Cost: "favoriteproduct.cost",
}

// Generated where

type whereHelpernull_Int struct{ field string }

func (w whereHelpernull_Int) EQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_Int) NEQ(x null.Int) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_Int) LT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_Int) LTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_Int) GT(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_Int) GTE(x null.Int) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}
func (w whereHelpernull_Int) IN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelpernull_Int) NIN(slice []int) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

func (w whereHelpernull_Int) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_Int) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }

var FavoriteproductWhere = struct {
	URL  whereHelperstring
	Jan  whereHelperstring
	Cost whereHelpernull_Int
}{
	URL:  whereHelperstring{field: "\"favoriteproduct\".\"url\""},
	Jan:  whereHelperstring{field: "\"favoriteproduct\".\"jan\""},
	Cost: whereHelpernull_Int{field: "\"favoriteproduct\".\"cost\""},
}

// FavoriteproductRels is where relationship names are stored.
var FavoriteproductRels = struct {
}{}

// favoriteproductR is where relationships are stored.
type favoriteproductR struct {
}

// NewStruct creates a new relationship struct
func (*favoriteproductR) NewStruct() *favoriteproductR {
	return &favoriteproductR{}
}

// favoriteproductL is where Load methods for each relationship are stored.
type favoriteproductL struct{}

var (
	favoriteproductAllColumns            = []string{"url", "jan", "cost"}
	favoriteproductColumnsWithoutDefault = []string{"url", "jan"}
	favoriteproductColumnsWithDefault    = []string{"cost"}
	favoriteproductPrimaryKeyColumns     = []string{"url", "jan"}
	favoriteproductGeneratedColumns      = []string{}
)

type (
	// FavoriteproductSlice is an alias for a slice of pointers to Favoriteproduct.
	// This should almost always be used instead of []Favoriteproduct.
	FavoriteproductSlice []*Favoriteproduct
	// FavoriteproductHook is the signature for custom Favoriteproduct hook methods
	FavoriteproductHook func(context.Context, boil.ContextExecutor, *Favoriteproduct) error

	favoriteproductQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	favoriteproductType                 = reflect.TypeOf(&Favoriteproduct{})
	favoriteproductMapping              = queries.MakeStructMapping(favoriteproductType)
	favoriteproductPrimaryKeyMapping, _ = queries.BindMapping(favoriteproductType, favoriteproductMapping, favoriteproductPrimaryKeyColumns)
	favoriteproductInsertCacheMut       sync.RWMutex
	favoriteproductInsertCache          = make(map[string]insertCache)
	favoriteproductUpdateCacheMut       sync.RWMutex
	favoriteproductUpdateCache          = make(map[string]updateCache)
	favoriteproductUpsertCacheMut       sync.RWMutex
	favoriteproductUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var favoriteproductAfterSelectHooks []FavoriteproductHook

var favoriteproductBeforeInsertHooks []FavoriteproductHook
var favoriteproductAfterInsertHooks []FavoriteproductHook

var favoriteproductBeforeUpdateHooks []FavoriteproductHook
var favoriteproductAfterUpdateHooks []FavoriteproductHook

var favoriteproductBeforeDeleteHooks []FavoriteproductHook
var favoriteproductAfterDeleteHooks []FavoriteproductHook

var favoriteproductBeforeUpsertHooks []FavoriteproductHook
var favoriteproductAfterUpsertHooks []FavoriteproductHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Favoriteproduct) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Favoriteproduct) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Favoriteproduct) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Favoriteproduct) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Favoriteproduct) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Favoriteproduct) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Favoriteproduct) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Favoriteproduct) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Favoriteproduct) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range favoriteproductAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddFavoriteproductHook registers your hook function for all future operations.
func AddFavoriteproductHook(hookPoint boil.HookPoint, favoriteproductHook FavoriteproductHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		favoriteproductAfterSelectHooks = append(favoriteproductAfterSelectHooks, favoriteproductHook)
	case boil.BeforeInsertHook:
		favoriteproductBeforeInsertHooks = append(favoriteproductBeforeInsertHooks, favoriteproductHook)
	case boil.AfterInsertHook:
		favoriteproductAfterInsertHooks = append(favoriteproductAfterInsertHooks, favoriteproductHook)
	case boil.BeforeUpdateHook:
		favoriteproductBeforeUpdateHooks = append(favoriteproductBeforeUpdateHooks, favoriteproductHook)
	case boil.AfterUpdateHook:
		favoriteproductAfterUpdateHooks = append(favoriteproductAfterUpdateHooks, favoriteproductHook)
	case boil.BeforeDeleteHook:
		favoriteproductBeforeDeleteHooks = append(favoriteproductBeforeDeleteHooks, favoriteproductHook)
	case boil.AfterDeleteHook:
		favoriteproductAfterDeleteHooks = append(favoriteproductAfterDeleteHooks, favoriteproductHook)
	case boil.BeforeUpsertHook:
		favoriteproductBeforeUpsertHooks = append(favoriteproductBeforeUpsertHooks, favoriteproductHook)
	case boil.AfterUpsertHook:
		favoriteproductAfterUpsertHooks = append(favoriteproductAfterUpsertHooks, favoriteproductHook)
	}
}

// One returns a single favoriteproduct record from the query.
func (q favoriteproductQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Favoriteproduct, error) {
	o := &Favoriteproduct{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for favoriteproduct")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Favoriteproduct records from the query.
func (q favoriteproductQuery) All(ctx context.Context, exec boil.ContextExecutor) (FavoriteproductSlice, error) {
	var o []*Favoriteproduct

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Favoriteproduct slice")
	}

	if len(favoriteproductAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Favoriteproduct records in the query.
func (q favoriteproductQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count favoriteproduct rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q favoriteproductQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if favoriteproduct exists")
	}

	return count > 0, nil
}

// Favoriteproducts retrieves all the records using an executor.
func Favoriteproducts(mods ...qm.QueryMod) favoriteproductQuery {
	mods = append(mods, qm.From("\"favoriteproduct\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"favoriteproduct\".*"})
	}

	return favoriteproductQuery{q}
}

// FindFavoriteproduct retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFavoriteproduct(ctx context.Context, exec boil.ContextExecutor, uRL string, jan string, selectCols ...string) (*Favoriteproduct, error) {
	favoriteproductObj := &Favoriteproduct{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"favoriteproduct\" where \"url\"=$1 AND \"jan\"=$2", sel,
	)

	q := queries.Raw(query, uRL, jan)

	err := q.Bind(ctx, exec, favoriteproductObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from favoriteproduct")
	}

	if err = favoriteproductObj.doAfterSelectHooks(ctx, exec); err != nil {
		return favoriteproductObj, err
	}

	return favoriteproductObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Favoriteproduct) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no favoriteproduct provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(favoriteproductColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	favoriteproductInsertCacheMut.RLock()
	cache, cached := favoriteproductInsertCache[key]
	favoriteproductInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			favoriteproductAllColumns,
			favoriteproductColumnsWithDefault,
			favoriteproductColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(favoriteproductType, favoriteproductMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(favoriteproductType, favoriteproductMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"favoriteproduct\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"favoriteproduct\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into favoriteproduct")
	}

	if !cached {
		favoriteproductInsertCacheMut.Lock()
		favoriteproductInsertCache[key] = cache
		favoriteproductInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Favoriteproduct.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Favoriteproduct) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	favoriteproductUpdateCacheMut.RLock()
	cache, cached := favoriteproductUpdateCache[key]
	favoriteproductUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			favoriteproductAllColumns,
			favoriteproductPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update favoriteproduct, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"favoriteproduct\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, favoriteproductPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(favoriteproductType, favoriteproductMapping, append(wl, favoriteproductPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update favoriteproduct row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for favoriteproduct")
	}

	if !cached {
		favoriteproductUpdateCacheMut.Lock()
		favoriteproductUpdateCache[key] = cache
		favoriteproductUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q favoriteproductQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for favoriteproduct")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for favoriteproduct")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FavoriteproductSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), favoriteproductPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"favoriteproduct\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, favoriteproductPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in favoriteproduct slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all favoriteproduct")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Favoriteproduct) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no favoriteproduct provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(favoriteproductColumnsWithDefault, o)

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

	favoriteproductUpsertCacheMut.RLock()
	cache, cached := favoriteproductUpsertCache[key]
	favoriteproductUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			favoriteproductAllColumns,
			favoriteproductColumnsWithDefault,
			favoriteproductColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			favoriteproductAllColumns,
			favoriteproductPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert favoriteproduct, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(favoriteproductPrimaryKeyColumns))
			copy(conflict, favoriteproductPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"favoriteproduct\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(favoriteproductType, favoriteproductMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(favoriteproductType, favoriteproductMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert favoriteproduct")
	}

	if !cached {
		favoriteproductUpsertCacheMut.Lock()
		favoriteproductUpsertCache[key] = cache
		favoriteproductUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Favoriteproduct record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Favoriteproduct) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Favoriteproduct provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), favoriteproductPrimaryKeyMapping)
	sql := "DELETE FROM \"favoriteproduct\" WHERE \"url\"=$1 AND \"jan\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from favoriteproduct")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for favoriteproduct")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q favoriteproductQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no favoriteproductQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from favoriteproduct")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for favoriteproduct")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FavoriteproductSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(favoriteproductBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), favoriteproductPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"favoriteproduct\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, favoriteproductPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from favoriteproduct slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for favoriteproduct")
	}

	if len(favoriteproductAfterDeleteHooks) != 0 {
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
func (o *Favoriteproduct) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindFavoriteproduct(ctx, exec, o.URL, o.Jan)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FavoriteproductSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := FavoriteproductSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), favoriteproductPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"favoriteproduct\".* FROM \"favoriteproduct\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, favoriteproductPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in FavoriteproductSlice")
	}

	*o = slice

	return nil
}

// FavoriteproductExists checks if the Favoriteproduct row exists.
func FavoriteproductExists(ctx context.Context, exec boil.ContextExecutor, uRL string, jan string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"favoriteproduct\" where \"url\"=$1 AND \"jan\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, uRL, jan)
	}
	row := exec.QueryRowContext(ctx, sql, uRL, jan)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if favoriteproduct exists")
	}

	return exists, nil
}

// Exists checks if the Favoriteproduct row exists.
func (o *Favoriteproduct) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return FavoriteproductExists(ctx, exec, o.URL, o.Jan)
}