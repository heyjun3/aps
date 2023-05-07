package pc4u

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/testutil"
)

func TestGetPc4uProductsByProductCode(t *testing.T) {
	conn, ctx := testutil.DatabaseFactory()
	conn.ResetModel(ctx, (*Pc4uProduct)(nil))
	f := scrape.GetByProductCodes([]*Pc4uProduct{})
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
			codes: []string{"test_code", "test1", "test2"},
		},
		want: scrape.Products{
			NewPc4uProduct("test", "test1", "https://google.com", "", 1111),
			NewPc4uProduct("test", "test2", "https://google.com", "", 2222),
			NewPc4uProduct("test", "test_code", "https://google.com", "", 7777),
		},
		wantErr: false,
	}}

	ps := scrape.Products{
		NewPc4uProduct("test", "test_code", "https://google.com", "", 7777),
		NewPc4uProduct("test", "code", "https://google.com", "", 7777),
		NewPc4uProduct("test", "test", "https://google.com", "", 7777),
		NewPc4uProduct("test", "test1", "https://google.com", "", 1111),
		NewPc4uProduct("test", "test2", "https://google.com", "", 2222),
	}
	err := ps.BulkUpsert(conn, ctx)
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
	conn, ctx := testutil.DatabaseFactory()
	conn.ResetModel(ctx, (*Pc4uProduct)(nil))

	t.Run("upsert pc4u product", func(t *testing.T) {
		p := NewPc4uProduct("test", "test", "test url", "1111", 9000)

		err := scrape.Products{p}.BulkUpsert(conn, ctx)

		assert.Equal(t, nil, err)
		expectd, _ := scrape.GetByProductCodes([]*Pc4uProduct{})(conn, ctx, "test")
		assert.Equal(t, (expectd[0]).(*Pc4uProduct), p)
	})
}
