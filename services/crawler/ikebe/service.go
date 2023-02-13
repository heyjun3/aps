package ikebe

import (
	"context"
	"crawler/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/volatiletech/null/v8"
)

const (
	scheme = "https"
	host = "www.ikebe-gakki.com"
)

type ScrapeService struct {}

func (s ScrapeService) StartScrape(url string) {

	repo := IkebeProductRepository{}
	client := Client{&http.Client{}}
	products := []*models.IkebeProduct{}
	for url != "" {
		logger.Info("http request", "url", url)
		res, err := client.request("GET", url, nil)
		if err != nil {
			logger.Error("http request error", err)
			break
		}
		var product []*models.IkebeProduct
		product, url = parseProducts(res)
		products = append(products, product...)
	}

	var codes []string
	for _, p := range products {
		codes = append(codes, p.ProductCode)
	}

	ctx := context.Background()
	conn, _ := NewDBconnection(cfg.dsn())
	productsInDB, err := repo.getByProductCodes(ctx, conn, codes...)
	if err != nil {
		logger.Error("get ikebe products error", err)
	}

	products = mappingIkebeProducts(products, productsInDB)
	for _, product := range products {
		if product.Jan.Valid == true {
			continue
		}
		logger.Info("http request", "url", product.URL.String)
		res, err := client.request("GET", product.URL.String, nil)
		if err != nil {
			logger.Error("http request error", err)
		}
		jan, err := parseProduct(res)
		if err != nil {
			logger.Error("jancode is valid", err)
			continue
		}
		product.Jan = null.StringFrom(jan)
	}

	err = repo.bulkUpsert(conn, products...)
	if err != nil {
		logger.Error("bulk upsert is failed", err)
	}
	var messages [][]byte
	filename := "ikebe_" + timeToStr(time.Now())
	for _, p := range products {
		m, err := generateMessage(p, filename)
		if err != nil {
			logger.Error("generate message error", err)
			continue
		}
		messages = append(messages, m)
	}
	mqClient := NewMQClient(cfg.MQDsn(), "mws")
	mqClient.batchPublish(messages...)
}

type httpClient interface {
	request(string, string, io.Reader) (*http.Response, error)
}

type Client struct {
	httpClient *http.Client
}

func (c Client) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 3; i++ {
		res, err := c.httpClient.Do(req)
		time.Sleep(time.Second * 2)
		if err != nil && res.StatusCode > 200 {
			logger.Error("http request error", err)
			continue
		}
		return res, err
	}
	return nil, err
}

func mappingIkebeProducts(products, productsInDB []*models.IkebeProduct) []*models.IkebeProduct{
	inDB := map[string]*models.IkebeProduct{}
	for _, p := range productsInDB {
		inDB[p.ProductCode] = p
	}

	for _, product := range products {
		p := inDB[product.ProductCode]
		if p == nil {
			continue
		}
		product.Jan = p.Jan
	}
	return products
}

func generateMessage(p *models.IkebeProduct, filename string) ([]byte, error) {
	if !p.Jan.Valid {
		return nil, fmt.Errorf("Jan code isn't valid %s", p.ProductCode)
	}
	if !p.Price.Valid {
		return nil, fmt.Errorf("price isn't valid %s", p.ProductCode)
	}
	if !p.URL.Valid {
		return nil, fmt.Errorf("url isn't valid %s", p.ProductCode)
	}
	m := NewMWSSchema(filename, p.Jan.String, p.URL.String, p.Price.Int64)
	message, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return message, err
}

func timeToStr(t time.Time) string {
	return t.Format("20060102_150405")
}

func (s ScrapeService) scrapeProductsList(client httpClient, url string) chan []*models.IkebeProduct{
	c := make(chan []*models.IkebeProduct)
	go func() {
		defer close(c)
		for url != "" {
			res, err := client.request("GET", url, nil)
			if err != nil {
				logger.Error("http request error", err)
				break
			}
			var products []*models.IkebeProduct
			products, url = parseProducts(res)
			c <- products
		}
	}()
	return c
}

func (s ScrapeService) getIkebeProduct(c chan []*models.IkebeProduct, dsn string) chan []*models.IkebeProduct{
	send := make(chan []*models.IkebeProduct)
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
			products := mappingIkebeProducts(p, dbProduct)
			send <- products
		}
	}()
	return send
}

func (s ScrapeService) scrapeProduct(
	ch chan []*models.IkebeProduct, client httpClient)(
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

					res, err := client.request("GET", product.URL.String, nil)
					if err != nil {
						logger.Error("http request error", err, "action", "scrapeProduct")
						continue
					}
					jan, err := parseProduct(res)
					if err != nil {
						logger.Error("jan code isn't valid", err)
						continue
					}
					product.Jan = null.StringFrom(jan)
					send <- product
				}
			}
		}()
		return send
	}
