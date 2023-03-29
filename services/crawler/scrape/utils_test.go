package scrape

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestTimeToStr(t *testing.T) {
	t.Run("format time to str", func(t *testing.T) {
		d := time.Date(2023, 2, 9, 22, 59, 0, 0, time.Local)

		s := timeToStr(d)
		fmt.Println(s)
		assert.Equal(t, "20230209_225900", s)
	})
}

func TestPullOutPrice(t *testing.T) {
	t.Run("pull out price", func(t *testing.T) {
		s := " 199,800円"

		price, err := PullOutNumber(s)

		assert.Equal(t, nil, err)
		assert.Equal(t, int64(199800), price)
	})

	t.Run("pull out price not digits", func(t *testing.T) {
		s := "aaa  fdsagfda"

		price, err := PullOutNumber(s)

		assert.Error(t, err)
		assert.Equal(t, int64(0), price)
	})

	t.Run("blank string", func(t *testing.T) {
		s := ""

		price, err := PullOutNumber(s)

		assert.Error(t, err)
		assert.Equal(t, int64(0), price)
	})
}
