package main

import (
	"api-server/spapi"
	"log"
	"fmt"

	"os"
)

func main() {
	u := os.Getenv("SPAPI_SERVICE_URL")
	c, err := spapi.NewSpapiClient(u)
	if err != nil {
		log.Fatal(err)
	}
	res, err := c.GetLowestPricing([]string{"4528678022729-N-6426-210307", "850005352686-N-6980-20231104"})
	if err != nil {
		panic(err)
	}
	offer := res.Responses[0].Body.Payload.Offers[100]
	fmt.Println(offer)
}
