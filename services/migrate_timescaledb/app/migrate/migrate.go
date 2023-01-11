package migrate

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"

	"github.com/volatiletech/null/v8"
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

func ConvKeepaProductToAsinsInfo(p *models.KeepaProduct) ([]models.AsinsInfoTime, error) {

	var asinInfos []models.AsinsInfoTime
	prices := make(map[string]float64)
	if err := json.Unmarshal(p.PriceData.JSON, &prices); err != nil {
		fmt.Println("price data unmarshal error")
		return nil, err
	}
	for time, price := range prices {
		t, err := convKeepaTimeToTime(time)
		if err != nil {
			fmt.Printf("action=ConvKeepaProductToAsinsInfo keepa time convert error value %v", err)
			return nil, err
		}
		data := models.AsinsInfoTime{Time: *t, Asin: p.Asin, Price: null.NewInt(int(price), true)}
		asinInfos = append(asinInfos, data)
	}

	ranks := make(map[string]float64)
	if err := json.Unmarshal(p.RankData.JSON, &ranks); err != nil {
		fmt.Println("rank data unmarshal error")
		return nil, err
	}
	for time, rank := range ranks {
		t, err := convKeepaTimeToTime(time)
		if err != nil {
			fmt.Printf("action=ConvKeepaProductToAsinsInfo keepa time convert error value %v", err)
			return nil, err
		}
		data := models.AsinsInfoTime{Time: *t, Asin: p.Asin, Rank: null.NewInt(int(rank), true)}
		asinInfos = append(asinInfos, data)
	}
	return asinInfos, nil
}
