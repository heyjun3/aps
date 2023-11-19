package main

import (
	"api-server/database"
	"api-server/inventory"
	"context"
	"fmt"
	"os"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	db := database.OpenDB(dsn)
	repo := inventory.InventoryRepository{}
	quantity := 0
	iv, err := repo.GetByCondition(context.Background(), db, inventory.Condition{Quantity: &quantity, IsNotOnlyLowestPrice: true})
	fmt.Println(err)
	fmt.Println(len(iv))
}
