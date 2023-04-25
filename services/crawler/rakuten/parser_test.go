package rakuten

import (
	"crawler/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductList(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		count int
		url   string
		first *RakutenProduct
		last  *RakutenProduct
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "parse product list page",
		args: args{
			"html/test_product_list.html",
		},
		want: want{
			count: 45,
			url:   "https://search.rakuten.co.jp/search/mall/?p=2&sid=212220",
			first: NewRakutenProduct(
				"【Z16L0005G】 Apple Mac mini 2023年CTOモデル（ベースモデル MMFK3J/A)",
				"397925",
				"https://item.rakuten.co.jp/jtus/397925/?variantId=397925",
				"",
				"jtus",
				140800,
				11794,
			),
			last: NewRakutenProduct(
				"ELECOM 外付けHDD ELD-FTV020UBK",
				"173391",
				"https://item.rakuten.co.jp/jtus/173391/?variantId=173391",
				"",
				"jtus",
				8480,
				852,
			),
		},
	}}
	parser := RakutenParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := testutil.CreateHttpResponse(tt.args.filename)
			if err != nil {
				logger.Error("file open error", err)
				panic(err)
			}
			defer resp.Body.Close()

			products, url := parser.ProductList(resp.Body)

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.url, url)
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
		})
	}
}

func TestProduct(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		jan string
		err error
	}

	tests := []struct {
		name string
		args args
		want want
		wantErr bool
	}{{
		name: "parse product",
		args: args{"html/test_product.html"},
		want: want{
			jan: "4526541041112",
			err: nil,
		},
	}}
	p := RakutenParser{}

	for _, tt := range tests {
		res, err := testutil.CreateHttpResponse(tt.args.filename)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		jan, err := p.Product(res.Body)

		assert.Equal(t, tt.want.jan, jan)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
