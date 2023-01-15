package migrate

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	
	"migrate_timescaledb/app/models"
)


func upsertAsinsInfoTimes(ctx context.Context, db boil.ContextExecutor, p []models.AsinsInfoTime) error {
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

func balkUpsertAsinsInfoTimes(ctx context.Context, db boil.ContextExecutor, p []models.AsinsInfoTime) error {
	limit := 30000
	for i := 0; i < len(p); i += limit {
		end := i + limit
		if len(p) < end {
			end = len(p)
		}
		strs := []string{}
		args := []interface{}{}
		for i, v := range p[i:end] {
			t := i * 4
			strs = append(strs, fmt.Sprintf("($%d, $%d, $%d, $%d)", t+1, t+2, t+3, t+4))

			args = append(args, v.Time.Format(time.RFC3339))
			args = append(args, v.Asin)

			if p, _ := v.Price.Value(); true {
				args = append(args, p)
			}
			if r, _ := v.Rank.Value(); true {
				args = append(args, r)
			}
		}

		stmt := fmt.Sprintf(
			`INSERT INTO asins_info_time(time, asin, price, rank) VALUES %s 
			ON CONFLICT (time, asin) DO UPDATE SET price = excluded.price, 
			rank = excluded.rank`, strings.Join(strs, ", "))
		_, err := db.Exec(stmt, args...)
		if err != nil {
			fmt.Printf("upsert error: %v", err)
			return err
		}
	}
	return nil
}

func getAllAsinsFromKeepaProduct(ctx context.Context, db boil.ContextExecutor) (models.KeepaProductSlice, error) {
	asins, err := models.KeepaProducts(
		qm.Select("asin"),
	).All(ctx, db)
	if err != nil {
		fmt.Printf("Get All Asins error: %v", err)
		return nil, err
	}
	return asins, nil
}
