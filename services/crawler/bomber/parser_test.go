package bomber

import (
	"net/http"
	"net/url"
	"testing"

	"crawler/product"
	"crawler/test/util"

	"github.com/stretchr/testify/assert"
)

func NewTestbomberProduct(name, productCode, url, jan string, price int64) *product.Product {
	p, err := NewBomberProduct(name, productCode, url, jan, price)
	if err != nil {
		panic(err)
	}
	return p
}

func NewTestRequest(u string) *http.Request {
	URL := util.OmitError(url.Parse(u))
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
			args: args{filename: "html/test_product_list.html", url: "https://www.pc-bomber.co.jp/shop/pcbomber/c/c1204_s!0.price/?search_shop=1002&0_price_to=50000"},
			want: want{
				count: 60,
				req:   NewTestRequest("https://www.pc-bomber.co.jp/shop/pcbomber/c/c1204_s!0.price_p2/?search_shop=1002&0_price_to=50000"),
				first: NewTestbomberProduct("外でもドアホン VL-SVD505KS", "g1002-2546600202102", "https://www.pc-bomber.co.jp/shop/pcbomber/g/g1002-2546600202102/", "", 49980),
				last:  NewTestbomberProduct("F-VXU70-W ホワイト", "g1002-2546600237685", "https://www.pc-bomber.co.jp/shop/pcbomber/g/g1002-2546600237685/", "", 43800),
			},
		},
		{
			name: "parse last page",
			args: args{filename: "html/test_product_list_last.html", url: "https://www.pc-bomber.co.jp/shop/pcbomber/c/c1204_s!0.price_p22/?search_shop=1002&0_price_to=50000"},
			want: want{
				count: 3,
				req:   nil,
				first: NewTestbomberProduct("[取寄10]ソフトクリームしゃぼん玉 [1個][4902923152032]", "g1002-2523180000052", "https://www.pc-bomber.co.jp/shop/pcbomber/g/g1002-2523180000052/", "", 145),
				last:  NewTestbomberProduct("[取寄10]換気扇スイッチひも DZ-SH902 DZ-SH902 ホワイト [1本入り][4971275018488]", "g1002-2565300015249", "https://www.pc-bomber.co.jp/shop/pcbomber/g/g1002-2565300015249/", "", 74),
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
