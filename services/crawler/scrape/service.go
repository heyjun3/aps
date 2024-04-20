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
	"crawler/product"
)

var logger = config.Logger

type IParser interface {
	ProductListByReq(io.ReadCloser, *http.Request) (product.Products, *http.Request)
	Product(io.ReadCloser) (string, error)
}

type Parser struct{}

func (p Parser) ConvToReq(products product.Products, url string) (product.Products, *http.Request) {
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

type Service struct {
	Parser            IParser
	Repo              ProductRepositoryInterface
	HistoryRepository RunServiceHistoryRepository
	EntryReq          *http.Request
	httpClient        HttpClient
	mqClient          RabbitMQClient
	fileId            string
}

func NewService[T product.IProduct](
	parser IParser, p T, ps []T, opts ...Option[T]) Service {
	s := &Service{
		Parser: parser,
		// Repo:              NewProductRepository(p, ps),
		HistoryRepository: RunServiceHistoryRepository{},
		httpClient:        NewClient(),
		mqClient:          NewMQClient(config.MQDsn, "mws"),
	}
	for _, opt := range opts {
		opt(s)
	}
	return *s
}

type Option[T product.IProduct] func(*Service)

func WithHttpClient[T product.IProduct](c HttpClient) func(*Service) {
	return func(s *Service) {
		s.httpClient = c
	}
}

func WithMQClient[T product.IProduct](c RabbitMQClient) func(*Service) {
	return func(s *Service) {
		s.mqClient = c
	}
}

func WithFileId[T product.IProduct](fileId string) func(*Service) {
	return func(s *Service) {
		s.fileId = fileId
	}
}

func WithCustomRepository(
	repo ProductRepositoryInterface) func(*Service) {
	return func(s *Service) {
		s.Repo = repo
	}
}

func (s Service) StartScrape(url, shopName string) {
	ctx := context.Background()
	db := CreateDBConnection(config.DBDsn)
	history := NewRunServiceHistory(shopName, url, "PROGRESS")
	s.HistoryRepository.Save(ctx, db, history)

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

	var fileId string
	if s.fileId == "" {
		fileId = shopName + "_" + TimeToStr(time.Now())
	} else {
		fileId = s.fileId
	}
	s.SendMessage(c4, fileId, &wg)

	wg.Wait()
	history.Status = "DONE"
	history.EndedAt = time.Now()
	s.HistoryRepository.Save(ctx, db, history)
}

func (s Service) ScrapeProductsList(req *http.Request) chan product.Products {
	c := make(chan product.Products, 100)
	go func() {
		defer close(c)
		for req != nil {
			logger.Info("request product list", "url", req.URL.String())
			res, err := s.httpClient.Request(req)
			if err != nil {
				logger.Error("http request error", err)
				break
			}
			var products product.Products
			products, req = s.Parser.ProductListByReq(res.Body, req)
			res.Body.Close()

			c <- products
		}
	}()
	return c
}

func (s Service) GetProductsBatch(ctx context.Context, db *bun.DB, c chan product.Products) chan product.Products {
	send := make(chan product.Products, 100)
	go func() {
		defer close(send)

		for p := range c {
			dbProduct, err := s.Repo.GetByCodes(ctx, db, p.GetCodes()...)
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

func (s Service) ScrapeProduct(
	ch chan product.Products) chan product.Products {

	send := make(chan product.Products, 100)
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

func (s Service) SaveProduct(ctx context.Context, db *bun.DB, ch chan product.Products) chan product.IProduct {

	send := make(chan product.IProduct, 100)
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

func (s Service) SendMessage(
	ch chan product.IProduct,
	fileId string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for p := range ch {
			m, err := p.GenerateMessage(fileId)
			if err != nil {
				logger.Error("generate message error", err)
				continue
			}

			err = s.mqClient.Publish(m)
			if err != nil {
				logger.Error("message publish error", err)
			}
		}
	}()
}
