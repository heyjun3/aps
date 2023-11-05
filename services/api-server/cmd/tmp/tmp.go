package main

import (
	// "api-server/inventory"
	"api-server/spapi"
	"fmt"
)

func main() {
	// s := inventory.InventoryService{}
	// s.UpdatePricing()
	c := spapi.SpapiClient{}
	fmt.Println(c.URL)
}
