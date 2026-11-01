package hsl

import (
	"net/http"
	"time"
)

const timeout = 10 * time.Second

type HslClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *HslClient {
	return &HslClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		apiKey: apiKey,
	}
}
