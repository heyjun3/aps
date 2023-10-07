package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"

	shopv1 "api-server/shop/gen/shop/v1"
	"api-server/shop/gen/shop/v1/shopv1connect"
)

func main() {
	client := shopv1connect.NewShopServiceClient(
		http.DefaultClient,
		"http://crawler_dev:8080",
	)
	res, err := client.ShopList(
		context.Background(),
		connect.NewRequest(&shopv1.ShopListRequest{}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.GetShops().GetShop())
	log.Println(res.Header())
}
