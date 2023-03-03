package ikebe

import (
	"context"
	"crawler/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func IkebeProductTableFactory(conn boil.ContextExecutor) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS ikebe_product (
        name VARCHAR, 
        jan VARCHAR, 
        price BIGINT, 
        shop_code VARCHAR NOT NULL, 
        product_code VARCHAR NOT NULL, 
        url VARCHAR, 
        PRIMARY KEY (shop_code, product_code));`)
	if err != nil {
		return err
	}
	return err
}

func TestGetIkebeProductsByProductCode(t *testing.T) {
	ctx := context.Background()
	conf, _ := NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn, err := NewDBconnection(conf.dsn())
	if err != nil {
		fmt.Println(err)
	}
	err = IkebeProductTableFactory(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	models.IkebeProducts().DeleteAll(ctx, conn)
	p := NewIkebeProduct("test", "test_code", "https://test.com", "", 1111)
	pro := models.IkebeProduct(*p)
	pro.Insert(ctx, conn, boil.Infer())

	t.Run("get products", func(t *testing.T) {
		r := IkebeProductRepository{}

		products, err := r.getByProductCodes(ctx, conn, "test_code", "test", "code")

		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(products))
		assert.Equal(t, p, products[0])
	})
}

func TestBulkUpsertIkebeProducts(t *testing.T) {
	conf, _ := NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn, err := NewDBconnection(conf.dsn())
	if err != nil {
		fmt.Println(err)
		return
	}
	err = IkebeProductTableFactory(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	models.IkebeProducts().DeleteAll(ctx, conn)

	t.Run("upsert ikebe products", func(t *testing.T) {
		p := IkebeProducts{
			NewIkebeProduct("test", "test", "https://test.jp", "1111", 9000),
			NewIkebeProduct("test", "test1", "https://test.jp", "1111", 9000),
			NewIkebeProduct("test", "test2", "https://test.jp", "1111", 9000),
		}

		err = p.bulkUpsert(conn)

		assert.Equal(t, nil, err)
		i, _ := models.FindIkebeProduct(ctx, conn, "ikebe", "test1")
		assert.Equal(t, p[1], i)
	})
}
