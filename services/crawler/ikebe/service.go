package ikebe

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	// "strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeService() {
	url := "https://www.ikebe-gakki.com/p/search?sort=latest&keyword=&tag=&tag=&tag=&minprice=&maxprice=100000&cat1=&value2=&cat2=&value3=&cat3=&tag=%E6%96%B0%E5%93%81&detailRadio=%E6%96%B0%E5%93%81&detailShop=null"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	parseProducts(res)
	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// doc, err := goquery.NewDocumentFromResponse(res)
	// if err != nil {
		// log.Fatal(err)
	// }

	// tags := doc.Find(".product-card")
	// fmt.Println(len(tags.Nodes))
}

func tmp() {
	r := http.Response{
		Body: ioutil.NopCloser(bytes.NewReader(nil)),
	}
	goquery.NewDocumentFromResponse(&r)
}
