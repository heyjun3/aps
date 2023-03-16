package ikebe

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"crawler/config"
)

type IkebeItem struct {
	bun.BaseModel `bun:"table:ikebe_product"`
	Name string
	Jan string
	Price int64
	ShopCode string `bun:"shop_code,pk"`
	ProductCode string `bun:"product_code,pk"`
	URL string
}

func (i *IkebeItem) Upsert(conn *bun.DB, ctx context.Context) error {
	_, err := conn.NewInsert().
		Model(i).
		On("CONFLICT (shop_code, product_code) DO UPDATE").
		Set(`
		name = EXCLUDED.name,
		jan = EXCLUDED.jan,
		price = EXCLUDED.price,
		url = EXCLUDED.url
		`).
		Exec(ctx)
	return err
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID int64 `bun:",pk,autoincrement"`
	Name string
}

func Tmp() {
	conn := createDBConnection(config.DBDsn)
	ctx := context.Background()
	i := &IkebeItem{
		ShopCode: "test",
		ProductCode: "test",
		Name: "name2",
	}
	err := i.Upsert(conn, ctx)
	if err != nil {
		logger.Error("error", err)
	}
	Test(test1)
	Test(test2("a", "aa"))
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
		Set("? = ?", bun.Ident("name"), "name").
		Exec(ctx)
	return err
}

func test1(s string) {
	logger.Info(s)
}

func test2(s1, s2 string) func(string){
	return func(s string) {
		logger.Info(s1)
		logger.Info(s2)
		logger.Info(s)
	}
}

func Test(f func(string)) {
	f("aaaa")
}