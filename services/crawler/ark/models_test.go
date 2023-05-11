package ark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/testutil"
)

func TestArkGetByProductCodes(t *testing.T) {
	conn, ctx := testutil.DatabaseFactory()
	conn.ResetModel(ctx, (*ArkProduct)(nil))
	p := scrape.Products{NewArkProduct("test", "test_code", "https://google.com", "", 1111)}
	f := scrape.GetByProductCodes([]*ArkProduct{})

	type args struct {
		conn  *bun.DB
		ctx   context.Context
		codes []string
		f     func(*bun.DB, context.Context, ...string) (scrape.Products, error)
	}
	tests := []struct {
		name    string
		args    args
		want    scrape.Products
		wantErr bool
	}{{
		name: "get procuets",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"test_code"},
			f:     f,
		},
		want:    p,
		wantErr: false,
	}, {
		name: "get products none",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"code", "test"},
			f:     f,
		},
		want:    scrape.Products(nil),
		wantErr: false,
	}}

	p.BulkUpsert(conn, ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := tt.args.f(tt.args.conn, tt.args.ctx, tt.args.codes...)

			assert.Equal(t, tt.want, products)
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestBulkUpsert(t *testing.T) {
	conn, ctx := testutil.DatabaseFactory()
	conn.ResetModel(ctx, ArkProduct{})
	type args struct {
		conn     *bun.DB
		ctx      context.Context
		products scrape.Products
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: "bulkupsert",
		args: args{
			conn: conn,
			ctx:  ctx,
			products: scrape.Products{
				NewArkProduct("test", "test_code", "https://google.com", "", 3333),
				NewArkProduct("test", "test", "https://google.com", "", 11111),
				NewArkProduct("test", "code", "https://google.com", "", 9999),
				NewArkProduct("test1", "code1", "https://google.com", "", 9999),
				NewArkProduct("test2", "code3", "https://google.com", "", 9999),
				NewArkProduct("test3", "code4", "https://google.com", "", 9999),
				NewArkProduct("test4", "code6", "https://google.com", "", 9999),
			},
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.products.BulkUpsert(tt.args.conn, tt.args.ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
