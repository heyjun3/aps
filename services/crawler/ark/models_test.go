package ark

import (
	"context"
	"crawler/config"
	"crawler/scrape"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func ArkDatabaseFactory() (*bun.DB, context.Context, error) {
	ctx := context.Background()
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn := scrape.CreateDBConnection(conf.Dsn())
	conn.ResetModel(ctx, (*ArkProduct)(nil))
	return conn, ctx, nil
}

func TestArkGetByProductCodes(t *testing.T) {
	conn, ctx, _ := ArkDatabaseFactory()
	p := NewArkProduct("test", "test_code", "https://google.com", "", 1111)
	type args struct {
		conn *bun.DB
		ctx context.Context
		codes []string
	}
	tests := []struct{
		name string
		args args
		want scrape.Products
		wantErr bool
	}{{
		name: "get procuets",
		args: args{
			conn: conn,
			ctx: ctx,
			codes: []string{"test_code"},
		},
		want: scrape.Products{p},
		wantErr: false,
	},{
		name: "get products none",
		args: args {
			conn: conn,
			ctx: ctx,
			codes: []string{"code", "test"},
		},
		want: scrape.Products(nil),
		wantErr: false,
	}}

	p.Upsert(conn, ctx)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			products, err := GetByProductCodes(tt.args.conn, tt.args.ctx, tt.args.codes...)
			
			assert.Equal(t, tt.want, products)
			if err != nil {
				assert.Error(t, err)
			}
		})
	}
}
