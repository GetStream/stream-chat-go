package stream_chat

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

const (
	PushProviderAPNS     = pushProvider("apn")
	PushProviderFirebase = pushProvider("firebase")
)

type pushProvider = string

type Device struct {
	ID           string       `json:"id"`            // The device ID.
	UserID       string       `json:"user_id"`       // The user ID for this device.
	PushProvider pushProvider `json:"push_provider"` // The push provider for this device. One of constants PushProvider*
}

type DevicesResponse struct {
	Devices []*Device `json:"devices"`
	Response
}

// GetDevices retrieves the list of devices for user.
func (c *Client) GetDevices(ctx context.Context, userID string) (*DevicesResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	params := url.Values{}
	params.Set("user_id", userID)

	var resp DevicesResponse

	err := c.makeRequest(ctx, http.MethodGet, "devices", params, nil, &resp)
	return &resp, err
}

// AddDevice adds new device.
func (c *Client) AddDevice(ctx context.Context, device *Device) (*Response, error) {
	switch {
	case device == nil:
		return nil, errors.New("device is nil")
	case device.ID == "":
		return nil, errors.New("device ID is empty")
	case device.UserID == "":
		return nil, errors.New("device user ID is empty")
	case device.PushProvider == "":
		return nil, errors.New("device push provider is empty")
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "devices", nil, device, &resp)
	return &resp, err
}

// DeleteDevice deletes a device from the user.
func (c *Client) DeleteDevice(ctx context.Context, userID, deviceID string) (*Response, error) {
	switch {
	case userID == "":
		return nil, errors.New("user ID is empty")
	case deviceID == "":
		return nil, errors.New("device ID is empty")
	}

	params := url.Values{}
	params.Set("id", deviceID)
	params.Set("user_id", userID)

	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, "devices", params, nil, &resp)
	return &resp, err
}
