package ikebe

import (
	"fmt"
	"log"
	URL "net/url"
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

	isSold := false
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

		path, exist := s.Find(".fs-c-productListItem__image.fs-c-productImage a[href]").Attr("href")
		url, err := URL.Parse(path)
		if exist == false || err != nil {
			fmt.Println("Not Found url")
			return
		}
		url.Scheme = scheme
		url.Host = host

		price := s.Find(".fs-c-productPrice__addon__price.fs-c-price .fs-c-price__value").Text()
		p, err := strconv.Atoi(strings.ReplaceAll(price, ",", ""))
		if err != nil {
			fmt.Println("Not Founc price")
			return
		}

		sold := s.Find(".fs-c-productListItem__outOfStock.fs-c-productListItem__notice.fs-c-productStock").Text()
		if sold == "SOLD" {
			fmt.Println("this product is sold out")
			isSold = true
			return
		}

		products = append(products, NewIkebeProduct(name, productId, url.String(), int64(p)))
	})

	nextPath, exist := doc.Find(".fs-c-pagination__item.fs-c-pagination__item--next[href]").First().Attr("href")
	u, err := URL.Parse(nextPath)
	if exist == false || err != nil || isSold == true {
		fmt.Println("Next Page URL is Not Found")
		return products, ""
	}
	u.Scheme = scheme
	u.Host = host

	return products, u.String()
}