package scrape

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"crawler/config"
)

var logger = config.Logger

type Service[T IProduct] struct {
	Parser   IParser
	Repo     ProductRepository[T]
	EntryReq *http.Request
}

func NewService[T IProduct](parser IParser, p T, ps []T) Service[T] {
	return Service[T]{
		Parser: parser,
		Repo:   NewProductRepository(p, ps),
	}
}

type IParser interface {
	ProductListByReq(io.ReadCloser, *http.Request) (Products, *http.Request)
	Product(io.ReadCloser) (string, error)
}

type Parser struct{}

func (p Parser) ConvToReq(products Products, url string) (Products, *http.Request) {
	if url == "" {
		return products, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("create new request error", err)
		return products, nil
	}
	return products, req
}

func (s Service[T]) StartScrape(url, shopName string) {
	client := NewClient()
	mqClient := NewMQClient(config.MQDsn, "mws")
	wg := sync.WaitGroup{}
	wg.Add(1)

	var err error
	if s.EntryReq == nil {
		s.EntryReq, err = http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatalln(err)
		}
	}

	c1 := s.ScrapeProductsList(client, s.EntryReq)
	c2 := s.GetProductsBatch(c1, config.DBDsn)
	c3 := s.ScrapeProduct(c2, client)
	c4 := s.SaveProduct(c3, config.DBDsn)
	s.SendMessage(c4, mqClient, shopName, &wg)

	wg.Wait()
}

func (s Service[T]) ScrapeProductsList(
	client httpClient, req *http.Request) chan Products {
	c := make(chan Products, 100)
	go func() {
		defer close(c)
		for req != nil {
			logger.Info("request product list", "url", req.URL.String())
			res, err := client.Request(req)
			if err != nil {
				logger.Error("http request error", err)
				break
			}
			var products Products
			products, req = s.Parser.ProductListByReq(res.Body, req)
			res.Body.Close()

			c <- products
		}
	}()
	return c
}

func (s Service[T]) GetProductsBatch(c chan Products, dsn string) chan Products {
	send := make(chan Products, 100)
	go func() {
		defer close(send)
		ctx := context.Background()
		conn := CreateDBConnection(dsn)

		for p := range c {
			dbProduct, err := s.Repo.GetByProductAndShopCodes(ctx, conn, p.getProductAndShopCodes()...)
			if err != nil {
				logger.Error("db get product error", err)
				continue
			}
			products := p.MapProducts(dbProduct)
			send <- products
		}
	}()
	return send
}

func (s Service[T]) ScrapeProduct(
	ch chan Products, client httpClient) chan Products {

	send := make(chan Products, 100)
	go func() {
		defer close(send)
		for products := range ch {
			for _, product := range products {
				if product.IsValidJan() {
					continue
				}

				logger.Info("product request url", "url", product.GetURL())
				res, err := client.RequestURL("GET", product.GetURL(), nil)
				if err != nil {
					logger.Error("http request error", err, "action", "scrapeProduct")
					continue
				}
				jan, _ := s.Parser.Product(res.Body)
				res.Body.Close()
				product.SetJan(jan)
			}
			send <- products
		}
	}()
	return send
}

func (s Service[T]) SaveProduct(ch chan Products, dsn string) chan IProduct {

	send := make(chan IProduct, 100)
	go func() {
		defer close(send)
		ctx := context.Background()
		conn := CreateDBConnection(dsn)
		for p := range ch {
			err := s.Repo.BulkUpsert(ctx, conn, p)
			if err != nil {
				logger.Error("product upsert error", err)
				continue
			}
			for _, product := range p {
				send <- product
			}
		}
	}()
	return send
}

func (s Service[T]) SendMessage(
	ch chan IProduct, client RabbitMQClient,
	shop_name string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		filename := shop_name + "_" + timeToStr(time.Now())
		for p := range ch {
			m, err := p.GenerateMessage(filename)
			if err != nil {
				logger.Error("generate message error", err)
				continue
			}

			err = client.Publish(m)
			if err != nil {
				logger.Error("message publish error", err)
			}
		}
	}()
}
