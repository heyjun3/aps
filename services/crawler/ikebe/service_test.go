package ikebe

import (
	"crawler/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func ikebeProductFactor(name, jan, shopCode, productCode, url string, price int64) *models.IkebeProduct{
	return &models.IkebeProduct{
		Name: null.StringFrom(name),
		Jan: null.StringFrom(jan),
		Price: null.Int64From(price),
		ShopCode: shopCode,
		ProductCode: productCode,
		URL: null.StringFrom(url),
	}
}

func TestMappingIkebeProducts(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		p := []*models.IkebeProduct{
			NewIkebeProduct("test", "test", "http://test.jp", 1111),
			NewIkebeProduct("test1", "test1", "http://test.jp", 1111),
			NewIkebeProduct("test2", "test2", "http://test.jp", 1111),
		}

		dbp := []*models.IkebeProduct{
			ikebeProductFactor("test", "4444", "test", "test", "", 4000),
			ikebeProductFactor("test", "555", "test", "test1", "", 4000),
			ikebeProductFactor("test", "77777", "test", "test2", "", 4000),
		}

		result := mappingIkebeProducts(p, dbp)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, ikebeProductFactor("test", "4444", "ikebe", "test", "http://test.jp", 1111), result[0])
		assert.Equal(t, ikebeProductFactor("test1", "555", "ikebe", "test1", "http://test.jp", 1111), result[1])
		assert.Equal(t, ikebeProductFactor("test2", "77777", "ikebe", "test2", "http://test.jp", 1111), result[2])
	})

	t.Run("product is empty", func(t *testing.T) {
		p := []*models.IkebeProduct{}
		dbp := []*models.IkebeProduct{
			ikebeProductFactor("test", "4444", "test", "test", "", 4000),
			ikebeProductFactor("test", "555", "test", "test1", "", 4000),
		}

		result := mappingIkebeProducts(p, dbp)

		assert.Equal(t, 0, len(result))
		assert.Equal(t, p, result)
	})

	t.Run("db product is empty", func(t *testing.T) {
		p := []*models.IkebeProduct{
			NewIkebeProduct("test", "test", "http://test.jp", 1111),
			NewIkebeProduct("test1", "test1", "http://test.jp", 1111),
			NewIkebeProduct("test2", "test2", "http://test.jp", 1111),
		}
		db := []*models.IkebeProduct{}

		result := mappingIkebeProducts(p, db)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, p, result)
	})
}

func TestGenerateMessage(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		p := NewIkebeProduct("test", "test", "https://test.com", 6000)
		p.Jan = null.StringFrom("4444")
		f := "ikebe_20220301_120303"

		m, err := generateMessage(p, f)

		assert.Equal(t, nil, err)
		ex := `{"filename":"ikebe_20220301_120303","jan":"4444","price":6000,"url":"https://test.com"}`
		assert.Equal(t, ex, string(m))
	})

	t.Run("Jan code isn't Valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", 5000)
		f := "ikebe_20220202_020222"

		m, err := generateMessage(p, f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("Price isn't valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", 5000)
		p.Price = null.Int64FromPtr(nil)
		f := "ikebe_20220202_020222"

		m, err := generateMessage(p, f)

		assert.Error(t, err)
		assert.Equal(t, []byte(nil), m)
	})

	t.Run("URL isn't valid", func(t *testing.T) {
		p := NewIkebeProduct("TEST", "test", "https://test.com", 5000)
		p.URL = null.StringFromPtr(nil)
		f := "ikebe_20220202_020222"

		m, err := generateMessage(p, f)

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
