package ikebe

import (
	"fmt"
	"regexp"
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
		logger.Error("response parse error", err)
		return nil, ""
	}

	isSold := false
	var products []*models.IkebeProduct
	doc.Find(".fs-c-productList__list__item.fs-c-productListItem").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".fs-c-productName__name").Text()
		if name == "" {
			logger.Info("Not Found product name")
			return
		}

		productId, exist := s.Find("input[name=staffStartSkuCode]").Attr("value")
		if exist == false {
			logger.Info("Not Found productId")
			return
		}

		path, exist := s.Find(".fs-c-productListItem__image.fs-c-productImage a[href]").Attr("href")
		url, err := URL.Parse(path)
		if exist == false || err != nil {
			logger.Info("Not Found url")
			return
		}
		url.Scheme = scheme
		url.Host = host

		price := s.Find(".fs-c-productPrice__addon__price.fs-c-price .fs-c-price__value").Text()
		p, err := strconv.Atoi(strings.ReplaceAll(price, ",", ""))
		if err != nil {
			logger.Info("Not Founc price")
			return
		}

		sold := s.Find(".fs-c-productListItem__outOfStock.fs-c-productListItem__notice.fs-c-productStock").Text()
		if sold == "SOLD" {
			logger.Info("product is sold out")
			isSold = true
			return
		}

		products = append(products, NewIkebeProduct(name, productId, url.String(), "", int64(p)))
	})

	nextPath, exist := doc.Find(".fs-c-pagination__item.fs-c-pagination__item--next[href]").First().Attr("href")
	u, err := URL.Parse(nextPath)
	if exist == false || err != nil || isSold == true {
		logger.Info("Next Page URL is Not Found")
		return products, ""
	}
	u.Scheme = scheme
	u.Host = host

	return products, u.String()
}

func parseProduct(r *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(r)
	if err != nil {
		logger.Error("response parse error", err)
		return "", err
	}


	jan := doc.Find(".janCode").Text()
	if jan == "" {
		err = fmt.Errorf("Not Found jan code")
		return "", err
	}

	rex := regexp.MustCompile("[0-9]{13}")
	janCode := rex.FindString(jan)
	if janCode == "" {
		err = fmt.Errorf("Not Found jan code")
		return "", err
	}

	return janCode, nil
}