package tutorial

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	point := float64(58)
	price := float64(5811)
	percent := int(math.Round(point / price * 100))

	assert.Equal(t, 1, percent)
}
