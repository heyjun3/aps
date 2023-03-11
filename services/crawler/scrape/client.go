package scrape

import (
	"io"
	"net/http"
	"time"
)

type httpClient interface {
	Request(string, string, io.Reader) (*http.Response, error)
}

type Client struct {
	httpClient *http.Client
}

func (c Client) Request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 3; i++ {
		res, err := c.httpClient.Do(req)
		time.Sleep(time.Second * 2)
		if err != nil && res.StatusCode > 200 {
			logger.Error("http request error", err)
			continue
		}
		return res, err
	}
	return nil, err
}
