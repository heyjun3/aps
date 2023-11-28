package point

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

type IUpdatePointInput interface {
	UpdatePoints() []UpdatePointInput
}

type UpdatePointInput struct {
	Sku          string `json:"sku"`
	PercentPoint int    `json:"percent_point"`
}

func UpdatePoints(URL *url.URL, input IUpdatePointInput) error {
	inputs := input.UpdatePoints()
	if len(inputs) == 0 {
		return errors.New("expect at least one input")
	}
	buf, err := json.Marshal(inputs)
	if err != nil {
		return err
	}
	URL.Path = "points"
	res, err := http.Post(URL.String(), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, _ := io.ReadAll(res.Body)
	slog.Info("update points response", "res", resBody)
	return nil
}
