package ikebe

import (
	"context"
	"crawler/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/volatiletech/null/v8"
)

const (
	scheme = "https"
	host = "www.ikebe-gakki.com"
)

func ScrapeService(url string) {
	url = "https://www.ikebe-gakki.com/p/search?sort=latest&keyword=&tag=&tag=&tag=&minprice=&maxprice=100000&cat1=&value2=&cat2=&value3=&cat3=&tag=%E6%96%B0%E5%93%81&detailRadio=%E6%96%B0%E5%93%81&detailShop=null"

	httpClient := &http.Client{}
	products := []*models.IkebeProduct{}
	for url != "" {
		res, err := request(httpClient, "GET", url, nil)
		if err != nil {
			log.Fatal(err)
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
	productsInDB, err := getIkebeProductsByProductCode(ctx, conn, codes...)
	if err != nil {
		log.Fatalln(err)
	}

	products = mappingIkebeProducts(products, productsInDB)
	for _, product := range products {
		if product.Jan.Valid == true {
			continue
		}
		res, err := request(httpClient, "GET", product.URL.String, nil)
		if err != nil {
			log.Fatalln(err)
		}
		jan := parseProduct(res)
		product.Jan = null.StringFrom(*jan)
	}
	
	var messages [][]byte
	filename := "ikebe_" + timeToStr(time.Now())
	for _, p := range products {
		m, err := generateMessage(p, filename)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messages = append(messages, m)
	}
	mqClient := NewMQClient(cfg.MQDsn(), "mws")
	mqClient.batchPublish(messages...)
}

func request(client *http.Client, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalln("action=request message=new request error")
		log.Fatalln(err)
		return nil, err
	}

	for i := 0; i < 3; i++ {
		res, err := client.Do(req)
		time.Sleep(time.Second * 2)
		if err != nil && res.StatusCode > 200 {
			log.Fatalln(err)
			log.Fatalf("status code: %d %s", res.StatusCode, res.Status)
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
		log.Fatalln(err)
		return nil, err
	}
	return message, err
}

func timeToStr(t time.Time) string {
	return t.Format("20060102_150405")
}

// func scrapeProductsList(url string) chan<- *models.IkebeProduct{
// 	c := make(chan<- *models.IkebeProduct)
// 	go func() {
// 		defer close(c)
// 		httpClient := &http.Client{}
// 		for url != "" {
// 			res, err := request(httpClient, "GET", url, nil)
// 			if err != nil {
// 				log.Fatal(err)
// 				break
// 			}
// 			var product []*models.IkebeProduct
// 			product, url = parseProducts(res)
// 			for _, p := range product {
// 				c <- p
// 			}
// 		}
// 	}()
// 	return c
// }

// func getIkebeProduct(c <-chan *models.IkebeProduct) chan<- *models.IkebeProduct{
// 	send := make(chan *models.IkebeProduct)
// 	go func() {
// 		defer close(send)
// 		ctx := context.Background()
// 		conn, err := NewDBconnection(cfg.dsn())
// 		if err != nil {
// 			log.Fatalln(err)
// 			return
// 		}

// 		for p := range c {
// 			ikebe, err := models.FindIkebeProduct(ctx, conn, p.ShopCode, p.ProductCode)
// 			if err != nil {
// 				log.Fatalln(err)
// 				continue
// 			}
// 			if ikebe.Jan.Valid {
// 				p.Jan = ikebe.Jan
// 			}
// 			send <- p
// 		}
// 	}()
// 	return send
// }

func Tmp() {
	c := NewMQClient(cfg.MQDsn(), "mws")
	c.batchPublish([]byte("TEST"), []byte("HELLO"))
}
