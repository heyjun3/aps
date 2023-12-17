package main

import (
	"context"
	// "time"
	// "encoding/json"
	"fmt"
	"log"
	"os"

	"api-server/database"
	"api-server/product"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("dsn null value")
	}
	db := database.OpenDB(dsn, true)
	repo := product.KeepaRepository{DB: db}
	k, err := repo.Get(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// s := make([]map[string]float64, 0, len(k.Prices))
	// for k, v := range k.Prices {
	// 	s = append(s, map[string]float64{k: v})
	// }
	fmt.Println(k.Asin)
	// k.Prices = map[string]float64{"a": 3.0}
	// k.Ranks = map[string]float64{"a": 3.0}
	// k.Charts = product.ChartData{
	// 	Data: []product.Chart{
	// 		{Date: time.Now().Format("2006-01-02"), Rank: 3.0, Price: 3.0},
	// 	},
	// }
	repo.Save(context.Background(), []*product.Keepa{k})
}
