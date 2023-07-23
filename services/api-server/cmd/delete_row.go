package main

import (
	"context"
	"log"
	"os"

	"api-server/product"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("db dsn is null")
	}
	db := product.OpenDB(dsn)
	repo := product.ProductRepository{DB: db}
	condition := product.NewCondition(200, 2, 0.1)
	if err := repo.DeleteIfCondition(context.Background(), condition); err != nil {
		log.Fatal(err)
	}
}
