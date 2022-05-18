package stream_chat

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SegmentData struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Filter      *SegmentFilter `json:"filter,omitempty"`
}

type Segment struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	UserCount    int    `json:"user_count"`
	ChannelCount int    `json:"channel_count"`
	SegmentData

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateSegmentRequest struct {
	Segment SegmentData `json:"segment"`
}

type UpdateSegmentRequest struct {
	SegmentUpdateableFields `json:"segment"`
}

type SegmentUpdateableFields struct {
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Filter      *SegmentFilter `json:"filter,omitempty"`
}

type ListSegmentsOptions struct {
	Limit  int
	Offset int
}

type CreateSegmentResponse struct {
	Segment Segment `json:"segment"`
	Response
}

type GetSegmentResponse struct {
	Segment Segment `json:"segment"`
	Response
}

type ListSegmentsResponse struct {
	Segments []Segment `json:"segments"`
	Response
}

type UpdateSegmentResponse struct {
	Segment Segment `json:"segment"`
	Response
}

// CreateSegment creates a new segment.
func (c *Client) CreateSegment(ctx context.Context, segment *CreateSegmentRequest) (*CreateSegmentResponse, error) {
	var response CreateSegmentResponse
	err := c.makeRequest(ctx, http.MethodPost, "segments", nil, segment, &response)
	return &response, err
}

// GetSegment returns a segment.
func (c *Client) GetSegment(ctx context.Context, id string) (*GetSegmentResponse, error) {
	var response GetSegmentResponse
	err := c.makeRequest(ctx, http.MethodGet, "segments/"+id, nil, nil, &response)
	return &response, err
}

// ListSegments returns a list of segments.
func (c *Client) ListSegments(ctx context.Context, opts *ListSegmentsOptions) (*ListSegmentsResponse, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
	}

	var response ListSegmentsResponse
	err := c.makeRequest(ctx, http.MethodGet, "segments", params, nil, &response)
	return &response, err
}

// UpdateSegment updates a segment.
func (c *Client) UpdateSegment(ctx context.Context, id string, segment *UpdateSegmentRequest) (*UpdateSegmentResponse, error) {
	var response UpdateSegmentResponse
	err := c.makeRequest(ctx, http.MethodPut, "segments/"+id, nil, segment, &response)
	return &response, err
}

// DeleteSegment deletes a segment.
func (c *Client) DeleteSegment(ctx context.Context, id string) error {
	return c.makeRequest(ctx, http.MethodDelete, "segments/"+id, nil, nil, nil)
}
