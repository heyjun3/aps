package scrape

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type httpClient interface {
	Request(string, string, io.Reader) (*http.Response, error)
}

func NewClient() Client{
	return Client{
		httpClient: &http.Client{
			Transport: &crawlerRoundTripper{
				base: http.DefaultTransport,
				logger: logger,
				attempts: 3,
				waitTime: time.Second * 2,
			},
		},
	}
}

type Client struct {
	httpClient *http.Client
}

func (c Client) Request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

type crawlerRoundTripper struct {
	base http.RoundTripper
	logger *slog.Logger
	attempts int
	waitTime time.Duration
}

func (t *crawlerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		res *http.Response
		err error
	)
	for count := 0; count < t.attempts; count++ {
		res, err = t.base.RoundTrip(req)
		time.Sleep(t.waitTime)
		if !t.shouldRetry(res, err) {
			logger.Info("http request error", "statuCode", res.StatusCode, "url", req.URL.String())
			return res, err
		}
	}
	return res, err
}

func (t *crawlerRoundTripper) shouldRetry(res *http.Response, err error) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			return true
		}
	}

	if res != nil {
		if res.StatusCode > http.StatusBadRequest {
			return true
		}
	}
	return false
}
