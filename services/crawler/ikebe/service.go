package ikebe

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	scheme = "https"
	host   = "www.ikebe-gakki.com"
)

type ScrapeService struct{
	repo Repository
	parser Parser
}

func NewScrapeService(repo Repository, parser Parser) *ScrapeService{
	return &ScrapeService{
		repo: repo,
		parser: parser,
	}
}

type Repository interface {
	getByProductCodes(ctx context.Context, conn boil.ContextExecutor, codes ...string) (Products, error)
}

type Parser interface {
	productList(io.ReadCloser) (Products, string)
	product(io.ReadCloser) (string, error)
}

func (s ScrapeService) StartScrape(url, shopName string) {

	client := Client{&http.Client{}}
	mqClient := NewMQClient(cfg.MQDsn(), "mws")
	wg := sync.WaitGroup{}
	wg.Add(1)

	c1 := s.scrapeProductsList(client, url)
	c2 := s.getProducts(c1, cfg.dsn())
	c3 := s.scrapeProduct(c2, client)
	c4 := s.saveProduct(c3, cfg.dsn())
	s.sendMessage(c4, mqClient, shopName, &wg)

	wg.Wait()
}

func (s ScrapeService) scrapeProductsList(
	client httpClient, url string) chan Products {
	c := make(chan Products, 10)
	go func() {
		defer close(c)
		for url != "" {
			logger.Info("product list request url", "url", url)
			res, err := client.request("GET", url, nil)
			if err != nil {
				logger.Error("http request error", err)
				break
			}
			var products Products
			products, url = s.parser.productList(res.Body)
			res.Body.Close()

			c <- products
		}
	}()
	return c
}

func (s ScrapeService) getProducts(c chan Products, dsn string) chan Products {
	send := make(chan Products, 10)
	go func() {
		defer close(send)
		ctx := context.Background()
		conn, err := NewDBconnection(dsn)
		if err != nil {
			logger.Error("db open error", err)
			return
		}

		for p := range c {
			dbProduct, err := s.repo.getByProductCodes(ctx, conn, p.getProductCodes()...)
			if err != nil {
				logger.Error("db get product error", err)
				continue
			}
			products := p.mapProducts(dbProduct)
			send <- products
		}
	}()
	return send
}

func (s ScrapeService) scrapeProduct(
	ch chan Products, client httpClient) chan Product {

	send := make(chan Product)
	go func() {
		defer close(send)
		for products := range ch {
			for _, product := range products {
				if product.isValidJan() {
					send <- product
					continue
				}

				logger.Info("product request url", "url", product.getURL())
				res, err := client.request("GET", product.getURL(), nil)
				if err != nil {
					logger.Error("http request error", err, "action", "scrapeProduct")
					continue
				}
				jan, err := s.parser.product(res.Body)
				res.Body.Close()
				if err != nil {
					logger.Error("jan code isn't valid", err, "url", res.Request.URL)
					continue
				}
				product.setJan(jan)
				send <- product
			}
		}
	}()
	return send
}

func (s ScrapeService) saveProduct(ch chan Product, dsn string) chan Product {

	send := make(chan Product)
	go func() {
		defer close(send)
		ctx := context.Background()
		conn, err := NewDBconnection(dsn)
		if err != nil {
			logger.Error("db open error", err)
			return
		}
		for p := range ch {
			err := p.Upsert(ctx, conn, true, []string{"shop_code", "product_code"}, boil.Infer(), boil.Infer())
			if err != nil {
				logger.Error("product upsert error", err)
				continue
			}
			send <- p
		}
	}()
	return send
}

func (s ScrapeService) sendMessage(
	ch chan Product, client RabbitMQClient,
	shop_name string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		filename := shop_name + "_" + timeToStr(time.Now())
		for p := range ch {
			m, err := p.generateMessage(filename)
			if err != nil {
				logger.Error("generate message error", err)
				continue
			}

			err = client.publish(m)
			if err != nil {
				logger.Error("message publish error", err)
			}
		}
	}()
}
