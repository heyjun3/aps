package rakuten

import (
	_ "embed"

	"gopkg.in/yaml.v3"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService() scrape.Service[*RakutenProduct] {
	return scrape.NewService(RakutenParser{}, &RakutenProduct{}, []*RakutenProduct{})
}

func RunServices() {
	shops, err := getShopList()
	if err != nil {
		logger.Error("error", err)
		return
	}
	for _, s := range shops.List {
		logger.Info("run service", "shop", s.ID, "url", s.URL)
		NewScrapeService().StartScrape(s.URL, "rakuten")
	}
}

type shop struct {
	ID  string `yaml:"id"`
	URL string `yaml:"url"`
}

type shops struct {
	List []shop `yaml:"rakuten"`
}

//go:embed shop.yaml
var contents []byte

func getShopList() (*shops, error) {
	s := shops{}
	if err := yaml.Unmarshal(contents, &s); err != nil {
		return nil, err
	}

	return &s, nil
}
