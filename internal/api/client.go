package api

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string, timeoutSeconds int) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("API URL is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}, nil
}
