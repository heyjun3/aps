package main

import (
	"context"
	"flag"
	"fmt"

	"crawler/ark"
	"crawler/bomber"
	"crawler/config"
	"crawler/hikaritv"
	"crawler/ikebe"
	"crawler/kaago"
	"crawler/murauchi"
	"crawler/nojima"
	"crawler/pc4u"
	"crawler/rakuten"
	"crawler/scrape"
	"crawler/shop"
)

func init() {
	models := []interface{}{
		&ark.ArkProduct{},
		&bomber.BomberProduct{},
		&hikaritv.HikaritvProduct{},
		&ikebe.IkebeProduct{},
		&kaago.KaagoProduct{},
		&murauchi.MurauchiProduct{},
		&nojima.NojimaProduct{},
		&pc4u.Pc4uProduct{},
		&rakuten.RakutenProduct{},
		&shop.Shop{},
		&scrape.RunServiceHistory{},
	}
	conn := scrape.CreateDBConnection(config.DBDsn)
	ctx := context.Background()

	for _, model := range models {
		if err := scrape.CreateTable(conn, ctx, model); err != nil {
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
	case shop == "bomber" && url != "":
		bomber.NewScrapeService().StartScrape(url, shop)
	case shop == "hikaritv" && url != "":
		hikaritv.NewScrapeService().StartScrape(url, shop)
	case shop == "ikebe" && url != "":
		ikebe.NewScrapeService().StartScrape(url, shop)
	case shop == "kaago" && url != "":
		kaago.NewScrapeService(url).StartScrape(url, shop)
	case shop == "murauchi" && category == "all":
		murauchi.RunAllCategories()
	case shop == "murauchi" && category != "":
		murauchi.NewScrapeService(category).StartScrape("", shop)
	case shop == "nojima" && category == "all":
		nojima.ScrapeAll()
	case shop == "pc4u" && url != "":
		pc4u.NewScrapeService().StartScrape(url, shop)
	case shop == "pc4u" && category == "all":
		pc4u.ScrapeAll(shop)
	case shop == "rakuten" && url != "":
		rakuten.NewScrapeService().StartScrape(url, shop)
	case shop == "rakuten" && id == "all":
		rakuten.RunServices()
	case shop == "rakuten" && id == "daily":
		rakuten.RunServicesByDaily()
	case shop == "move" && url == "":
		scrape.MoveMessages("mws", "mws")
	default:
		fmt.Printf("argument error: s=%s, u=%s", shop, url)
	}
}
