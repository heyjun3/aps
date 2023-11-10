package main

import (
	"api-server/database"
	"api-server/inventory"

	"context"
	"os"
)

func main() {
	in := inventory.NewInventory("asin", "fnsku", "sku", "new", "name", 1)
	in.CurrentPrice = &inventory.CurrentPrice{
		Price: inventory.Price{
			Amount: 100,
		},
	}
	inventories := inventory.Inventories{
		in,
	}
	dsn := os.Getenv("DB_DSN")
	db := database.OpenDB(dsn)
	err := inventory.InventoryRepository{}.Save(context.Background(), db, inventories)
	if err != nil {
		panic(err)
	}
}
