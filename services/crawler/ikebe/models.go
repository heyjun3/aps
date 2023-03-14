package ikebe

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"crawler/config"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID int64 `bun:",pk,autoincrement"`
	Name string
}

func CreateTable() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.DBDsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		logger.Info("erro")
	}
}