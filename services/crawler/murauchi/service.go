package murauchi

import (
	"log"

	"crawler/config"
	"crawler/scrape"
)

var logger = config.Logger

func NewScrapeService(category string) scrape.Service[*MurauchiProduct] {
	service := scrape.NewService(MurauchiParser{}, &MurauchiProduct{}, []*MurauchiProduct{})
	req, err := generateRequest(0, category)
	if err != nil {
		log.Fatalln(err)
	}
	service.EntryReq = req
	return service
}

func RunAllCategories() {
	logger.Info("start all categories")
	categories, err := GetAllCategories()
	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		logger.Info("start scrape", "category", category)
		NewScrapeService(category).StartScrape("", "murauchi")
		logger.Info("end scrape", "category", category)
	}
	logger.Info("end all categories")
}

func GetAllCategories() ([]string, error) {
	res, err := scrape.NewClient().RequestURL("GET", "https://www.murauchi.com/MCJ-front-web/index.html", nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return MurauchiParser{}.FindCategories(res.Body)
}
