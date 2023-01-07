package connection

import(
	"fmt"
	"database/sql"
	"migrate_timescaledb/app/config"

	_ "github.com/lib/pq"
)

var DbConnection *sql.DB

func init() {
	var err error
	DbConnection, err = CreateDBConnection(config.Cfg.Dsn())
	if err != nil {
		fmt.Printf("sql doesn't open database %v", err)
		return
	}
}

func CreateDBConnection(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("sql doesn't open database %v", err)
		return nil, err
	}
	return conn, nil
}