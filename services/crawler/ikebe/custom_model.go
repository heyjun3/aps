package ikebe

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"golang.org/x/net/context"

	"crawler/models"
)

func NewDBconnection(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conn, nil
}

func NewIkebeProduct(name, productCode, URL string, price int64) *models.IkebeProduct {
	return &models.IkebeProduct{
		Name: null.StringFrom(name),
		ProductCode: productCode,
		URL: null.StringFrom(URL),
		Price: null.Int64From(price),
		ShopCode: "ikebe",
	}
}

func getIkebeProductsByProductCode(ctx context.Context, conn boil.ContextExecutor, codes ...string) ([]*models.IkebeProduct, error){
	var i []interface{}
	for _, code := range codes {
		i = append(i, code)
	}
	return models.IkebeProducts(
		qm.WhereIn("product_code in ?", i...),
	).All(ctx, conn)
}

func bulkUpsertIkebeProducts(conn *sql.DB, products ...*models.IkebeProduct) error{
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
						url = excluded.url RETURNING url, name;`, strings.Join(strs, ",") )
	_, err := conn.Exec(stmt, args...)
	if err != nil {
		return err
	}
	return err
}
