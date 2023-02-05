package ikebe

import (
	"context"
	"crawler/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)


func TestGetIkebeProductsByProductCode(t *testing.T) {
	ctx := context.Background()
	conf := NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn, err := NewDBconnection(conf.dsn())
	if err != nil {
		fmt.Println(err)
	}
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS ikebe_product (
        name VARCHAR, 
        jan VARCHAR, 
        price BIGINT, 
        shop_code VARCHAR NOT NULL, 
        product_code VARCHAR NOT NULL, 
        url VARCHAR, 
        PRIMARY KEY (shop_code, product_code));`)
	if err != nil {
		fmt.Println(err)
	}
	p := NewIkebeProduct("test", "test_code", "https://test.com", 1111)
	p.Insert(ctx, conn, boil.Infer())
	defer models.IkebeProducts().DeleteAll(ctx, conn)

	t.Run("get products", func(t *testing.T) {
		products, err := getIkebeProductsByProductCode(ctx, conn, "test_code", "test", "code")

		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(products))
		assert.Equal(t, p, products[0])
	})
}