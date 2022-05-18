package stream_chat

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CampaignStatusName string

const (
	StatusDraft      CampaignStatusName = "draft"
	StatusReady      CampaignStatusName = "ready"
	StatusStopped    CampaignStatusName = "stopped"
	StatusScheduled  CampaignStatusName = "scheduled"
	StatusCompleted  CampaignStatusName = "completed"
	StatusFailed     CampaignStatusName = "failed"
	StatusInProgress CampaignStatusName = "in_progress"
)

type SegmentFilter struct {
	Channel map[string]interface{} `json:"channel,omitempty"`
	User    map[string]interface{} `json:"user,omitempty"`
}

type CampaignData struct {
	SegmentID         string            `json:"segment_id,omitempty"`
	SenderID          string            `json:"sender_id,omitempty"`
	ChannelType       string            `json:"channel_type,omitempty"`
	Text              string            `json:"text,omitempty"`
	Defaults          map[string]string `json:"defaults,omitempty"`
	Filter            SegmentFilter     `json:"filter,omitempty"`
	Attachments       []*Attachment     `json:"attachments,omitempty"`
	PushNotifications bool              `json:"push_notifications,omitempty"`
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
}

type Campaign struct {
	CampaignData

	Status       CampaignStatusName `json:"status"`
	ScheduledFor time.Time          `json:"scheduled_for"`
	ScheduledAt  time.Time          `json:"scheduled_at"`
	CompletedAt  time.Time          `json:"completed_at"`
	FailedAt     time.Time          `json:"failed_at"`
	StoppedAt    time.Time          `json:"stopped_at"`
	ResumedAt    time.Time          `json:"resumed_at"`
	Progress     int                `json:"progress"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CampaignUpdateableFields struct {
	Name              string            `json:"name,omitempty"`
	Description       string            `json:"description,omitempty"`
	Text              string            `json:"text,omitempty"`
	SegmentID         string            `json:"segment_id,omitempty"`
	SenderID          string            `json:"sender_id,omitempty"`
	Defaults          map[string]string `json:"defaults"`
	Attachments       []*Attachment     `json:"attachments,omitempty"`
	PushNotifications *bool             `json:"push_notifications,omitempty"`
	ChannelType       string            `json:"channel_type,omitempty"`
}

type CampaignRequest struct {
	Campaign *CampaignData `json:"campaign"`
}

type UpdateCampaignRequest struct {
	Campaign *CampaignUpdateableFields `json:"campaign"`
}

type ListCampaignOptions struct {
	Limit  int
	Offset int
}

type GetCampaignResponse struct {
	Campaign Campaign `json:"campaign"`
	Response
}

type CreateCampaignResponse struct {
	Campaign Campaign `json:"campaign"`
	Response
}

type UpdateCampaignResponse struct {
	Campaign Campaign `json:"campaign"`
	Response
}

type ListCampaignsResponse struct {
	Campaigns []Campaign `json:"campaigns"`
	Response
}

// CreateCampaign creates a new campaign.
func (c *Client) CreateCampaign(ctx context.Context, req *CampaignRequest) (*CreateCampaignResponse, error) {
	var resp CreateCampaignResponse
	err := c.makeRequest(ctx, http.MethodPost, "campaigns", nil, req, &resp)
	return &resp, err
}

// GetCampaign retrieves a campaign by ID.
func (c *Client) GetCampaign(ctx context.Context, id string) (*GetCampaignResponse, error) {
	var resp GetCampaignResponse
	err := c.makeRequest(ctx, http.MethodGet, "campaigns/"+id, nil, nil, &resp)
	return &resp, err
}

// ListCampaigns retrieves a list of campaigns.
func (c *Client) ListCampaigns(ctx context.Context, opts *ListCampaignOptions) (*ListCampaignsResponse, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Limit > 0 {
			params.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Set("offset", strconv.Itoa(opts.Offset))
		}
	}

	var resp ListCampaignsResponse
	err := c.makeRequest(ctx, http.MethodGet, "campaigns", params, nil, &resp)
	return &resp, err
}

// UpdateCampaign updates a campaign by ID.
func (c *Client) UpdateCampaign(ctx context.Context, id string, req *UpdateCampaignRequest) (*UpdateCampaignResponse, error) {
	var resp UpdateCampaignResponse
	err := c.makeRequest(ctx, http.MethodPut, "campaigns/"+id, nil, req, &resp)
	return &resp, err
}

// DeleteCampaign deletes a campaign by ID.
func (c *Client) DeleteCampaign(ctx context.Context, id string) error {
	return c.makeRequest(ctx, http.MethodDelete, "campaigns/"+id, nil, nil, nil)
}
