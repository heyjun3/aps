package bomber

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"crawler/scrape"

	"github.com/PuerkitoBio/goquery"
)

type BomberParser struct{}

func (p BomberParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	return nil, nil
}

func (p BomberParser) Product(r io.ReadCloser) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("[0-9]{12,13}")

	codes := re.FindAllString(doc.Find(".detail_goods_name2").Text(), -1)
	if len(codes) > 0 {
		return codes[0], nil
	}

	return "", fmt.Errorf("not found jan code")
}
