package ikebe

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"

	"crawler/models"
	"crawler/scrape"
)



func NewIkebeProduct(name, productCode, url, jan string, price int64) *IkebeProduct {
	isJan := true
	if jan == "" {
		isJan = false
	}
	return &IkebeProduct{
		models.IkebeProduct{
			Name:        null.StringFrom(name),
			Jan:         null.NewString(jan, isJan),
			Price:       null.Int64From(price),
			ShopCode:    "ikebe",
			ProductCode: productCode,
			URL:         null.StringFrom(url),
		},
	}
}

type IkebeProductRepository struct{}

func (r IkebeProductRepository) GetByProductCodes(ctx context.Context, conn boil.ContextExecutor, codes ...string) (scrape.Products, error) {
	var i []interface{}
	for _, code := range codes {
		i = append(i, code)
	}
	ikebeProducts, err := models.IkebeProducts(
		qm.WhereIn("product_code in ?", i...),
	).All(ctx, conn)

	if err != nil {
		return nil, fmt.Errorf("getByProductCodes is failed")
	}

	var products scrape.Products
	for _, p := range ikebeProducts {
		products = append(products, &IkebeProduct{*p})
	}
	return products, nil
}

func (repo IkebeProductRepository) bulkUpsert(conn *sql.DB, products ...*IkebeProduct) error {
	strs := []string{}
	args := []interface{}{}
	for i, p := range products {
		d := i * 6
		str := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", d+1, d+2, d+3, d+4, d+5, d+6)
		strs = append(strs, str)

		args = append(args, p.Name.String)
		args = append(args, p.Jan.String)
		args = append(args, p.Price.Int64)
		args = append(args, p.ShopCode)
		args = append(args, p.ProductCode)
		args = append(args, p.URL.String)
	}
	stmt := fmt.Sprintf(`INSERT INTO ikebe_product (name, jan, price, shop_code, product_code, url) 
						VALUES %s ON CONFLICT (shop_code, product_code) DO UPDATE SET 
						name = excluded.name, jan = excluded.jan, price = excluded.price, 
						url = excluded.url RETURNING url, name;`, strings.Join(strs, ","))
	_, err := conn.Exec(stmt, args...)
	if err != nil {
		return err
	}
	return err
}

type IkebeProduct struct {
	models.IkebeProduct
}

func (p *IkebeProduct) GenerateMessage(filename string) ([]byte, error) {
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

func (p *IkebeProduct) GetProductCode() string {
	return p.ProductCode
}

func (p *IkebeProduct) GetJan() string {
	return p.Jan.String
}

func (p *IkebeProduct) GetURL() string {
	return p.URL.String
}

func (p *IkebeProduct) IsValidJan() bool {
	return p.Jan.Valid
}

func (p *IkebeProduct) SetJan(jan string) {
	p.Jan = null.StringFrom(jan)
}
