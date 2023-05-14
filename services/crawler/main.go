package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/uptrace/bun"

	"crawler/ark"
	"crawler/config"
	"crawler/ikebe"
	"crawler/pc4u"
	"crawler/rakuten"
	"crawler/scrape"
)

func init() {
	fs := []func(*bun.DB, context.Context) error{
		ark.CreateTable,
		ikebe.CreateTable,
		pc4u.CreateTable,
	}
	conn := scrape.CreateDBConnection(config.Config.Dsn())
	ctx := context.Background()

	for _, f := range fs {
		if err := f(conn, ctx); err != nil {
			panic(err)
		}
	}
}

func main() {
	var (
		shop string
		url  string
		id   string
	)
	flag.StringVar(&shop, "s", "", "expect crawle shop")
	flag.StringVar(&url, "u", "", "expect crawle url")
	flag.StringVar(&id, "i", "", "expect crawle shop id")
	flag.Parse()

	switch {
	case shop == "ark" && url != "":
		ark.NewScrapeService().StartScrape(url, shop)
	case shop == "ikebe" && url != "":
		ikebe.NewScrapeService().StartScrape(url, shop)
	case shop == "pc4u" && url != "":
		pc4u.NewScrapeService().StartScrape(url, shop)
	case shop == "rakuten" && url != "":
		rakuten.NewScrapeService().StartScrapeBySeries(url, shop)
	case shop == "rakuten" && id == "all":
		rakuten.RunServices()
	case shop == "move" && url == "":
		scrape.MoveMessages("mws", "mws")
	default:
		fmt.Printf("argument error: s=%s, u=%s", shop, url)
	}
}
