package pc4u

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/config"
	"crawler/scrape"
)

func Pc4uDatabaseFactory() (*bun.DB, context.Context, error){
	ctx := context.Background()
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn := scrape.CreateDBConnection(conf.Dsn())
	_, err := conn.NewCreateTable().Model((*Pc4uProduct)(nil)).Exec(ctx)
	if err != nil {
		return nil, nil, err
	}
	conn.NewDelete().Model((*Pc4uProduct)(nil)).Exec(ctx)
	return conn, ctx, nil
}

func TestGetPc4uProductsByProductCode(t *testing.T) {
	conn, ctx, err := Pc4uDatabaseFactory()
	if err != nil {
		return
	}
	p := NewPc4uProduct("test", "test_code", "https://google.com", "", 7777)
	repo := Pc4uProductRepository{}
	if err := repo.Upsert(conn, ctx, p); err != nil {
		logger.Error("insert error", err)
	}

	t.Run("get products", func(t *testing.T) {
		
		products, err := repo.GetByProductCodes(conn, ctx, "test_code")

		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(products))
		assert.Equal(t, p, products[0])
	})
}

func TestUpsert(t *testing.T) {
	conn, ctx, err := Pc4uDatabaseFactory()
	if err != nil {
		return 
	}
	repo := Pc4uProductRepository{}
	t.Run("upsert pc4u product", func(t *testing.T) {
		p := NewPc4uProduct("test", "test", "test url", "1111", 9000)

		err := repo.Upsert(conn, ctx, p)

		assert.Equal(t, nil, err)
		expectd, _ := repo.GetByProductCodes(conn, ctx, "test")
		assert.Equal(t, (expectd[0]).(*Pc4uProduct), p)
	})
}
