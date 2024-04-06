package product

import (
	"crawler/scrape"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"crawler.products"`
	SiteCode      string `bun:"site_code,pk"`
	ShopCode      string `bun:"shop_code,pk"`
	ProductCode   string `bun:"product_code,pk"`
	Name          string
	Jan           *string
	Price         int64
	URL           string
}

var _ scrape.IProduct = &Product{}

func New(product Product) (*Product, error) {
	if err := product.validateZeroValues(); err != nil {
		return nil, err
	}
	return &product, nil
}

func (p Product) validateZeroValues() (err error) {
	structType := reflect.TypeOf(p)
	structValue := reflect.ValueOf(p)
	fieldsNum := structValue.NumField()

	for i := 0; i < fieldsNum; i++ {
		field := structValue.Field(i)
		fieldName := structType.Field(i).Name

		if fieldName == "Jan" || fieldName == "BaseModel" {
			continue
		}

		if isSet := field.IsValid() && !field.IsZero(); !isSet {
			err = fmt.Errorf("%s is not set; ", fieldName)
			return err
		}
	}
	return nil
}

func (p Product) GenerateMessage(filename string) ([]byte, error) {
	message, err := scrape.NewMessage(filename, p.URL, p.Jan, p.Price)
	if err != nil {
		return nil, err
	}
	return json.Marshal(message)
}

func (p Product) GetName() string {
	return p.Name
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

func (p Product) GetShopCode() string {
	return p.ShopCode
}

func (p Product) GetProductAndShopCode() []string {
	return []string{p.ProductCode, p.ShopCode}
}

func (p Product) IsValidJan() bool {
	return p.Jan != nil
}

func (p *Product) SetJan(jan string) {
	if jan != "" {
		p.Jan = &jan
	}
}
