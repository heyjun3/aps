package rakuten

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService(shopCode string) *scrape.Service {
	return &scrape.Service{
		FetchProductByProductCodes: getByProductCodes(shopCode),
		Parser:                     RakutenParser{},
	}
}

func RunServices() {
	shops, err := getShopList()
	if err != nil {
		logger.Error("error", err)
		return
	}
	for _, s := range shops.List {
		logger.Info("run service", "shop", s.ID, "url", s.URL)
		NewScrapeService(s.ID).StartScrape(s.URL, "rakuten")
	}
}

type shop struct {
	ID  string `yaml:"id"`
	URL string `yaml:"url"`
}

type shops struct {
	List []shop `yaml:"rakuten"`
}

func getShopList() (*shops, error) {
	s := shops{}
	path := os.Getenv("ROOT_PATH")
	f, err := os.ReadFile(filepath.Join(path, "shop.yaml"))
	if err != nil {
		logger.Error("err", err)
	}
	if err := yaml.Unmarshal(f, &s); err != nil {
		return nil, err
	}

	return &s, nil
}
