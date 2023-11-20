package main

import (
	"api-server/database"
	"api-server/inventory"
	"api-server/spapi"
	"context"
	"fmt"
	"os"
)

func main() {
	tmpHttp()
}

func tmpHttp() {
	URL := os.Getenv("SPAPI_SERVICE_URL")
	client, err := spapi.NewSpapiClient(URL)
	if err != nil {
		panic(err)
	}
	if err := client.UpdatePricing("I1-JGFK-MS72", 55000); err != nil {
		panic(err)
	}
}

func tmpDatabase() {
	dsn := os.Getenv("DB_DSN")
	db := database.OpenDB(dsn)
	repo := inventory.InventoryRepository{}
	quantity := 0
	iv, err := repo.GetByCondition(context.Background(), db, inventory.Condition{Quantity: &quantity, IsNotOnlyLowestPrice: true})
	fmt.Println(err)
	fmt.Println(len(iv))
}
