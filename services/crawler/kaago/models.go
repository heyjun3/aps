package kaago

import (
	"fmt"
	"reflect"

	"crawler/product"
)

const (
	siteCode = "kaago"
)

func NewKaagoProduct(name, productCode, url, jan, shopCode string,
	price int64) (*product.Product, error) {
	return product.New(
		siteCode,
		shopCode,
		productCode,
		name,
		jan,
		url,
		price,
	)
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
