package murauchi

import (
	"io"
	"net/http"
	"testing"

	"crawler/test/util"

	"github.com/stretchr/testify/assert"
)

func TestProductListbyReq(t *testing.T) {
	type args struct {
		filename string
		req      *http.Request
	}
	type want struct {
		count int
		first *MurauchiProduct
		last  *MurauchiProduct
		url   string
		body  string
	}

	tests := []struct {
		name   string
		args   args
		want   want
		isLast bool
	}{{
		name: "parse murauchi product list",
		args: args{
			filename: "html/test_product_list.html",
			req:      util.OmitError(generateRequest(0, "1000000000001")),
		},
		want: want{
			count: 112,
			first: util.OmitError(NewMurauchiProduct(
				"NEC\n超値下げ 15.6型ノートPC (Core i5/8GBメモリ/256GB SSD/Officeなし) PC-VKM44XDFHB8CSEZZY",
				"0000025905156",
				"https://www.murauchi.com/MCJ-front-web/CoD/0000025905156",
				"",
				69999,
			)),
			last: util.OmitError(NewMurauchiProduct(
				"NEC\n超値下げ フルHD対応21.5型ワイド液晶ディスプレイ 5年保証 LCD-L221F ホワイト",
				"0000025912888",
				"https://www.murauchi.com/MCJ-front-web/CoD/0000025912888",
				"",
				7980,
			)),
			url:  "https://www.murauchi.com/MCJ-front-web/WH/front/Default.do",
			body: "categoryNo=1000000000001&handlingType=0&keyword=%81%40&listCount=120&mode=graphic&pageNumber=1&searchType=keyword&sortOrder=1&type=COMMODITY_LIST",
		},
		isLast: false,
	}, {
		name: "parse murauchi last page",
		args: args{
			filename: "html/test_product_list_last_page.html",
			req:      util.OmitError(generateRequest((0), "0000")),
		},
		want: want{
			count: 14,
			first: util.OmitError(NewMurauchiProduct(
				"Lenovo レノボ\nThinkPad 90W ACアダプター (X1 Carbon用) 0B46997",
				"0000012968837",
				"https://www.murauchi.com/MCJ-front-web/CoD/0000012968837",
				"",
				5636,
			)),
			last: util.OmitError(NewMurauchiProduct(
				"サンワサプライ\nPDA-PEN16 入力ペン 3本セット",
				"0000001629549",
				"https://www.murauchi.com/MCJ-front-web/CoD/0000001629549",
				"",
				481,
			)),
		},
		isLast: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponseOnSjis(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			products, req := MurauchiParser{}.ProductListByReq(res.Body, tt.args.req)

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
			if tt.isLast {
				assert.Equal(t, (*http.Request)(nil), req)
			} else {
				assert.Equal(t, tt.want.url, req.URL.String())
				body, _ := io.ReadAll(req.Body)
				assert.Equal(t, tt.want.body, string(body))
			}
		})
	}
}

func TestParseProduct(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		jan string
	}
	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{{
		name:  "parse product page",
		args:  args{"html/test_product.html"},
		want:  want{"4989027022188"},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponseOnSjis(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			jan, err := MurauchiParser{}.Product(res.Body)

			logger.Info(jan)
			assert.Equal(t, tt.want.jan, jan)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindCategories(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		first string
		last string
		count int
	}
	tests := []struct{
		name string
		args args
		want want
		isErr bool
	}{{
		name: "parse find categories",
		args: args{"html/test_top_page.html"},
		want: want{"1000000020895", "1000000021188", 481},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponseOnSjis(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			categories, err := MurauchiParser{}.FindCategories(res.Body)

			assert.Equal(t, tt.want.first, categories[0])
			assert.Equal(t, tt.want.last, categories[len(categories)-1])
			assert.Equal(t, tt.want.count, len(categories))
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
