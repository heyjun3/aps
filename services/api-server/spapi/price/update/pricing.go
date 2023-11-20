package update

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func Pricing(URL *url.URL, sku string, price int) error {
	type input struct {
		Sku   string `json:"sku"`
		Price int    `json:"price"`
	}
	buf, err := json.Marshal(input{Sku: sku, Price: price})
	if err != nil {
		return err
	}
	URL.Path = "price"
	res, err := http.Post(URL.String(), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	slog.Info("update price response", "res", resBody)
	return nil
}
