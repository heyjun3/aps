package main

import (
	"api-server/product"
)

func main() {
	product.KeepaService{}.Consume()
}
