package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"api-server/database"
	"api-server/product"

	"github.com/spf13/cobra"
)

func main() {
	var cmdUpdateDrops = &cobra.Command{
		Use: "updateDrops",
		Run: func(cmd *cobra.Command, args []string) {
			updateDrops()
		},
	}

	var rootCmd = &cobra.Command{
		Use: "keepa",
	}
	rootCmd.AddCommand(cmdUpdateDrops)
	rootCmd.Execute()
}

func updateDrops() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("dsn null value")
	}
	db := database.OpenDB(dsn, false)
	repo := product.KeepaRepository{DB: db}
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	var asin string
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
		go updateDropsMA(ctx, repo, keepas, wg)

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

func updateDropsMA(ctx context.Context, repo product.KeepaRepository, keepa []*product.Keepa, wg *sync.WaitGroup) error {
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
