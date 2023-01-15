package main

import (
	"fmt"
	"migrate_timescaledb/app/migrate"
	"sort"
	"strconv"

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
	migrate.StartMigrate()
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