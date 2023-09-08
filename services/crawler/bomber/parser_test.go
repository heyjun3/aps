package bomber

import (
	"crawler/test/util"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewTestbomberProduct(name, productCode, url, jan string, price int64) *BomberProduct {
	p, err := NewBomberProduct(name, productCode, url, jan, price)
	if err != nil {
		panic(err)
	}
	return p
}

func NewTestRequest(u string) *http.Request {
	URL := util.OmitError(url.Parse(u))
	// URL.RawQuery = URL.Query().Encode()
	return util.OmitError(http.NewRequest("GET", URL.String(), nil))
}

func TestProductListByReq(t *testing.T) {
	type args struct {
		filename string
		url      string
	}
	type want struct {
		count int
		req   *http.Request
		first *BomberProduct
		last  *BomberProduct
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "parse product list page",
			args: args{filename: "html/test_product_list.html", url: "https://www.pc-bomber.co.jp/shop/pcbomber/c/c1204_s!0.price/?search_shop=1002&0_price_to=50000"},
			want: want{
				count: 60,
				req:   NewTestRequest("https://www.pc-bomber.co.jp/shop/pcbomber/c/c1204_s!0.price/?search_shop=1002&0_price_to=50000"),
				first: NewTestbomberProduct("外でもドアホン VL-SVD505KS", "g1002-2546600202102", "https://www.pc-bomber.co.jp/shop/pcbomber/g/g1002-2546600202102/", "", 49980),
				last:  NewTestbomberProduct("F-VXU70-W ホワイト", "g1002-2546600237685", "https://www.pc-bomber.co.jp/shop/pcbomber/g/g1002-2546600237685/", "", 43800),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.args.url, nil)
			if err != nil {
				panic(err)
			}
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			products, r := BomberParser{}.ProductListByReq(res.Body, req)

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
			assert.Equal(t, tt.want.req, r)
		})
	}
}

func TestProduct(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		code string
	}

	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{{
		name:  "parse product page",
		args:  args{filename: "html/test_product.html"},
		want:  want{code: "4550161577867"},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			code, err := BomberParser{}.Product(res.Body)

			assert.Equal(t, tt.want.code, code)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
