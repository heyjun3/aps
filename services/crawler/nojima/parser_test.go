package nojima

import (
	"testing"
	"net/url"

	"crawler/testutil"

	"github.com/stretchr/testify/assert"
)

func TestProductList(t *testing.T) {
	type args struct {
		filename string
		URL string
	}
	type want struct {
		count int
		url string
		first *NojimaProduct
		last *NojimaProduct
	}
	u, _ := url.Parse("https://online.nojima.co.jp/app/catalog/list/init?searchCategoryCode=0&mode=image&pageSize=60&currentPage=4&alignmentSequence=9&searchDispFlg=true&immediateDeliveryDispFlg=1&searchWord=%E3%82%A4%E3%83%B3%E3%82%AF")
	u.RawQuery = u.Query().Encode()
	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "parse product list page",
		args: args{
			filename: "html/test_product_list.html",
			URL: "https://online.nojima.co.jp/app/catalog/list/init?searchCategoryCode=0&mode=image&pageSize=60&currentPage=3&alignmentSequence=9&searchDispFlg=true&immediateDeliveryDispFlg=1&searchWord=%E3%82%A4%E3%83%B3%E3%82%AF",
		},
		want: want{
			count: 60,
			url: u.String(),
			first: NewNojimaProduct(
				"Canon純正インクイエロー",
				"4960999905549",
				"https://online.nojima.co.jp/commodity/1/4960999905549/",
				"4960999905549",
				873,
			),
			last: NewNojimaProduct(
				"brother純正インクカートリッジ大容量タイプシアン",
				"4977766788977",
				"https://online.nojima.co.jp/commodity/1/4977766788977/",
				"4977766788977",
				2601,
			),
		},
	},{
		name: "parse last page",
		args: args{
			filename: "html/test_product_list_last.html",
		},
		want: want{
			count: 26,
			url: "",
			first: NewNojimaProduct(
				"BANDAIひろがるスカイ!プリキュアふしぎなミラージュペン",
				"4549660880431",
				"https://online.nojima.co.jp/commodity/1/4549660880431/",
				"4549660880431",
				1380,
			),
			last: NewNojimaProduct(
				"ELSONICDVD-R【1回録画用/4.7GB/1-16倍速/10枚/5mmプラケース】",
				"0479960012833",
				"https://online.nojima.co.jp/commodity/1/0479960012833/",
				"0479960012833",
				1078,
			),
		},
	}}
	parser := NojimaParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		res, err := testutil.CreateHttpResponse(tt.args.filename)
		if err != nil {
			logger.Error("file open error", err)
			panic(err)
		}
		defer res.Body.Close()

		products, url := parser.ProductList(res.Body, tt.args.URL)

		assert.Equal(t, tt.want.count, len(products))
		assert.Equal(t, tt.want.url, url)
		assert.Equal(t, tt.want.first, products[0])
		assert.Equal(t, tt.want.last, products[len(products)-1])
		})
	}
}