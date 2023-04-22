package rakuten

import (
	"context"
	"testing"

	"github.com/uptrace/bun"
	"github.com/stretchr/testify/assert"

	"crawler/scrape"
	"crawler/testutils"
)

func TestGetRakutenProductsByProductCode(t *testing.T) {
	conn, ctx := testutils.DatabaseFactory()
	conn.ResetModel(ctx, (*RakutenProduct)(nil))

	type args struct {
		conn *bun.DB
		ctx context.Context
		codes []string
	}
	tests := []struct {
		name string
		args args
		want scrape.Products
		wantError bool
	}{{
		name: "get products",
		args: args{
			conn: conn,
			ctx: ctx,
			codes: []string{"test", "code", "test_code"},
		},
		want: scrape.Products{
			NewRakutenProduct("name", "test", "http://", "4444", "rakuten", 9900, 0),
			NewRakutenProduct("name", "code", "http://", "4444444", "rakuten", 9900, 0),
		},
		wantError: false,
	}, {
		name: "get no products",
		args: args{
			conn: conn,
			ctx: ctx,
			codes: []string{"aaaate", "dddd"},
		},
		want: scrape.Products(nil),
		wantError: false,
	}}

	preProducts := scrape.Products{
		NewRakutenProduct("name", "test", "http://", "4444", "rakuten", 9900, 1),
		NewRakutenProduct("name", "code", "http://", "4444444", "rakuten", 9900, 1),
	}
	err := preProducts.BulkUpsert(conn, ctx)
	if err != nil {
		panic("prerequisites error")
	}

	for _, tt := range tests {
		products, err := GetByProductCodes(tt.args.conn, tt.args.ctx, tt.args.codes...)

		assert.Equal(t, tt.want, products)
		if tt.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
