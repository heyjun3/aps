package ikebe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProducts(t *testing.T) {
	t.Run("parse product list page", func(t *testing.T) {
		b, err := ioutil.ReadFile("html/product_list.html")
		if err != nil {
			fmt.Println("file open error")
		}
		res := http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(b)),
			Request: &http.Request{},
		}
		result := parseProducts(&res)

		assert.Equal(t, "hello", result)
	})
}
