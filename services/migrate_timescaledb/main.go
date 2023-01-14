package main

import (
	"context"
	"fmt"
	"migrate_timescaledb/app/connection"
	"migrate_timescaledb/app/migrate"
	"migrate_timescaledb/app/models"
	"sort"
	"strconv"

	_ "github.com/lib/pq"
	// "github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
	// infos[0].Rank = null.IntFrom(1000)
	var upCol []string
	if infos[0].Rank.IsZero() == false{
		upCol = append(upCol, "rank")
	}
	if infos[0].Price.IsZero() == false {
		upCol = append(upCol, "price")
	}
	updateColumns := boil.Whitelist(upCol...)
	infos[0].Upsert(context.Background(), connection.DbConnection, true, []string{"time", "asin"}, updateColumns, boil.Infer())
	if err != nil {
		fmt.Printf("insert error: %v", err)
	}
}

func tmp2() {
	keys := []string{"81111", "222", "999"}
	var is = []int{}
	for _, k := range keys {
		i, _ := strconv.Atoi(k)
		is = append(is, i)
	}
	sort.Ints(is)
	fmt.Println(is)
}