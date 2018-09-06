// Package pubg is a wrapper with helper functions for accessing pubg
// servers. Only thing that is required is a developer API key.
package pubg

import (
	"net/http"
	"net/url"
)

// Client is the main struct for pubg
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    *url.URL
}

// New returns a new defaulted Client struct.
func New(key string, httpClient *http.Client) (s *Client, err error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	url, err := url.Parse("https://api.pubg.com")
	if err != nil {
		return nil, err
	}

	return &Client{
		apiKey:     key,
		httpClient: httpClient,
		baseURL:    url,
	}, nil
}

func GetShards() (shards []string) {
	shards = append(shards, XboxNorthAmerica)
	shards = append(shards, XboxOceania)
	shards = append(shards, XboxEurope)
	shards = append(shards, XboxAsia)
	shards = append(shards, PCAsia)
	shards = append(shards, PCKAKAO)
	shards = append(shards, PCKorea)
	shards = append(shards, PCJapan)
	shards = append(shards, PCEurope)
	shards = append(shards, PCOceania)
	shards = append(shards, PCSouthAsia)
	shards = append(shards, PCKoreaJapan)
	shards = append(shards, PCNorthAmerica)
	return
}
