package main

import (
	"crawler/ikebe"
)

func main() {
	url := "https://www.ikebe-gakki.com/p/search?sort=latest&keyword=&tag=&tag=&tag=&minprice=100000&maxprice=200000&cat1=&value2=&cat2=&value3=&cat3=&tag=%E3%82%A2%E3%82%A6%E3%83%88%E3%83%AC%E3%83%83%E3%83%88&detailRadio=%E3%82%A2%E3%82%A6%E3%83%88%E3%83%AC%E3%83%83%E3%83%88&tag=&detailShop="
	ikebe.ScrapeService(url)
}
