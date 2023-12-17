package competitive

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnMarshalResponse(t *testing.T) {
	buf, err := os.ReadFile("./res.json")
	if err != nil {
		panic(err)
	}
	var res GetCompetitivePricingResponse
	if err := json.Unmarshal(buf, &res); err != nil {
		panic(err)
	}

	prices := res.LandedPrices()

	assert.Equal(t, "ASIN1", prices[0].Asin)
	assert.Equal(t, 18315, prices[0].LandedPrice.Amount)
	assert.Equal(t, 18315, prices[0].ListingPrice.Amount)
	assert.Equal(t, 62046, prices[0].SalesRankings[0].Rank)
	assert.Equal(t, "Success", prices[0].Status)

	assert.Equal(t, "ASIN2", prices[1].Asin)
	assert.Equal(t, 11662, prices[1].LandedPrice.Amount)
	assert.Equal(t, 11780, prices[1].ListingPrice.Amount)
	assert.Equal(t, 10003, prices[1].SalesRankings[0].Rank)
	assert.Equal(t, "Success", prices[1].Status)
}
