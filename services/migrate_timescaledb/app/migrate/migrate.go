package migrate

import (
	"encoding/json"
	"context"
	"fmt"
	"time"
	"strconv"

	con "migrate_timescaledb/app/connection"
	"migrate_timescaledb/app/models"
)

func convKeepaTimeToTime(keepa_time string) (*time.Time, error) {
	k, err := strconv.Atoi(keepa_time)
	if err != nil {
		fmt.Printf("keepa time arg is valid %v", keepa_time)
		return &time.Time{}, err
	}
	unix_time := (k + 21564000) * 60
	t := time.Unix(int64(unix_time), 0)
	return &t, nil
}

func ConvKeepaProductToAsinsInfo(asin string) {
	p, err := models.FindKeepaProduct(context.Background(), con.DbConnection, asin)
	if err != nil {
		fmt.Printf("get keepa product failed. argument is %v\n", asin)
		return
	}

	prices := make(map[string]float64)
	if err := json.Unmarshal(p.PriceData.JSON, &prices); err != nil {
		fmt.Println("price data unmarshal error")
		return
	}

	ranks := make(map[string]float64)
	if err := json.Unmarshal(p.RankData.JSON, &ranks); err != nil {
		fmt.Println("rank data unmarshal error")
		return
	}
}
