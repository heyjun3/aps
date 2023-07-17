package murauchi

import (
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"crawler/scrape"
)

type MurauchiParser struct{}

func (p MurauchiParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		logger.Error("response parse error", err)
		return nil, nil
	}

	var products scrape.Products
	doc.Find("")

	return products, nil
}
