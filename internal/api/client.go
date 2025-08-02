package api

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	token    string
	accountId     string
	httpClient *http.Client
}

func NewClient(token, accountid string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("Harvest token is required")
	}
	if accountid == "" {
		return nil, fmt.Errorf("Harvest account id is required")
	}

	return &Client{
		token: token,
		accountId:  accountid,
		httpClient: &http.Client{
			Timeout: time.Duration(30) * time.Second,
		},
	}, nil
}
