package rakuten

import (
	"crawler/scrape"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	scheme = ""
	host = ""
)

type RakutenParser struct{}

func (p RakutenParser) ProductList(r io.ReadCloser) (scrape.Products, string) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, ""
	}

	var products scrape.Products
	doc.Find(".dui-card.searchresultitem").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".title-link--3Ho6z").Text()
		if name == "" {
			logger.Info("Not Found product name")
			return
		}

		path, exist := s.Find("image a").Attr("href")
		URL, err := url.Parse(path)
		if !exist || err != nil {
			logger.Info("Not Found url", "name", name)
			return
		}
		URL.Scheme = scheme
		URL.Host = host

		paths := strings.Split(URL.Path, "/")
		productId := paths[len(paths)-1]
		shopId := paths[len(paths)-2]

		price, err := scrape.PullOutNumber(s.Find("price--OX_YW").Text())
		if err != nil {
			logger.Info("Not Found price", "name", name, "url", URL)
			return
		}

		point, err := scrape.PullOutNumber(s.Find(".points--AHzKn span").Text())
		if err != nil {
			logger.Info("Not Found point", "name", name, "url", URL)
		}

		products = append(products,
			NewRakutenProduct(name, productId, URL.String(), "", shopId, price, point))
	})

	nextURL, exist := doc.Find(".item.-next.nextPage").Attr("href")
	if !exist {
		logger.Info("Not Found Next Page URL")
		return products, ""
	}
	return products, nextURL
}
