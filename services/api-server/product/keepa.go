package product

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/uptrace/bun"

	"api-server/spapi/price/competitive"
)

const (
	timeFormat string = "2006-01-02"
)

type Chart struct {
	Date    string  `json:"date" bun:"type:date"`
	Rank    float64 `json:"rank"`
	Price   float64 `json:"price"`
	RankMA7 int     `json:"rank_ma7"`
}

type ChartData struct {
	Data []Chart `json:"data"`
}

func (c *ChartData) filteringPastDays(days int) {
	pastDate := time.Now().Add(-time.Hour * time.Duration(24*days))
	charts := make([]Chart, 0, len(c.Data))
	for _, chart := range c.Data {
		date, err := time.Parse("2006-01-02", chart.Date)
		if err != nil {
			slog.Error("date parse error", "err", err)
			continue
		}
		if date.After(pastDate) {
			charts = append(charts, chart)
		}
	}
	c.Data = charts
}

type Keepa struct {
	bun.BaseModel `bun:"keepa_products"`
	Asin          string             `bun:"asin,pk"`
	Drops         int                `bun:"sales_drops_90"`
	DropsMA7      int                `bun:"drops_ma_7" json:"drops_ma_7"`
	Prices        map[string]float64 `bun:"price_data,type:jsonb"`
	Ranks         map[string]float64 `bun:"rank_data,type:jsonb"`
	Charts        ChartData          `bun:"render_data,type:jsonb"`
	Created       time.Time          `bun:",type:date,nullzero,notnull,default:current_timestamp"`
	Modified      time.Time          `bun:",type:date,nullzero,notnull,default:current_timestamp"`
}

func (k *Keepa) updateRenderData(data renderData, keepaTime, date string) *Keepa {
	if k.Prices == nil {
		k.Prices = make(map[string]float64)
	}
	if k.Ranks == nil {
		k.Ranks = make(map[string]float64)
	}

	price := float64(data.price)
	rank := float64(data.rank)
	k.Prices[keepaTime] = price
	k.Ranks[keepaTime] = rank
	chart := Chart{Date: date, Price: price, Rank: rank}
	k.Charts.Data = append(k.Charts.Data, chart)
	k.Charts.filteringPastDays(90)
	return k
}

func (k *Keepa) calculateDropsMA7() {
	dropsMA7 := 0
	for i := 1; i < len(k.Charts.Data)-1; i++ {
		if k.Charts.Data[i-1].RankMA7 > k.Charts.Data[i].RankMA7 &&
			k.Charts.Data[i].RankMA7 < k.Charts.Data[i+1].RankMA7 {
			dropsMA7 += 1
		}
	}
	k.DropsMA7 = dropsMA7
}

// calculate rank of moving average for a period of days.
func (k *Keepa) CalculateRankMA(days int) error {
	sort.Slice(k.Charts.Data, func(i, j int) bool {
		ti, err := time.Parse(timeFormat, k.Charts.Data[i].Date)
		if err != nil {
			logger.Warn("date string parse error", "date", k.Charts.Data[i].Date)
			return false
		}
		tj, err := time.Parse(timeFormat, k.Charts.Data[j].Date)
		if err != nil {
			logger.Warn("date string parse error", "date", k.Charts.Data[j].Date)
			return false
		}
		return ti.Before(tj)
	})
	for i := days; i < len(k.Charts.Data); i++ {
		var rankSum int
		data := k.Charts.Data[i-days : i]
		for _, chart := range data {
			rankSum += int(chart.Rank)
		}
		k.Charts.Data[i].RankMA7 = int(rankSum / days)
	}
	k.calculateDropsMA7()
	return nil
}

type Keepas []*Keepa

func (k Keepas) Asins() []string {
	asins := make([]string, len(k))
	for i, keepa := range k {
		asins[i] = keepa.Asin
	}
	return asins
}
func (k Keepas) UpdateRenderData(renderDatas renderDatas) Keepas {
	keepaTime := fmt.Sprint(UnixTimeToKeepaTime(time.Now().Unix()))
	date := time.Now().Format("2006-01-02")
	m := renderDatas.Map()
	for _, keepa := range k {
		data := m[keepa.Asin]
		if data == nil {
			continue
		}
		keepa.updateRenderData(*data, keepaTime, date)
	}
	return k
}

type renderData struct {
	asin  string
	price int
	rank  int
}
type renderDatas []*renderData

func (r renderDatas) Map() map[string]*renderData {
	m := make(map[string]*renderData)
	for _, data := range r {
		m[data.asin] = data
	}
	return m
}

func (r renderDatas) Asins() []string {
	asins := make([]string, 0, len(r))
	for _, data := range r {
		asins = append(asins, data.asin)
	}
	return asins
}

func ConvertLandedProducts(products competitive.LandedProducts) renderDatas {
	re := regexp.MustCompile("^[0-9]+$")
	renderDatas := make(renderDatas, 0, len(products))

	for _, product := range products {
		asin := product.Asin
		price := func() int {
			for _, p := range []*competitive.Price{product.LandedPrice, product.ListingPrice} {
				if p != nil {
					return p.Amount
				}
			}
			return -1
		}()
		rank := func() int {
			for _, r := range product.SalesRankings {
				if !re.MatchString(r.ProductCategoryId) && r.ProductCategoryId != "" {
					return r.Rank
				}
			}
			return -1
		}()
		renderDatas = append(renderDatas, &renderData{
			asin:  asin,
			price: price,
			rank:  rank,
		})
	}
	return renderDatas
}

func UnixTimeToKeepaTime(unix int64) int64 {
	return (unix/60 - 21564000)
}

func KeepaTimeToUnix(keepaTime int64) int64 {
	return (keepaTime + 21564000) * 60
}

type KeepaRepository struct {
	DB *bun.DB
}

func (k KeepaRepository) Save(ctx context.Context, keepas []*Keepa) error {
	if len(keepas) == 0 {
		return errors.New("expect at least on keepa object")
	}
	_, err := k.DB.
		NewInsert().
		Model(&keepas).
		On("CONFLICT (asin) DO UPDATE").
		Set(strings.Join([]string{
			"sales_drops_90 = EXCLUDED.sales_drops_90",
			"drops_ma_7 = EXCLUDED.drops_ma_7",
			"price_data = EXCLUDED.price_data",
			"rank_data =  EXCLUDED.rank_data",
			"render_data = EXCLUDED.render_data",
			"modified = now()",
		}, ",")).
		Exec(ctx)
	return err
}

func (k KeepaRepository) UpdateDropsMA(ctx context.Context, keepas []*Keepa) error {
	if len(keepas) == 0 {
		return fmt.Errorf("expect at least on keepa ojbects")
	}
	_, err := k.DB.NewInsert().
		Model(&keepas).
		On("CONFLICT (asin) DO UPDATE").
		Set(strings.Join([]string{
			"drops_ma_7 = EXCLUDED.drops_ma_7",
		}, ",")).
		Exec(ctx)
	return err
}

func (k KeepaRepository) UpdateDropsMAV2(ctx context.Context, keepas []*Keepa) error {
	_, err := k.DB.NewUpdate().
		Model(&keepas).
		Column("drops_ma_7").
		Bulk().
		Exec(ctx)
	return err
}

func (k KeepaRepository) GetByAsins(ctx context.Context, asins []string) (Keepas, error) {
	keepas := make([]*Keepa, 0, len(asins))
	err := k.DB.NewSelect().
		Model(&keepas).
		Where("asin IN (?)", bun.In(asins)).
		Order("asin").
		Scan(ctx)
	return keepas, err
}

func (k KeepaRepository) GetCounts(ctx context.Context) (map[string]int, error) {
	now := time.Now().Format("2006-01-02")

	var total, modified int
	err := k.DB.NewSelect().
		Model((*Keepa)(nil)).
		ColumnExpr("count(*)").
		ColumnExpr("count(? = ? or NULL)", bun.Ident("modified"), now).
		Scan(ctx, &total, &modified)

	return map[string]int{"total": total, "modified": modified}, err
}

type Cursor struct {
	Start string
	End   string
}

func NewCursor(keepas Keepas) Cursor {
	if len(keepas) == 0 {
		return Cursor{}
	}
	return Cursor{
		Start: keepas[0].Asin,
		End:   keepas[len(keepas)-1].Asin,
	}
}

type LoadingData struct {
	Price bool
	Rank  bool
}

type SelectQueryOption = func(*bun.SelectQuery) *bun.SelectQuery

func WithOutPriceData() SelectQueryOption {
	return func(qb *bun.SelectQuery) *bun.SelectQuery {
		qb.ExcludeColumn("price_data")
		return qb
	}
}
func WithOutRankData() SelectQueryOption {
	return func(qb *bun.SelectQuery) *bun.SelectQuery {
		qb.ExcludeColumn("rank_data")
		return qb
	}
}
func OnlyDropsMA7IsNull() SelectQueryOption {
	return func(qb *bun.SelectQuery) *bun.SelectQuery {
		qb.Where("drops_ma_7 IS NULL")
		return qb
	}
}

func (k KeepaRepository) GetKeepaWithPaginate(
	ctx context.Context, cursor string, limit int,
	opts ...SelectQueryOption) (Keepas, Cursor, error) {
	var keepas Keepas
	qb := k.DB.NewSelect().
		Model(&keepas).
		Where("asin > ?", cursor).
		Order("asin ASC").
		Limit(limit)

	for _, opt := range opts {
		qb = opt(qb)
	}

	if err := qb.Scan(ctx); err != nil {
		return nil, Cursor{}, err
	}
	return keepas, NewCursor(keepas), nil
}
