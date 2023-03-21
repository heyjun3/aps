package ikebe

import (
	"context"
	"crawler/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"crawler/config"
	"crawler/scrape"
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
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn := scrape.CreateDBConnection(conf.Dsn())
	err := IkebeProductTableFactory(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	models.IkebeProducts().DeleteAll(ctx, conn)
	p := NewIkebeProduct("test", "test_code", "https://test.com", "", 1111)
	repo := IkebeProductRepository{}
	if err := repo.Upsert(conn, ctx, p); err != nil {
		logger.Error("insert error", err)
	}

	t.Run("get products", func(t *testing.T) {
		r := IkebeProductRepository{}

		products, err := r.GetByProductCodes(conn, ctx, "test_code")

		assert.Equal(t, nil, err)
		assert.Equal(t, 1, len(products))
		assert.Equal(t, p, products[0])
	})
}

func TestBulkUpsertIkebeProducts(t *testing.T) {
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn := scrape.CreateDBConnection(conf.Dsn())
	err := IkebeProductTableFactory(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	models.IkebeProducts().DeleteAll(ctx, conn)

	t.Run("upsert ikebe products", func(t *testing.T) {
		p := []*IkebeProduct{
			NewIkebeProduct("test", "test", "https://test.jp", "1111", 9000),
			NewIkebeProduct("test", "test1", "https://test.jp", "1111", 9000),
			NewIkebeProduct("test", "test2", "https://test.jp", "1111", 9000),
		}
		repo := IkebeProductRepository{}
		for _, product := range p {
			repo.Upsert(conn, ctx, product)
		}

		assert.Equal(t, nil, err)
		var i IkebeProduct
		conn.NewSelect().Model(&i).Where("product_code = ? and shop_code = ?", "test1", "ikebe").Scan(ctx)
		assert.Equal(t, *p[1], i)
	})
}

func TestMappingIkebeProducts(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		p := scrape.Products{
			NewIkebeProduct("test", "test", "http://test.jp", "", 1111),
			NewIkebeProduct("test1", "test1", "http://test.jp", "", 1111),
			NewIkebeProduct("test2", "test2", "http://test.jp", "", 1111),
		}

		dbp := scrape.Products{
			NewIkebeProduct("test", "test", "test", "4444", 4000),
			NewIkebeProduct("test", "test1", "test1", "555", 4000),
			NewIkebeProduct("test", "test2", "test2", "7777", 4000),
		}

		result := p.MapProducts(dbp)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, NewIkebeProduct("test", "test", "http://test.jp", "4444", 1111), result[0])
		assert.Equal(t, NewIkebeProduct("test1", "test1", "http://test.jp", "555", 1111), result[1])
		assert.Equal(t, NewIkebeProduct("test2", "test2", "http://test.jp", "7777", 1111), result[2])
	})

	t.Run("product is empty", func(t *testing.T) {
		p := scrape.Products{}
		dbp := scrape.Products{
			NewIkebeProduct("test", "test", "test", "11111", 4000),
			NewIkebeProduct("test", "test", "test1", "55555", 4000),
		}

		result := p.MapProducts(dbp)

		assert.Equal(t, 0, len(result))
		assert.Equal(t, p, result)
	})

	t.Run("db product is empty", func(t *testing.T) {
		p := scrape.Products{
			NewIkebeProduct("test", "test", "http://test.jp", "", 1111),
			NewIkebeProduct("test1", "test1", "http://test.jp", "", 1111),
			NewIkebeProduct("test2", "test2", "http://test.jp", "", 1111),
		}
		db := scrape.Products{}

		result := p.MapProducts(db)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, p, result)
	})
}

func TestGenerateMessage(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		p := NewIkebeProduct("test", "test", "https://test.com", "4444", 6000)
		f := "ikebe_20220301_120303"

		m, err := p.GenerateMessage(f)

		assert.Equal(t, nil, err)
		ex := `{"filename":"ikebe_20220301_120303","jan":"4444","cost":6000,"url":"https://test.com"}`
		assert.Equal(t, ex, string(m))
	})

	t.Run("Jan code isn't Valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", "", 5000)
		f := "ikebe_20220202_020222"

		m, err := p.GenerateMessage(f)
		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("Price isn't valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", "", 5000)
		p.Price = 0
		f := "ikebe_20220202_020222"

		m, err := p.GenerateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("URL isn't valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", "", 5000)
		p.URL = ""
		f := "ikebe_20220202_020222"

		m, err := p.GenerateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})
}