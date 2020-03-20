package stream_chat //nolint:golint

import (
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

type devicesResponse struct {
	Devices []*Device `json:"devices"`
}

// GetDevices retrieves the list of devices for user
func (c *Client) GetDevices(userID string) (devices []*Device, err error) {
	if userID == "" {
		return nil, ErrorMissingUserID
	}

	params := url.Values{}
	params.Set("user_id", userID)

	var resp devicesResponse

	err = c.makeRequest(http.MethodGet, "devices", params, nil, &resp)

	return resp.Devices, err
}

// AddDevice adds new device.
func (c *Client) AddDevice(device *Device) error {
	switch {
	case device == nil:
		return errors.New("device is nil")
	case device.ID == "":
		return ErrorMissingDeviceID
	case device.UserID == "":
		return errors.New("device user ID is empty")
	case device.PushProvider == "":
		return errors.New("device push provider is empty")
	}

	return c.makeRequest(http.MethodPost, "devices", nil, device, nil)
}

// DeleteDevice deletes a device from the user
func (c *Client) DeleteDevice(userID, deviceID string) error {
	switch {
	case userID == "":
		return ErrorMissingUserID
	case deviceID == "":
		return ErrorMissingDeviceID
	}

	options := []Option{
		NewOption("id", deviceID),
		NewOption("user_id", userID),
	}

	return c.makeRequestWithOptions(http.MethodDelete, "devices", options, nil, nil)
}
