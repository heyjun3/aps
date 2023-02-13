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
			return
		}
		res := http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(b)),
			Request: &http.Request{},
		}

		products, url := parseProducts(&res)

		assert.Equal(t, 40, len(products))
		assert.Equal(t, "https://www.ikebe-gakki.com/p/search?maxprice=100000&tag=%E6%96%B0%E5%93%81&page=2&sort=latest", url)

		p1 := NewIkebeProduct(
			"Apocalyptica(オンライン納品専用)※代引きはご利用いただけません",
			"750802",
			"https://www.ikebe-gakki.com/c/c-/dt/dt03/dt031548/dt031548002/750802",
			"",
			34300,
		)
		p40 := NewIkebeProduct(
			"DRUM MIDI - EDM GROOVES(オンライン納品専用)※代引きはご利用いただけません",
			"750628",
			"https://www.ikebe-gakki.com/c/c-/dt/dt03/dt031598/dt031598002/750628",
			"",
			3630,
		)
		assert.Equal(t, p1, products[0])
		assert.Equal(t, p40, products[len(products)-1])
	})

	t.Run("parse last product list page", func(t *testing.T) {
		b, err := ioutil.ReadFile("html/last_product_list.html")
		if err != nil {
			fmt.Println("file open err")
			return
		}
		res := http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(b)),
			Request: &http.Request{},
		}

		products, url := parseProducts(&res)

		assert.Equal(t, 17, len(products))
		assert.Equal(t, "", url)

		p1 := NewIkebeProduct(
			"SR-SK30【次回3月入荷予定】",
			"124704",
			"https://www.ikebe-gakki.com/c/c-/pr/pr09/pr092127/124704",
			"",
			3267,
		)
		p17 := NewIkebeProduct(
			"SS-6B 【6口電源タップ】(SS6B)",
			"100469",
			"https://www.ikebe-gakki.com/c/c-/am/am09/am090814/100469",
			"",
			6050,
		)
		assert.Equal(t, p1, products[0])
		assert.Equal(t, p17, products[len(products)-1])
	})
}

func TestParseProduct(t *testing.T) {
	t.Run("parse product page", func(t *testing.T) {
		f, err := ioutil.ReadFile("html/product.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		res := &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(f)),
			Request: &http.Request{},
		}

		jan, err := parseProduct(res)

		assert.Equal(t, nil, err)
		assert.Equal(t, "2500140008600", jan)
	})
}
