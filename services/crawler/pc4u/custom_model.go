package pc4u

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"crawler/models"
	"crawler/scrape"
)

type Pc4uProduct struct {
	models.Pc4uProduct
}

func NewPc4uProduct(name, productCode, url, jan string, price int64) *Pc4uProduct {
	isJan := true
	if jan == "" {
		isJan = false
	}
	return &Pc4uProduct{
		models.Pc4uProduct{
			Name:        null.StringFrom(name),
			Jan:         null.NewString(jan, isJan),
			Price:       null.Int64From(price),
			ShopCode:    "pc4u",
			ProductCode: productCode,
			URL:         null.StringFrom(url),
		},
	}
}

func (p *Pc4uProduct) GenerateMessage(filename string) ([]byte, error) {
	if !p.Jan.Valid {
		return nil, fmt.Errorf("jan code isn't valid %s", p.ProductCode)
	}
	if !p.Price.Valid {
		return nil, fmt.Errorf("price isn't valid %s", p.ProductCode)
	}
	if !p.URL.Valid {
		return nil, fmt.Errorf("url isn't valid %s", p.ProductCode)
	}
	m := scrape.NewMWSSchema(filename, p.Jan.String, p.URL.String, p.Price.Int64)
	message, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return message, err
}

func (p *Pc4uProduct) GetProductCode() string {
	return p.ProductCode
}

func (p *Pc4uProduct) GetJan() string {
	return p.Jan.String
}

func (p *Pc4uProduct) GetURL() string {
	return p.URL.String
}

func (p *Pc4uProduct) IsValidJan() bool {
	return p.Jan.Valid
}

func (p *Pc4uProduct) SetJan(jan string) {
	p.Jan = null.StringFrom(jan)
}

type Pc4uProductRepository struct{}

func (r Pc4uProductRepository) GetByProductCodes(ctx context.Context,
	conn boil.ContextExecutor, codes ...string) (scrape.Products, error) {

	var i []interface{}
	for _, code := range codes {
		i = append(i, code)
	}
	pc4uProducts, err := models.Pc4uProducts(
		qm.WhereIn("product_code in ?", i...)).All(ctx, conn)

	if err != nil {
		return nil, err
	}
	var products scrape.Products
	for _, p := range pc4uProducts {
		products = append(products, &Pc4uProduct{*p})
	}
	return products, nil
}
