package main

import (
	"context"
	// "time"
	// "encoding/json"
	// "fmt"
	"flag"
	"log"
	"os"

	"api-server/database"
	"api-server/product"
)

func main() {
	var asin string
	var scanningRows bool

	flag.StringVar(&asin, "a", "", "start asin")
	flag.BoolVar(&scanningRows, "s", false, "scanning rows")
	flag.Parse()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("dsn null value")
	}
	db := database.OpenDB(dsn, true)
	if scanningRows {
		rows, err := db.NewSelect().Model((*product.Keepa)(nil)).Order("asin").Rows(context.Background())
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			keepa := new(product.Keepa)
			if err := db.ScanRow(context.Background(), rows, keepa); err != nil {
				panic(err)
			}
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}
	} else {
		repo := product.KeepaRepository{DB: db}
		var keepas product.Keepas
		var err error
		cursor := product.Cursor{
			End: asin,
		}
		limit := 100
		for {
			keepas, cursor, err = repo.GetPageNate(context.Background(), cursor.End, limit)
			if err != nil {
				log.Print(cursor.End)
				log.Fatal(err)
			}
			if len(keepas) != limit {
				log.Print(len(keepas))
				log.Print("return data len not equal limit")
				return
			}
			log.Print(len(keepas))
		}

	}
}
