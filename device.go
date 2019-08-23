package stream_chat

import (
	"errors"
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

	err = c.Get("devices", params, &resp)

	return resp.Devices, err
}

// AddDevice adds new device.
func (c *Client) AddDevice(device *Device) error {
	switch {
	case device == nil:
		return errors.New("device is nil")
	case device.ID == "":
		return errors.New("device ID is empty")
	case device.UserID == "":
		return errors.New("device user ID is empty")
	case device.PushProvider == "":
		return errors.New("device push provider is empty")
	}

	return c.Post("devices", nil, device, nil)
}

// Delete a device from the user
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

	return c.Delete("devices", params, nil)
}
