package product

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
)

type IProduct interface {
	GenerateMessage(filename string) ([]byte, error)
	GetName() string
	GetProductCode() string
	GetJan() string
	GetURL() string
	GetPrice() int64
	GetShopCode() string
	GetProductAndShopCode() []string
	IsValidJan() bool
	SetJan(string)
}

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

var _ IProduct = &Product{}

func New(product Product) (*Product, error) {
	if *product.Jan == "" {
		product.Jan = nil
	}
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
	message, err := NewMessage(filename, p.URL, p.Jan, p.Price)
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

type Products []IProduct

func ConvToProducts[T IProduct](products []T) Products {
	var result Products
	for i := 0; i < len(products); i++ {
		result = append(result, products[i])
	}
	return result
}

func (p Products) GetProductAndShopCodes() [][]string {
	codes := make([][]string, 0, len(p))
	for _, product := range p {
		codes = append(codes, product.GetProductAndShopCode())
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

func NewMessage(filename, url string, jan *string, price int64) (*message, error) {
	m := message{
		Filename: filename,
		Jan:      jan,
		Price:    price,
		URL:      url,
	}
	if err := m.validation(); err != nil {
		return nil, err
	}
	return &m, nil
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
