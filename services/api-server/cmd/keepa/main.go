package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

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
	db := database.OpenDB(dsn, false)
	repo := product.KeepaRepository{DB: db}
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	updateDropsMA := func(keepa []*product.Keepa, wg *sync.WaitGroup) error {
		defer wg.Done()
		fmt.Println(keepa[0].Asin)
		for _, k := range keepa {
			if err := k.CalculateRankMA(7); err != nil {
				return err
			}
		}
		if err := repo.UpdateDropsMAV2(ctx, keepa); err != nil {
			return err
		}
		return nil
	}

	if scanningRows {
		rows, err := db.NewSelect().Model((*product.Keepa)(nil)).Order("asin").Rows(context.Background())
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		for rows.Next() {
			keepa := new(product.Keepa)
			if err := db.ScanRow(ctx, rows, keepa); err != nil {
				panic(err)
			}
			if err := updateDropsMA([]*product.Keepa{keepa}, wg); err != nil {
				panic(err)
			}
		}
		if err := rows.Err(); err != nil {
			panic(err)
		}
	} else {
		var keepas product.Keepas
		var err error
		cursor := product.Cursor{
			End: asin,
		}
		limit := 500
		for {
			keepas, cursor, err = repo.GetPageNate(context.Background(), cursor.End, limit, product.LoadingData{})
			if err != nil {
				log.Print(cursor.End)
				log.Fatal(err)
			}
			wg.Add(1)
			go updateDropsMA(keepas, wg)

			if len(keepas) != limit {
				log.Print(len(keepas))
				log.Print("return data len not equal limit")
				break
			}
			log.Print(len(keepas))
		}
		wg.Wait()
		log.Println("done")
	}
}
