package api

import "time"

type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
	Private   bool      `json:"private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEntryRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Status  string   `json:"status,omitempty"`
	Private bool     `json:"private,omitempty"`
}

type UpdateEntryRequest struct {
	Title   *string   `json:"title,omitempty"`
	Content *string   `json:"content,omitempty"`
	Tags    *[]string `json:"tags,omitempty"`
	Status  *string   `json:"status,omitempty"`
	Private *bool     `json:"private,omitempty"`
}

type ListParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type CreateEntryResponse struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type ListEntriesResponse struct {
	Entries []*Entry `json:"entries"`
	Total   int      `json:"total"`
}
