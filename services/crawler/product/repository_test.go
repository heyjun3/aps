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
		test func(*testing.T, context.Context, *bun.DB)
	}{{
		name: "test get product",
		test: testGetProduct,
	}}

	for _, tt := range tests {
		db, ctx := setupTest()
		tt.test(t, ctx, db)
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
	repo := Repository{}
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
	repo := Repository{}

	result, err := repo.GetProduct(ctx, db, "productCode_1", "testShop")
	cast, ok := result.(*Product)

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, want, cast)
}
