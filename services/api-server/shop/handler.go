package shop

import (
	"context"
	"net/http"
	"os"

	shopv1 "api-server/shop/gen/shop/v1"
	"api-server/shop/gen/shop/v1/shopv1connect"
	"github.com/go-playground/validator/v10"

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

func AddValidateRules(validate *validator.Validate) *validator.Validate {
	rule := map[string]string{
		"Id":       "required",
		"SiteName": "required",
		"Name":     "required",
		"Url":      "required",
		"Interval": "required",
	}
	validate.RegisterStructValidationMapRules(map[string]string{"Shop": "dive"}, shopv1.Shops{})
	validate.RegisterStructValidationMapRules(rule, shopv1.Shop{})
	return validate
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
		slog.Error("bad request", "detail", err)
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

func DeleteShop(c echo.Context) error {
	type Ids struct {
		Id []string `json:"ids" validate:"required,min=1"`
	}
	ids := new(Ids)
	if err := c.Bind(ids); err != nil {
		slog.Error("bad request", "detail", err)
		return c.JSON(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(ids); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	res, err := client.DeleteShop(context.Background(), connect.NewRequest(&shopv1.DeleteShopRequest{
		Ids: ids.Id,
	}))
	if err != nil {
		slog.Error("failed delete shop", "detail", err)
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	slog.Info("success delete shop", "response", res.Msg)
	return c.JSON(http.StatusOK, nil)
}
