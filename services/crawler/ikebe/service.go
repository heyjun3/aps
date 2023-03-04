package ikebe

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	scheme = "https"
	host   = "www.ikebe-gakki.com"
)

type Product interface {
	Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns boil.Columns, insertColumns boil.Columns) error
	generateMessage(filename string) ([]byte, error)
}

type Products interface {
	getProductCodes() []string
	mappingIkebeProducts(IkebeProducts) IkebeProducts
	slice() IkebeProducts
}

type ScrapeService struct{}

func (s ScrapeService) StartScrape(url string) {

	client := Client{&http.Client{}}
	mqClient := NewMQClient(cfg.MQDsn(), "mws")
	wg := sync.WaitGroup{}
	wg.Add(1)

	c1 := s.scrapeProductsList(client, url)
	c2 := s.getIkebeProduct(c1, cfg.dsn())
	c3 := s.scrapeProduct(c2, client)
	c4 := s.saveProduct(c3, cfg.dsn())
	s.sendMessage(c4, mqClient, "ikebe", &wg)

	wg.Wait()
}

func (s ScrapeService) scrapeProductsList(client httpClient, url string) chan Products {
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
			var products IkebeProducts
			products, url = parseProducts(res.Body)
			res.Body.Close()
			c <- products
		}
	}()
	return c
}

func (s ScrapeService) getIkebeProduct(c chan Products, dsn string) chan Products {
	send := make(chan Products, 10)
	go func() {
		defer close(send)
		ctx := context.Background()
		conn, err := NewDBconnection(dsn)
		if err != nil {
			logger.Error("db open error", err)
			return
		}

		repo := IkebeProductRepository{}
		for p := range c {
			dbProduct, err := repo.getByProductCodes(ctx, conn, p.getProductCodes()...)
			if err != nil {
				logger.Error("db get product error", err)
				continue
			}
			products := p.mappingIkebeProducts(dbProduct)
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
			for _, product := range products.slice() {
				if product.Jan.Valid {
					send <- product
					continue
				}

				logger.Info("product request url", "url", product.URL.String)
				res, err := client.request("GET", product.URL.String, nil)
				if err != nil {
					logger.Error("http request error", err, "action", "scrapeProduct")
					continue
				}
				jan, err := parseProduct(res.Body)
				res.Body.Close()
				if err != nil {
					logger.Error("jan code isn't valid", err, "url", res.Request.URL)
					continue
				}
				product.Jan = null.StringFrom(jan)
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
				logger.Error("ikebe product upsert error", err)
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
