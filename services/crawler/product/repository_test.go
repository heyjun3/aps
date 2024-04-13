package product

import (
	"context"
	"fmt"
	"testing"

	"crawler/scrape"
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
		name: "test get product and shop code",
		fn:   testGetByProductAndShopCodes,
	}, {
		name: "test buld upsert",
		fn:   testBulkUpsert,
	}}

	for _, tt := range tests {
		db, ctx := setupTest()
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
	products := make(scrape.Products, 0, count)
	for i := 0; i < count; i++ {
		p, err := New(Product{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: fmt.Sprintf("productCode_%d", i),
			Name:        fmt.Sprintf("productName_%d", i),
			Jan:         ptr(fmt.Sprintf("jan%d", i)),
			Price:       int64((i + 1) * 1000),
			URL:         fmt.Sprintf("testURL%d", i),
		})
		if err != nil {
			panic(err)
		}
		products = append(products, p)
	}
	repo := NewRepository[*Product]("testSite")
	if err := repo.BulkUpsert(ctx, db, products); err != nil {
		panic(err)
	}
}

func testGetProduct(t *testing.T, ctx context.Context, db *bun.DB) {
	want := &Product{
		SiteCode:    "testSite",
		ShopCode:    "testShop",
		ProductCode: "productCode_1",
		Name:        "productName_1",
		Jan:         ptr("jan1"),
		Price:       int64(2000),
		URL:         "testURL1",
	}
	repo := NewRepository[*Product]("testSite")

	result, err := repo.GetProduct(ctx, db, "productCode_1", "testShop")

	assert.NoError(t, err)
	assert.Equal(t, want, result)

	result, err = repo.GetProduct(ctx, db, "nonExistsCode", "testShop")

	assert.Error(t, err)
	assert.Equal(t, &Product{}, result)
}

func testGetByProductAndShopCodes(t *testing.T, ctx context.Context,
	db *bun.DB) {
	want := scrape.Products{
		&Product{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: "productCode_1",
			Name:        "productName_1",
			Jan:         ptr("jan1"),
			Price:       int64(2000),
			URL:         "testURL1",
		},
		&Product{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: "productCode_10",
			Name:        "productName_10",
			Jan:         ptr("jan10"),
			Price:       int64(11000),
			URL:         "testURL10",
		},
	}
	repo := NewRepository[*Product]("testSite")

	result, err := repo.GetByProductAndShopCodes(ctx, db,
		[][]string{{"productCode_1", "testShop"}, {"productCode_10", "testShop"}}...)

	assert.NoError(t, err)
	assert.Equal(t, want, result)

	result, err = repo.GetByProductAndShopCodes(ctx, db,
		[][]string{{"nonExistsProductCode", "testShop"}}...)

	assert.NoError(t, err)
	assert.Equal(t, scrape.Products(nil), result)
}

func testBulkUpsert(t *testing.T, ctx context.Context, db *bun.DB) {
	repo := NewRepository[*Product]("testSite")
	products := scrape.Products{
		&Product{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: "productCode_1",
			Name:        "productName_1",
			Jan:         ptr("jan1"),
			Price:       int64(2000),
			URL:         "testURL1",
		},
		&Product{
			SiteCode:    "testSite",
			ShopCode:    "testShop",
			ProductCode: "productCode_10",
			Name:        "productName_10",
			Jan:         ptr("jan10"),
			Price:       int64(11000),
			URL:         "testURL10",
		},
	}
	err := repo.BulkUpsert(ctx, db, products)

	assert.NoError(t, err)
}
