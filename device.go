package stream_chat

import (
	"net/http"
)

const (
	PushProviderAPNS     = pushProvider("apn")
	PushProviderFirebase = pushProvider("firebase")
)

type (
	DeviceID string

	pushProvider string
)

type DeviceAPI interface {
	// Get list of devices for user
	GetDevices(userId string) ([]Device, error)
	// Add device to a user. Provider should be one of PushProvider* constant
	AddDevice(userId string, deviceID DeviceID, provider pushProvider) error
	// Delete a device for a user
	DeleteDevice(userId string, deviceID DeviceID) error
}

type Device struct {
	//The device ID.
	ID string `json:"id"`
	//The user ID for this device.
	UserID string `json:"user_id"`
	//The push provider for this device. One of constants PushProvider*
	PushProvider pushProvider `json:"push_provider"`
}

func (d *Device) fromHash(hash map[string]string) {
	d.UserID = hash["user_id"]
	d.ID = hash["id"]
	d.PushProvider = pushProvider(hash["provider"])
}

func (d *Device) toHash() map[string]interface{} {
	return map[string]interface{}{
		"user_id":  d.UserID,
		"id":       d.ID,
		"provider": d.PushProvider,
	}
}

// Get list of devices for user
func (c *Client) GetDevices(userId string) (devices []Device, err error) {
	params := map[string][]string{
		"user_id": {userId},
	}

	var resp struct {
		Devices []Device `json:"devices"`
	}

	err = c.makeRequest(http.MethodGet, "devices", params, nil, &resp)

	return resp.Devices, err
}

// Add device to a user. Provider should be one of PushProvider* constant
func (c *Client) AddDevice(device Device) error {
	return c.makeRequest(http.MethodPost, "devices", nil, device, nil)
}

// Delete a device for a user
func (c *Client) DeleteDevice(userId string, deviceID string) error {
	params := map[string][]string{
		"id":      {deviceID},
		"user_id": {userId},
	}

	return c.makeRequest(http.MethodDelete, "devices", params, nil, nil)
}
