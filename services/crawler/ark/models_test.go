package ark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/test/util"
)

func TestArkGetByProductCodes(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*ArkProduct)(nil))
	s := NewScrapeService()
	p := scrape.Products{util.OmitError(NewArkProduct("test", "test_code", "https://google.com", "", 1111))}

	type args struct {
		conn  *bun.DB
		ctx   context.Context
		codes []string
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
		},
		want:    p,
		wantErr: false,
	}, {
		name: "get products none",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"code", "test"},
		},
		want:    scrape.Products(nil),
		wantErr: false,
	}}

	s.Repo.BulkUpsert(ctx, conn, p)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := s.Repo.GetByProductCodes(tt.args.ctx, tt.args.conn, tt.args.codes...)

			assert.Equal(t, tt.want, products)
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestBulkUpsert(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, ArkProduct{})
	s := NewScrapeService()

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
				util.OmitError(NewArkProduct("test", "test_code", "https://google.com", "", 3333)),
				util.OmitError(NewArkProduct("test", "test", "https://google.com", "", 11111)),
				util.OmitError(NewArkProduct("test", "code", "https://google.com", "", 9999)),
				util.OmitError(NewArkProduct("test1", "code1", "https://google.com", "", 9999)),
				util.OmitError(NewArkProduct("test2", "code3", "https://google.com", "", 9999)),
				util.OmitError(NewArkProduct("test3", "code4", "https://google.com", "", 9999)),
				util.OmitError(NewArkProduct("test4", "code6", "https://google.com", "", 9999)),
			},
		},
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Repo.BulkUpsert(tt.args.ctx, tt.args.conn, tt.args.products)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
