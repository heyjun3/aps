package ikebe

import (
	"encoding/json"
	"fmt"
	"time"

	"crawler/models"
)

type ikebeProducts []*models.IkebeProduct

func NewIkebeProducts(products ...*models.IkebeProduct) ikebeProducts{
	p := make(ikebeProducts, len(products))
	copy(p, products)
	return p
}

func (p ikebeProducts) mappingIkebeProducts(productsInDB ikebeProducts) ikebeProducts{
	inDB := map[string]*models.IkebeProduct{}
	for _, v := range productsInDB {
		inDB[v.ProductCode] = v
	}

	for _, v := range p {
		product := inDB[v.ProductCode]
		if product == nil {
			continue
		}
		v.Jan = product.Jan
	}
	return p
}

func generateMessage(p *models.IkebeProduct, filename string) ([]byte, error) {
	if !p.Jan.Valid {
		return nil, fmt.Errorf("jan code isn't valid %s", p.ProductCode)
	}
	if !p.Price.Valid {
		return nil, fmt.Errorf("price isn't valid %s", p.ProductCode)
	}
	if !p.URL.Valid {
		return nil, fmt.Errorf("url isn't valid %s", p.ProductCode)
	}
	m := NewMWSSchema(filename, p.Jan.String, p.URL.String, p.Price.Int64)
	message, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return message, err
}

func timeToStr(t time.Time) string {
	return t.Format("20060102_150405")
}
