package ikebe

import (
	"fmt"
	"database/sql"

	"github.com/volatiletech/null/v8"
	_ "github.com/lib/pq"

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

func bulkUpsertIkebeProducts(conn *sql.DB) {
	stmt := fmt.Sprintf(`INSERT INTO ikebe_product (name, jan, price, shop_code, product_code, url) 
						VALUES %s ON CONFLICT (shop_code, product_code) DO UPDATE SET 
						name = excluded.name, jan = excluded.jan, price = excluded.price, 
						url = excluded.url RETURNING url, name;`, "($1, $2, $3, $4, $5, $6)")
	// p, err := conn.Exec(stmt, "test", "4444", 4444, "ikebe", "test", "url")
	rows, err := conn.Query(stmt, "test", sql.NullString{}, 1111, "t", "t", "url")
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var name *string
		rows.Scan(&name)
		fmt.Println(*name)
	}
}
