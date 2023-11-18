package main

import (
	"api-server/spapi"
	"log"

	"os"
)

func main() {
	u := os.Getenv("SPAPI_SERVICE_URL")
	c, err := spapi.NewSpapiClient(u)
	if err != nil {
		log.Fatal(err)
	}
	c.GetLowestPricing([]string{"4528678022729-N-6426-210307"})
}
