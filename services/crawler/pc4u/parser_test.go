package pc4u

import (
	"bytes"
	"io"
	"os"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProducts(t *testing.T) {
	parser := Pc4uParser{}

	t.Run("parse product list page", func(t *testing.T) {
		b, err := os.ReadFile("html/test_product_list.html")
		if err != nil {
			logger.Error("file open error", err)
			return
		}
		res := &http.Response{
			Body: io.NopCloser(bytes.NewReader(b)),
			Request: &http.Request{},
		}
		defer res.Body.Close()

		products, url := parser.ProductList(res.Body)

		assert.Equal(t, 50, len(products))
		assert.Equal(t, "https://www.pc4u.co.jp/view/search?page=2", url)

		first := NewPc4uProduct(
			"GIGABYTE B760I AORUS PRO DDR4 第13 & 12世代 Intel Core プロセッサー対応 Mini-ITX マザーボード｜B760I AORUS PRO DDR4",
			"000000081834",
			"https://www.pc4u.co.jp/view/item/000000081834",
			"",
			36960,
		)
		last := NewPc4uProduct(
			"ADATA XPG GAMMIX D20 16GB(16GBx1) DDR4 3600MHz(PC4-28800) U-DIMM SINGLE COLOR BOX ブラック｜AX4U360016G18I-CBK20",
			"000000081728",
			"https://www.pc4u.co.jp/view/item/000000081728",
			"",
			7429,
		)
		assert.Equal(t, first, products[0])
		assert.Equal(t, last, products[len(products)-1])
	})

	t.Run("parse last page", func(t *testing.T) {
		b, err := os.ReadFile("html/test_product_list_last_page.html")	
		if err != nil {
			logger.Error("file open error", err)
			return
		}
		res := &http.Response{
			Body: io.NopCloser(bytes.NewReader(b)),
			Request: &http.Request{},
		}
		defer res.Body.Close()

		products, url := parser.ProductList(res.Body)
		
		assert.Equal(t, 17, len(products))
		assert.Equal(t, "", url)

		first := NewPc4uProduct(
			"【アウトレット特価・新品】Corning Thunderbolt 3, 50m Optical Cable Thunderboltケーブル｜AOC-CCU6JPN050M20",
			"000000072144",
			"https://www.pc4u.co.jp/view/item/000000072144?category_page_id=outlet",
			"",
			55990,
		)
		last := NewPc4uProduct(
			"【アウトレット特価・新品】Keyspan USBシリアルアダプタ  USB Serial Adapter (USA-19HS)",
			"014004000005",
			"https://www.pc4u.co.jp/view/item/014004000005?category_page_id=outlet",
			"",
			5940,
		)
		assert.Equal(t, first, products[0])
		assert.Equal(t, last, products[len(products)-1])
	})
}

func TestPullOutPrice(t *testing.T) {
	t.Run("pull out price", func(t *testing.T) {
		s := " 199,800円"

		price, err := pullOutPrice(s)

		assert.Equal(t, nil, err)
		assert.Equal(t, int64(199800), price)
	})

	t.Run("pull out price not digits", func(t *testing.T) {
		s := "aaa  fdsagfda"

		price, err := pullOutPrice(s)

		assert.Error(t, err)
		assert.Equal(t, int64(0), price)
	})

	t.Run("blank string", func(t *testing.T) {
		s := ""

		price, err := pullOutPrice(s)

		assert.Error(t, err)
		assert.Equal(t, int64(0), price)
	})
}
