package product

import (
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

		if fieldName == "Jan" {
			continue
		}

		if isSet := field.IsValid() && !field.IsZero(); !isSet {
			err = fmt.Errorf("%s is not set; ", fieldName)
			return err
		}
	}
	return nil
}
