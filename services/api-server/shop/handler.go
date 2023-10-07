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

func CreateShop(c echo.Context) error {
	shops := new(shopv1.Shops)
	if err := c.Bind(shops); err != nil {
		slog.Error("bad request", "detail", err, "request body", c.Request().Body)
		return c.JSON(http.StatusBadRequest, "bad request")
	}
	
	if err := c.Validate(shops); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := client.CreateShop(context.Background(), connect.NewRequest(&shopv1.CreateShopRequest{
		Shops: shops,
	}))
	if err != nil {
		slog.Error("failed create shop", "detail", err)
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	slog.Info("success create shop", "response", res.Msg)
	return c.JSON(http.StatusOK, nil)
}
