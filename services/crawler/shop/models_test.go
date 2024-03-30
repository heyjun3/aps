package shop

import (
	"testing"

	"crawler/test/util"
	"github.com/stretchr/testify/assert"
)

func TestShopSave(t *testing.T) {
	db, ctx := util.DatabaseFactory()
	if err := db.ResetModel(ctx, (*Shop)(nil)); err != nil {
		panic(err)
	}
	repo := ShopRepository{}

	type args struct {
		shops []*Shop
	}
	tests := []struct {
		name  string
		args  args
		isErr bool
	}{{
		name: "save shops",
		args: args{
			shops: []*Shop{
				{ID: "test", SiteName: "site_test", Name: "test", URL: "http://test.com"},
				{ID: "test1", SiteName: "site_test", Name: "test", URL: "http://test.com"},
			},
		},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(db, ctx, tt.args.shops)

			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestShopGetAll(t *testing.T) {
	db, ctx := util.DatabaseFactory()
	if err := db.ResetModel(ctx, (*Shop)(nil)); err != nil {
		panic(err)
	}
	repo := ShopRepository{}
	shops := []*Shop{
		{ID: "test", SiteName: "site_test", Name: "test", URL: "http://test.com"},
		{ID: "test1", SiteName: "site_test", Name: "test", URL: "http://test.com"},
	}
	if err := repo.Save(db, ctx, shops); err != nil {
		logger.Error("test", err)
		panic(err)
	}

	tests := []struct {
		name  string
		shops []*Shop
		isErr bool
	}{{
		"get all shops",
		shops,
		false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shops, err := repo.GetAll(db, ctx)

			assert.Equal(t, tt.shops, shops)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetByInterval(t *testing.T) {
	db, ctx := util.DatabaseFactory()
	if err := db.ResetModel(ctx, (*Shop)(nil)); err != nil {
		panic(err)
	}
	repo := ShopRepository{}
	shops := []*Shop{
		{ID: "test", SiteName: "site_test", Name: "test", URL: "http://test.com", Interval: "daily"},
		{ID: "test1", SiteName: "site_test", Name: "test", URL: "http://test.com", Interval: "weekly"},
	}
	if err := repo.Save(db, ctx, shops); err != nil {
		logger.Error("test", err)
		panic(err)
	}

	tests := []struct {
		name     string
		interval Interval
		shops    []Shop
		isErr    bool
	}{{
		name:     "get shop by interval string",
		interval: Daily,
		shops: []Shop{
			{ID: "test", SiteName: "site_test", Name: "test", URL: "http://test.com", Interval: "daily"},
		},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shops, err := ShopRepository{}.GetByInterval(db, ctx, tt.interval)

			assert.Equal(t, tt.shops, shops)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
