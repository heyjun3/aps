package config

import (
	"os"

	"golang.org/x/exp/slog"
)

type config struct {
	Psql     Psql
	RabbitMQ RabbitMQ
	Http     HttpConf
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

type HttpConf struct {
	UserAgent string `toml:"useragent"`
}

var Config config
var DBDsn string
var MQDsn string
var DstMQDsn string
var Logger *slog.Logger
var Http HttpConf

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout))
	DBDsn = os.Getenv("DB_DSN")
	if DBDsn == "" {
		panic("DB DSN isn't empty string")
	}
	MQDsn = os.Getenv("MQ_DSN")
	if MQDsn == "" {
		panic("MQ DSN isn't empty string")
	}

	Http = HttpConf{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"}
}
