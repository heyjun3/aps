package ikebe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/test/util"
)

func TestGetIkebeProductsByProductCode(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*IkebeProduct)(nil))

	s := NewScrapeService()
	p := scrape.Products{
		NewTestIkebeProduct("test", "test_code", "https://test.com", "", 1111),
	}
	err := s.Repo.BulkUpsert(ctx, conn, p)
	if err != nil {
		logger.Error("insert error", err)
	}

	type args struct {
		conn  *bun.DB
		ctx   context.Context
		codes [][]string
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
			codes: [][]string{{"test_code", "ikebe"}},
		},
		want:    p,
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := s.Repo.GetByProductAndShopCodes(tt.args.ctx, tt.args.conn, tt.args.codes...)

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
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*IkebeProduct)(nil))
	s := NewScrapeService()

	t.Run("upsert ikebe product", func(t *testing.T) {
		p := NewTestIkebeProduct("test", "test", "test url", "1111", 9000)

		err := s.Repo.BulkUpsert(ctx, conn, scrape.Products{p})

		assert.Equal(t, nil, err)
		expectd, _ := s.Repo.GetByProductAndShopCodes(ctx, conn, []string{"test", "ikebe"})
		assert.Equal(t, (expectd[0]).(*IkebeProduct), p)
	})
}
