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

func TestConvertPercentToPoint(t *testing.T) {
	percent := 2
	price := 5811
	
	point := int(math.Round((float64(price) / 100 * float64(percent))))

	assert.Equal(t, 116, point)
}

func TestMathRound(t *testing.T) {
	price := 5.5

	assert.Equal(t, 5, int(price))
	assert.Equal(t, float64(6), math.Round(price))
}
