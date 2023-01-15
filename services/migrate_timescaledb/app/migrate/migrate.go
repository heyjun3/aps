package migrate

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"migrate_timescaledb/app/models"
	"migrate_timescaledb/app/connection"

	"github.com/volatiletech/null/v8"
)

func StartMigrate() {
	ctx := context.Background()
	db := connection.DbConnection
	asins, err := getAllAsinsFromKeepaProduct(ctx, db)
	if err != nil {
		return
	}
	for _, asin := range asins[:10] {
		product, err := models.FindKeepaProduct(context.Background(), connection.DbConnection, asin.Asin)
		if err != nil {
			fmt.Printf("get keepa product failed: %v", err)
			return
		}

		infos, err := convKeepaProductToAsinsInfo(product)
		if err != nil {
			fmt.Printf("convert asins info failed: %v", err)
			return
		}
		p := deleteDuplicateAsinsInfoTimes(infos)
		err = balkUpsertAsinsInfoTimes(ctx, connection.DbConnection, p)
	}
}

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


func deleteDuplicateAsinsInfoTimes(p []models.AsinsInfoTime) []models.AsinsInfoTime {
	type key struct {
		time time.Time
		asin string
	}
	m := make(map[key]models.AsinsInfoTime)
	for _, model := range p {
		if v, b := m[key{model.Time, model.Asin}]; b {
			if model.Rank.IsZero() == false {
				v.Rank = model.Rank
			}
			if model.Price.IsZero() == false {
				v.Price = model.Price
			}
			m[key{model.Time, model.Asin}] = v
		} else {
			m[key{model.Time, model.Asin}] = model
		}
	}

	keys := []key{}
	for k := range m {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].time.Before(keys[j].time)
	})

	asinsInfoTimes := []models.AsinsInfoTime{}
	for _, k := range keys {
		asinsInfoTimes = append(asinsInfoTimes, m[k])
	}

	return asinsInfoTimes
}

func convKeepaProductToAsinsInfo(p *models.KeepaProduct) ([]models.AsinsInfoTime, error) {

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

func TmpUpsert() {
	ctx := context.Background()
	p := []models.AsinsInfoTime{
		{
			Time: time.Date(2023, 1, 15, 21, 57, 0, 0, time.Local),
			Asin: "AAAA",
			Price: null.IntFrom(2000),
			Rank: null.IntFrom(1100),
		},
		{
			Time: time.Date(2024, 1, 15, 21, 57, 0, 0, time.Local),
			Asin: "BBBB",
			Price: null.IntFrom(1000),
			Rank: null.IntFrom(1000),
		},
	}
	err := balkUpsertAsinsInfoTimes(ctx, connection.DbConnection, p)
	fmt.Println(err)
}
