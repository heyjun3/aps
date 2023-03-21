package main

import (
	"flag"
	"fmt"

	"crawler/ikebe"
	"crawler/pc4u"
	"crawler/scrape"
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
	case shop == "move" && url == "":
		scrape.MoveMessages("test", "mws")
	default:
		fmt.Printf("argument error: s=%s, u=%s", shop, url)
	}
}
