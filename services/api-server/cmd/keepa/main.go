package main

import (
	"context"
	"log"
	"os"
	"sync"

	"api-server/database"
	"api-server/product"

	"github.com/spf13/cobra"
)

func main() {
	var isAll bool
	var cmdUpdateDrops = &cobra.Command{
		Use: "updateDrops",
		Run: func(cmd *cobra.Command, args []string) {
			updateDrops(isAll)
		},
	}
	cmdUpdateDrops.Flags().BoolVarP(&isAll, "all", "a", false, "update all record")

	var rootCmd = &cobra.Command{
		Use: "keepa",
	}
	rootCmd.AddCommand(cmdUpdateDrops)
	rootCmd.Execute()
}

func updateDrops(isAll bool) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("dsn null value")
	}
	db := database.OpenDB(dsn, false)
	repo := product.KeepaRepository{DB: db}
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	opts := []product.SelectQueryOption{
		product.WithOutPriceData(),
		product.WithOutRankData(),
	}
	if !isAll {
		opts = append(opts, product.OnlyDropsMA7IsNull())
	}

	var asin string
	var keepas product.Keepas
	var err error
	cursor := product.Cursor{
		End: asin,
	}
	limit := 500
	for {
		keepas, cursor, err = repo.GetKeepaWithPaginate(
			context.Background(),
			cursor.End,
			limit,
			opts...,
		)
		log.Printf("cursor end: %s", cursor.End)
		if err != nil {
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
	if len(keepa) == 0 {
		return nil
	}
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
