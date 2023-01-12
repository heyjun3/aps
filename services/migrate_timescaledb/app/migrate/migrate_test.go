package migrate

import (
	"fmt"
	"testing"
	"time"
	"strconv"

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

		result, err := ConvKeepaProductToAsinsInfo(&p)
		
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

		result, err := ConvKeepaProductToAsinsInfo(&p)

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

		result, err := ConvKeepaProductToAsinsInfo(&p)

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

		result, err := ConvKeepaProductToAsinsInfo(&p)

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

		result, err := ConvKeepaProductToAsinsInfo(&p)

		assert.Error(t, err)
		assert.Equal(t, []models.AsinsInfoTime(nil), result)
	})
}
