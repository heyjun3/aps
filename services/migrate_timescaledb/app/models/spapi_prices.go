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

// SpapiPrice is an object representing the database table.
type SpapiPrice struct {
	Asin     string     `boil:"asin" json:"asin" toml:"asin" yaml:"asin"`
	Price    null.Int64 `boil:"price" json:"price,omitempty" toml:"price" yaml:"price,omitempty"`
	Modified null.Time  `boil:"modified" json:"modified,omitempty" toml:"modified" yaml:"modified,omitempty"`

	R *spapiPriceR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L spapiPriceL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SpapiPriceColumns = struct {
	Asin     string
	Price    string
	Modified string
}{
	Asin:     "asin",
	Price:    "price",
	Modified: "modified",
}

var SpapiPriceTableColumns = struct {
	Asin     string
	Price    string
	Modified string
}{
	Asin:     "spapi_prices.asin",
	Price:    "spapi_prices.price",
	Modified: "spapi_prices.modified",
}

// Generated where

var SpapiPriceWhere = struct {
	Asin     whereHelperstring
	Price    whereHelpernull_Int64
	Modified whereHelpernull_Time
}{
	Asin:     whereHelperstring{field: "\"spapi_prices\".\"asin\""},
	Price:    whereHelpernull_Int64{field: "\"spapi_prices\".\"price\""},
	Modified: whereHelpernull_Time{field: "\"spapi_prices\".\"modified\""},
}

// SpapiPriceRels is where relationship names are stored.
var SpapiPriceRels = struct {
	AsinAsinsInfo string
}{
	AsinAsinsInfo: "AsinAsinsInfo",
}

// spapiPriceR is where relationships are stored.
type spapiPriceR struct {
	AsinAsinsInfo *AsinsInfo `boil:"AsinAsinsInfo" json:"AsinAsinsInfo" toml:"AsinAsinsInfo" yaml:"AsinAsinsInfo"`
}

// NewStruct creates a new relationship struct
func (*spapiPriceR) NewStruct() *spapiPriceR {
	return &spapiPriceR{}
}

func (r *spapiPriceR) GetAsinAsinsInfo() *AsinsInfo {
	if r == nil {
		return nil
	}
	return r.AsinAsinsInfo
}

// spapiPriceL is where Load methods for each relationship are stored.
type spapiPriceL struct{}

var (
	spapiPriceAllColumns            = []string{"asin", "price", "modified"}
	spapiPriceColumnsWithoutDefault = []string{"asin"}
	spapiPriceColumnsWithDefault    = []string{"price", "modified"}
	spapiPricePrimaryKeyColumns     = []string{"asin"}
	spapiPriceGeneratedColumns      = []string{}
)

type (
	// SpapiPriceSlice is an alias for a slice of pointers to SpapiPrice.
	// This should almost always be used instead of []SpapiPrice.
	SpapiPriceSlice []*SpapiPrice
	// SpapiPriceHook is the signature for custom SpapiPrice hook methods
	SpapiPriceHook func(context.Context, boil.ContextExecutor, *SpapiPrice) error

	spapiPriceQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	spapiPriceType                 = reflect.TypeOf(&SpapiPrice{})
	spapiPriceMapping              = queries.MakeStructMapping(spapiPriceType)
	spapiPricePrimaryKeyMapping, _ = queries.BindMapping(spapiPriceType, spapiPriceMapping, spapiPricePrimaryKeyColumns)
	spapiPriceInsertCacheMut       sync.RWMutex
	spapiPriceInsertCache          = make(map[string]insertCache)
	spapiPriceUpdateCacheMut       sync.RWMutex
	spapiPriceUpdateCache          = make(map[string]updateCache)
	spapiPriceUpsertCacheMut       sync.RWMutex
	spapiPriceUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var spapiPriceAfterSelectHooks []SpapiPriceHook

var spapiPriceBeforeInsertHooks []SpapiPriceHook
var spapiPriceAfterInsertHooks []SpapiPriceHook

var spapiPriceBeforeUpdateHooks []SpapiPriceHook
var spapiPriceAfterUpdateHooks []SpapiPriceHook

var spapiPriceBeforeDeleteHooks []SpapiPriceHook
var spapiPriceAfterDeleteHooks []SpapiPriceHook

var spapiPriceBeforeUpsertHooks []SpapiPriceHook
var spapiPriceAfterUpsertHooks []SpapiPriceHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SpapiPrice) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SpapiPrice) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SpapiPrice) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SpapiPrice) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SpapiPrice) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SpapiPrice) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SpapiPrice) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SpapiPrice) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SpapiPrice) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range spapiPriceAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSpapiPriceHook registers your hook function for all future operations.
func AddSpapiPriceHook(hookPoint boil.HookPoint, spapiPriceHook SpapiPriceHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		spapiPriceAfterSelectHooks = append(spapiPriceAfterSelectHooks, spapiPriceHook)
	case boil.BeforeInsertHook:
		spapiPriceBeforeInsertHooks = append(spapiPriceBeforeInsertHooks, spapiPriceHook)
	case boil.AfterInsertHook:
		spapiPriceAfterInsertHooks = append(spapiPriceAfterInsertHooks, spapiPriceHook)
	case boil.BeforeUpdateHook:
		spapiPriceBeforeUpdateHooks = append(spapiPriceBeforeUpdateHooks, spapiPriceHook)
	case boil.AfterUpdateHook:
		spapiPriceAfterUpdateHooks = append(spapiPriceAfterUpdateHooks, spapiPriceHook)
	case boil.BeforeDeleteHook:
		spapiPriceBeforeDeleteHooks = append(spapiPriceBeforeDeleteHooks, spapiPriceHook)
	case boil.AfterDeleteHook:
		spapiPriceAfterDeleteHooks = append(spapiPriceAfterDeleteHooks, spapiPriceHook)
	case boil.BeforeUpsertHook:
		spapiPriceBeforeUpsertHooks = append(spapiPriceBeforeUpsertHooks, spapiPriceHook)
	case boil.AfterUpsertHook:
		spapiPriceAfterUpsertHooks = append(spapiPriceAfterUpsertHooks, spapiPriceHook)
	}
}

// One returns a single spapiPrice record from the query.
func (q spapiPriceQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SpapiPrice, error) {
	o := &SpapiPrice{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for spapi_prices")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all SpapiPrice records from the query.
func (q spapiPriceQuery) All(ctx context.Context, exec boil.ContextExecutor) (SpapiPriceSlice, error) {
	var o []*SpapiPrice

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SpapiPrice slice")
	}

	if len(spapiPriceAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all SpapiPrice records in the query.
func (q spapiPriceQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count spapi_prices rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q spapiPriceQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if spapi_prices exists")
	}

	return count > 0, nil
}

// AsinAsinsInfo pointed to by the foreign key.
func (o *SpapiPrice) AsinAsinsInfo(mods ...qm.QueryMod) asinsInfoQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"asin\" = ?", o.Asin),
	}

	queryMods = append(queryMods, mods...)

	return AsinsInfos(queryMods...)
}

// LoadAsinAsinsInfo allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (spapiPriceL) LoadAsinAsinsInfo(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSpapiPrice interface{}, mods queries.Applicator) error {
	var slice []*SpapiPrice
	var object *SpapiPrice

	if singular {
		var ok bool
		object, ok = maybeSpapiPrice.(*SpapiPrice)
		if !ok {
			object = new(SpapiPrice)
			ok = queries.SetFromEmbeddedStruct(&object, &maybeSpapiPrice)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", object, maybeSpapiPrice))
			}
		}
	} else {
		s, ok := maybeSpapiPrice.(*[]*SpapiPrice)
		if ok {
			slice = *s
		} else {
			ok = queries.SetFromEmbeddedStruct(&slice, maybeSpapiPrice)
			if !ok {
				return errors.New(fmt.Sprintf("failed to set %T from embedded struct %T", slice, maybeSpapiPrice))
			}
		}
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &spapiPriceR{}
		}
		args = append(args, object.Asin)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &spapiPriceR{}
			}

			for _, a := range args {
				if a == obj.Asin {
					continue Outer
				}
			}

			args = append(args, obj.Asin)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`asins_info`),
		qm.WhereIn(`asins_info.asin in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load AsinsInfo")
	}

	var resultSlice []*AsinsInfo
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice AsinsInfo")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for asins_info")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for asins_info")
	}

	if len(asinsInfoAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.AsinAsinsInfo = foreign
		if foreign.R == nil {
			foreign.R = &asinsInfoR{}
		}
		foreign.R.AsinSpapiPrice = object
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.Asin == foreign.Asin {
				local.R.AsinAsinsInfo = foreign
				if foreign.R == nil {
					foreign.R = &asinsInfoR{}
				}
				foreign.R.AsinSpapiPrice = local
				break
			}
		}
	}

	return nil
}

// SetAsinAsinsInfo of the spapiPrice to the related item.
// Sets o.R.AsinAsinsInfo to related.
// Adds o to related.R.AsinSpapiPrice.
func (o *SpapiPrice) SetAsinAsinsInfo(ctx context.Context, exec boil.ContextExecutor, insert bool, related *AsinsInfo) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"spapi_prices\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"asin"}),
		strmangle.WhereClause("\"", "\"", 2, spapiPricePrimaryKeyColumns),
	)
	values := []interface{}{related.Asin, o.Asin}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.Asin = related.Asin
	if o.R == nil {
		o.R = &spapiPriceR{
			AsinAsinsInfo: related,
		}
	} else {
		o.R.AsinAsinsInfo = related
	}

	if related.R == nil {
		related.R = &asinsInfoR{
			AsinSpapiPrice: o,
		}
	} else {
		related.R.AsinSpapiPrice = o
	}

	return nil
}

// SpapiPrices retrieves all the records using an executor.
func SpapiPrices(mods ...qm.QueryMod) spapiPriceQuery {
	mods = append(mods, qm.From("\"spapi_prices\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"spapi_prices\".*"})
	}

	return spapiPriceQuery{q}
}

// FindSpapiPrice retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSpapiPrice(ctx context.Context, exec boil.ContextExecutor, asin string, selectCols ...string) (*SpapiPrice, error) {
	spapiPriceObj := &SpapiPrice{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"spapi_prices\" where \"asin\"=$1", sel,
	)

	q := queries.Raw(query, asin)

	err := q.Bind(ctx, exec, spapiPriceObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from spapi_prices")
	}

	if err = spapiPriceObj.doAfterSelectHooks(ctx, exec); err != nil {
		return spapiPriceObj, err
	}

	return spapiPriceObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SpapiPrice) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no spapi_prices provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(spapiPriceColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	spapiPriceInsertCacheMut.RLock()
	cache, cached := spapiPriceInsertCache[key]
	spapiPriceInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			spapiPriceAllColumns,
			spapiPriceColumnsWithDefault,
			spapiPriceColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(spapiPriceType, spapiPriceMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(spapiPriceType, spapiPriceMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"spapi_prices\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"spapi_prices\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into spapi_prices")
	}

	if !cached {
		spapiPriceInsertCacheMut.Lock()
		spapiPriceInsertCache[key] = cache
		spapiPriceInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the SpapiPrice.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SpapiPrice) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	spapiPriceUpdateCacheMut.RLock()
	cache, cached := spapiPriceUpdateCache[key]
	spapiPriceUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			spapiPriceAllColumns,
			spapiPricePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update spapi_prices, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"spapi_prices\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, spapiPricePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(spapiPriceType, spapiPriceMapping, append(wl, spapiPricePrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update spapi_prices row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for spapi_prices")
	}

	if !cached {
		spapiPriceUpdateCacheMut.Lock()
		spapiPriceUpdateCache[key] = cache
		spapiPriceUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q spapiPriceQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for spapi_prices")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for spapi_prices")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SpapiPriceSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), spapiPricePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"spapi_prices\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, spapiPricePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in spapiPrice slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all spapiPrice")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SpapiPrice) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no spapi_prices provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(spapiPriceColumnsWithDefault, o)

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

	spapiPriceUpsertCacheMut.RLock()
	cache, cached := spapiPriceUpsertCache[key]
	spapiPriceUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			spapiPriceAllColumns,
			spapiPriceColumnsWithDefault,
			spapiPriceColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			spapiPriceAllColumns,
			spapiPricePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert spapi_prices, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(spapiPricePrimaryKeyColumns))
			copy(conflict, spapiPricePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"spapi_prices\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(spapiPriceType, spapiPriceMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(spapiPriceType, spapiPriceMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert spapi_prices")
	}

	if !cached {
		spapiPriceUpsertCacheMut.Lock()
		spapiPriceUpsertCache[key] = cache
		spapiPriceUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single SpapiPrice record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SpapiPrice) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no SpapiPrice provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), spapiPricePrimaryKeyMapping)
	sql := "DELETE FROM \"spapi_prices\" WHERE \"asin\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from spapi_prices")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for spapi_prices")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q spapiPriceQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no spapiPriceQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from spapi_prices")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for spapi_prices")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SpapiPriceSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(spapiPriceBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), spapiPricePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"spapi_prices\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, spapiPricePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from spapiPrice slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for spapi_prices")
	}

	if len(spapiPriceAfterDeleteHooks) != 0 {
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
func (o *SpapiPrice) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSpapiPrice(ctx, exec, o.Asin)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SpapiPriceSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SpapiPriceSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), spapiPricePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"spapi_prices\".* FROM \"spapi_prices\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, spapiPricePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SpapiPriceSlice")
	}

	*o = slice

	return nil
}

// SpapiPriceExists checks if the SpapiPrice row exists.
func SpapiPriceExists(ctx context.Context, exec boil.ContextExecutor, asin string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"spapi_prices\" where \"asin\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, asin)
	}
	row := exec.QueryRowContext(ctx, sql, asin)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if spapi_prices exists")
	}

	return exists, nil
}

// Exists checks if the SpapiPrice row exists.
func (o *SpapiPrice) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return SpapiPriceExists(ctx, exec, o.Asin)
}
