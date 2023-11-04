package main

import (
	// "api-server/inventory"
	"fmt"
	"api-server/spapi"
)

func main() {
	// s := inventory.InventoryService{}
	// s.UpdatePricing()
	c := spapi.SpapiClient{}
	fmt.Println(c.URL)
}
