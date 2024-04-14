package rakuten

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/product"
	"crawler/test/util"
)

func NewTestRakutenProduct(
	name, productCode, url, jan, shopCode string, price, point int64) *product.Product {
	p, _ := NewRakutenProduct(name, productCode, url, jan, shopCode, price, point)
	return p
}

func TestProductList(t *testing.T) {
	type args struct {
		filename string
	}
	type want struct {
		count int
		url   string
		first *product.Product
		last  *product.Product
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
			url:   "https://search.rakuten.co.jp/search/mall/?p=2&sid=233405",
			first: NewTestRakutenProduct(
				"CORSAIR コルセアCorsair SF750 PLATINUM 750W PC電源ユニット 80PLUS PLATINUM CP-9020186-JP(2476065)代引不可 送料無料",
				"itm0015775592",
				"https://item.rakuten.co.jp/e-zoa/itm0015775592/?variantId=itm0015775592",
				"",
				"e-zoa",
				26044,
				2630,
			),
			last: NewTestRakutenProduct(
				"CORSAIR コルセアK70 RGB TKL CHAMPION MX Cherry MX Speed 日本レイアウト ゲーミングキーボード CH9119014JP(2509921)送料無料",
				"itm0015786850",
				"https://item.rakuten.co.jp/e-zoa/itm0015786850/?variantId=itm0015786850",
				"",
				"e-zoa",
				19434,
				1964,
			),
		},
	}, {
		name: "parse last page",
		args: args{
			filename: "html/test_last_product_list.html",
		},
		want: want{
			count: 21,
			url:   "",
			first: NewTestRakutenProduct(
				"IO DATA UD-RPCASE1　Raspberry Pi 2/3用ケース",
				"1000-01530328-00000001",
				"https://item.rakuten.co.jp/ioplaza/1000-01530328-00000001/?variantId=1000-01530328-00000001",
				"",
				"ioplaza",
				2178,
				209,
			),
			last: NewTestRakutenProduct(
				"【税込み】【メーカー保証】三菱ケミカルメディア SR80SP50V1",
				"1000-00007675-00000001",
				"https://item.rakuten.co.jp/ioplaza/1000-00007675-00000001/?variantId=1000-00007675-00000001",
				"",
				"ioplaza",
				2508,
				245,
			),
		},
	}, {
		name: "parse last page and get next url",
		args: args{
			filename: "html/test_last_products.html",
		},
		want: want{
			count: 45,
			url:   "https://search.rakuten.co.jp/search/mall/?max=45412&p=1&s=12&sid=206032",
			first: NewTestRakutenProduct(
				"449276 シマノ リミテッドプロガードタイツ LLAサイズ(TFイエロー) SHIMANO FI-014U",
				"4969363449276-36-58834-n",
				"https://item.rakuten.co.jp/jism/4969363449276-36-58834-n/?variantId=4969363449276-36-58834-n",
				"",
				"jism",
				45470,
				6234,
			),
			last: NewTestRakutenProduct(
				"16MMF1.4_DCDN_C_EF-M シグマ 16mm F1.4 DC DN ※EF-Mレンズ（APS-Cサイズミラーレス用）",
				"0085126402716-34-52183-n",
				"https://item.rakuten.co.jp/jism/0085126402716-34-52183-n/?variantId=0085126402716-34-52183-n",
				"",
				"jism",
				45412,
				4986,
			),
		},
	}}
	parser := RakutenParser{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				logger.Error("file open error", err)
				panic(err)
			}
			defer resp.Body.Close()

			products, url := parser.ProductList(resp.Body, "")

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
		name    string
		args    args
		want    want
		wantErr bool
	}{{
		name: "parse product",
		args: args{"html/test_product.html"},
		want: want{
			jan: "4526541041112",
			err: nil,
		},
	}, {
		name: "parse product in jan length 12",
		args: args{"html/test_product_12.html"},
		want: want{
			jan: "718037891460",
			err: nil,
		},
	}}
	p := RakutenParser{}

	for _, tt := range tests {
		res, err := util.CreateHttpResponse(tt.args.filename)
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
