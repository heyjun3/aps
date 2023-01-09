package migrate

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)


func TestConvKeepaTimeToTime(t *testing.T) {
	t.Run("conver keepa time to time.Time", func(t *testing.T) {
		dt, err := ConvKeepaTimeToTime("5831784")

		assert.Equal(t, nil, err)
		assert.Equal(t, dt, time.Date(2022, 2, 2, 5, 24, 0, 0, time.Local))
	})
}
