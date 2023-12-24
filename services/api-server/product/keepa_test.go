package product

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"api-server/spapi/price/competitive"
	"api-server/test"
)

func keepaSeed(repo KeepaRepository) error {
	keepas := make([]*Keepa, 100)
	for i := 0; i < 100; i++ {
		keepas[i] = &Keepa{Asin: "asin_" + strconv.Itoa(i)}
	}
	return repo.Save(context.Background(), keepas)
}

func TestKeepaGetByAsins(t *testing.T) {
	db := test.CreateTestDBConnection()
	if err := db.ResetModel(context.Background(), &Keepa{}); err != nil {
		panic(err)
	}
	repo := KeepaRepository{DB: db}
	if err := keepaSeed(repo); err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{name: "keepas get by asins", args: []string{"asin_1"}, want: []string{"asin_1"}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ks, err := repo.GetByAsins(context.Background(), tt.args)

			assert.Equal(t, tt.want, ks.Asins())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestKeepaGetCounts(t *testing.T) {
	db := test.CreateTestDBConnection()
	if err := db.ResetModel(context.Background(), &Keepa{}); err != nil {
		panic(err)
	}
	ks := []*Keepa{
		{Asin: "aaa", Modified: time.Date(2022, 4, 1, 0, 0, 0, 0, time.Local)},
		{Asin: "bbb"},
		{Asin: "ccc"},
	}
	repo := KeepaRepository{DB: db}
	if err := repo.Save(context.Background(), ks); err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		ctx     context.Context
		want    map[string]int
		wantErr bool
	}{
		{"get counts", context.Background(), map[string]int{"total": 3, "modified": 2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetCounts(tt.ctx)

			assert.Equal(t, tt.want, result)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateRenderData(t *testing.T) {
	type args struct {
		ks   Keepas
		data renderDatas
	}
	type want struct {
		ks Keepas
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Chartが追加される",
			args: args{
				ks:   Keepas{{Asin: "asin_1"}},
				data: renderDatas{{asin: "asin_1", price: 1000, rank: 2000}},
			},
			want: want{
				ks: Keepas{{Asin: "asin_1", Charts: ChartData{Data: []Chart{{Date: time.Now().Format("2006-01-02"), Price: 1000, Rank: 2000}}}}},
			},
		},
		{
			name: "Chartが複数になること",
			args: args{
				ks:   Keepas{{Asin: "asin_2", Charts: ChartData{Data: []Chart{{Date: time.Now().Add(-time.Hour * 24).Format("2006-01-02"), Price: 1000, Rank: 2000}}}}},
				data: renderDatas{{asin: "asin_2", price: 3000, rank: 4000}},
			},
			want: want{
				ks: Keepas{{Asin: "asin_2", Charts: ChartData{Data: []Chart{{Date: time.Now().Add(-time.Hour * 24).Format("2006-01-02"), Price: 1000, Rank: 2000}, {Date: time.Now().Format("2006-01-02"), Price: 3000, Rank: 4000}}}}},
			},
		},
		{
			name: "対応するAsinがなければChart追加されない",
			args: args{
				ks:   Keepas{{Asin: "asin_3", Prices: make(map[string]float64), Ranks: make(map[string]float64)}},
				data: renderDatas{},
			},
			want: want{
				ks: Keepas{{Asin: "asin_3"}},
			},
		},
		{
			name: "90日より前のChartは削除される",
			args: args{
				ks:   Keepas{{Asin: "asin_4", Charts: ChartData{Data: []Chart{{Date: time.Now().Add(-time.Hour * 24 * 91).Format("2006-01-02"), Price: 1000, Rank: 2000}}}}},
				data: renderDatas{{asin: "asin_4", price: 4000, rank: 5000}},
			},
			want: want{
				ks: Keepas{{Asin: "asin_4", Charts: ChartData{Data: []Chart{{Date: time.Now().Format("2006-01-02"), Price: 4000, Rank: 5000}}}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.ks.UpdateRenderData(tt.args.data)
			for i := range result {
				assert.Equal(t, tt.want.ks[i].Asin, result[i].Asin)
				assert.Equal(t, tt.want.ks[i].Charts, result[i].Charts)
				assert.NotNil(t, result[i].Prices)
				assert.NotNil(t, result[i].Ranks)
			}
		})
	}
}

func TestConvertLandedProducts(t *testing.T) {
	type args struct {
		p competitive.LandedProducts
	}
	type want struct {
		data renderDatas
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "rankとpriceがない場合、-1になる",
			args: args{p: competitive.LandedProducts{{Asin: "asin_1"}}},
			want: want{data: renderDatas{{asin: "asin_1", price: -1, rank: -1}}},
		},
		{
			name: "LandedPriceがあればそのままpriceになる",
			args: args{p: competitive.LandedProducts{{Asin: "asin_2", LandedPrice: &competitive.Price{Amount: 2000}}}},
			want: want{data: renderDatas{{asin: "asin_2", price: 2000, rank: -1}}},
		},
		{
			name: "LandedPriceがなくListingPriceがある時、ListingPriceがpriceになる",
			args: args{p: competitive.LandedProducts{{Asin: "asin_3", ListingPrice: &competitive.Price{Amount: 3000}}}},
			want: want{data: renderDatas{{asin: "asin_3", price: 3000, rank: -1}}},
		},
		{
			name: "salesRankがあればrankになる",
			args: args{p: competitive.LandedProducts{{Asin: "asin_4", SalesRankings: []competitive.SalesRank{{ProductCategoryId: "pc", Rank: 40}}}}},
			want: want{data: renderDatas{{asin: "asin_4", price: -1, rank: 40}}},
		},
		{
			name: "salesRankのProductCateogryIdが数値のみの場合、無視する",
			args: args{p: competitive.LandedProducts{{Asin: "asin_5", SalesRankings: []competitive.SalesRank{{ProductCategoryId: "5555", Rank: 50}}}}},
			want: want{data: renderDatas{{asin: "asin_5", price: -1, rank: -1}}},
		},
		{
			name: "salesRankのProductCateogryIdが空文字の場合、無視する",
			args: args{p: competitive.LandedProducts{{Asin: "asin_6", SalesRankings: []competitive.SalesRank{{ProductCategoryId: "", Rank: 60}}}}},
			want: want{data: renderDatas{{asin: "asin_6", price: -1, rank: -1}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertLandedProducts(tt.args.p)

			assert.Equal(t, tt.want.data, result)
		})
	}
}

func TestKeepatimeToUnix(t *testing.T) {
	keepaTime := int64(6815860)

	unix := KeepaTimeToUnix(keepaTime)

	date := time.Unix(unix, 0).Format("2006-01-02")

	assert.Equal(t, "2023-12-17", date)
}

func TestUnixTimeToKeepaTime(t *testing.T) {
	ttime, err := time.Parse("2006-01-02", "2023-12-17")
	if err != nil {
		panic(err)
	}

	keepaTime := UnixTimeToKeepaTime(ttime.Unix())

	assert.Equal(t, int64(6815520), keepaTime)
}
