// Code generated by SQLBoiler 4.14.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testNetseaProducts(t *testing.T) {
	t.Parallel()

	query := NetseaProducts()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testNetseaProductsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testNetseaProductsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := NetseaProducts().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testNetseaProductsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := NetseaProductSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testNetseaProductsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := NetseaProductExists(ctx, tx, o.ShopCode, o.ProductCode)
	if err != nil {
		t.Errorf("Unable to check if NetseaProduct exists: %s", err)
	}
	if !e {
		t.Errorf("Expected NetseaProductExists to return true, but got false.")
	}
}

func testNetseaProductsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	netseaProductFound, err := FindNetseaProduct(ctx, tx, o.ShopCode, o.ProductCode)
	if err != nil {
		t.Error(err)
	}

	if netseaProductFound == nil {
		t.Error("want a record, got nil")
	}
}

func testNetseaProductsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = NetseaProducts().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testNetseaProductsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := NetseaProducts().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testNetseaProductsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	netseaProductOne := &NetseaProduct{}
	netseaProductTwo := &NetseaProduct{}
	if err = randomize.Struct(seed, netseaProductOne, netseaProductDBTypes, false, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}
	if err = randomize.Struct(seed, netseaProductTwo, netseaProductDBTypes, false, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = netseaProductOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = netseaProductTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := NetseaProducts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testNetseaProductsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	netseaProductOne := &NetseaProduct{}
	netseaProductTwo := &NetseaProduct{}
	if err = randomize.Struct(seed, netseaProductOne, netseaProductDBTypes, false, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}
	if err = randomize.Struct(seed, netseaProductTwo, netseaProductDBTypes, false, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = netseaProductOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = netseaProductTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func netseaProductBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func netseaProductAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaProduct) error {
	*o = NetseaProduct{}
	return nil
}

func testNetseaProductsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &NetseaProduct{}
	o := &NetseaProduct{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, netseaProductDBTypes, false); err != nil {
		t.Errorf("Unable to randomize NetseaProduct object: %s", err)
	}

	AddNetseaProductHook(boil.BeforeInsertHook, netseaProductBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	netseaProductBeforeInsertHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.AfterInsertHook, netseaProductAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	netseaProductAfterInsertHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.AfterSelectHook, netseaProductAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	netseaProductAfterSelectHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.BeforeUpdateHook, netseaProductBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	netseaProductBeforeUpdateHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.AfterUpdateHook, netseaProductAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	netseaProductAfterUpdateHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.BeforeDeleteHook, netseaProductBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	netseaProductBeforeDeleteHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.AfterDeleteHook, netseaProductAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	netseaProductAfterDeleteHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.BeforeUpsertHook, netseaProductBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	netseaProductBeforeUpsertHooks = []NetseaProductHook{}

	AddNetseaProductHook(boil.AfterUpsertHook, netseaProductAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	netseaProductAfterUpsertHooks = []NetseaProductHook{}
}

func testNetseaProductsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testNetseaProductsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(netseaProductColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testNetseaProductsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testNetseaProductsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := NetseaProductSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testNetseaProductsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := NetseaProducts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	netseaProductDBTypes = map[string]string{`Name`: `character varying`, `Jan`: `character varying`, `Price`: `bigint`, `ShopCode`: `character varying`, `ProductCode`: `character varying`, `URL`: `character varying`}
	_                    = bytes.MinRead
)

func testNetseaProductsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(netseaProductPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(netseaProductAllColumns) == len(netseaProductPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testNetseaProductsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(netseaProductAllColumns) == len(netseaProductPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &NetseaProduct{}
	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, netseaProductDBTypes, true, netseaProductPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(netseaProductAllColumns, netseaProductPrimaryKeyColumns) {
		fields = netseaProductAllColumns
	} else {
		fields = strmangle.SetComplement(
			netseaProductAllColumns,
			netseaProductPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := NetseaProductSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testNetseaProductsUpsert(t *testing.T) {
	t.Parallel()

	if len(netseaProductAllColumns) == len(netseaProductPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := NetseaProduct{}
	if err = randomize.Struct(seed, &o, netseaProductDBTypes, true); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert NetseaProduct: %s", err)
	}

	count, err := NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, netseaProductDBTypes, false, netseaProductPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize NetseaProduct struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert NetseaProduct: %s", err)
	}

	count, err = NetseaProducts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}