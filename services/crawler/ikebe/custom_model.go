package ikebe

import (
	"database/sql"
	"fmt"

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

func getIkebeProductsByProductCode(ctx context.Context, conn boil.ContextExecutor, codes... string) ([]*models.IkebeProduct, error){
	return models.IkebeProducts(
		qm.WhereIn("product_code in ?", codes),
	).All(ctx, conn)
}

func bulkUpsertIkebeProducts(conn *sql.DB) {
	stmt := fmt.Sprintf(`INSERT INTO ikebe_product (name, jan, price, shop_code, product_code, url) 
						VALUES %s ON CONFLICT (shop_code, product_code) DO UPDATE SET 
						name = excluded.name, jan = excluded.jan, price = excluded.price, 
						url = excluded.url RETURNING url, name;`, "($1, $2, $3, $4, $5, $6)")
	_, err := conn.Exec(stmt, "test", "4444", 4444, "ikebe", "test", "url")
	if err != nil {
		fmt.Println(err)
	}
}
