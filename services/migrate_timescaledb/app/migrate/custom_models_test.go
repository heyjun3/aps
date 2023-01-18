package migrate

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"migrate_timescaledb/app/config"
	"migrate_timescaledb/app/connection"
	"migrate_timescaledb/app/models"
)
func TestUpsertAsinsInfoTimes(t *testing.T) {
	c, _ := config.NewConfig("../../sqlboiler.yaml")
	db, err := connection.CreateDBConnection(c.Dsn())
	if err != nil {
		fmt.Printf("Doesn't connection database")
		return
	}

	t.Run("upsert asins_info_time records", func(t *testing.T) {
		tx, _ := db.Begin()
		defer tx.Rollback()
		ctx := context.Background()
		p := []models.AsinsInfoTime{
			{
				Time: time.Date(2023, 1, 15, 1, 4, 0, 0, time.Local),
				Asin: "XXXXXXX",
				Price: null.IntFrom(1000),
			},
			{
				Time: time.Date(2023, 1, 15, 1, 4, 0, 0, time.Local),
				Asin: "XXXXXXX",
				Rank: null.IntFrom(9999),
			},
		}

		err := upsertAsinsInfoTimes(ctx, tx, p)

		assert.Equal(t, nil, err)
	})
}

func TestBalkUpsertAsinsInfoTimes(t *testing.T) {
	c, _ := config.NewConfig("../../sqlboiler.yaml")
	db, _ := connection.CreateDBConnection(c.Dsn())

	t.Run("balk upsert asinsInfoTimes", func(t *testing.T) {
		m := []models.AsinsInfoTime{
			{
				Time: time.Now(),
				Asin: "XXXX",
				Price: null.IntFrom(1111),
				Rank: null.IntFrom(1111),
			},
			{
				Time: time.Now(),
				Asin: "YYYY",
				Price: null.IntFrom(2222),
				Rank: null.IntFrom(2222),
			},
		}
		ctx := context.Background()
		tx, _ := db.Begin()
		defer tx.Rollback()

		err := balkUpsertAsinsInfoTimes(ctx, tx, m)

		assert.Equal(t, nil, err)
	})
}

func TestGetAllAsinsFromKeepaProduct(t *testing.T) {
	ct := context.Background()
	c, _ := config.NewConfig("../../sqlboiler.yaml")
	c.DB.Name = "test"
	db, err := connection.CreateDBConnection(c.Dsn())
	if err != nil {
		fmt.Printf("Doesn't connection database")
		return
	}
	m := models.KeepaProduct{
		Asin: "test",
	}
	m.Upsert(ct, db, false, []string{"asin"}, boil.Infer(), boil.Infer())

	defer models.KeepaProducts().DeleteAll(ct, db)

	t.Run("get all asins", func(t *testing.T) {
		tx, _ := db.Begin()
		defer tx.Rollback()
		ctx := context.Background()

		asins, err := getAllAsinsFromKeepaProduct(ctx, db)

		assert.Equal(t, nil, err)
		ext := models.KeepaProductSlice{&models.KeepaProduct{Asin: "test"}}
		assert.Equal(t, ext, asins)
	})
}