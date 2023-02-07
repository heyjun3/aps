package ikebe

import (
	"fmt"

	"github.com/BurntSushi/toml"
)



type Config struct {
	Psql Psql
	RabbitMQ RabbitMQ
}

type Psql struct {
	DBname string `toml:"dbname"`
	Host string `toml:"host"`
	Port string `toml:"port"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
	SSLmode string `toml:"sslmode"`
	Blacklist []string `toml:"blacklist"`
}

type RabbitMQ struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
	Host string `toml:"host"`
	Port string `toml:"port"`
}

var cfg Config
func init() {
	cfg = NewConfig("sqlboiler.toml")
}

func NewConfig(path string) Config {
	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		fmt.Println(err)
		return cfg
	}
	return cfg
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