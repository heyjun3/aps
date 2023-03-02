package ikebe

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"golang.org/x/exp/slog"
)

type Config struct {
	Psql     Psql
	RabbitMQ RabbitMQ
}

type Psql struct {
	DBname    string   `toml:"dbname"`
	Host      string   `toml:"host"`
	Port      string   `toml:"port"`
	User      string   `toml:"user"`
	Pass      string   `toml:"pass"`
	SSLmode   string   `toml:"sslmode"`
	Blacklist []string `toml:"blacklist"`
}

type RabbitMQ struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
	Host string `toml:"host"`
	Port string `toml:"port"`
}

var cfg Config
var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout))
	path := os.Getenv("ROOT_PATH")
	if path == "" {
		panic("Not set env ROOT_PATH")
	}
	var err error
	cfg, err = NewConfig(filepath.Join(path, "sqlboiler.toml"))
	if err != nil {
		panic(err)
	}
}

func NewConfig(path string) (Config, error) {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (c Config) dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Psql.User,
		c.Psql.Pass,
		c.Psql.Host,
		c.Psql.Port,
		c.Psql.DBname,
		c.Psql.SSLmode,
	)
}

func (c Config) MQDsn() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Pass,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)
}
