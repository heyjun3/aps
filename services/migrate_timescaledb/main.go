package main

import (
	"context"
	"fmt"
	"migrate_timescaledb/app/connection"
	"migrate_timescaledb/app/models"
	"migrate_timescaledb/app/migrate"

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
	tmp()
}

func tmp() {
	asin := "B07MTRXVR7"
	product, err := models.FindKeepaProduct(context.Background(), connection.DbConnection, asin)
	if err != nil {
		fmt.Printf("get keepa product failed: %v", err)
		return
	}

	infos, err := migrate.ConvKeepaProductToAsinsInfo(product)
	if err != nil {
		fmt.Printf("convert asins info failed: %v", err)
		return
	}
	fmt.Println(infos)
}