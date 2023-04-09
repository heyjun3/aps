package pc4u

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/testutils"
)

func TestGetPc4uProductsByProductCode(t *testing.T) {
	conn, ctx := testutils.DatabaseFactory()
	conn.ResetModel(ctx, (*Pc4uProduct)(nil))
	f := GetByProductCodes
	type args struct {
		conn  *bun.DB
		ctx   context.Context
		f     func(*bun.DB, context.Context, ...string) (scrape.Products, error)
		codes []string
	}
	tests := []struct {
		name    string
		args    args
		want    scrape.Products
		wantErr bool
	}{{
		name: "get product",
		args: args{
			conn:  conn,
			ctx:   ctx,
			f:     f,
			codes: []string{"test_code"},
		},
		want: scrape.Products{
			NewPc4uProduct("test", "test_code", "https://google.com", "", 7777),
		},
		wantErr: false,
	}}

	p := NewPc4uProduct("test", "test_code", "https://google.com", "", 7777)
	err := scrape.Products{p}.BulkUpsert(conn, ctx)
	if err != nil {
		logger.Error("insert error", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pros, err := tt.args.f(tt.args.conn, tt.args.ctx, tt.args.codes...)

			assert.Equal(t, tt.want, pros)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpsert(t *testing.T) {
	conn, ctx := testutils.DatabaseFactory()
	conn.ResetModel(ctx, (*Pc4uProduct)(nil))

	t.Run("upsert pc4u product", func(t *testing.T) {
		p := NewPc4uProduct("test", "test", "test url", "1111", 9000)

		err := scrape.Products{p}.BulkUpsert(conn, ctx)

		assert.Equal(t, nil, err)
		expectd, _ := GetByProductCodes(conn, ctx, "test")
		assert.Equal(t, (expectd[0]).(*Pc4uProduct), p)
	})
}
