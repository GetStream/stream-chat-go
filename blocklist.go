package stream_chat

import (
	"context"
	"net/http"
	"time"
)

type BlocklistBase struct {
	Name  string   `json:"name"`
	Words []string `json:"words"`
}

type Blocklist struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	BlocklistBase
}

type BlocklistCreateRequest struct {
	BlocklistBase
}

type GetBlocklistResponse struct {
	Blocklist *Blocklist `json:"blocklist"`
	Response
}

type ListBlocklistsResponse struct {
	Blocklists []*Blocklist `json:"blocklists"`
	Response
}

// CreateBlocklist creates a blocklist.
func (c *Client) CreateBlocklist(ctx context.Context, blocklist *BlocklistCreateRequest) (*Response, error) {
	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "blocklists", nil, blocklist, &resp)
	return &resp, err
}

// GetBlocklist gets a blocklist.
func (c *Client) GetBlocklist(ctx context.Context, name string) (*GetBlocklistResponse, error) {
	var resp GetBlocklistResponse
	err := c.makeRequest(ctx, http.MethodGet, "blocklists/"+name, nil, nil, &resp)
	return &resp, err
}

// UpdateBlocklist updates a blocklist.
func (c *Client) UpdateBlocklist(ctx context.Context, name string, words []string) (*Response, error) {
	var resp Response
	err := c.makeRequest(ctx, http.MethodPut, "blocklists/"+name, nil, map[string][]string{"words": words}, &resp)
	return &resp, err
}

// ListBlocklists lists all blocklists.
func (c *Client) ListBlocklists(ctx context.Context) (*ListBlocklistsResponse, error) {
	var resp ListBlocklistsResponse
	err := c.makeRequest(ctx, http.MethodGet, "blocklists", nil, nil, &resp)
	return &resp, err
}

// DeleteBlocklist deletes a blocklist.
func (c *Client) DeleteBlocklist(ctx context.Context, name string) (*Response, error) {
	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, "blocklists/"+name, nil, nil, &resp)
	return &resp, err
}
