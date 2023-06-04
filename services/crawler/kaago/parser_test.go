package kaago

import (
	"testing"

	"crawler/test/util"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		count int
		url   string
		first *KaagoProduct
		last  *KaagoProduct
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "parse product list",
		args: args{
			filename: "json/test_product_list.json",
		},
		want: want{
			count: 36,
			url:   "",
			first: NewKaagoProduct(
				"S223ATES-W ダイキン ルームエアコン6畳 ホワイト",
				"9900000004842",
				"/seicaplus/S223ATES-W-%E3%83%80%E3%82%A4%E3%82%AD%E3%83%B3-%E3%83%AB%E3%83%BC%E3%83%A0%E3%82%A8%E3%82%A2%E3%82%B3%E3%83%B36%E7%95%B3-%E3%83%9B%E3%83%AF%E3%82%A4%E3%83%88/?itemcode=9900000004842",
				"9900000004842",
				"seicaplus",
				46990,
			),
			last: NewKaagoProduct(
				"S71ZTCXP-W ダイキン ルームエアコン23畳 ホワイト 200V",
				"9900000004160",
				"/seica/S71ZTCXP-W-%E3%83%80%E3%82%A4%E3%82%AD%E3%83%B3-%E3%83%AB%E3%83%BC%E3%83%A0%E3%82%A8%E3%82%A2%E3%82%B3%E3%83%B323%E7%95%B3-%E3%83%9B%E3%83%AF%E3%82%A4%E3%83%88-200V/?itemcode=9900000004160",
				"9900000004160",
				"seica",
				146000,
			),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			products, url := KaagoParser{}.ProductList(res.Body, "")

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.url, url)
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
		})
	}

}
