package main

import (
	"api-server/database"
	"api-server/inventory"
	"api-server/spapi"
	"api-server/spapi/point"
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
	inputs := []point.UpdatePointInput{{Sku: "4957054511319-B-35800-20230811", PercentPoint: 7}}
	if err := client.UpdatePoints(inputs); err != nil {
		panic(err)
	}
	// if err := client.UpdatePricing("I1-JGFK-MS72", 55000); err != nil {
		// panic(err)
	// }
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
