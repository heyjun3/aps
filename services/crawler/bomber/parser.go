package bomber

import (
	"io"
	"net/http"

	"crawler/scrape"
)

type BomberParser struct{}

func (p BomberParser) ProductListByReq(r io.ReadCloser, req *http.Request) (scrape.Products, *http.Request) {
	return nil, nil
}

func (p BomberParser) Product(r io.ReadCloser) (string, error) {
	return "", nil
}
