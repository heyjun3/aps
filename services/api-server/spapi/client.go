package spapi

import (
	"net/url"
)

type SpapiClient struct {
	URL *url.URL
}

func NewSpapiClient(URL string) (*SpapiClient, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	return &SpapiClient{
		URL: u,
	}, nil
}

