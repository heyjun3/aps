package main

import (
	"os"

	"crawler/config"
	"crawler/product"
	"crawler/scrape"

	"github.com/uptrace/bun"
)

func modelsToByte(db *bun.DB, models []interface{}) []byte {
	var buf []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model).WithForeignKeys().PartitionBy("LIST (site_code)")
		rawQuery, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		buf = append(buf, rawQuery...)
		buf = append(buf, ";\n"...)
	}
	return buf
}

func main() {
	db := scrape.CreateDBConnection(config.DBDsn)

	models := []interface{}{
		(*product.Product)(nil),
	}

	var buf []byte
	buf = append(buf, modelsToByte(db, models)...)

	os.WriteFile("schema.sql", buf, 0644)
}
