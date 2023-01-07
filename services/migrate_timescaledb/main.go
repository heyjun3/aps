package main

import (
	"fmt"
	"context"
	"strconv"
	"migrate_timescaledb/app/models"
	"migrate_timescaledb/app/connection"
	_ "github.com/lib/pq"
)

func main() {
	count, err := models.AsinsInfos().Count(context.Background(), connection.DbConnection)
	if err != nil {
		fmt.Printf("Not get count %v", err)
		return
	}
	fmt.Println(fmt.Sprintf("Count: %s", strconv.FormatInt(count, 10)))
}