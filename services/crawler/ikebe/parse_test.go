package ikebe

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/test/util"
)

func TestParseProducts(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		count int
		url   string
		first *IkebeProduct
		last  *IkebeProduct
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "parse product list page",
		args: args{
			filename: "html/test_product_list.html",
		},
		want: want{
			count: 40,
			url:   "https://www.ikebe-gakki.com/Form/Product/ProductList.aspx?shop=0&cat=&bid=ec&dpcnt=40&img=1&sort=07&udns=1&fpfl=0&sfl=0&pno=2",
			first: util.OmitError(NewIkebeProduct(
				"BLUE LAVA Touch wIdeal Bag (Ice Sail White) 【特価】",
				"755076",
				"https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=755076&bid=ec",
				"",
				99800,
			)),
			last: util.OmitError(NewIkebeProduct(
				"FENDER TONE SAVER 250K (#7706416049)",
				"755032",
				"https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=755032&bid=ec",
				"",
				6600,
			)),
		},
	}, {
		name: "parse last product list page",
		args: args{
			filename: "html/test_last_product_list.html",
		},
		want: want{
			count: 17,
			url:   "",
			first: util.OmitError(NewIkebeProduct(
				"CRY BABY 95Q WAH",
				"529",
				"https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=529&bid=ec",
				"",
				23100,
			)),
			last: util.OmitError(NewIkebeProduct(
				"PO-5S",
				"42",
				"https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=42&bid=ec",
				"",
				1925,
			)),
		},
	}, {
		name: "parse sale price products",
		args: args{
			filename: "html/test_sale_product_list.html",
		},
		want: want{
			count: 40,
			url:   "https://www.ikebe-gakki.com/Form/Product/ProductList.aspx?shop=0&cat=&bid=ec&cicon=6&dpcnt=40&img=1&sort=07&udns=1&fpfl=0&sfl=0&_cblCampaignIcon=6&_class=%e3%82%a2%e3%82%a6%e3%83%88%e3%83%ac%e3%83%83%e3%83%88&pno=2",
			first: util.OmitError(NewIkebeProduct(
				"Hydra Elite PRO 7 Trans Black Satine (T-BLK-S)【在庫処分超特価】",
				"754420",
				"https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=754420&bid=ec",
				"",
				473000,
			)),
			last: util.OmitError(NewIkebeProduct(
				"Marine Band 1896/20 (キー：Fm) 【特価】",
				"737193",
				"https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=737193&bid=ec",
				"",
				3850,
			)),
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := IkebeParser{}
			b, err := os.ReadFile(tt.args.filename)
			if err != nil {
				fmt.Println("file open error")
				return
			}
			res := &http.Response{
				Body:    io.NopCloser(bytes.NewReader(b)),
				Request: &http.Request{},
			}
			defer res.Body.Close()
			products, url := parser.ProductList(res.Body, "")

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.url, url)
			assert.Equal(t, tt.want.first, products[0])
			assert.Equal(t, tt.want.last, products[len(products)-1])
		})
	}
}

func TestParseProduct(t *testing.T) {
	t.Run("parse product page", func(t *testing.T) {
		parser := IkebeParser{}
		f, err := os.ReadFile("html/test_product.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		res := &http.Response{
			Body:    io.NopCloser(bytes.NewReader(f)),
			Request: &http.Request{},
		}
		defer res.Body.Close()
		jan, err := parser.Product(res.Body)

		assert.Equal(t, nil, err)
		assert.Equal(t, "4515515829030", jan)
	})
}
