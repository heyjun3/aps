package migrate

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"migrate_timescaledb/app/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
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

func getMapKeys(m map[string]float64) ([]int, error) {
	keys := []int{}
	for k := range m {
		i, err := strconv.Atoi(k)
		if err != nil {
			fmt.Printf("conver key error isn't valid value: %v", err)
			return nil, err
		}

		keys = append(keys, i)
	}
	sort.Ints(keys)
	return keys, nil
}

func ConvKeepaProductToAsinsInfo(p *models.KeepaProduct) ([]models.AsinsInfoTime, error) {

	var asinInfos []models.AsinsInfoTime
	prices := make(map[string]float64)
	if err := json.Unmarshal(p.PriceData.JSON, &prices); err != nil {
		fmt.Println("price data unmarshal error")
		return nil, err
	}
	times, err := getMapKeys(prices)
	if err != nil {
		fmt.Printf("prices map conver key error")
		return nil, err
	}
	for _, time := range times {
		ts := strconv.Itoa(time)
		t, err := convKeepaTimeToTime(ts)
		if err != nil {
			fmt.Printf("action=ConvKeepaProductToAsinsInfo keepa time convert error value %v", err)
			return nil, err
		}
		data := models.AsinsInfoTime{Time: *t, Asin: p.Asin, Price: null.NewInt(int(prices[ts]), true)}
		asinInfos = append(asinInfos, data)
	}

	ranks := make(map[string]float64)
	if err := json.Unmarshal(p.RankData.JSON, &ranks); err != nil {
		fmt.Println("rank data unmarshal error")
		return nil, err
	}
	times, err = getMapKeys(ranks)
	if err != nil {
		fmt.Printf("ranks map conver key error")
		return nil, err
	}
	for _, time := range times{
		ts := strconv.Itoa(time)
		t, err := convKeepaTimeToTime(ts)
		if err != nil {
			fmt.Printf("action=ConvKeepaProductToAsinsInfo keepa time convert error value %v", err)
			return nil, err
		}
		data := models.AsinsInfoTime{Time: *t, Asin: p.Asin, Rank: null.NewInt(int(ranks[ts]), true)}
		asinInfos = append(asinInfos, data)
	}
	return asinInfos, nil
}

func UpsertAsinsInfoTimes(ctx context.Context, db boil.ContextExecutor, p []models.AsinsInfoTime) error {
	conflictColums := []string{"time", "asin"}

	for _, r := range p {
		upCol := []string{}
		if r.Rank.IsZero() == false {
			upCol = append(upCol, "rank")
		}
		if r.Price.IsZero() == false {
			upCol = append(upCol, "price")
		}
		updateColumns := boil.Whitelist(upCol...)

		err := r.Upsert(ctx, db, true, conflictColums, updateColumns, boil.Infer())
		if err != nil {
			fmt.Printf("AsinsInfoTime Upsert error: %v, value: %v", err, r)
			return err
		}
	}
	return nil
}
