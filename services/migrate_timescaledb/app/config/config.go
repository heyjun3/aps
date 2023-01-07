package config

import (
	"os"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB DBConfig `yaml:"psql"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string 
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name   string `yaml:"dbname"`
}

var Cfg Config
func init() {
	var err error
	Cfg, err = NewConfig("sqlboiler.yaml")
	if err != nil {
		fmt.Printf("Initialize config error")
	}
}

func NewConfig(path string) (Config, error) {
	var cfg Config
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Doesn't read config file %v", err)
		return cfg, err
	}
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		fmt.Printf("yaml unmarshal error %v", err)
		return cfg, err
	}
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	return cfg, nil
}

func (c *Config) Dsn() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}
