package shop

import (
	"context"
	"testing"

	"crawler/test/util"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func shopSeeder(repo ShopRepository, db *bun.DB) {
	shops := []*Shop{
		{ID: "test", SiteName: "site_test", Name: "test", URL: "http://test.com", Interval: "daily"},
		{ID: "test1", SiteName: "site_test", Name: "test", URL: "http://test.com", Interval: "weekly"},
	}
	if err := repo.Save(context.Background(), db, shops); err != nil {
		panic(err)
	}
}

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
			err := repo.Save(ctx, db, tt.args.shops)

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
	if err := repo.Save(ctx, db, shops); err != nil {
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
			shops, err := repo.GetAll(ctx, db)

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
	if err := repo.Save(ctx, db, shops); err != nil {
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
			shops, err := ShopRepository{}.GetByInterval(ctx, db, tt.interval)

			assert.Equal(t, tt.shops, shops)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetBySiteName(t *testing.T) {
	db, ctx := util.DatabaseFactory()
	repo := ShopRepository{}
	if err := db.ResetModel(ctx, (*Shop)(nil)); err != nil {
		panic(err)
	}
	shopSeeder(repo, db)

	t.Run("get by site name", func(t *testing.T) {
		shops, err := repo.GetBySiteName(ctx, db, "site_test")
		assert.NoError(t, err)

		assert.Greater(t, len(shops), 0)
		for _, shop := range shops {
			assert.Equal(t, shop.SiteName, "site_test")
		}
	})
}

func TestGetBySiteNameAndInterval(t *testing.T) {
	db, ctx := util.DatabaseFactory()
	repo := ShopRepository{}
	if err := db.ResetModel(ctx, (*Shop)(nil)); err != nil {
		panic(err)
	}
	shopSeeder(repo, db)

	t.Run("get by site name and interval", func(t *testing.T) {
		shops, err := repo.GetBySiteNameAndInterval(
			ctx, db, "site_test", Daily)

		assert.NoError(t, err)
		assert.Greater(t, len(shops), 0)
		for _, shop := range shops {
			assert.Equal(t, shop.SiteName, "site_test")
			assert.Equal(t, shop.Interval, Daily.String())
		}
	})
}
