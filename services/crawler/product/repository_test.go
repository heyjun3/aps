package product

import (
	"context"
	"fmt"
	"testing"

	"crawler/test/util"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func TestRepository(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T, context.Context, *bun.DB)
	}{{
		name: "test get product",
		fn:   testGetProduct,
	}, {
		name: "test get product by codes",
		fn:   testGetByCodes,
	}, {
		name: "test buld upsert",
		fn:   testBulkUpsert,
	}}

	db, ctx := setupTest()
	for _, tt := range tests {
		tt.fn(t, ctx, db)
	}
}

func setupTest() (*bun.DB, context.Context) {
	conn, ctx := util.DatabaseFactory()
	conn.NewTruncateTable().Model((*Product)(nil)).Exec(ctx)
	seedProducts(ctx, conn)
	return conn, ctx
}

func ptr[T any](a T) *T {
	return &a
}

func seedProducts(ctx context.Context, db *bun.DB) {
	count := 100
	products := make(Products, 0, count)
	for i := 0; i < count; i++ {
		p := NewTestProduct(
			fmt.Sprintf("productName_%d", i),
			fmt.Sprintf("productCode_%d", i),
			fmt.Sprintf("testURL%d", i),
			fmt.Sprintf("jan%d", i),
			"testShop",
			int64((i+1)*1000),
		)
		products = append(products, p)
	}
	repo := NewRepository()
	if err := repo.BulkUpsert(ctx, db, products); err != nil {
		panic(err)
	}
}

func testGetProduct(t *testing.T, ctx context.Context, db *bun.DB) {
	want := &Product{
		Code: Code{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: "productCode_1",
		},
		Name:  "productName_1",
		Jan:   ptr("jan1"),
		Price: int64(2000),
		URL:   "testURL1",
	}
	repo := NewRepository()

	result, err := repo.GetProduct(ctx, db, "testSite", "testShop", "productCode_1")

	assert.NoError(t, err)
	assert.Equal(t, want, result)

	result, err = repo.GetProduct(ctx, db, "testSite", "testShop", "nonExistsCode")

	assert.Error(t, err)
	assert.Equal(t, &Product{}, result)
}

func testGetByCodes(t *testing.T, ctx context.Context,
	db *bun.DB) {
	want := Products{
		&Product{
			Code: Code{
				SiteCode:    "testSite",
				ShopCode:    "testShop",
				ProductCode: "productCode_1",
			},
			Name:  "productName_1",
			Jan:   ptr("jan1"),
			Price: int64(2000),
			URL:   "testURL1",
		},
		&Product{
			Code: Code{
				SiteCode:    "testSite",
				ShopCode:    "testShop",
				ProductCode: "productCode_10",
			},
			Name:  "productName_10",
			Jan:   ptr("jan10"),
			Price: int64(11000),
			URL:   "testURL10",
		},
	}
	repo := NewRepository()

	result, err := repo.GetByCodes(ctx, db,
		[]Code{
			{
				SiteCode:    "testSite",
				ShopCode:    "testShop",
				ProductCode: "productCode_1",
			},
			{
				SiteCode:    "testSite",
				ShopCode:    "testShop",
				ProductCode: "productCode_10",
			},
		})

	assert.NoError(t, err)
	assert.Equal(t, want, result)

	result, err = repo.GetByCodes(ctx, db,
		[]Code{{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: "nonExistsProductCode",
		}})

	assert.NoError(t, err)
	assert.Equal(t, Products(nil), result)
}

func testBulkUpsert(t *testing.T, ctx context.Context, db *bun.DB) {
	repo := NewRepository()
	products := Products{
		&Product{
			Code: Code{
				SiteCode:    "testSite",
				ShopCode:    "testShop",
				ProductCode: "productCode_1",
			},
			Name:  "productName_1",
			Jan:   ptr("jan1"),
			Price: int64(2000),
			URL:   "testURL1",
		},
		&Product{
			Code: Code{
				SiteCode:    "testSite",
				ShopCode:    "testShop",
				ProductCode: "productCode_10",
			},
			Name:  "productName_10",
			Jan:   ptr("jan10"),
			Price: int64(11000),
			URL:   "testURL10",
		},
	}
	err := repo.BulkUpsert(ctx, db, products)

	assert.NoError(t, err)
}
