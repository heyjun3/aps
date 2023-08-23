package kaago

import (
	"fmt"
	"reflect"

	"github.com/uptrace/bun"

	"crawler/scrape"
)

type KaagoProduct struct {
	bun.BaseModel `bun:"table:kaago_products"`
	scrape.Product
}

func NewKaagoProduct(name, productCode, url, jan, shopCode string, price int64) (*KaagoProduct, error) {
	p, err := scrape.NewProduct(name, productCode, url, jan, shopCode, price)
	if err != nil {
		return nil, err
	}
	return &KaagoProduct{
		Product: *p,
	}, nil
}

type KaagoResp struct {
	CurrentPage int64              `json:"currentPage"`
	ProductList []KaagoRespProduct `json:"resultList"`
}

type KaagoRespProduct struct {
	Name        string `json:"commodityName"`
	Price       int64  `json:"retailPrice"`
	ProductCode string `json:"commodityCode"`
	ShopCode    string `json:"shopCode"`
	URL         string `json:"commodityDetailUrl"`
}

func ValidateKaagoRespProduct(k KaagoRespProduct) (err error) {

	structType := reflect.TypeOf(k)
	structVal := reflect.ValueOf(k)
	fieldsNum := structVal.NumField()

	for i := 0; i < fieldsNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		if isSet := field.IsValid() && !field.IsZero(); !isSet {
			err = fmt.Errorf("%v%s is not set; ", err, fieldName)
		}
	}
	return err
}
