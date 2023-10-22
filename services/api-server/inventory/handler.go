package inventory

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

func RefreshInventory(c echo.Context) error {
	host := os.Getenv("SPAPI_HOST")
	if host == "" {
		slog.Error("No set SPAPI_HOST")
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	query := url.Values{}
	query.Set("next_token", "")
	url := url.URL{
		Scheme: "http",
		Host: fmt.Sprintf("%s:8000", host),
		Path: "inventory-summaries",
		RawQuery: query.Encode(),
	}
	res, err := http.Get(url.String())
	if err != nil {
		slog.Error("get inventory summaries error", err)
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	return c.JSON(http.StatusOK, "success")
}
