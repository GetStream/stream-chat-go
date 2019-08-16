package stream_chat

import (
	"errors"
	"net/http"
)

const (
	PushProviderAPNS     = pushProvider("apn")
	PushProviderFirebase = pushProvider("firebase")
)

type pushProvider = string

type Device struct {
	ID           string       `json:"id"`            //The device ID.
	UserID       string       `json:"user_id"`       //The user ID for this device.
	PushProvider pushProvider `json:"push_provider"` //The push provider for this device. One of constants PushProvider*
}

type devicesResponse struct {
	Devices []*Device `json:"devices"`
}

// Get list of devices for user
func (c *Client) GetDevices(userId string) (devices []*Device, err error) {
	if userId == "" {
		return nil, errors.New("user ID is empty")
	}

	params := map[string][]string{
		"user_id": {userId},
	}

	var resp devicesResponse

	err = c.makeRequest(http.MethodGet, "devices", params, nil, &resp)

	return resp.Devices, err
}

// Add device to a user. Provider should be one of PushProvider* constant
func (c *Client) AddDevice(device *Device) error {
	return c.makeRequest(http.MethodPost, "devices", nil, device, nil)
}

// Delete a device for a user
func (c *Client) DeleteDevice(userID string, deviceID string) error {
	switch {
	case userID == "":
		return errors.New("user ID is empty")
	case deviceID == "":
		return errors.New("device ID is empty")
	}

	params := map[string][]string{
		"id":      {deviceID},
		"user_id": {userID},
	}

	return c.makeRequest(http.MethodDelete, "devices", params, nil, nil)
}
