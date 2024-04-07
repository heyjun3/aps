package scrape

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"crawler/test/util"
)

type ClientMock struct {
	path string
}

func (c ClientMock) RequestURL(method, url string, body io.Reader) (*http.Response, error) {
	return util.CreateHttpResponse(c.path)
}

func (c ClientMock) Request(req *http.Request) (*http.Response, error) {
	return util.CreateHttpResponse(c.path)
}

type ParserMock struct {
	products Products
	URL      string
	jan      string
	err      error
}

func (p ParserMock) ProductListByReq(r io.ReadCloser, req *http.Request) (Products, *http.Request) {
	if p.URL == "" {
		return p.products, nil
	}
	u, _ := url.Parse(p.URL)
	return p.products, &http.Request{URL: u}
}

func (p ParserMock) ProductList(doc io.ReadCloser, url string) (Products, string) {
	return p.products, p.URL
}

func (p ParserMock) Product(doc io.ReadCloser) (string, error) {
	return p.jan, p.err
}

func TestScrapeProductsList(t *testing.T) {
	type args struct {
		service Service[*Product]
		URL     string
	}
	type want struct {
		first IProduct
		last  IProduct
		len   int
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: Service[*Product]{
				Parser: ParserMock{
					products: Products{
						(NewTestProduct("test", "test1", "http://test.jp", "1111", "test", 1111)),
						(NewTestProduct("test", "test3", "http://test.jp", "3333", "test", 3333)),
						(NewTestProduct("test", "test2", "http://test.jp", "2222", "test", 2222)),
					},
					URL: "",
				},
				httpClient: ClientMock{"html/test_scrape_products_list.html"},
			},
			URL: "https://google.com",
		},
		want: want{
			first: (NewTestProduct("test", "test1", "http://test.jp", "1111", "test", 1111)),
			last:  (NewTestProduct("test", "test2", "http://test.jp", "2222", "test", 2222)),
			len:   3,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			req, _ := http.NewRequest("GET", tt.args.URL, nil)
			ch := tt.args.service.ScrapeProductsList(req)

			for p := range ch {
				assert.Equal(t, tt.want.first, p[0])
				assert.Equal(t, tt.want.last, p[len(p)-1])
				assert.Equal(t, tt.want.len, len(p))
			}
		})
	}
}

func TestGetProductsBatch(t *testing.T) {
	db, ctx := util.DatabaseFactory()

	type args struct {
		service  Service[*Product]
		products Products
		ctx      context.Context
		db       *bun.DB
	}
	type want struct {
		products Products
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: Service[*Product]{
				Repo: NewProductRepository(&Product{}, []*Product{}),
			},
			products: Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "", "test", 2222)),
				(NewTestProduct("test3", "test3", "http://test.jp", "", "test", 3333)),
				(NewTestProduct("test4", "test4", "http://test.jp", "", "test", 4444)),
			},
			ctx: ctx,
			db:  db,
		},
		want: want{
			Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "1111", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "2222", "test", 2222)),
				(NewTestProduct("test3", "test3", "http://test.jp", "3333", "test", 3333)),
				(NewTestProduct("test4", "test4", "http://test.jp", "4444", "test", 4444)),
			},
		},
	}, {
		name: "get products return null",
		args: args{
			service: NewService(ParserMock{}, &Product{}, []*Product{}, WithHttpClient[*Product](ClientMock{})),
			products: Products{
				(NewTestProduct("test11", "test11", "http://test.jp", "", "test", 1111)),
				(NewTestProduct("test22", "test22", "http://test.jp", "", "test", 2222)),
				(NewTestProduct("test33", "test33", "http://test.jp", "", "test", 3333)),
			},
			ctx: ctx,
			db:  db,
		},
		want: want{
			Products{
				(NewTestProduct("test11", "test11", "http://test.jp", "", "test", 1111)),
				(NewTestProduct("test22", "test22", "http://test.jp", "", "test", 2222)),
				(NewTestProduct("test33", "test33", "http://test.jp", "", "test", 3333)),
			},
		},
	}}

	db.ResetModel(ctx, (*Product)(nil))
	ps := Products{
		(NewTestProduct("test1", "test1", "http://test.jp", "1111", "test", 1111)),
		(NewTestProduct("test2", "test2", "http://test.jp", "2222", "test", 2222)),
		(NewTestProduct("test3", "test3", "http://test.jp", "3333", "test", 3333)),
		(NewTestProduct("test4", "test4", "http://test.jp", "4444", "test", 4444)),
	}
	ProductRepository[*Product]{}.BulkUpsert(ctx, db, ps)

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.GetProductsBatch(tt.args.ctx, tt.args.db, ch)

			assert.Equal(t, tt.want.products, <-c)
		})
	}
}

func TestScrapeProduct(t *testing.T) {
	type args struct {
		service  Service[*Product]
		products Products
	}
	type want struct {
		products Products
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: Service[*Product]{
				Parser: ParserMock{
					jan: "99999",
				},
				httpClient: ClientMock{"html/test_scrape_products_list.html"},
			},
			products: Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "", "test", 2222)),
			},
		},
		want: want{
			products: Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.ScrapeProduct(ch)

			assert.Equal(t, tt.want.products, <-c)
		})
	}
}

func TestSaveProduct(t *testing.T) {
	db, ctx := util.DatabaseFactory()

	type args struct {
		service  Service[*Product]
		products Products
		ctx      context.Context
		db       *bun.DB
	}
	type want struct {
		products Products
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: NewService(ParserMock{}, &Product{}, []*Product{}),
			products: Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
				(NewTestProduct("test3", "test3", "http://test.jp", "99999", "test", 3333)),
			},
			ctx: ctx,
			db:  db,
		},
		want: want{
			products: Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
				(NewTestProduct("test3", "test3", "http://test.jp", "99999", "test", 3333)),
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.SaveProduct(tt.args.ctx, tt.args.db, ch)

			var products Products
			for p := range c {
				products = append(products, p)
			}

			assert.Equal(t, tt.want.products, products)
		})
	}
}

type MQMock struct{}

func (m MQMock) Publish(message []byte) error {
	fmt.Println(string(message))
	return nil
}

func TestSendMessage(t *testing.T) {
	type args struct {
		service  Service[*Product]
		products Products
		siteName string
	}

	tests := []struct {
		name string
		args args
	}{{
		name: "happy path",
		args: args{
			service: Service[*Product]{
				mqClient: MQMock{},
			},
			products: Products{
				(NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
				(NewTestProduct("test3", "test3", "http://test.jp", "99999", "test", 3333)),
			},
			siteName: "test",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan IProduct, 10)
			for _, v := range tt.args.products {
				ch <- v
			}
			close(ch)
			wg := sync.WaitGroup{}
			wg.Add(1)

			tt.args.service.SendMessage(ch, tt.args.siteName, &wg)
			wg.Wait()
		})
	}
}
