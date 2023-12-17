package main

import (
	"context"
	"log"
	"os"

	"api-server/database"
	"api-server/product"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("db dsn is null")
	}
	db := database.OpenDB(dsn, true)
	repo := product.ProductRepository{DB: db}

	for _, fn := range []func(ctx context.Context) error{repo.RefreshGeneratedColumns, repo.DeleteIfConditionWithKeepa} {
		if err := fn(context.Background()); err != nil {
			panic(err)
		}
	}

	condition := product.NewCondition(200, 2, 0.1)
	if err := repo.DeleteIfCondition(context.Background(), condition); err != nil {
		log.Fatal(err)
	}
}
