package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"

	shopv1 "crawler/server/gen/shop/v1"
	"crawler/server/gen/shop/v1/shopv1connect"
)

func main() {
	client := shopv1connect.NewShopServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	res, err := client.CreateShop(
		context.Background(),
		connect.NewRequest(&shopv1.CreateShopRequest{
			Shop: []*shopv1.Shop{
				{Id: "test", SiteName: "test", Name: "test", Url: "http://test.com", Interval: "daily"},
			},
		}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg)
	log.Println(res.Header())
}
