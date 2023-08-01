package scrape

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/test/util"
)

func NewTestProduct(name, productCode, url, jan, shopCode string, price int64) (*Product) {
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

	t.Run("happy path", func(t *testing.T) {
		p := Products{
			(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
			(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
			(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
		}

		dbp := Products{
			(NewTestProduct("test", "test", "test", "4444", "test_shop", 4444)),
			(NewTestProduct("test", "test1", "test1", "555", "test_shop", 4444)),
			(NewTestProduct("test", "test2", "test2", "7777", "test_shop", 4444)),
		}

		result := p.MapProducts(dbp)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, (NewTestProduct("test", "test", "http://test.jp", "4444", "test_shop", 1111)), result[0])
		assert.Equal(t, (NewTestProduct("test1", "test1", "http://test.jp", "555", "test_shop", 1111)), result[1])
		assert.Equal(t, (NewTestProduct("test2", "test2", "http://test.jp", "7777", "test_shop", 1111)), result[2])
	})

	t.Run("product is empty", func(t *testing.T) {
		p := Products{}
		dbp := Products{
			(NewTestProduct("test", "test", "test", "11111", "test_shop", 4444)),
			(NewTestProduct("test", "test", "test1", "55555", "test_shop", 4444)),
		}

		result := p.MapProducts(dbp)

		assert.Equal(t, 0, len(result))
		assert.Equal(t, p, result)
	})

	t.Run("db product is empty", func(t *testing.T) {
		p := Products{
			(NewTestProduct("test", "test", "http://test.jp", "", "test_shop", 1111)),
			(NewTestProduct("test1", "test1", "http://test.jp", "", "test_shop", 1111)),
			(NewTestProduct("test2", "test2", "http://test.jp", "", "test_shop", 1111)),
		}
		db := Products{}

		result := p.MapProducts(db)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, p, result)
	})
}

func TestGenerateMessage(t *testing.T) {
	f := "ikebe_20220301_120303"
	t.Run("generate message", func(t *testing.T) {
		p := (NewTestProduct("test", "test", "https://test.com", "4444", "test_shop", 6000))

		m, err := p.GenerateMessage(f)

		assert.Equal(t, nil, err)
		ex := `{"filename":"ikebe_20220301_120303","jan":"4444","cost":6000,"url":"https://test.com"}`
		assert.Equal(t, ex, string(m))
	})

	t.Run("Jan code isn't Valid", func(t *testing.T) {
		p := (NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000))

		m, err := p.GenerateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("Price isn't valid", func(t *testing.T) {
		p := (NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000))
		p.Price = 0

		m, err := p.GenerateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("URL isn't valid", func(t *testing.T) {
		p := (NewTestProduct("TEST", "test", "https://test.com", "", "test_shop", 5000))
		p.URL = ""

		m, err := p.GenerateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})
}
