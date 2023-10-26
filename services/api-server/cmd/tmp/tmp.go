package main

import (
	"api-server/inventory"
)

func main() {
	s := inventory.InventoryService{}
	s.UpdatePricing()
}
