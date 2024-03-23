package scrape

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/uptrace/bun"

	"crawler/config"
)

var logger = config.Logger

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

type Service[T IProduct] struct {
	Parser     IParser
	Repo       ProductRepository[T]
	EntryReq   *http.Request
	httpClient httpClient
}

func NewService[T IProduct](parser IParser, p T, ps []T) Service[T] {
	return Service[T]{
		Parser:     parser,
		Repo:       NewProductRepository(p, ps),
		httpClient: NewClient(),
	}
}

func (s Service[T]) StartScrape(url, shopName string) {
	db := CreateDBConnection(config.DBDsn)
	ctx := context.Background()
	history := NewRunServiceHistory(shopName, url, "PROGRESS")
	RunServiceHistoryRepository{}.Save(ctx, db, history)

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

	c1 := s.ScrapeProductsList(s.EntryReq)
	c2 := s.GetProductsBatch(ctx, db, c1)
	c3 := s.ScrapeProduct(c2)
	c4 := s.SaveProduct(ctx, db, c3)
	s.SendMessage(c4, mqClient, shopName, &wg)

	wg.Wait()
	history.Status = "DONE"
	history.EndedAt = time.Now()
	RunServiceHistoryRepository{}.Save(ctx, db, history)
}

func (s Service[T]) ScrapeProductsList(req *http.Request) chan Products {
	c := make(chan Products, 100)
	go func() {
		defer close(c)
		for req != nil {
			logger.Info("request product list", "url", req.URL.String())
			res, err := s.httpClient.Request(req)
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

func (s Service[T]) GetProductsBatch(ctx context.Context, db *bun.DB, c chan Products) chan Products {
	send := make(chan Products, 100)
	go func() {
		defer close(send)

		for p := range c {
			dbProduct, err := s.Repo.GetByProductAndShopCodes(ctx, db, p.getProductAndShopCodes()...)
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
	ch chan Products) chan Products {

	send := make(chan Products, 100)
	go func() {
		defer close(send)
		for products := range ch {
			for _, product := range products {
				if product.IsValidJan() {
					continue
				}

				logger.Info("product request url", "url", product.GetURL())
				res, err := s.httpClient.RequestURL("GET", product.GetURL(), nil)
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

func (s Service[T]) SaveProduct(ctx context.Context, db *bun.DB, ch chan Products) chan IProduct {

	send := make(chan IProduct, 100)
	go func() {
		defer close(send)
		for p := range ch {
			err := s.Repo.BulkUpsert(ctx, db, p)
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
