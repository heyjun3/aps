package scrape

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func CreateDBConnection(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	conn := bun.NewDB(sqldb, pgdialect.New())
	return conn
}

type IProduct interface {
	GenerateMessage(filename string) ([]byte, error)
	GetProductCode() string
	GetJan() string
	GetURL() string
	GetPrice() int64
	GetShopCode() string
	IsValidJan() bool
	SetJan(string)
}

type Product struct {
	Name        string
	Jan         *string
	Price       int64
	ShopCode    string `bun:"shop_code,pk"`
	ProductCode string `bun:"product_code,pk"`
	URL         string
}

func NewProduct(name, productCode, url, jan, shopCode string, price int64) *Product {
	janPtr := &jan
	if jan == "" {
		janPtr = nil
	}
	return &Product{
		Name:        name,
		Jan:         janPtr,
		Price:       price,
		ShopCode:    shopCode,
		ProductCode: productCode,
		URL:         url,
	}
}

func (p Product) GenerateMessage(filename string) ([]byte, error) {
	message := message{
		Filename: filename,
		Jan:      p.Jan,
		Price:    p.Price,
		URL:      p.URL,
	}
	if err := message.validation(); err != nil {
		return nil, err
	}
	return json.Marshal(message)
}

func (p Product) GetProductCode() string {
	return p.ProductCode
}

func (p Product) GetJan() string {
	if p.Jan == nil {
		return ""
	}
	return *p.Jan
}

func (p Product) GetURL() string {
	return p.URL
}

func (p Product) GetPrice() int64 {
	return p.Price
}

func (p Product) GetShopCode() string{
	return p.ShopCode
}

func (p Product) IsValidJan() bool {
	return p.Jan != nil
}

func (p *Product) SetJan(jan string) {
	if jan != "" {
		p.Jan = &jan
	}
}

func GetProduct(p IProduct) func(*bun.DB, context.Context, string, string) (IProduct, error) {
	return func(conn *bun.DB, ctx context.Context, productCode, shopCode string) (IProduct, error) {
		product := reflect.New(reflect.ValueOf(p).Elem().Type()).Interface().(IProduct)
		err := conn.NewSelect().
			Model(product).
			Where("product_code = ?", productCode).
			Where("shop_code = ?", shopCode).
			Scan(ctx, product)

		return product, err
	}
}

type Products []IProduct

func GetByProductCodes[T IProduct](ps []T) func(*bun.DB, context.Context, ...string) (Products, error) {
	return func(conn *bun.DB, ctx context.Context, codes ...string) (Products, error) {
		products := ps

		err := conn.NewSelect().
			Model(&products).
			Where("product_code IN (?)", bun.In(codes)).
			Order("product_code ASC").
			Scan(ctx, &products)

		var result Products
		for _, p := range products {
			result = append(result, p)
		}

		return result, err
	}
}

func (p Products) BulkUpsert(conn *bun.DB, ctx context.Context) error {
	mapProduct := map[string]IProduct{}
	for _, v := range p {
		mapProduct[v.GetProductCode()] = v
	}
	var products Products
	for _, v := range mapProduct {
		products = append(products, v)
	}

	_, err := conn.NewInsert().
		Model(&products).
		On("CONFLICT (shop_code, product_code) DO UPDATE").
		Set(`
			name = EXCLUDED.name,
			jan = EXCLUDED.jan,
			price = EXCLUDED.price,
			url = EXCLUDED.url
		`).
		Returning("NULL").
		Exec(ctx)

	return err
}

func (p Products) getProductCodes() []string {
	var codes []string
	for _, pro := range p {
		codes = append(codes, pro.GetProductCode())
	}
	return codes
}

func (p Products) MapProducts(products Products) Products {
	mapped := map[string]IProduct{}
	for _, v := range products {
		code := v.GetProductCode()
		mapped[code] = v
	}

	for _, v := range p {
		product, exist := mapped[v.GetProductCode()]
		if !exist {
			continue
		}
		v.SetJan(product.GetJan())
	}
	return p
}

type message struct {
	Filename string  `json:"filename"`
	Jan      *string `json:"jan"`
	Price    int64   `json:"cost"`
	URL      string  `json:"url"`
}

func (m *message) validation() error {
	if m.Jan == nil || *m.Jan == "" {
		return fmt.Errorf("jan is zero value. url: %s", m.URL)
	}
	if m.Price == 0 {
		return fmt.Errorf("price is zero value. url:%s", m.URL)
	}
	if m.URL == "" {
		return fmt.Errorf("url is zero value")
	}
	return nil
}
