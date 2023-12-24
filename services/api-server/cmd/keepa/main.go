package main

import (
	"context"
	// "time"
	// "encoding/json"
	// "fmt"
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
	var keepas product.Keepas
	var err error
	cursor := product.Cursor{}
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
