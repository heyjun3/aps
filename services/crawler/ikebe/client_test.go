package ikebe

import (
	"io"
	"os"
	"net/http"
	"bytes"
)

type clientMock struct {
	path string
}

func (c clientMock) request(method, url string, body io.Reader) (*http.Response, error) {
	b, err := os.ReadFile(c.path)
	if err != nil {
		logger.Error("open file error", err)
		return nil, err
	}
	res := &http.Response{
		Body:    io.NopCloser(bytes.NewReader(b)),
		Request: &http.Request{},
	}
	return res, nil
}
