package scrape

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/test/util"
)

func NewTestProduct(name, productCode, url, jan, shopCode string, price int64) *Product {
	p, _ := NewProduct(name, productCode, url, jan, shopCode, price)
	return p
}

func TestGetProduct(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Product)(nil))
	repo := NewProductRepository(&Product{}, []*Product{})

	type args struct {
		conn        *bun.DB
		ctx         context.Context
		productCode string
		shopCode    string
	}

	tests := []struct {
		name    string
		args    args
		want    *Product
		wantErr bool
	}{{
		name: "test get same product",
		args: args{
			conn:        conn,
			ctx:         ctx,
			productCode: "p1",
			shopCode:    "shop1",
		},
		want:    (NewTestProduct("test", "p1", "google.com", "111", "shop1", 9999)),
		wantErr: false,
	}, {
		name: "get none product",
		args: args{
			conn:        conn,
			ctx:         ctx,
			productCode: "ppp",
			shopCode:    "shop11",
		},
		want:    &Product{},
		wantErr: true,
	}}

	pre := Products{
		(NewTestProduct("name", "test", "https://test.com", "1111", "shop", 1111)),
		(NewTestProduct("name", "test1", "https://test.com", "2222", "shop", 11)),
		(NewTestProduct("name", "test2", "https://test.com", "", "shop", 2)),
		(NewTestProduct("test", "p1", "google.com", "111", "shop1", 9999)),
	}
	repo.BulkUpsert(ctx, conn, pre)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := repo.GetProduct(tt.args.ctx, tt.args.conn, tt.args.productCode, tt.args.shopCode)

			assert.Equal(t, tt.want, p)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBulkUpsert(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Product)(nil))
	repo := NewProductRepository(&Product{}, []*Product{})

	type args struct {
		conn     *bun.DB
		ctx      context.Context
		products Products
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: "success upsert",
		args: args{
			conn: conn,
			ctx:  ctx,
			products: Products{
				(NewTestProduct("test", "test", "https://test.com", "1111", "test", 1111)),
				(NewTestProduct("test", "test1", "https://test.com", "1111", "test", 1111)),
				(NewTestProduct("test", "test2", "https://test.com", "1111", "test", 1111)),
				(NewTestProduct("test", "test3", "https://test.com", "1111", "test", 1111)),
				(NewTestProduct("test", "test4", "https://test.com", "1111", "test", 1111)),
			},
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.BulkUpsert(tt.args.ctx, tt.args.conn, tt.args.products)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Product)(nil))
	repo := NewProductRepository(&Product{}, []*Product{})

	type args struct {
		conn  *bun.DB
		ctx   context.Context
		codes [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    Products
		wantErr bool
	}{{
		name: "get product",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: [][]string{{"test", "shop"}, {"test1", "shop"}, {"test2", "shop"}, {"test3", "shop"}, {"test4", "shop"}},
		},
		want: Products{
			(NewTestProduct("name", "test", "https://test.com", "1111", "shop", 1111)),
			(NewTestProduct("name", "test1", "https://test.com", "2222", "shop", 11)),
			(NewTestProduct("name", "test2", "https://test.com", "", "shop", 2)),
			(NewTestProduct("name", "test3", "https://test.com", "", "shop", 2)),
			(NewTestProduct("name", "test4", "https://test.com", "", "shop", 2)),
		},
		wantErr: false,
	}, {
		name: "get another product",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: [][]string{{"test2", "shop"}},
		},
		want: Products{
			(NewTestProduct("name", "test2", "https://test.com", "", "shop", 2)),
		},
		wantErr: false,
	}, {
		name: "get none product",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: [][]string{{"ttttt", "shop"}, {"eeeee", "shop"}},
		},
		want:    Products(nil),
		wantErr: false,
	}}

	pre := Products{
		(NewTestProduct("name", "test", "https://test.com", "1111", "shop", 1111)),
		(NewTestProduct("name", "test1", "https://test.com", "2222", "shop", 11)),
		(NewTestProduct("name", "test2", "https://test.com", "", "shop", 2)),
		(NewTestProduct("name", "test3", "https://test.com", "", "shop", 2)),
		(NewTestProduct("name", "test4", "https://test.com", "", "shop", 2)),
		(NewTestProduct("name", "test5", "https://test.com", "", "shop", 2)),
		(NewTestProduct("name", "test6", "https://test.com", "", "shop", 2)),
		(NewTestProduct("name", "test7", "https://test.com", "", "shop", 2)),
		(NewTestProduct("name", "test8", "https://test.com", "", "shop", 2)),
	}
	repo.BulkUpsert(ctx, conn, pre)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := repo.GetByProductAndShopCodes(tt.args.ctx, tt.args.conn, tt.args.codes...)

			assert.Equal(t, tt.want, products)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMappingProducts(t *testing.T) {
	type args struct {
		mergeProducts  Products
		targetProducts Products
	}
	type want struct {
		mergedProducts Products
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "merge products",
			args: args{
				mergeProducts: Products{
					(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
					(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
					(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
				},
				targetProducts: Products{
					(NewTestProduct("test", "test", "test", "4444", "test_shop", 4444)),
					(NewTestProduct("test", "test1", "test1", "555", "test_shop", 4444)),
					(NewTestProduct("test", "test2", "test2", "7777", "test_shop", 4444)),
				},
			},
			want: want{
				mergedProducts: Products{
					(NewTestProduct("test", "test", "http://test.jp", "4444", "test_shop", 1111)),
					(NewTestProduct("test1", "test1", "http://test.jp", "555", "test_shop", 1111)),
					(NewTestProduct("test2", "test2", "http://test.jp", "7777", "test_shop", 1111)),
				},
			},
		},
		{
			name: "merge product is empty",
			args: args{
				mergeProducts: Products{},
				targetProducts: Products{
					(NewTestProduct("test", "test", "test", "11111", "test_shop", 4444)),
					(NewTestProduct("test", "test", "test1", "55555", "test_shop", 4444)),
				},
			},
			want: want{mergedProducts: Products{}},
		},
		{
			name: "target product is empty",
			args: args{
				mergeProducts: Products{
					(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
					(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
					(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
				},
				targetProducts: Products{},
			},
			want: want{
				mergedProducts: Products{
					(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
					(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
					(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := tt.args.mergeProducts.MapProducts(tt.args.targetProducts)

			assert.Equal(t, tt.want.mergedProducts, merged)
		})
	}
}

func TestGenerateMessage(t *testing.T) {
	f := "ikebe_20220301_120303"
	type args struct {
		product  *Product
		filename string
	}
	type want struct {
		message string
	}
	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{
		{
			name:  "generate message",
			args:  args{product: NewTestProduct("test", "test", "https://test.com", "4444", "test_shop", 6000), filename: f},
			want:  want{message: `{"filename":"ikebe_20220301_120303","jan":"4444","cost":6000,"url":"https://test.com"}`},
			isErr: false,
		},
		{
			name:  "Jan code isn't Valid",
			args:  args{product: NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000), filename: f},
			want:  want{message: ""},
			isErr: true,
		},
		{
			name:  "Price isn't valid",
			args:  args{product: NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000), filename: f},
			want:  want{message: ""},
			isErr: true,
		},
		{
			name:  "URL isn't valid",
			args:  args{product: NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000), filename: f},
			want:  want{message: ""},
			isErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := tt.args.product.GenerateMessage(tt.args.filename)

			assert.Equal(t, tt.want.message, string(message))
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
