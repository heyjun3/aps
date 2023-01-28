package buffalo

import (
	"fmt"
	"log"
	"net/http"
	// "strings"

	"github.com/PuerkitoBio/goquery"
)

func BuffaloScrapeAllPageService() {
	url := "https://buffalo-direct.com/collections/broadband"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// doc, err := goquery.NewDocumentFromReader(res.Body)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatal(err)
	}

	tags := doc.Find(".product-card")
	fmt.Println(len(tags.Nodes))
}
