package ikebe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/testutils"
)

func TestGetIkebeProductsByProductCode(t *testing.T) {
	conn, ctx := testutils.DatabaseFactory()
	conn.ResetModel(ctx, (*IkebeProduct)(nil))
	f := GetByProductCodes
	p := scrape.Products{
		NewIkebeProduct("test", "test_code", "https://test.com", "", 1111),
	}

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
		name: "get products",
		args: args{
			conn:  conn,
			ctx:   ctx,
			f:     f,
			codes: []string{"test_code"},
		},
		want:    p,
		wantErr: false,
	}}

	err := p.BulkUpsert(conn, ctx)
	if err != nil {
		logger.Error("insert error", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := tt.args.f(tt.args.conn, tt.args.ctx, tt.args.codes...)

			assert.Equal(t, tt.want, products)
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
	conn.ResetModel(ctx, (*IkebeProduct)(nil))

	t.Run("upsert ikebe product", func(t *testing.T) {
		p := NewIkebeProduct("test", "test", "test url", "1111", 9000)

		err := scrape.Products{p}.BulkUpsert(conn, ctx)

		assert.Equal(t, nil, err)
		expectd, _ := GetByProductCodes(conn, ctx, "test")
		assert.Equal(t, (expectd[0]).(*IkebeProduct), p)
	})
}
