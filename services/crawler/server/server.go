package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"github.com/uptrace/bun"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"crawler/config"
	"crawler/rakuten"
	"crawler/scrape"
	greetv1 "crawler/server/gen/greet/v1"
	"crawler/server/gen/greet/v1/greetv1connect"
	shopv1 "crawler/server/gen/shop/v1"
	"crawler/server/gen/shop/v1/shopv1connect"
)

var logger = config.Logger
var db *bun.DB

func init() {
	db = scrape.CreateDBConnection(config.DBDsn)
}

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Msg.Name)
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func convertShopsIntoShopv1Shops(shops []rakuten.Shop) *shopv1.Shops {
	shopsv1 := []*shopv1.Shop{}
	for _, shop := range shops {
		shopv1 := shopv1.Shop{Id: shop.ID, SiteName: shop.SiteName, Name: shop.Name, Url: shop.URL, Interval: shop.Interval}
		shopsv1 = append(shopsv1, &shopv1)
	}
	return &shopv1.Shops{Shop: shopsv1}
}

type ShopServer struct{}

func (s *ShopServer) ShopList(ctx context.Context, req *connect.Request[shopv1.ShopListRequest]) (*connect.Response[shopv1.ShopListResponse], error) {
	logger.Info("ShopList", "status", "run")
	repo := rakuten.ShopRepository{}
	shops, err := repo.GetAll(db, context.Background())
	if err != nil {
		return nil, err
	}
	shopsv1 := convertShopsIntoShopv1Shops(shops)
	res := connect.NewResponse(&shopv1.ShopListResponse{
		Shops: shopsv1,
	})
	logger.Info("ShopList", "status", "done")
	return res, nil
}

func StartServer() {
	greeter := &GreetServer{}
	shopServer := &ShopServer{}

	greetPath, greetHandler := greetv1connect.NewGreetServiceHandler(greeter)
	shopPath, shopHandler := shopv1connect.NewShopServiceHandler(shopServer)

	mux := http.NewServeMux()
	mux.Handle(greetPath, greetHandler)
	mux.Handle(shopPath, shopHandler)

	http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
