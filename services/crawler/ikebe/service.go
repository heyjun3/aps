package ikebe

import (
	"crawler/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const (
	scheme = "https"
	host = "www.ikebe-gakki.com"
)

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
		var p []*models.IkebeProduct
		p, url = parseProducts(res)
		products = append(products, p...)
	}

	var codes []string
	for _, p := range products {
		codes = append(codes, p.ProductCode)
	}
	productsInDB := models.IkebeProducts(
		qm.WhereIn("")
	)

	// parseProducts(res)
	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// doc, err := goquery.NewDocumentFromResponse(res)
	// if err != nil {
		// log.Fatal(err)
	// }

	// tags := doc.Find(".product-card")
	// fmt.Println(len(tags.Nodes))
}

func Tmp() {
	url := true
	for url {
		fmt.Println("test")
		url = false
	}
}
