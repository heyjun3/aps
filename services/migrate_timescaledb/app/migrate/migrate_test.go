package migrate

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"migrate_timescaledb/app/models"
)


func TestConvKeepaTimeToTime(t *testing.T) {
	t.Run("convert keepa time to time.Time", func(t *testing.T) {
		dt, err := convKeepaTimeToTime("5831784")

		assert.Equal(t, nil, err)
		assert.Equal(t, time.Date(2022, 2, 2, 5, 24, 0, 0, time.Local), *dt)
		tmp, _ := convKeepaTimeToTime("22222")
		fmt.Println(*tmp)
	})

	t.Run("a-z argument is not valid", func(t *testing.T) {
		dt, err := convKeepaTimeToTime("aazzz")
		
		assert.IsType(t, &strconv.NumError{}, err)
		assert.Equal(t, time.Time{}, *dt)
	})
}

func TestGetMapKeys(t *testing.T) {
	t.Run("get keys for map", func(t *testing.T) {
		m := map[string]float64{
			"1111": 5000,
			"0999": 1000,
			"5555": 5555,
		}

		keys, err := getMapKeys(m)

		assert.Equal(t, nil, err)
		assert.Equal(t, []int{999, 1111, 5555}, keys)
	})

	t.Run("keys specified 0 to 9 only", func(t *testing.T) {
		m := map[string]float64{
			"1111": 5000,
			"0999": 1000,
			"aaa": 5555,
		}

		keys, err := getMapKeys(m)

		assert.Error(t, err)
		assert.Equal(t, []int(nil), keys)
	})
}

func TestDeleteDuplicateAsinsInfoTimes(t *testing.T) {
	t.Run("delete duplicate AsinsInfoTimes", func(t *testing.T) {
		p := []models.AsinsInfoTime{
			{
				Time: time.Date(2023, 1, 16, 0, 35, 0, 0, time.Local),
				Asin: "XXXX",
				Price: null.IntFrom(1000),
				Rank: null.IntFrom(1000),
			},
			{
				Time: time.Date(2023, 1, 16, 0, 35, 0, 0, time.Local),
				Asin: "XXXX",
				Price: null.IntFrom(2000),
				Rank: null.IntFrom(2000),
			},
			{
				Time: time.Date(2023, 1, 17, 0, 35, 0, 0, time.Local),
				Asin: "XXXX",
				Rank: null.IntFrom(4000),
			},
			{
				Time: time.Date(2023, 1, 17, 0, 35, 0, 0, time.Local),
				Asin: "XXXX",
				Price: null.IntFrom(3000),
			},
		}

		d := deleteDuplicateAsinsInfoTimes(p)

		ext := []models.AsinsInfoTime{
			{
				Time: time.Date(2023, 1, 16, 0, 35, 0, 0, time.Local),
				Asin: "XXXX",
				Price: null.IntFrom(2000),
				Rank: null.IntFrom(2000),
			},
			{
				Time: time.Date(2023, 1, 17, 0, 35, 0, 0, time.Local),
				Asin: "XXXX",
				Price: null.IntFrom(3000),
				Rank: null.IntFrom(4000),
			},
		}
		assert.Equal(t, ext, d)
	})
}

func TestConvKeepaProductToAsinsInfo(t *testing.T) {
	t.Run("convert keepa product to asins infos", func (t *testing.T) {
		p := models.KeepaProduct{
			Asin: "AAAABBBB",
		    SalesDrops90: null.NewInt(1, true),
			Created: null.NewTime(time.Date(2023, 1, 12, 1, 40, 0, 0, time.Local), true),	
			Modified: null.NewTime(time.Date(2024, 1, 12, 1, 40, 0, 0, time.Local), true),	
			PriceData: null.JSONFrom([]byte(`{"10000": 10000, "20000": 20000}`)),
			RankData: null.JSONFrom([]byte(`{"11111": 11111, "22222": 22222}`)),
		}

		result, err := convKeepaProductToAsinsInfo(&p)
		
		assert.Equal(t, nil, err)
		assert.IsType(t, []models.AsinsInfoTime{}, result)
		assert.Equal(t, 4, len(result))
		k1 := models.AsinsInfoTime{
			Time: time.Date(2011, 1, 8, 7, 40, 0, 0, time.Local),
			Asin: "AAAABBBB",
			Price: null.NewInt(10000, true),
		}
		k2 := models.AsinsInfoTime{
			Time: time.Date(2011, 1, 16, 19, 22, 0, 0, time.Local),
			Asin: "AAAABBBB",
			Rank: null.NewInt(22222, true),
		}
		assert.Equal(t, k1, result[0])
		assert.Equal(t, k2, result[len(result)-1])
	})

	t.Run("price is required", func (t *testing.T) {
		p := models.KeepaProduct{
			Asin: "AAAABBBB",
		    SalesDrops90: null.NewInt(1, true),
			Created: null.NewTime(time.Date(2023, 1, 12, 1, 40, 0, 0, time.Local), true),	
			Modified: null.NewTime(time.Date(2024, 1, 12, 1, 40, 0, 0, time.Local), true),	
			RankData: null.JSONFrom([]byte(`{"11111": 11111, "22222": 22222}`)),
		}

		result, err := convKeepaProductToAsinsInfo(&p)

		assert.Error(t, err)
		assert.Equal(t, []models.AsinsInfoTime(nil), result)
	})

	t.Run("rank is required", func (t *testing.T) {
		p := models.KeepaProduct{
			Asin: "AAAABBBB",
		    SalesDrops90: null.NewInt(1, true),
			Created: null.NewTime(time.Date(2023, 1, 12, 1, 40, 0, 0, time.Local), true),	
			Modified: null.NewTime(time.Date(2024, 1, 12, 1, 40, 0, 0, time.Local), true),	
			PriceData: null.JSONFrom([]byte(`{"11111": 11111, "22222": 22222}`)),
		}

		result, err := convKeepaProductToAsinsInfo(&p)

		assert.Error(t, err)
		assert.Equal(t, []models.AsinsInfoTime(nil), result)
	})

	t.Run("price keepa time is 0 to 9 strings", func(t *testing.T) {
		p := models.KeepaProduct{
			Asin: "AAAABBBB",
		    SalesDrops90: null.NewInt(1, true),
			Created: null.NewTime(time.Date(2023, 1, 12, 1, 40, 0, 0, time.Local), true),	
			Modified: null.NewTime(time.Date(2024, 1, 12, 1, 40, 0, 0, time.Local), true),	
			PriceData: null.JSONFrom([]byte(`{"aaaa": 10000, "20000": 20000}`)),
			RankData: null.JSONFrom([]byte(`{"11111": 11111, "22222": 22222}`)),
		}

		result, err := convKeepaProductToAsinsInfo(&p)

		assert.Error(t, err)
		assert.Equal(t, []models.AsinsInfoTime(nil), result)
	})

	t.Run("rank keepa time is 0 to 9 strings", func(t *testing.T) {
		p := models.KeepaProduct{
			Asin: "AAAABBBB",
		    SalesDrops90: null.NewInt(1, true),
			Created: null.NewTime(time.Date(2023, 1, 12, 1, 40, 0, 0, time.Local), true),	
			Modified: null.NewTime(time.Date(2024, 1, 12, 1, 40, 0, 0, time.Local), true),	
			PriceData: null.JSONFrom([]byte(`{"10000": 10000, "20000": 20000}`)),
			RankData: null.JSONFrom([]byte(`{"11111": 11111, "LLLLLL": 22222}`)),
		}

		result, err := convKeepaProductToAsinsInfo(&p)

		assert.Error(t, err)
		assert.Equal(t, []models.AsinsInfoTime(nil), result)
	})
}


