package price

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type IUpdatePriceInput interface {
	UpdatePrice() UpdatePriceInput
}

type UpdatePriceInput struct {
	Sku   string `json:"sku"`
	Price int    `json:"price"`
}

func UpdatePricing(URL *url.URL, input IUpdatePriceInput) error {
	buf, err := json.Marshal(input.UpdatePrice())
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
