package config

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewConfig(t *testing.T) {
	pass := os.Getenv("DB_PASSWORD")
	config := Config{
		DB: DBConfig{
			User: "postgres",
			Password: pass,
			Host: "timescaledb",
			Port: "5432",
			Name: "aps",
		},
	}

	t.Run("test new config", func(t *testing.T) {
		cfg, err := NewConfig("../../sqlboiler.yaml")

		assert.Equal(t, nil, err)
		assert.Equal(t, config, cfg)
	})
}