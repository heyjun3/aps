package scrape

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	// "reflect"

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

func (i Product) GenerateMessage(filename string) ([]byte, error) {
	message := message{
		Filename: filename,
		Jan:      i.Jan,
		Price:    i.Price,
		URL:      i.URL,
	}
	if err := message.validation(); err != nil {
		return nil, err
	}
	return json.Marshal(message)
}

func (i Product) GetProductCode() string {
	return i.ProductCode
}

func (i Product) GetJan() string {
	if i.Jan == nil {
		return ""
	}
	return *i.Jan
}

func (i Product) GetURL() string {
	return i.URL
}

func (i Product) IsValidJan() bool {
	return i.Jan != nil
}

func (i *Product) SetJan(jan string) {
	i.Jan = &jan
}

type Products []IProduct

func GetByProductCodes(p IProduct) (func(*bun.DB, context.Context, ...string)(Products, error)) {
	return func(conn *bun.DB, ctx context.Context, codes ...string) (Products, error) {
		s := len(codes) * 2
		products := make(Products, s)
		for i := 0; i < s; i++ {
			products[i] = reflect.New(reflect.ValueOf(p).Elem().Type()).Interface().(IProduct)
		}
		
		err := conn.NewSelect().
			Model(&products).
			Where("product_code IN (?)", bun.In(codes)).
			Scan(ctx)
		
		return products, err
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
		return fmt.Errorf("jan is zero value")
	}
	if m.Price == 0 {
		return fmt.Errorf("price is zero value")
	}
	if m.URL == "" {
		return fmt.Errorf("url is zero value")
	}
	return nil
}
