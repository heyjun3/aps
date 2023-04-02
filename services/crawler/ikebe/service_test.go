package ikebe

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"os"
	"net/http"
	"io"
	"bytes"

	"github.com/stretchr/testify/assert"

	"crawler/config"
	"crawler/scrape"
)

type ClientMock struct {
	path string
}

func (c ClientMock) Request(method, url string, body io.Reader) (*http.Response, error) {
	b, err := os.ReadFile(c.path)
	if err != nil {
		logger.Error("open file error", err)
		return nil, err
	}
	res := &http.Response{
		Body:    io.NopCloser(bytes.NewReader(b)),
		Request: &http.Request{},
	}
	return res, nil
}

func TestScrapeProductsList(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c := ClientMock{"html/test_last_product_list.html"}
		s := NewScrapeService()

		ch := s.ScrapeProductsList(c, "https://google.com")

		p1 := NewIkebeProduct(
			"SR-SK30【次回3月入荷予定】",
			"124704",
			"https://www.ikebe-gakki.com/c/c-/pr/pr09/pr092127/124704",
			"",
			3267,
		)
		p17 := NewIkebeProduct(
			"SS-6B 【6口電源タップ】(SS6B)",
			"100469",
			"https://www.ikebe-gakki.com/c/c-/am/am09/am090814/100469",
			"",
			6050,
		)
		for p := range ch {
			assert.Equal(t, 17, len(p))
			assert.Equal(t, p1, p[0])
			assert.Equal(t, p17, p[len(p)-1])
		}
	})
}

func TestGetIkebeProduct(t *testing.T) {
	ctx := context.Background()
	conf, _ := config.NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn := scrape.CreateDBConnection(conf.Dsn())
	conn.NewDelete().Model((*IkebeProduct)(nil)).Exec(ctx)
	ps := []*IkebeProduct{
		NewIkebeProduct("test1", "test1", "http://", "1111", 1111),
		NewIkebeProduct("test2", "test2", "http://", "2222", 2222),
		NewIkebeProduct("test3", "test3", "http://", "3333", 3333),
	}
	for _, p := range ps {
		p.Upsert(conn, ctx)
		logger.Info(p.Name)
	}

	t.Run("happy path", func(t *testing.T) {
		s := NewScrapeService()
		p := scrape.Products{
			NewIkebeProduct("test1", "test1", "http://", "", 1111),
			NewIkebeProduct("test2", "test2", "http://", "", 2222),
			NewIkebeProduct("test3", "test3", "http://", "", 3333),
		}
		ch := make(chan scrape.Products)
		go func() {
			defer close(ch)
			ch <- p
		}()

		c := s.GetProducts(ch, conf.Dsn())

		var expected scrape.Products
		for _, p := range ps {
			expected = append(expected, p)
		}
		for product := range c {
			assert.Equal(t, expected, product)
		}
	})

	t.Run("get products return null", func(t *testing.T) {
		s := NewScrapeService()
		p := scrape.Products{
			NewIkebeProduct("test1", "test4", "http://", "", 1111),
			NewIkebeProduct("test2", "test5", "http://", "", 2222),
			NewIkebeProduct("test3", "test6", "http://", "", 3333),
		}

		ch := make(chan scrape.Products)
		go func() {
			defer close(ch)
			ch <- p
		}()

		c := s.GetProducts(ch, conf.Dsn())

		for product := range c {
			assert.Equal(t, p, product)
			assert.Equal(t, "", product[0].GetJan())
		}
	})
}

func TestScrapeProduct(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		s := NewScrapeService()
		c := ClientMock{"html/test_product.html"}
		p := scrape.Products{
			NewIkebeProduct("test1", "test4", "http://", "", 1111),
			NewIkebeProduct("test3", "test6", "http://", "", 3333),
		}
		ch := make(chan scrape.Products)
		go func() {
			defer close(ch)
			ch <- p
		}()

		channel := s.ScrapeProduct(ch, c)

		expectProduct := []scrape.IProduct{
			NewIkebeProduct("test1", "test4", "http://", "2500140008600", 1111),
			NewIkebeProduct("test3", "test6", "http://", "2500140008600", 3333),
		}
		var products []scrape.IProduct
		for product := range channel {
			products = append(products, product)
		}

		assert.Equal(t, expectProduct, products)
	})
}

func TestSaveProduct(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		conf, _ := config.NewConfig("../sqlboiler.toml")
		conf.Psql.DBname = "test"
		ch := make(chan scrape.IProduct)
		p := []*IkebeProduct{
			NewIkebeProduct("test1", "test4", "http://", "", 1111),
			NewIkebeProduct("test2", "test5", "http://", "", 2222),
			NewIkebeProduct("test3", "test6", "http://", "", 3333),
		}
		go func() {
			defer close(ch)
			for _, pro := range p {
				ch <- pro
			}
		}()
		s := NewScrapeService()

		channel := s.SaveProduct(ch, conf.Dsn())

		var ps []scrape.IProduct
		for p := range channel {
			ps = append(ps, p)
			fmt.Println(p)
		}
		var extProducts []scrape.IProduct
		for _, pro := range p {
			extProducts = append(extProducts, scrape.IProduct(pro))
		}
		assert.Equal(t, extProducts, ps)
	})
}

type MQMock struct{}

func (m MQMock) Publish(message []byte) error {
	fmt.Println(string(message))
	return nil
}

func TestSendMessage(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ch := make(chan scrape.IProduct)
		p := []*IkebeProduct{
			NewIkebeProduct("test1", "test4", "http://", "1111", 1111),
			NewIkebeProduct("test2", "test5", "http://", "2222", 2222),
			NewIkebeProduct("test3", "test6", "http://", "3333", 3333),
		}
		go func() {
			defer close(ch)
			for _, t := range p {
				ch <- t
			}
		}()
		c := MQMock{}
		s := NewScrapeService()
		wg := sync.WaitGroup{}
		wg.Add(1)

		s.SendMessage(ch, c, "ikebe", &wg)
	})
}
