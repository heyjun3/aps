package pc4u

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/test/util"
)

func TestGetPc4uProductsByProductCode(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Pc4uProduct)(nil))
	s := NewScrapeService()

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
		name: "get product",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: [][]string{{"test_code", "pc4u"}, {"test1", "pc4u"}, {"test2", "pc4u"}},
		},
		want: scrape.Products{
			util.OmitError(NewPc4uProduct("test", "test1", "https://google.com", "", 1111)),
			util.OmitError(NewPc4uProduct("test", "test2", "https://google.com", "", 2222)),
			util.OmitError(NewPc4uProduct("test", "test_code", "https://google.com", "", 7777)),
		},
		wantErr: false,
	}}

	ps := scrape.Products{
		util.OmitError(NewPc4uProduct("test", "test_code", "https://google.com", "", 7777)),
		util.OmitError(NewPc4uProduct("test", "code", "https://google.com", "", 7777)),
		util.OmitError(NewPc4uProduct("test", "test", "https://google.com", "", 7777)),
		util.OmitError(NewPc4uProduct("test", "test1", "https://google.com", "", 1111)),
		util.OmitError(NewPc4uProduct("test", "test2", "https://google.com", "", 2222)),
	}
	err := s.Repo.BulkUpsert(ctx, conn, ps)
	if err != nil {
		logger.Error("insert error", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pros, err := s.Repo.GetByProductAndShopCodes(tt.args.ctx, tt.args.conn, tt.args.codes...)

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
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Pc4uProduct)(nil))
	s := NewScrapeService()

	t.Run("upsert pc4u product", func(t *testing.T) {
		p := util.OmitError(NewPc4uProduct("test", "test", "test url", "1111", 9000))

		err := s.Repo.BulkUpsert(ctx, conn, scrape.Products{p})

		assert.Equal(t, nil, err)
		expectd, _ := s.Repo.GetByProductAndShopCodes(ctx, conn, []string{"test", "pc4u"})
		assert.Equal(t, (expectd[0]).(*Pc4uProduct), p)
	})
}
