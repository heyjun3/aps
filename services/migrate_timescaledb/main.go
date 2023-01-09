package main

import (
	"context"
	"encoding/json"
	"fmt"
	"migrate_timescaledb/app/connection"
	"migrate_timescaledb/app/models"

	_ "github.com/lib/pq"
)

type detailData struct {
	Date string `json:"date"`
	Rank float64 `json:"rank"`
	Price float64 `json:"price"`
}

type RenderData struct {
	Data []detailData `json:"data"`
}

type Data struct {
	Price map[string]float64
}

func main() {
	asin := "B07MTRXVR7"
	product, err := models.FindKeepaProduct(context.Background(), connection.DbConnection, asin)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(product.Asin)
	fmt.Println(product.SalesDrops90)
	fmt.Println(product.Created)
	fmt.Println(product.Modified)
	d := make(map[string]float64)
	err = json.Unmarshal(product.PriceData.JSON, &d)
	for k, v := range d {
		fmt.Println(k, v)
	}
	m := 11111
	fmt.Println(&m)
	fmt.Println(*&m)
}