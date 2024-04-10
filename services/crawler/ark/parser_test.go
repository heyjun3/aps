package ark

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/product"
	"crawler/test/util"
)

func newTestArkProduct(name, productCode, url, jan string,
	price int64) *product.Product {
	p, err := newArkProduct(name, productCode, url, jan, price)
	if err != nil {
		panic(err)
	}
	return p
}

func TestProductList(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		count int
		url   string
		first *product.Product
		last  *product.Product
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "parse product list page",
			args: args{filename: "html/test_product_list.html"},
			want: want{
				count: 50,
				url:   "https://www.ark-pc.co.jp/search/?offset=50&limit=50&nouki=1",
				first: newTestArkProduct(
					"CMSX16GX5M1A4800C40",
					"11755303",
					"https://www.ark-pc.co.jp/i/11755303/",
					"",
					8980,
				),
				last: newTestArkProduct(
					"Loupedeck Live",
					"50284987",
					"https://www.ark-pc.co.jp/i/50284987/",
					"",
					39600,
				),
			},
		},
		{
			name: "parse last page",
			args: args{filename: "html/test_product_list_last.html"},
			want: want{
				count: 4,
				url:   "",
				first: newTestArkProduct(
					"AINEX YH-3020A チップ用ヒートシンク30mm角",
					"40000501",
					"https://www.ark-pc.co.jp/i/40000501/",
					"",
					440,
				),
				last: newTestArkProduct(
					"SteelSeries QcK+ (Qck L)",
					"50190022",
					"https://www.ark-pc.co.jp/i/50190022/",
					"",
					1840,
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			products, url := ArkParser{}.ProductList(res.Body, "")

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.url, url)
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
		})
	}
}

func TestProductListInCouponPrice(t *testing.T) {
	t.Run("parse coupon price", func(t *testing.T) {
		res, err := util.CreateHttpResponse("html/test_coupon_price.html")
		if err != nil {
			logger.Error("file open error", err)
			panic(err)
		}
		defer res.Body.Close()

		products, url := ArkParser{}.ProductList(res.Body, "")

		assert.Equal(t, 50, len(products))
		assert.Equal(t, "https://www.ark-pc.co.jp/search/?offset=2900&limit=50&nouki=1", url)

		for _, v := range products {
			if v.GetProductCode() == "50283314" {
				assert.Equal(t, int64(11980), v.(*product.Product).Price)
			}
		}
	})
}

func TestProduct(t *testing.T) {
	parser := ArkParser{}

	t.Run("parse product page", func(t *testing.T) {
		res, err := util.CreateHttpResponse("html/test_product_page.html")
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
