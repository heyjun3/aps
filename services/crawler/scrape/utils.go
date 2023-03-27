package scrape

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func timeToStr(t time.Time) string {
	return t.Format("20060102_150405")
}

func PullOutNumber(s string) (int64, error) {
	r := regexp.MustCompile(`[0-9]+`)
	strs := r.FindAllString(s, -1)
	if len(strs) == 0 {
		return 0, fmt.Errorf("pull out number error: %s", s)
	}

	price, err := strconv.Atoi(strings.Join(strs, ""))
	if err != nil {
		return 0, err
	}
	return int64(price), nil
}

// テスト時のみ使用している。
func CreateHttpResponse(path string) (*http.Response, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &http.Response{
		Body:    io.NopCloser(bytes.NewReader(b)),
		Request: &http.Request{},
	}
	return res, nil
}
