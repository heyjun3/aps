package ikebe

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	
	"crawler/models"
)


func parseProducts(r *http.Response) ([]*models.IkebeProduct, string) {
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		log.Fatal(err)
		return nil, ""
	}

	var products []*models.IkebeProduct
	doc.Find(".fs-c-productList__list__item.fs-c-productListItem").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".fs-c-productName__name").Text()
		if name == "" {
			fmt.Println("Not Found product name")
			return
		}

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
		p, err := strconv.Atoi(strings.ReplaceAll(price, ",", ""))
		if err != nil {
			fmt.Println("Not Founc price")
			return
		}

		products = append(products, NewIkebeProduct(name, productId, url, int64(p)))
	})

	nextURL, exist := doc.Find(".fs-c-pagination__item.fs-c-pagination__item--nex[href]").First().Attr("href")
	if exist == false {
		fmt.Println("Next Page URL is Not Found")
	}

	return products, nextURL
}