package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/uptrace/bun"

	"crawler/ark"
	"crawler/config"
	"crawler/ikebe"
	"crawler/kaago"
	"crawler/murauchi"
	"crawler/nojima"
	"crawler/pc4u"
	"crawler/rakuten"
	"crawler/scrape"
)

func init() {
	fs := []func(*bun.DB, context.Context) error{
		ark.CreateTable,
		ikebe.CreateTable,
		pc4u.CreateTable,
		nojima.CreateTable,
		kaago.CreateTable,
		murauchi.CreateTable,
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
		category string
		id       string
		shop     string
		url      string
	)
	flag.StringVar(&category, "c", "", "expect category")
	flag.StringVar(&id, "i", "", "expect crawle shop id")
	flag.StringVar(&shop, "s", "", "expect crawle shop")
	flag.StringVar(&url, "u", "", "expect crawle url")
	flag.Parse()

	switch {
	case shop == "ark" && url != "":
		ark.NewScrapeService().StartScrape(url, shop)
	case shop == "ikebe" && url != "":
		ikebe.NewScrapeService().StartScrape(url, shop)
	case shop == "kaago" && url != "":
		kaago.NewScrapeService(url).StartScrape(url, shop)
	case shop == "murauchi" && category != "":
		murauchi.NewScrapeService(category).StartScrape("", shop)
	case shop == "nojima" && url != "":
		nojima.NewScrapeService().StartScrape(url, shop)
	case shop == "pc4u" && url != "":
		pc4u.NewScrapeService().StartScrape(url, shop)
	case shop == "rakuten" && url != "":
		rakuten.NewScrapeService().StartScrape(url, shop)
	case shop == "rakuten" && id == "all":
		rakuten.RunServices()
	case shop == "move" && url == "":
		scrape.MoveMessages("mws", "mws")
	default:
		fmt.Printf("argument error: s=%s, u=%s", shop, url)
	}
}
