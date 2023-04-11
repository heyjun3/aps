package ikebe

import (
	"fmt"
	"io"
	URL "net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"crawler/scrape"
)

type IkebeParser struct{}

func (parser IkebeParser) ProductList(r io.ReadCloser) (scrape.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, ""
	}

	var products scrape.Products
	doc.Find(".grid_item.product_item").Each(func(i int, s *goquery.Selection) {
		nameElem := s.Find(".item-information_productName.restrictTarget")
		name := nameElem.Text()
		if name == "" {
			logger.Info("Not Found product name")
			return
		}

		path, exist := nameElem.Attr("href")
		url, err := URL.Parse(path)
		if !exist || err != nil {
			logger.Info("Not Found url")
			return
		}
		url.Scheme = scheme
		url.Host = host

		productId := url.Query().Get("pid")
		if productId == "" {
			logger.Info("Not Found productId")
			return
		}

		price := s.Find(".price-bold.price-nomal").Text()
		if price == "" {
			price = s.Find(".price-bold.price-sale").Text()
		}
		p, err := scrape.PullOutNumber(price)
		if err != nil {
			logger.Info("Not Found price")
			return
		}

		products = append(products, NewIkebeProduct(name, productId, url.String(), "", int64(p)))
	})

	nextPath, exist := doc.Find(".product_pager-bottom .next a").Attr("href")
	if !exist || nextPath == "" || nextPath == "#" {
		logger.Info("Next Page URL is Not Found")
		return products, ""
	}
	u, err := URL.Parse(nextPath)
	if err != nil {
		logger.Error("url parse error", err)
		return products, ""
	}
	u.Scheme = scheme
	u.Host = host

	return products, u.String()
}

func (parser IkebeParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return "", err
	}

	jan := doc.Find(".removeSection_targetElm-jan").Text()
	if jan == "" {
		return "", fmt.Errorf("not found jan code")
	}

	rex := regexp.MustCompile("[0-9]{13}")
	janCode := rex.FindString(jan)
	if janCode == "" {
		return "", fmt.Errorf("not found jan code")
	}

	return janCode, nil
}
