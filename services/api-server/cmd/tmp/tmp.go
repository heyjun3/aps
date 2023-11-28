package main

import (
	"api-server/database"
	"api-server/inventory"
	"api-server/spapi"
	"context"
	"encoding/json"
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
	res, err := client.GetLowestPricing([]string{"4562312235052-N-6980-20231105"})
	if err != nil {
		panic(err)
	}
	flag := false
	if flag {
		buf, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(buf))
	}
}

func updatePoint() {
	URL := os.Getenv("SPAPI_SERVICE_URL")
	client, err := spapi.NewSpapiClient(URL)
	if err != nil {
		panic(err)
	}
	sku := "4957054511319-B-35800-20230811"
	price := 44000
	point := 7
	low, err := inventory.NewLowestPrice(sku, float64(price), float64(point))
	if err != nil {
		panic(err)
	}
	des, err := inventory.NewDesiredPrice(&sku, &price, &point, *low)
	if err != nil {
		panic(err)
	}
	fmt.Println(des)
	inputs := inventory.DesiredPrices{des}
	// inputs := []point.UpdatePointInput{{Sku: "4957054511319-B-35800-20230811", PercentPoint: 7}}
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
