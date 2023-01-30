package ikebe

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)


func parseProducts(r *http.Response) string{
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	doc.Find(".fs-c-productList__list__item.fs-c-productListItem").Each(func(i int, s *goquery.Selection) {
		productId, exist := s.Find("input[name=staffStartSkuCode]").Attr("value")
		if exist == false {
			fmt.Println("Not Found productId")
			return
		}

		url, exist := s.Find(".fs-c-productListItem__image.fs-c-productImage a[href]").Attr("href")
		if exist == false {
			fmt.Println("Not Found url")
			return
		}

		price := s.Find(".fs-c-productPrice__addon__price.fs-c-price .fs-c-price__value").Text()
		if price == "" {
			fmt.Println("Not Founc price")
			return
		}
		p, _ := strconv.Atoi(strings.ReplaceAll(price, ",", ""))
		fmt.Println(productId, url, p)
	})
	return "hello"
}