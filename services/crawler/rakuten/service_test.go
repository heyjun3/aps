package rakuten

import (
	"testing"
)

func TestGetShopList(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		s, _ := getShopList()
		for _, v := range s.List {
			logger.Info("show shop code", "code", v.ID, "url", v.URL)
		}
	})
}
