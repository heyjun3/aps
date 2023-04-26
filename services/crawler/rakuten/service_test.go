package rakuten

import (
	"testing"
)

func TestGetShopList(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		s, _ := getShopList()
		logger.Info("shoplist", "value", s)
	})
}
