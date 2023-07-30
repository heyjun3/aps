package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// "github.com/uptrace/bun/extra/bundebug"

	"api-server/product"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("db dsn is null")
	}
	db := product.OpenDB(dsn)
	// db.AddQueryHook(bundebug.NewQueryHook())
	// bundebug.NewQueryHook(bundebug.WithVerbose(true))
	repo := product.ProductRepository{DB: db}

	count := 1500000
	products := make([]product.Product, 0, count)
	for i := 0; i < count; i++ {
		product := product.Product{Asin: "asin_" + fmt.Sprint(i), Filename: "file"}
		products = append(products, product)
	}

	if err := repo.Save(context.Background(), products); err != nil {
		log.Fatal(err)
	}
}
