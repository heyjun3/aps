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