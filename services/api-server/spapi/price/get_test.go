package price

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Ptr[T any](v T) *T {
	return &v
}

func TestFilterCondition(t *testing.T) {
	buf, err := os.ReadFile("./res.json")
	if err != nil {
		panic(err)
	}
	var response GetLowestPricingResponse
	if err := json.Unmarshal(buf, &response); err != nil {
		panic(err)
	}

	offers := response.Responses[0].Body.Payload.Offers
	c := Condition{}
	filterd := offers.FilterCondition(c)
	assert.Equal(t, len(offers), len(filterd))

	c = Condition{
		MyOffer: Ptr(true),
	}
	myOffers := offers.FilterCondition(c)
	for _, o := range myOffers {
		assert.True(t, o.MyOffer)
	}

	c = Condition{
		IsFullfilledByAmazon: Ptr(true),
		IsBuyBoxWinner:       Ptr(true),
	}
	filterd = offers.FilterCondition(c)
	for _, o := range filterd {
		assert.True(t, o.IsBuyBoxWinner)
		assert.True(t, o.IsFulfilledByAmazon)
	}
}
