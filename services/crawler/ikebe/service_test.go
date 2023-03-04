package ikebe

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"crawler/models"
)

func TestMappingIkebeProducts(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		p := IkebeProducts{
			NewIkebeProduct("test", "test", "http://test.jp", "", 1111),
			NewIkebeProduct("test1", "test1", "http://test.jp", "", 1111),
			NewIkebeProduct("test2", "test2", "http://test.jp", "", 1111),
		}

		dbp := IkebeProducts{
			NewIkebeProduct("test", "test", "test", "4444", 4000),
			NewIkebeProduct("test", "test1", "test1", "555", 4000),
			NewIkebeProduct("test", "test2", "test2", "7777", 4000),
		}

		result := p.mappingIkebeProducts(dbp)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, NewIkebeProduct("test", "test", "http://test.jp", "4444", 1111), result[0])
		assert.Equal(t, NewIkebeProduct("test1", "test1", "http://test.jp", "555", 1111), result[1])
		assert.Equal(t, NewIkebeProduct("test2", "test2", "http://test.jp", "7777", 1111), result[2])
	})

	t.Run("product is empty", func(t *testing.T) {
		p := IkebeProducts{}
		dbp := IkebeProducts{
			NewIkebeProduct("test", "test", "test", "11111", 4000),
			NewIkebeProduct("test", "test", "test1", "55555", 4000),
		}

		result := p.mappingIkebeProducts(dbp)

		assert.Equal(t, 0, len(result))
		assert.Equal(t, p, result)
	})

	t.Run("db product is empty", func(t *testing.T) {
		p := IkebeProducts{
			NewIkebeProduct("test", "test", "http://test.jp", "", 1111),
			NewIkebeProduct("test1", "test1", "http://test.jp", "", 1111),
			NewIkebeProduct("test2", "test2", "http://test.jp", "", 1111),
		}
		db := IkebeProducts{}

		result := p.mappingIkebeProducts(db)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, p, result)
	})
}

func TestGenerateMessage(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		p := NewIkebeProduct("test", "test", "https://test.com", "", 6000)
		p.Jan = null.StringFrom("4444")
		f := "ikebe_20220301_120303"

		m, err := p.generateMessage(f)

		assert.Equal(t, nil, err)
		ex := `{"filename":"ikebe_20220301_120303","jan":"4444","cost":6000,"url":"https://test.com"}`
		assert.Equal(t, ex, string(m))
	})

	t.Run("Jan code isn't Valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", "", 5000)
		f := "ikebe_20220202_020222"

		m, err := p.generateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("Price isn't valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", "", 5000)
		p.Price = null.Int64FromPtr(nil)
		f := "ikebe_20220202_020222"

		m, err := p.generateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("URL isn't valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", "", 5000)
		p.URL = null.StringFromPtr(nil)
		f := "ikebe_20220202_020222"

		m, err := p.generateMessage(f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})
}

func TestTimeToStr(t *testing.T) {
	t.Run("format time to str", func(t *testing.T) {
		d := time.Date(2023, 2, 9, 22, 59, 0, 0, time.Local)

		s := timeToStr(d)
		fmt.Println(s)
		assert.Equal(t, "20230209_225900", s)
	})
}

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

func TestScrapeProductsList(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		c := clientMock{"html/test_last_product_list.html"}
		s := ScrapeService{}

		ch := s.scrapeProductsList(c, "https://google.com")

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
		for p := range ch {
			assert.Equal(t, 17, len(p.slice()))
			assert.Equal(t, p1, p.slice()[0])
			assert.Equal(t, p17, p.slice()[len(p.slice())-1])
		}
	})
}

func TestGetIkebeProduct(t *testing.T) {
	ctx := context.Background()
	conf, _ := NewConfig("../sqlboiler.toml")
	conf.Psql.DBname = "test"
	conn, _ := NewDBconnection(conf.dsn())
	models.IkebeProducts().DeleteAll(ctx, conn)
	ps := IkebeProducts{
		NewIkebeProduct("test1", "test1", "http://", "1111", 1111),
		NewIkebeProduct("test2", "test2", "http://", "2222", 2222),
		NewIkebeProduct("test3", "test3", "http://", "3333", 3333),
	}
	ps.bulkUpsert(conn)

	t.Run("happy path", func(t *testing.T) {
		s := ScrapeService{}
		p := IkebeProducts{
			NewIkebeProduct("test1", "test1", "http://", "", 1111),
			NewIkebeProduct("test2", "test2", "http://", "", 2222),
			NewIkebeProduct("test3", "test3", "http://", "", 3333),
		}
		ch := make(chan Products)
		go func() {
			defer close(ch)
			ch <- p
		}()

		c := s.getIkebeProduct(ch, conf.dsn())

		for product := range c {
			assert.Equal(t, ps, product)
		}
	})

	t.Run("get products return null", func(t *testing.T) {
		s := ScrapeService{}
		p := IkebeProducts{
			NewIkebeProduct("test1", "test4", "http://", "", 1111),
			NewIkebeProduct("test2", "test5", "http://", "", 2222),
			NewIkebeProduct("test3", "test6", "http://", "", 3333),
		}

		ch := make(chan Products)
		go func() {
			defer close(ch)
			ch <- p
		}()

		c := s.getIkebeProduct(ch, conf.dsn())

		for product := range c {
			assert.Equal(t, p, product)
			assert.Equal(t, "", product.slice()[0].Jan.String)
		}
	})
}

func TestScrapeProduct(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		s := ScrapeService{}
		c := clientMock{"html/test_product.html"}
		p := IkebeProducts{
			NewIkebeProduct("test1", "test4", "http://", "", 1111),
			NewIkebeProduct("test3", "test6", "http://", "", 3333),
		}
		ch := make(chan Products)
		go func() {
			defer close(ch)
			ch <- p
		}()

		channel := s.scrapeProduct(ch, c)

		expectProduct := []Product{
			NewIkebeProduct("test1", "test4", "http://", "2500140008600", 1111),
			NewIkebeProduct("test3", "test6", "http://", "2500140008600", 3333),
		}
		var products []Product
		for product := range channel {
			products = append(products, product)
		}

		assert.Equal(t, expectProduct, products)
	})
}

func TestSaveProduct(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		conf, _ := NewConfig("../sqlboiler.toml")
		conf.Psql.DBname = "test"
		ch := make(chan Product)
		p := IkebeProducts{
			NewIkebeProduct("test1", "test4", "http://", "", 1111),
			NewIkebeProduct("test2", "test5", "http://", "", 2222),
			NewIkebeProduct("test3", "test6", "http://", "", 3333),
		}
		go func() {
			defer close(ch)
			for _, pro := range p {
				ch <- pro
			}
		}()
		s := ScrapeService{}

		channel := s.saveProduct(ch, conf.dsn())

		var ps []Product
		for p := range channel {
			ps = append(ps, p)
			fmt.Println(p)
		}
		var extProducts []Product
		for _, pro := range p {
			extProducts = append(extProducts, Product(pro))
		}
		assert.Equal(t, extProducts, ps)
	})
}

type MQMock struct{}

func (m MQMock) publish(message []byte) error {
	fmt.Println(string(message))
	return nil
}

func TestSendMessage(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ch := make(chan Product)
		p := IkebeProducts{
			NewIkebeProduct("test1", "test4", "http://", "1111", 1111),
			NewIkebeProduct("test2", "test5", "http://", "2222", 2222),
			NewIkebeProduct("test3", "test6", "http://", "3333", 3333),
		}
		go func() {
			defer close(ch)
			for _, t := range p {
				ch <- t
			}
		}()
		c := MQMock{}
		s := ScrapeService{}
		wg := sync.WaitGroup{}
		wg.Add(1)

		s.sendMessage(ch, c, "ikebe", &wg)
	})
}
