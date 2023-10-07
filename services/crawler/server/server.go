package server

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/uptrace/bun"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"crawler/config"
	"crawler/scrape"
	"crawler/server/gen/greet/v1/greetv1connect"
	"crawler/server/gen/shop/v1/shopv1connect"
)

var logger = config.Logger
var db *bun.DB

func init() {
	db = scrape.CreateDBConnection(config.DBDsn)
}

func StartServer() {
	interceptors := connect.WithInterceptors(NewLoggerInterceptor())

	greetPath, greetHandler := greetv1connect.NewGreetServiceHandler(&GreetServer{}, interceptors)
	shopPath, shopHandler := shopv1connect.NewShopServiceHandler(&ShopServer{}, interceptors)

	mux := http.NewServeMux()
	mux.Handle(greetPath, greetHandler)
	mux.Handle(shopPath, shopHandler)

	http.ListenAndServe(
		"crawler_dev:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
