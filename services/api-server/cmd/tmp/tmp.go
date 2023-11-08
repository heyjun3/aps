package main

import (
	"api-server/database"
	"api-server/inventory"

	"context"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db := database.OpenDB(dsn)
	if err := database.CreateTable(context.Background(), db, &inventory.TmpInventory{CurrentPrice: &inventory.CurrentPrice{}}); err != nil {
		panic(err)
	}
}
