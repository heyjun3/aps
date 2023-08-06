package pc4u

import (
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
		first *Pc4uProduct
		last  *Pc4uProduct
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "parse product list page",
			args: args{filename: "html/test_product_list.html"},
			want: want{
				count: 50,
				url:   "https://www.pc4u.co.jp/view/search?page=2",
				first: NewTestPc4uProduct(
					"GIGABYTE B760I AORUS PRO DDR4 第13 & 12世代 Intel Core プロセッサー対応 Mini-ITX マザーボード｜B760I AORUS PRO DDR4",
					"000000081834",
					"https://www.pc4u.co.jp/view/item/000000081834",
					"",
					36960,
				),
				last: NewTestPc4uProduct(
					"ADATA XPG GAMMIX D20 16GB(16GBx1) DDR4 3600MHz(PC4-28800) U-DIMM SINGLE COLOR BOX ブラック｜AX4U360016G18I-CBK20",
					"000000081728",
					"https://www.pc4u.co.jp/view/item/000000081728",
					"",
					7429,
				),
			},
		},
		{
			name: "parse last page",
			args: args{filename: "html/test_product_list_last_page.html"},
			want: want{
				count: 17,
				url:   "",
				first: NewTestPc4uProduct(
					"【アウトレット特価・新品】Corning Thunderbolt 3, 50m Optical Cable Thunderboltケーブル｜AOC-CCU6JPN050M20",
					"000000072144",
					"https://www.pc4u.co.jp/view/item/000000072144?category_page_id=outlet",
					"",
					55990,
				),
				last: NewTestPc4uProduct(
					"【アウトレット特価・新品】Keyspan USBシリアルアダプタ  USB Serial Adapter (USA-19HS)",
					"014004000005",
					"https://www.pc4u.co.jp/view/item/014004000005?category_page_id=outlet",
					"",
					5940,
				),
			},
		},
		{
			name: "parse next page url",
			args: args{filename: "html/test_product_list_next_URL.html"},
			want: want{
				count: 50,
				url:   "https://www.pc4u.co.jp/view/search?page=6",
				first: NewTestPc4uProduct(
					"バッファロー ホウジンムケ RAID1タイオウ ソトヅケHDD 2ドライブ 4TB｜HD-WHA4U3/R1",
					"000000081503",
					"https://www.pc4u.co.jp/view/item/000000081503",
					"",
					46590,
				),
				last: NewTestPc4uProduct(
					"バッファロー TS5410RN3204 TS5410RNシリーズ 4ドライブ ラツクマウ｜TS5410RN3204",
					"000000081560",
					"https://www.pc4u.co.jp/view/item/000000081560",
					"",
					454300,
				),
			},
		},
		{
			name: "parse sold page url",
			args: args{filename: "html/test_product_list_soldout_next_URL.html"},
			want: want{
				count: 2,
				url:   "",
				first: NewTestPc4uProduct(
					"Keyspan  USBシリアルアダプタ USB Serial Adapter (USA-19HS)",
					"014004000001",
					"https://www.pc4u.co.jp/view/item/014004000001",
					"",
					6490,
				),
				last: NewTestPc4uProduct(
					"Cremax ICYDOCK Tray for MB559US-1S Black (MB559TRAY-B)",
					"032002000009",
					"https://www.pc4u.co.jp/view/item/032002000009",
					"",
					2269,
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			products, url := Pc4uParser{}.ProductList(res.Body, "")

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
	}
	tests := []struct {
		name  string
		args  args
		want  want
		isErr bool
	}{
		{
			name:  "parse product",
			args:  args{filename: "html/test_product.html"},
			want:  want{jan: "4719331990053"},
			isErr: false,
		},
		{
			name:  "parser product no contain table",
			args:  args{filename: "html/test_product_no_table.html"},
			want:  want{jan: "4537694092371"},
			isErr: false,
		},
		{
			name:  "parser product on table",
			args:  args{filename: "html/test_product_on_table.html"},
			want:  want{jan: "4719512135716"},
			isErr: false,
		},
		{
			name:  "code is EAN",
			args:  args{filename: "html/test_product_ean.html"},
			want:  want{jan: "195553843515"},
			isErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			code, err := Pc4uParser{}.Product(res.Body)

			assert.Equal(t, tt.want.jan, code)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
