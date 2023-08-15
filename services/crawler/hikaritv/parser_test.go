package hikaritv

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/test/util"
)

func NewTestHikaritvProduct(name, jan, productCode, url string, price int64) *HikaritvProduct {
	p, err := NewHikaritvProduct(name, jan, productCode, url, price)
	if err != nil {
		panic(err)
	}
	return p
}

func NewTestRequest(u string) *http.Request {
	URL := util.OmitError(url.Parse(u))
	URL.RawQuery = URL.Query().Encode()
	req := util.OmitError(http.NewRequest("GET", URL.String(), nil))
	return req
}

func TestProductListByReq(t *testing.T) {
	type args struct {
		filename string
		url      string
	}
	type want struct {
		count int
		req   *http.Request
		first *HikaritvProduct
		last  *HikaritvProduct
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "parse product list page",
			args: args{filename: "html/test_product_list.html", url: "https://shop.hikaritv.net/shopping/app/catalog/list/init?searchCategoryCode=0&searchWord=&searchCommodityCode=&searchMethod=0&searchType=0&squeezeSerch=0&hideKeyWord=&hidePriceMin=&hidePriceMax=50000&keywordToggle=&alignmentSequence=3&pageSize=200&mode=image&pageLayout=window&searchMakerName=&pointFacet=&discountRateFacet=&searchPriceStart=&searchPriceEnd=&searchTagCode=&searchCouponCode=&fqGetPoint=&fqStartDateMin=&fqStartDateMax=&fqStartDateName=&fqAverageRating=&banner=&notDisplayFacet=&currentPage=1&fqStockStatus=1"},
			want: want{
				count: 200,
				req:   NewTestRequest("https://shop.hikaritv.net/shopping/app/catalog/list/init?searchCategoryCode=0&searchWord=&searchCommodityCode=&searchMethod=0&searchType=0&squeezeSerch=0&hideKeyWord=&hidePriceMin=&hidePriceMax=50000&keywordToggle=&alignmentSequence=3&pageSize=200&mode=image&pageLayout=window&searchMakerName=&pointFacet=&discountRateFacet=&searchPriceStart=&searchPriceEnd=&searchTagCode=&searchCouponCode=&fqGetPoint=&fqStartDateMin=&fqStartDateMax=&fqStartDateName=&fqAverageRating=&banner=&notDisplayFacet=&currentPage=2&fqStockStatus=1"),
				first: NewTestHikaritvProduct(
					"ReFa EPI GO RE-AR-02A",
					"",
					"2010122365",
					"https://shop.hikaritv.net/shopping/commodity/plala/2010122365/",
					50000,
				),
				last: NewTestHikaritvProduct(
					"◇Aruba AP-655-CVR-20 20-pack for AP-655 white n-g s-on covers R7J45A",
					"",
					"6519012245",
					"https://shop.hikaritv.net/shopping/commodity/plala/6519012245/",
					47000,
				),
			},
		},
		{
			name: "parse last page",
			args: args{filename: "html/test_last_product_list.html", url: "https://shop.hikaritv.net/shopping/app/catalog/list/init?searchCategoryCode=0&searchWord=&searchCommodityCode=&searchMethod=0&searchType=0&squeezeSerch=0&hideKeyWord=&hidePriceMin=&hidePriceMax=50000&keywordToggle=&alignmentSequence=3&pageSize=200&mode=image&pageLayout=window&searchMakerName=&pointFacet=&discountRateFacet=&searchPriceStart=&searchPriceEnd=&searchTagCode=&searchCouponCode=&fqGetPoint=&fqStartDateMin=&fqStartDateMax=&fqStartDateName=&fqAverageRating=&banner=&notDisplayFacet=&currentPage=177&fqStockStatus=1"},
			want: want{
				count: 159,
				req:   nil,
				first: NewTestHikaritvProduct(
					"ビオレu キッチンハンドジェルソープ 無香料 詰替 200ml 4901301336224",
					"",
					"2010120258",
					"https://shop.hikaritv.net/shopping/commodity/plala/2010120258/",
					378,
				),
				last: NewTestHikaritvProduct(
					"アルカリボタン電池 LR44EC",
					"",
					"5410008499",
					"https://shop.hikaritv.net/shopping/commodity/plala/5410008499/",
					123,
				),
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

			products, r := HikaritvParser{}.ProductListByReq(res.Body, req)

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
		name:  "parse product",
		args:  args{filename: "html/test_product.html"},
		want:  want{code: "4948872414975"},
		isErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := util.CreateHttpResponse(tt.args.filename)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			code, err := HikaritvParser{}.Product(res.Body)

			assert.Equal(t, tt.want.code, code)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
