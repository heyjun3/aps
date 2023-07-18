package murauchi

import (
	"net/http"
	"testing"

	"crawler/test/util"

	"github.com/stretchr/testify/assert"
)

func TestProductListbyReq(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		count int
		first *MurauchiProduct
		last *MurauchiProduct
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "parse murauchi product list",
		args: args{
			filename: "html/test_product_list.html",
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
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			products, req := MurauchiParser{}.ProductListByReq(res.Body, &http.Request{})

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
			assert.Equal(t, &http.Request{}, req)
		})
	}
}
