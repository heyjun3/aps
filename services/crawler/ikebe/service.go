package ikebe

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"crawler/models"
)

const (
	scheme = "https"
	host = "www.ikebe-gakki.com"
)

type ScrapeService struct {}

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

func (s ScrapeService) scrapeProductsList(client httpClient, url string) chan ikebeProducts{
	c := make(chan ikebeProducts, 10)
	go func() {
		defer close(c)
		for url != "" {
			logger.Info("product list request url", "url", url)
			res, err := client.request("GET", url, nil)
			if err != nil {
				logger.Error("http request error", err)
				break
			}
			var products ikebeProducts
			products, url = parseProducts(res.Body)
			res.Body.Close()
			c <- products
		}
	}()
	return c
}

func (s ScrapeService) getIkebeProduct(c chan ikebeProducts, dsn string) chan ikebeProducts{
	send := make(chan ikebeProducts, 10)
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
			var codes []string
			for _, pro := range p {
				codes = append(codes, pro.ProductCode)
			}
			dbProduct, err := repo.getByProductCodes(ctx, conn, codes...)
			if err != nil {
				logger.Error("db get product error", err)
				continue
			}
			products := NewIkebeProducts(p...).mappingIkebeProducts(dbProduct)
			send <- products
		}
	}()
	return send
}

func (s ScrapeService) scrapeProduct(
	ch chan ikebeProducts, client httpClient)(
	chan *models.IkebeProduct){

		send := make(chan *models.IkebeProduct)
		go func() {
			defer close(send)
			for products := range ch {
				for _, product := range products {
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

func (s ScrapeService) saveProduct(ch chan *models.IkebeProduct, dsn string) (
	chan *models.IkebeProduct) {
	
	send := make(chan *models.IkebeProduct)
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
	ch chan *models.IkebeProduct, client RabbitMQClient,
	shop_name string, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		filename := shop_name+ "_" + timeToStr(time.Now())
		for p := range ch {
			m, err := generateMessage(p, filename)
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
