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

func testNetseaShops(t *testing.T) {
	t.Parallel()

	query := NetseaShops()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testNetseaShopsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
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

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testNetseaShopsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := NetseaShops().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testNetseaShopsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := NetseaShopSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testNetseaShopsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := NetseaShopExists(ctx, tx, o.ShopID)
	if err != nil {
		t.Errorf("Unable to check if NetseaShop exists: %s", err)
	}
	if !e {
		t.Errorf("Expected NetseaShopExists to return true, but got false.")
	}
}

func testNetseaShopsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	netseaShopFound, err := FindNetseaShop(ctx, tx, o.ShopID)
	if err != nil {
		t.Error(err)
	}

	if netseaShopFound == nil {
		t.Error("want a record, got nil")
	}
}

func testNetseaShopsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = NetseaShops().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testNetseaShopsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := NetseaShops().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testNetseaShopsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	netseaShopOne := &NetseaShop{}
	netseaShopTwo := &NetseaShop{}
	if err = randomize.Struct(seed, netseaShopOne, netseaShopDBTypes, false, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}
	if err = randomize.Struct(seed, netseaShopTwo, netseaShopDBTypes, false, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = netseaShopOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = netseaShopTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := NetseaShops().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testNetseaShopsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	netseaShopOne := &NetseaShop{}
	netseaShopTwo := &NetseaShop{}
	if err = randomize.Struct(seed, netseaShopOne, netseaShopDBTypes, false, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}
	if err = randomize.Struct(seed, netseaShopTwo, netseaShopDBTypes, false, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = netseaShopOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = netseaShopTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func netseaShopBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func netseaShopAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *NetseaShop) error {
	*o = NetseaShop{}
	return nil
}

func testNetseaShopsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &NetseaShop{}
	o := &NetseaShop{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, netseaShopDBTypes, false); err != nil {
		t.Errorf("Unable to randomize NetseaShop object: %s", err)
	}

	AddNetseaShopHook(boil.BeforeInsertHook, netseaShopBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	netseaShopBeforeInsertHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.AfterInsertHook, netseaShopAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	netseaShopAfterInsertHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.AfterSelectHook, netseaShopAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	netseaShopAfterSelectHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.BeforeUpdateHook, netseaShopBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	netseaShopBeforeUpdateHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.AfterUpdateHook, netseaShopAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	netseaShopAfterUpdateHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.BeforeDeleteHook, netseaShopBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	netseaShopBeforeDeleteHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.AfterDeleteHook, netseaShopAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	netseaShopAfterDeleteHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.BeforeUpsertHook, netseaShopBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	netseaShopBeforeUpsertHooks = []NetseaShopHook{}

	AddNetseaShopHook(boil.AfterUpsertHook, netseaShopAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	netseaShopAfterUpsertHooks = []NetseaShopHook{}
}

func testNetseaShopsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testNetseaShopsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(netseaShopColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testNetseaShopsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
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

func testNetseaShopsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := NetseaShopSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testNetseaShopsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := NetseaShops().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	netseaShopDBTypes = map[string]string{`Name`: `character varying`, `ShopID`: `character varying`}
	_                 = bytes.MinRead
)

func testNetseaShopsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(netseaShopPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(netseaShopAllColumns) == len(netseaShopPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testNetseaShopsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(netseaShopAllColumns) == len(netseaShopPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &NetseaShop{}
	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, netseaShopDBTypes, true, netseaShopPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(netseaShopAllColumns, netseaShopPrimaryKeyColumns) {
		fields = netseaShopAllColumns
	} else {
		fields = strmangle.SetComplement(
			netseaShopAllColumns,
			netseaShopPrimaryKeyColumns,
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

	slice := NetseaShopSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testNetseaShopsUpsert(t *testing.T) {
	t.Parallel()

	if len(netseaShopAllColumns) == len(netseaShopPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := NetseaShop{}
	if err = randomize.Struct(seed, &o, netseaShopDBTypes, true); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert NetseaShop: %s", err)
	}

	count, err := NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, netseaShopDBTypes, false, netseaShopPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize NetseaShop struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert NetseaShop: %s", err)
	}

	count, err = NetseaShops().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}