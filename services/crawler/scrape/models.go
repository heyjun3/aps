package scrape

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func CreateDBConnection(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	conn := bun.NewDB(sqldb, pgdialect.New())
	return conn
}

type Product interface {
	GenerateMessage(filename string) ([]byte, error)
	GetProductCode() string
	GetJan() string
	GetURL() string
	IsValidJan() bool
	SetJan(string)
}

type Products []Product

func (p Products) getProductCodes() []string {
	var codes []string
	for _, pro := range p {
		codes = append(codes, pro.GetProductCode())
	}
	return codes
}

func (p Products) MapProducts(products Products) Products{
	mapped := map[string]Product{}
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

func NewProduct(name, productCode, url, jan, shopCode string, price int64) *BaseProduct {
	return &BaseProduct{
		Name:        name,
		Jan:         jan,
		Price:       price,
		ShopCode:    shopCode,
		ProductCode: productCode,
		URL:         url,
	}
}

type message struct {
	Filename string `json:"filename"`
	Jan      string `json:"jan"`
	Price    int64  `json:"cost"`
	URL      string `json:"url"`
}

func (m *message) validation() error {
	if m.Jan == "" {
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

type BaseProduct struct {
	Name        string
	Jan         string
	Price       int64
	ShopCode    string `bun:"shop_code,pk"`
	ProductCode string `bun:"product_code,pk"`
	URL         string
}

func (i *BaseProduct) GenerateMessage(filename string) ([]byte, error) {
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

func (i *BaseProduct) GetProductCode() string {
	return i.ProductCode
}

func (i *BaseProduct) GetJan() string {
	return i.Jan
}

func (i *BaseProduct) GetURL() string {
	return i.URL
}

func (i *BaseProduct) IsValidJan() bool {
	return i.Jan != ""
}

func (i *BaseProduct) SetJan(jan string) {
	i.Jan = jan
}
