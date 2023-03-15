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

func Tmp() {
	conn := createDBConnection(config.DBDsn)
	ctx := context.Background()
	u := &User{
		Name: "test",
	}
	err := u.Upsert(conn, ctx)
	if err != nil {
		logger.Error("error", err)
	}
}

func createDBConnection(dsn string) *bun.DB{
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func CreateTable() {
	db := createDBConnection(config.DBDsn)
	ctx := context.Background()
	_, err := db.NewCreateTable().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		logger.Info("erro")
	}
}

func (u *User) Insert(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewInsert().Model(u).Exec(ctx)
	return err
}

func (u *User) Upsert(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewInsert().
		Model(u).
		On("CONFLICT (id) DO UPDATE").
		Set("? = EXCLUDED.?", bun.Ident("name"), "name").
		Exec(ctx)
	return err
}