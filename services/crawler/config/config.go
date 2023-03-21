package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"golang.org/x/exp/slog"
)

type config struct {
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

var Config config
var DBDsn string
var MQDsn string
var DstMQDsn string
var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout))
	path := os.Getenv("ROOT_PATH")
	if path == "" {
		panic("Not set env ROOT_PATH")
	}
	var err error
	Config, err = NewConfig(filepath.Join(path, "sqlboiler.toml"))
	if err != nil {
		panic(err)
	}
	DBDsn = Config.Dsn()
	MQDsn = Config.MQDsn()

	Config.RabbitMQ.Host = "192.168.0.5"
	DstMQDsn = Config.MQDsn()
}

func NewConfig(path string) (config, error) {
	var Config config
	_, err := toml.DecodeFile(path, &Config)
	if err != nil {
		return Config, err
	}
	return Config, nil
}

func (c config) Dsn() string {
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

func (c config) MQDsn() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Pass,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)
}
