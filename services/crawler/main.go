package main

import (
	"flag"
	"fmt"

	"crawler/ikebe"
)

func main() {
	var (
		shop string
		url  string
	)
	flag.StringVar(&shop, "s", "", "expect crawle shop")
	flag.StringVar(&url, "u", "", "expect crawle url")
	flag.Parse()

	switch {
	case shop == "ikebe" && url != "":
		s := ikebe.NewScrapeService(ikebe.IkebeProductRepository{}, ikebe.IkebeParser{})
		s.StartScrape(url, shop)
	default:
		fmt.Printf("argument error: s=%s, u=%s", shop, url)
	}
}
