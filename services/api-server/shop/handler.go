package shop

import (
	"context"
	"net/http"
	"os"

	shopv1 "api-server/shop/gen/shop/v1"
	"api-server/shop/gen/shop/v1/shopv1connect"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

var client shopv1connect.ShopServiceClient

func init() {
	shopServiceURL := os.Getenv("SHOP_SERVICE_URL")
	if shopServiceURL == "" {
		panic("not found shop service url")
	}
	client = shopv1connect.NewShopServiceClient(
		http.DefaultClient,
		shopServiceURL,
	)
}

func GetShops(c echo.Context) error {
	res, err := client.ShopList(context.Background(), connect.NewRequest(&shopv1.ShopListRequest{}))
	if err != nil {
		slog.Error("error", "detail", err)
		return c.JSON(http.StatusInternalServerError, "error")
	}
	return c.JSON(http.StatusOK, res.Msg.GetShops())
}
