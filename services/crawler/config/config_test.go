package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigInit(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		u := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
		assert.Equal(t, u, Http.UserAgent)
	})
}
