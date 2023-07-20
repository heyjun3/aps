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
		name string
		args args
		want want
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
			url:  "https://www.murauchi.com/MCJ-front-web/WH/front/Default.do%3Ftype=COMMODITY_LIST",
			body: "categoryNo=1000000000001&handlingType=0&keyword=%E3%80%80&listCount=120&mode=graphic&pageNumber=1&searchType=keyword&sortOrder=&type=COMMODITY_LIST",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			products, req := MurauchiParser{}.ProductListByReq(res.Body, tt.args.req)

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
			assert.Equal(t, tt.want.url, req.URL.String())
			body, _ := io.ReadAll(req.Body)
			assert.Equal(t, tt.want.body, string(body))
		})
	}
}
