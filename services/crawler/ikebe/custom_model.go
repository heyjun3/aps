package ikebe

import (
	"github.com/volatiletech/null/v8"

	"crawler/models"
)

func NewIkebeProduct(name, productCode, URL string, price int64) *models.IkebeProduct {
	return &models.IkebeProduct{
		Name: null.StringFrom(name),
		ProductCode: productCode,
		URL: null.StringFrom(URL),
		Price: null.Int64From(price),
	}
}