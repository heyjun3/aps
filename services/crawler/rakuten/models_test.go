package rakuten

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/test/util"
)

func NewTestRakutenProduct(name, productCode, url, jan, shopCode string, price, point int64) *RakutenProduct {
	p, _ := NewRakutenProduct(name, productCode, url, jan, shopCode, price, point)
	return p
}

func TestGetRakutenProducts(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*RakutenProduct)(nil))
	s := NewScrapeService()

	type args struct {
		conn  *bun.DB
		ctx   context.Context
		codes []string
	}
	tests := []struct {
		name      string
		args      args
		want      scrape.Products
		wantError bool
	}{{
		name: "get products",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"test", "code", "test_code"},
		},
		want: scrape.Products{
			NewTestRakutenProduct("name", "test", "http://", "4444", "rakuten", 9900, 0),
			NewTestRakutenProduct("name", "code", "http://", "4444444", "rakuten", 9900, 0),
		},
		wantError: false,
	}, {
		name: "get no products",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"aaaate", "dddd"},
		},
		want:      scrape.Products(nil),
		wantError: false,
	}}

	preProducts := scrape.Products{
		NewTestRakutenProduct("name", "test", "http://", "4444", "rakuten", 9900, 0),
		NewTestRakutenProduct("name", "code", "http://", "4444444", "rakuten", 9900, 0),
	}
	err := s.Repo.BulkUpsert(ctx, conn, preProducts)
	if err != nil {
		panic("prerequisites error")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			var products scrape.Products
			for _, code := range tt.args.codes {
				product, err := s.Repo.GetProduct(tt.args.ctx, tt.args.conn, code, "rakuten")
				logger.Error("error", err)
				if err != nil {
					continue
				}
				products = append(products, product)
			}

			assert.Equal(t, tt.want, products)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
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
		interval: daily,
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
