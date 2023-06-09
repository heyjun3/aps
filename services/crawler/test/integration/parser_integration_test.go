//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/ark"
	"crawler/ikebe"
	"crawler/nojima"
	"crawler/pc4u"
	"crawler/rakuten"
	"crawler/scrape"
)

func TestProductListIntegration(t *testing.T) {
	type args struct {
		parser scrape.IParser
		url    string
	}
	type want struct {
		count int
		url   string
	}

	tests := []struct {
		name   string
		args   args
		want   want
		hasJan bool
	}{
		{
			name: "ark product list",
			args: args{
				parser: ark.ArkParser{},
				url:    "https://www.ark-pc.co.jp/search/?limit=50&nouki=1",
			},
			want: want{
				count: 50,
				url:   "https://www.ark-pc.co.jp/search/?offset=50&limit=50&nouki=1",
			},
		},
		{
			name: "ikebe product list",
			args: args{
				parser: ikebe.IkebeParser{},
				url:    "https://www.ikebe-gakki.com/Form/Product/ProductList.aspx?shop=0&cat=&bid=ec&dpcnt=40&img=1&sort=07&udns=1&fpfl=0&sfl=0&pno=1&cicon=1",
			},
			want: want{
				count: 40,
				url:   "https://www.ikebe-gakki.com/Form/Product/ProductList.aspx?shop=0&cat=&bid=ec&cicon=1&dpcnt=40&img=1&sort=07&udns=1&fpfl=0&sfl=0&pno=2",
			},
		},
		{
			name: "nojima product list",
			args: args{
				parser: nojima.NojimaParser{},
				url:    "https://online.nojima.co.jp/app/catalog/list/init?searchCategoryCode=0&mode=image&pageSize=60&currentPage=1&alignmentSequence=9&searchDispFlg=true&immediateDeliveryDispFlg=1&searchWord=+",
			},
			want: want{
				count: 60,
				url:   "https://online.nojima.co.jp/app/catalog/list/init?alignmentSequence=9&currentPage=2&immediateDeliveryDispFlg=1&mode=image&pageSize=60&searchCategoryCode=0&searchDispFlg=true&searchWord=+",
			},
			hasJan: true,
		},
		{
			name: "pc4u product list",
			args: args{
				parser: pc4u.Pc4uParser{},
				url:    "https://www.pc4u.co.jp/view/search?search_keyword=&search_price_low=&search_price_high=&search_category=&search_original_code=",
			},
			want: want{
				count: 50,
				url:   "https://www.pc4u.co.jp/view/search?page=2",
			},
		},
		{
			name: "rakuten product list",
			args: args{
				url:    "https://search.rakuten.co.jp/search/mall/?sid=197844",
				parser: rakuten.RakutenParser{},
			},
			want: want{
				count: 45,
				url:   "https://search.rakuten.co.jp/search/mall/?p=2&sid=197844",
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := scrape.NewClient().RequestURL("GET", tc.args.url, nil)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			products, url := tc.args.parser.ProductList(res.Body, tc.args.url)

			assert.Equal(t, tc.want.count, len(products))
			assert.Equal(t, tc.want.url, url)

			for _, p := range products {
				assert.NotEmpty(t, p.GetName())
				assert.NotEmpty(t, p.GetPrice())
				assert.NotEmpty(t, p.GetProductCode())
				assert.NotEmpty(t, p.GetShopCode())
				assert.NotEmpty(t, p.GetURL())

				if tc.hasJan {
					assert.NotEmpty(t, p.GetJan())
				}
			}
		})
	}
}

func TestProductIntegration(t *testing.T) {
	type args struct {
		url    string
		parser scrape.IParser
	}
	type want struct {
		jan string
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "ark product",
			args: args{
				url:    "https://www.ark-pc.co.jp/i/10170046/",
				parser: ark.ArkParser{},
			},
			want: want{
				jan: "0735858504447",
			},
		},
		{
			name: "ikebe product",
			args: args{
				url:    "https://www.ikebe-gakki.com/Form/Product/ProductDetail.aspx?shop=0&pid=758741&bid=ec&cat=dli003001",
				parser: ikebe.IkebeParser{},
			},
			want: want{
				jan: "4033653012218",
			},
		},
		{
			name: "nojima product",
			args: args{
				url:    "https://online.nojima.co.jp/commodity/1/4902370550733/",
				parser: nojima.NojimaParser{},
			},
			want: want{
				jan: "4902370550733",
			},
		},
		{
			name: "pc4u product",
			args: args{
				url:    "https://www.pc4u.co.jp/view/item/000000082708",
				parser: pc4u.Pc4uParser{},
			},
			want: want{
				jan: "0840006698692",
			},
		},
		{
			name: "rakuten product",
			args: args{
				url:    "https://item.rakuten.co.jp/dj/596022/?s-id=bk_pc_item_list_name_d",
				parser: rakuten.RakutenParser{},
			},
			want: want{
				jan: "4042477257071",
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := scrape.NewClient().RequestURL("GET", tc.args.url, nil)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			jan, err := tc.args.parser.Product(res.Body)

			assert.Equal(t, tc.want.jan, jan)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
