package rakuten

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/scrape"
	"crawler/test/util"
)

func NewTestRakutenProduct(name, productCode, url, jan, shopCode string, price, point int64) *RakutenProduct {
	p, _ := NewRakutenProduct(name, productCode, url, jan, shopCode, price, point)
	return p
}

func TestGetRakutenProducts(t *testing.T) {
	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*RakutenProduct)(nil))
	s := NewScrapeService()

	type args struct {
		conn  *bun.DB
		ctx   context.Context
		codes []string
	}
	tests := []struct {
		name      string
		args      args
		want      scrape.Products
		wantError bool
	}{{
		name: "get products",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"test", "code", "test_code"},
		},
		want: scrape.Products{
			NewTestRakutenProduct("name", "test", "http://", "4444", "rakuten", 9900, 0),
			NewTestRakutenProduct("name", "code", "http://", "4444444", "rakuten", 9900, 0),
		},
		wantError: false,
	}, {
		name: "get no products",
		args: args{
			conn:  conn,
			ctx:   ctx,
			codes: []string{"aaaate", "dddd"},
		},
		want:      scrape.Products(nil),
		wantError: false,
	}}

	preProducts := scrape.Products{
		NewTestRakutenProduct("name", "test", "http://", "4444", "rakuten", 9900, 0),
		NewTestRakutenProduct("name", "code", "http://", "4444444", "rakuten", 9900, 0),
	}
	err := s.Repo.BulkUpsert(ctx, conn, preProducts)
	if err != nil {
		panic("prerequisites error")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			var products scrape.Products
			for _, code := range tt.args.codes {
				product, err := s.Repo.GetProduct(tt.args.ctx, tt.args.conn, code, "rakuten")
				logger.Error("error", err)
				if err != nil {
					continue
				}
				products = append(products, product)
			}

			assert.Equal(t, tt.want, products)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
