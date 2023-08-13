package hikaritv

import (
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger


type HikaritvParser struct {
	scrape.Parser
}

func (p HikaritvParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	doc.Find(".w50p .inner").Each(func(i int, s *goquery.Selection) {
	})
	return nil, nil
}
