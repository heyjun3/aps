// +build integration

package rakuten

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/scrape"
)

func TestProductListIntegration(t *testing.T){
	type args struct {
		url string
	}
	type want struct {
		count int
		url string
	}

	tests := []struct{
		name string
		args args
		want want
	}{{
		name: "parse product list integration test",
		args: args{
			url: "https://search.rakuten.co.jp/search/mall/?sid=197844",
		},
		want: want{
			count: 45,
			url: "https://search.rakuten.co.jp/search/mall/?p=2&sid=197844",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := RakutenParser{}
			res, err := scrape.NewClient().Request("GET", tt.args.url, nil)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			products, url := p.ProductList(res.Body, tt.args.url)

			assert.Equal(t, tt.want.count, len(products))
			assert.Equal(t, tt.want.url, url)

			for _, p := range products {
				assert.NotEmpty(t, p.(*RakutenProduct).Name)
				assert.NotEmpty(t, p.(*RakutenProduct).Price)
				assert.NotEmpty(t, p.(*RakutenProduct).ProductCode)
				assert.NotEmpty(t, p.(*RakutenProduct).ShopCode)
				assert.NotEmpty(t, p.(*RakutenProduct).URL)
			}
		})
	}
}
