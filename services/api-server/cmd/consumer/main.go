package main

import (
	"log"
	"os"

	"api-server/database"
	"api-server/mq"
	"api-server/product"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("dsn null value")
	}
	db := database.OpenDB(dsn, false)
	keepaService := product.NewKeepaService(db)
	mq.Consume(keepaService.UpdateRenderData, "chart")
}
