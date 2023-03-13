package pc4u

import (
	"bytes"
	"io"
	"os"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createHttpResponse(path string) (*http.Response, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &http.Response{
		Body: io.NopCloser(bytes.NewReader(b)),
		Request: &http.Request{},
	}
	return res, nil
}

func TestParseProducts(t *testing.T) {
	parser := Pc4uParser{}

	t.Run("parse product list page", func(t *testing.T) {
		res, err := createHttpResponse("html/test_product_list.html")
		if err != nil {
			logger.Error("error", err)
			return
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
		res, err := createHttpResponse("html/test_product_list_last_page.html")	
		if err != nil {
			logger.Error("error", err)
			return
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
	t.Run("parse next page url", func(t *testing.T) {
		res, err := createHttpResponse("html/test_product_list_next_URL.html")
		if err != nil {
			logger.Error("error", err)
			return
		}
		defer res.Body.Close()

		_, url := parser.ProductList(res.Body)

		assert.Equal(t, "https://www.pc4u.co.jp/view/search?page=6", url)
	})
	t.Run("parse sold page url", func(t *testing.T) {
		res, err := createHttpResponse("html/test_product_list_soldout_next_URL.html")
		if err != nil {
			logger.Error("error", err)
			return
		}
		defer res.Body.Close()

		_, url := parser.ProductList(res.Body)

		assert.Equal(t, "", url)
	})
}

func TestProduct(t *testing.T) {
	parser := Pc4uParser{}
	t.Run("parse product", func(t *testing.T) {
		res, err := createHttpResponse("html/test_product.html")
		if err != nil {
			logger.Error("error", err)
			return
		}
		defer res.Body.Close()

		jan, err := parser.Product(res.Body)

		assert.Equal(t, nil, err)
		assert.Equal(t, "4719331990053", jan)
	})

	t.Run("parser product no contain table", func(t *testing.T) {
		res, err := createHttpResponse("html/test_product_no_table.html")
		if err != nil {
			logger.Error("error", err)
			return
		}
		defer res.Body.Close()

		jan, err := parser.Product(res.Body)

		assert.Equal(t, nil, err)
		assert.Equal(t, "4537694092371", jan)
	})

	t.Run("parser product on table", func(t *testing.T) {
		res, err := createHttpResponse("html/test_product_on_table.html")
		if err != nil {
			logger.Error("error", err)
			return
		}
		defer res.Body.Close()

		jan, err := parser.Product(res.Body)

		assert.Equal(t, nil, err)
		assert.Equal(t, "4719512135716", jan)
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
