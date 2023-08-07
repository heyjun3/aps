package kaago

import (
	"io"
	"net/http"
	"testing"

	"crawler/test/util"

	"github.com/stretchr/testify/assert"
)

func NewTestKaagoProduct(name, productCode, url, jan, shopCode string, price int64) *KaagoProduct {
	p, err := NewKaagoProduct(name, productCode, url, jan, shopCode, price)
	if err != nil {
		panic(err)
	}
	return p
}

func TestProductListByReq(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		body        string
		count       int
		currentPage int
		url         string
		first       *KaagoProduct
		last        *KaagoProduct
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
			body:        "categorycode=0&currentPage=2&hasStock=1&shopcode=",
			count:       36,
			currentPage: 2,
			url:         "https://kaago.com/ajax/catalog/list/init",
			first: NewTestKaagoProduct(
				"S223ATES-W ダイキン ルームエアコン6畳 ホワイト",
				"9900000004842",
				"https://kaago.com/seicaplus/S223ATES-W-%E3%83%80%E3%82%A4%E3%82%AD%E3%83%B3-%E3%83%AB%E3%83%BC%E3%83%A0%E3%82%A8%E3%82%A2%E3%82%B3%E3%83%B36%E7%95%B3-%E3%83%9B%E3%83%AF%E3%82%A4%E3%83%88/?itemcode=9900000004842",
				"9900000004842",
				"seicaplus",
				46990,
			),
			last: NewTestKaagoProduct(
				"S71ZTCXP-W ダイキン ルームエアコン23畳 ホワイト 200V",
				"9900000004160",
				"https://kaago.com/seica/S71ZTCXP-W-%E3%83%80%E3%82%A4%E3%82%AD%E3%83%B3-%E3%83%AB%E3%83%BC%E3%83%A0%E3%82%A8%E3%82%A2%E3%82%B3%E3%83%B323%E7%95%B3-%E3%83%9B%E3%83%AF%E3%82%A4%E3%83%88-200V/?itemcode=9900000004160",
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

			products, req := KaagoParser{}.ProductListByReq(res.Body, &http.Request{})

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.url, req.URL.String())
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
			body, _ := io.ReadAll(req.Body)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}

func TestProduct(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		jan      string
		wantErr  bool
	}{{
		name:     "parse product page",
		filename: "html/test_product.html",
		jan:      "4549980516713",
		wantErr:  false,
	}}
	p := KaagoParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := util.CreateHttpResponse(tt.filename)
			if err != nil {
				panic(err)
			}
			defer r.Body.Close()

			jan, err := p.Product(r.Body)

			assert.Equal(t, tt.jan, jan)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
