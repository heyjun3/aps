package main

import (
	"flag"
	"fmt"

	"crawler/ikebe"
	"crawler/pc4u"
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
		ikebe.NewScrapeService().StartScrape(url, shop)
	case shop == "pc4u" && url != "":
		pc4u.NewScrapeService().StartScrape(url, shop)
	default:
		ikebe.Tmp()
		fmt.Printf("argument error: s=%s, u=%s", shop, url)
	}
}
