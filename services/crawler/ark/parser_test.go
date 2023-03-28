package ark


import (
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/scrape"
)

func TestProductList(t *testing.T) {
	parser := ArkParser{}

	t.Run("parse product list page", func(t *testing.T) {
		res, err := scrape.CreateHttpResponse("html/test_product_list.html")
		if err != nil {
			logger.Error("error", err)
			panic(err)
		}
		defer res.Body.Close()

		products, url := parser.ProductList(res.Body)

		assert.Equal(t, 50, len(products))
		assert.Equal(t, "https://www.ark-pc.co.jp/search/?offset=50&limit=50&nouki=1", url)
		first := NewArkProduct(
			"CMSX16GX5M1A4800C40",
			"11755303",
			"https://www.ark-pc.co.jp/i/11755303/",
			"",
			8980,
		)
		last := NewArkProduct(
			"Loupedeck Live",
			"50284987",
			"https://www.ark-pc.co.jp/i/50284987/",
			"",
			39600,
		)
		
		assert.Equal(t, first, products[0])
		assert.Equal(t, last, products[len(products)-1])
	})

	t.Run("parse last page", func(t *testing.T) {
		res, err := scrape.CreateHttpResponse("html/test_product_list_last.html")
		if err != nil {
			logger.Error("file open error", err)
			panic(err)
		}
		defer res.Body.Close()

		products, url := parser.ProductList(res.Body)

		assert.Equal(t, 4, len(products))
		assert.Equal(t, "", url)

		first := NewArkProduct(
			"AINEX YH-3020A チップ用ヒートシンク30mm角",
			"40000501",
			"https://www.ark-pc.co.jp/i/40000501/",
			"",
			440,
		)
		last := NewArkProduct(
			"SteelSeries QcK+ (Qck L)",
			"50190022",
			"https://www.ark-pc.co.jp/i/50190022/",
			"",
			1840,
		)

		assert.Equal(t, first, products[0])
		assert.Equal(t, last, products[len(products)-1])
	})

	t.Run("parse coupon price", func(t *testing.T) {
		res, err := scrape.CreateHttpResponse("html/test_coupon_price.html")
		if err != nil {
			logger.Error("file open error", err)
			panic(err)
		}
		defer res.Body.Close()

		products, url := parser.ProductList(res.Body)

		assert.Equal(t, 50, len(products))
		assert.Equal(t, "https://www.ark-pc.co.jp/search/?offset=2900&limit=50&nouki=1", url)

		for _, v := range products {
			if v.GetProductCode() == "50283314" {
				assert.Equal(t, int64(11980), v.(*ArkProduct).Price)
			}
		}
	})
}

func TestProduct(t *testing.T) {
	parser := ArkParser{}
	
	t.Run("parse product page", func(t *testing.T) {
		res, err := scrape.CreateHttpResponse("html/test_product_page.html")
		if err != nil {
			logger.Error("file open error", err)
			panic(err)
		}
		defer res.Body.Close()

		jan, err := parser.Product(res.Body)

		assert.Equal(t, nil, err)
		assert.Equal(t, "0843591081313", jan)
	})
}
