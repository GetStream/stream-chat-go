package stream_chat

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ImportTaskHistory struct {
	CreatedAt time.Time `json:"created_at"`
	NextState string    `json:"next_state"`
	PrevState string    `json:"prev_state"`
}

type ImportTask struct {
	CreatedAt time.Time            `json:"created_at"`
	Filename  string               `json:"filename"`
	History   []*ImportTaskHistory `json:"history"`
	ID        string               `json:"id"`
	State     string               `json:"state"`
	UpdatedAt time.Time            `json:"updated_at"`
	Result    interface{}          `json:"result"`
	Size      *int                 `json:"size"`
}

type ListImportsOptions struct {
	Limit  int
	Offset int
}

type CreateImportResponse struct {
	ImportTask *ImportTask `json:"import_task"`
	UploadURL  string      `json:"upload_url"`
	Response
}

type GetImportResponse struct {
	ImportTask *ImportTask `json:"import_task"`
	Response
}

type ListImportsResponse struct {
	ImportTasks []*ImportTask `json:"import_tasks"`
	Response
}

// CreateImport creates a new import task.
// Note: Do not use this.
// It is present for internal usage only.
// This function can, and will, break and/or be removed at any point in time.
func (c *Client) CreateImport(ctx context.Context, filename string) (*CreateImportResponse, error) {
	var resp CreateImportResponse
	err := c.makeRequest(ctx, http.MethodPost, "imports", nil, map[string]string{"filename": filename}, &resp)

	return &resp, err
}

// GetImport returns an import task.
// Note: Do not use this.
// It is present for internal usage only.
// This function can, and will, break and/or be removed at any point in time.
func (c *Client) GetImport(ctx context.Context, id string) (*GetImportResponse, error) {
	var resp GetImportResponse
	err := c.makeRequest(ctx, http.MethodGet, "imports/"+id, nil, nil, &resp)

	return &resp, err
}

// ListImports returns all import tasks.
// Note: Do not use this.
// It is present for internal usage only.
// This function can, and will, break and/or be removed at any point in time.
func (c *Client) ListImports(ctx context.Context, opts *ListImportsOptions) (*ListImportsResponse, error) {
	params := url.Values{}
	if opts != nil {
		params.Set("limit", strconv.Itoa(opts.Limit))
		params.Set("offset", strconv.Itoa(opts.Offset))
	}

	var resp ListImportsResponse
	err := c.makeRequest(ctx, http.MethodGet, "imports", params, nil, &resp)

	return &resp, err
}
