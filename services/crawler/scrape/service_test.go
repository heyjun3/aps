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

	"crawler/product"
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
	products product.Products
	URL      string
	jan      string
	err      error
}

func (p ParserMock) ProductListByReq(r io.ReadCloser, req *http.Request) (product.Products, *http.Request) {
	if p.URL == "" {
		return p.products, nil
	}
	u, _ := url.Parse(p.URL)
	return p.products, &http.Request{URL: u}
}

func (p ParserMock) ProductList(doc io.ReadCloser, url string) (product.Products, string) {
	return p.products, p.URL
}

func (p ParserMock) Product(doc io.ReadCloser) (string, error) {
	return p.jan, p.err
}

func TestScrapeProductsList(t *testing.T) {
	type args struct {
		service Service
		URL     string
	}
	type want struct {
		first product.IProduct
		last  product.IProduct
		len   int
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: Service{
				Parser: ParserMock{
					products: product.Products{
						(product.NewTestProduct("test", "test1", "http://test.jp", "1111", "test", 1111)),
						(product.NewTestProduct("test", "test3", "http://test.jp", "3333", "test", 3333)),
						(product.NewTestProduct("test", "test2", "http://test.jp", "2222", "test", 2222)),
					},
					URL: "",
				},
				httpClient: ClientMock{"html/test_scrape_products_list.html"},
			},
			URL: "https://google.com",
		},
		want: want{
			first: (product.NewTestProduct("test", "test1", "http://test.jp", "1111", "test", 1111)),
			last:  (product.NewTestProduct("test", "test2", "http://test.jp", "2222", "test", 2222)),
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
		service  Service
		products product.Products
		ctx      context.Context
		db       *bun.DB
	}
	type want struct {
		products product.Products
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: Service{
				Repo: product.NewRepository(),
			},
			products: product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "", "test", 2222)),
				(product.NewTestProduct("test3", "test3", "http://test.jp", "", "test", 3333)),
				(product.NewTestProduct("test4", "test4", "http://test.jp", "", "test", 4444)),
			},
			ctx: ctx,
			db:  db,
		},
		want: want{
			product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "1111", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "2222", "test", 2222)),
				(product.NewTestProduct("test3", "test3", "http://test.jp", "3333", "test", 3333)),
				(product.NewTestProduct("test4", "test4", "http://test.jp", "4444", "test", 4444)),
			},
		},
	}, {
		name: "get products return null",
		args: args{
			service: NewService(
				ParserMock{},
				WithHttpClient(ClientMock{}),
				WithCustomRepository(product.NewRepository()),
			),
			products: product.Products{
				(product.NewTestProduct("test11", "test11", "http://test.jp", "", "test", 1111)),
				(product.NewTestProduct("test22", "test22", "http://test.jp", "", "test", 2222)),
				(product.NewTestProduct("test33", "test33", "http://test.jp", "", "test", 3333)),
			},
			ctx: ctx,
			db:  db,
		},
		want: want{
			product.Products{
				(product.NewTestProduct("test11", "test11", "http://test.jp", "", "test", 1111)),
				(product.NewTestProduct("test22", "test22", "http://test.jp", "", "test", 2222)),
				(product.NewTestProduct("test33", "test33", "http://test.jp", "", "test", 3333)),
			},
		},
	}}

	db.ResetModel(ctx, (*product.Product)(nil))
	ps := product.Products{
		(product.NewTestProduct("test1", "test1", "http://test.jp", "1111", "test", 1111)),
		(product.NewTestProduct("test2", "test2", "http://test.jp", "2222", "test", 2222)),
		(product.NewTestProduct("test3", "test3", "http://test.jp", "3333", "test", 3333)),
		(product.NewTestProduct("test4", "test4", "http://test.jp", "4444", "test", 4444)),
	}
	product.NewRepository().BulkUpsert(ctx, db, ps)

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan product.Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.GetProductsBatch(tt.args.ctx, tt.args.db, ch)

			assert.Equal(t, tt.want.products, <-c)
		})
	}
}

func TestScrapeProduct(t *testing.T) {
	type args struct {
		service  Service
		products product.Products
	}
	type want struct {
		products product.Products
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: Service{
				Parser: ParserMock{
					jan: "99999",
				},
				httpClient: ClientMock{"html/test_scrape_products_list.html"},
			},
			products: product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "", "test", 2222)),
			},
		},
		want: want{
			products: product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan product.Products, 10)
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
		service  Service
		products product.Products
		ctx      context.Context
		db       *bun.DB
	}
	type want struct {
		products product.Products
	}

	tests := []struct {
		name string
		args args
		want want
	}{{
		name: "happy path",
		args: args{
			service: NewService(
				ParserMock{},
				WithCustomRepository(product.NewRepository()),
			),
			products: product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
				(product.NewTestProduct("test3", "test3", "http://test.jp", "99999", "test", 3333)),
			},
			ctx: ctx,
			db:  db,
		},
		want: want{
			products: product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
				(product.NewTestProduct("test3", "test3", "http://test.jp", "99999", "test", 3333)),
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan product.Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.SaveProduct(tt.args.ctx, tt.args.db, ch)

			var products product.Products
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
		service  Service
		products product.Products
		siteName string
	}

	tests := []struct {
		name string
		args args
	}{{
		name: "happy path",
		args: args{
			service: Service{
				mqClient: MQMock{},
			},
			products: product.Products{
				(product.NewTestProduct("test1", "test1", "http://test.jp", "99999", "test", 1111)),
				(product.NewTestProduct("test2", "test2", "http://test.jp", "99999", "test", 2222)),
				(product.NewTestProduct("test3", "test3", "http://test.jp", "99999", "test", 3333)),
			},
			siteName: "test",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan product.IProduct, 10)
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
