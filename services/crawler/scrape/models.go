package scrape

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func CreateDBConnection(dsn string) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	conn := bun.NewDB(sqldb, pgdialect.New())
	return conn
}

func CreateTable(db *bun.DB, ctx context.Context, model interface{}) error {
	_, err := db.NewCreateTable().
		Model(model).
		IfNotExists().
		Exec(ctx)
	return err
}

type RunServiceHistory struct {
	bun.BaseModel
	ID        uuid.UUID `bun:",pk,type:uuid,default:gen_random_uuid()"`
	ShopName  string    `bun:",notnull"`
	URL       string    `bun:",notnull"`
	Status    string    `bun:",notnull"`
	StartedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	EndedAt   time.Time `bun:",nullzero"`
}

func NewRunServiceHistory(shopName, url, status string) *RunServiceHistory {
	return &RunServiceHistory{
		ShopName: shopName,
		URL:      url,
		Status:   status,
	}
}
