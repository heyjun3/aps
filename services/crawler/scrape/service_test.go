package scrape

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"crawler/test/util"
)

type ClientMock struct {
	path string
}

func (c ClientMock) Request(method, url string, body io.Reader) (*http.Response, error) {
	return util.CreateHttpResponse(c.path)
}

type ParserMock struct {
	products Products
	URL      string
	jan      string
	err      error
}

func (p ParserMock) ProductList(doc io.ReadCloser, url string) (Products, string) {
	return p.products, p.URL
}

func (p ParserMock) Product(doc io.ReadCloser) (string, error) {
	return p.jan, p.err
}

func TestScrapeProductsList(t *testing.T) {
	type args struct {
		client  httpClient
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
			client: ClientMock{"html/test_scrape_products_list.html"},
			service: Service[*Product]{
				Parser: ParserMock{
					products: Products{
						NewProduct("test", "test1", "http://test.jp", "1111", "test", 1111),
						NewProduct("test", "test3", "http://test.jp", "3333", "test", 3333),
						NewProduct("test", "test2", "http://test.jp", "2222", "test", 2222),
					},
					URL: "",
				},
			},
			URL: "https://google.com",
		},
		want: want{
			first: NewProduct("test", "test1", "http://test.jp", "1111", "test", 1111),
			last:  NewProduct("test", "test2", "http://test.jp", "2222", "test", 2222),
			len:   3,
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := tt.args.service.ScrapeProductsList(tt.args.client, tt.args.URL)

			for p := range ch {
				assert.Equal(t, tt.want.first, p[0])
				assert.Equal(t, tt.want.last, p[len(p)-1])
				assert.Equal(t, tt.want.len, len(p))
			}
		})
	}
}

func TestGetProductsBatch(t *testing.T) {
	type args struct {
		service  Service[*Product]
		products Products
		DSN      string
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
				NewProduct("test1", "test1", "http://test.jp", "", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "", "test", 3333),
				NewProduct("test4", "test4", "http://test.jp", "", "test", 4444),
			},
			DSN: util.TestDSN(),
		},
		want: want{
			Products{
				NewProduct("test1", "test1", "http://test.jp", "1111", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "2222", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "3333", "test", 3333),
				NewProduct("test4", "test4", "http://test.jp", "4444", "test", 4444),
			},
		},
	}, {
		name: "get products return null",
		args: args{
			service: Service[*Product]{
				Repo: NewProductRepository(&Product{}, []*Product{}),
			},
			products: Products{
				NewProduct("test11", "test11", "http://test.jp", "", "test", 1111),
				NewProduct("test22", "test22", "http://test.jp", "", "test", 2222),
				NewProduct("test33", "test33", "http://test.jp", "", "test", 3333),
			},
			DSN: util.TestDSN(),
		},
		want: want{
			Products{
				NewProduct("test11", "test11", "http://test.jp", "", "test", 1111),
				NewProduct("test22", "test22", "http://test.jp", "", "test", 2222),
				NewProduct("test33", "test33", "http://test.jp", "", "test", 3333),
			},
		},
	}}

	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Product)(nil))
	ps := Products{
		NewProduct("test1", "test1", "http://test.jp", "1111", "test", 1111),
		NewProduct("test2", "test2", "http://test.jp", "2222", "test", 2222),
		NewProduct("test3", "test3", "http://test.jp", "3333", "test", 3333),
		NewProduct("test4", "test4", "http://test.jp", "4444", "test", 4444),
	}
	ProductRepository[*Product]{}.BulkUpsert(ctx, conn, ps)

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.GetProductsBatch(ch, tt.args.DSN)

			assert.Equal(t, tt.want.products, <-c)
		})
	}
}

func TestGetProducts(t *testing.T) {
	type args struct {
		service  Service[*Product]
		products Products
		DSN      string
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
				NewProduct("test1", "test1", "http://test.jp", "", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "", "test", 3333),
			},
			DSN: util.TestDSN(),
		},
		want: want{
			products: Products{
				NewProduct("test1", "test1", "http://test.jp", "1111", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "2222", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "", "test", 3333),
			},
		},
	}}

	conn, ctx := util.DatabaseFactory()
	conn.ResetModel(ctx, (*Product)(nil))
	pre := Products{
		NewProduct("test1", "test1", "http://test.jp", "1111", "test", 1111),
		NewProduct("test2", "test2", "http://test.jp", "2222", "test", 2222),
		NewProduct("test3", "test3", "http://test.jp", "", "test", 3333),
	}
	ProductRepository[*Product]{}.BulkUpsert(ctx, conn, pre)

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			p := tt.args.service.GetProduct(ch, tt.args.DSN)

			for ps := range p {
				assert.Equal(t, tt.want.products, ps)
			}
		})
	}
}

func TestScrapeProduct(t *testing.T) {
	type args struct {
		service  Service[*Product]
		products Products
		client   httpClient
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
			},
			products: Products{
				NewProduct("test1", "test1", "http://test.jp", "", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "", "test", 2222),
			},
			client: ClientMock{"html/test_scrape_products_list.html"},
		},
		want: want{
			products: Products{
				NewProduct("test1", "test1", "http://test.jp", "99999", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "99999", "test", 2222),
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.ScrapeProduct(ch, tt.args.client)

			assert.Equal(t, tt.want.products, <-c)
		})
	}
}

func TestSaveProduct(t *testing.T) {
	type args struct {
		service  Service[*Product]
		products Products
		DSN      string
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
			service: Service[*Product]{},
			products: Products{
				NewProduct("test1", "test1", "http://test.jp", "99999", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "99999", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "99999", "test", 3333),
			},
			DSN: util.TestDSN(),
		},
		want: want{
			products: Products{
				NewProduct("test1", "test1", "http://test.jp", "99999", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "99999", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "99999", "test", 3333),
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan Products, 10)
			ch <- tt.args.products
			close(ch)

			c := tt.args.service.SaveProduct(ch, tt.args.DSN)

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
		client   RabbitMQClient
		siteName string
	}

	tests := []struct {
		name string
		args args
	}{{
		name: "happy path",
		args: args{
			service: Service[*Product]{},
			products: Products{
				NewProduct("test1", "test1", "http://test.jp", "99999", "test", 1111),
				NewProduct("test2", "test2", "http://test.jp", "99999", "test", 2222),
				NewProduct("test3", "test3", "http://test.jp", "99999", "test", 3333),
			},
			client:   MQMock{},
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

			tt.args.service.SendMessage(ch, tt.args.client, tt.args.siteName, &wg)
			wg.Wait()
		})
	}
}
