package stream_chat

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

type SharedLocation struct {
	MessageID  string `json:"message_id"`
	ChannelCID string `json:"channel_cid"`
	UserID     string `json:"user_id"`

	Latitude          *float64   `json:"latitude,omitempty"`
	Longitude         *float64   `json:"longitude,omitempty"`
	CreatedByDeviceID string     `json:"created_by_device_id"`
	EndAt             *time.Time `json:"end_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ActiveLiveLocationsResponse struct {
	ActiveLiveLocations []*SharedLocation `json:"active_live_locations"`
	Response
}

type SharedLocationResponse struct {
	SharedLocation

	Message *MessageResponse `json:"message,omitempty"`
	Channel *Channel         `json:"channel,omitempty"`
}

// GetUserActiveLocations returns all active live locations for a user
func (c *Client) GetUserActiveLocations(ctx context.Context, userID string) (*ActiveLiveLocationsResponse, error) {
	path := path.Join("users", "live_locations")
	var resp ActiveLiveLocationsResponse
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	params := url.Values{}
	params.Set("user_id", userID)

	err := c.makeRequest(ctx, http.MethodGet, path, params, nil, &resp)
	return &resp, err
}

// UpdateUserActiveLocation updates a location
func (c *Client) UpdateUserActiveLocation(ctx context.Context, userID string, location *SharedLocation) (*SharedLocationResponse, error) {
	path := path.Join("users", "live_locations")
	var resp SharedLocationResponse

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	params := url.Values{}
	params.Set("user_id", userID)

	err := c.makeRequest(ctx, http.MethodPut, path, params, location, &resp)
	return &resp, err
}
