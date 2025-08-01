package api

import (
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) CreateEntry(req CreateEntryRequest) (*CreateEntryResponse, error) {
	var response CreateEntryResponse
	err := c.makeRequest("POST", "/entries", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) GetEntry(id string) (*Entry, error) {
	var entry Entry
	endpoint := fmt.Sprintf("/entries/%s", id)
	err := c.makeRequest("GET", endpoint, nil, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (c *Client) ListEntries(params ListParams) ([]*Entry, error) {
	queryParams := url.Values{}
	if params.Limit > 0 {
		queryParams.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		queryParams.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Filter != "" {
		queryParams.Set("filter", params.Filter)
	}

	endpoint := "/entries"
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	var response ListEntriesResponse
	err := c.makeRequest("GET", endpoint, nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Entries, nil
}

func (c *Client) UpdateEntry(id string, req UpdateEntryRequest) (*Entry, error) {
	var entry Entry
	endpoint := fmt.Sprintf("/entries/%s", id)
	err := c.makeRequest("PUT", endpoint, req, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (c *Client) DeleteEntry(id string) error {
	endpoint := fmt.Sprintf("/entries/%s", id)
	return c.makeRequest("DELETE", endpoint, nil, nil)
}
