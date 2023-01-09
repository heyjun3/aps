package migrate

import (
	"testing"
	"time"
	"strconv"

	"github.com/stretchr/testify/assert"
)


func TestConvKeepaTimeToTime(t *testing.T) {
	t.Run("convert keepa time to time.Time", func(t *testing.T) {
		dt, err := ConvKeepaTimeToTime("5831784")

		assert.Equal(t, nil, err)
		assert.Equal(t, time.Date(2022, 2, 2, 5, 24, 0, 0, time.Local), *dt)
	})

	t.Run("a-z argument is not valid", func(t *testing.T) {
		dt, err := ConvKeepaTimeToTime("aazzz")
		
		assert.IsType(t, &strconv.NumError{}, err)
		assert.Equal(t, time.Time{}, *dt)
	})
}
