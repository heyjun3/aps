package ikebe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/config"
	"crawler/scrape"
)

func IkebeDatabaseFactory() (*bun.DB, context.Context, error) {
	ctx := context.Background()
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn := scrape.CreateDBConnection(conf.Dsn())
	conn.NewDelete().Model((*IkebeProduct)(nil)).Exec(ctx)
	return conn, ctx, nil
}

func TestGetIkebeProductsByProductCode(t *testing.T) {
	conn, ctx, err := IkebeDatabaseFactory()
	if err != nil {
		return
	}
	p := NewIkebeProduct("test", "test_code", "https://test.com", "", 1111)
	if err := p.Upsert(conn, ctx); err != nil {
		logger.Error("insert error", err)
	}

	t.Run("get products", func(t *testing.T) {

		products, err := GetByProductCodes(conn, ctx, "test_code")

		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(products))
		assert.Equal(t, p, products[0])
	})
}

func TestUpsert(t *testing.T) {
	conn, ctx, err := IkebeDatabaseFactory()
	if err != nil {
		return
	}

	t.Run("upsert ikebe product", func(t *testing.T) {
		p := NewIkebeProduct("test", "test", "test url", "1111", 9000)

		err := p.Upsert(conn, ctx)

		assert.Equal(t, nil, err)
		expectd, _ := GetByProductCodes(conn, ctx, "test")
		assert.Equal(t, (expectd[0]).(*IkebeProduct), p)
	})
}
