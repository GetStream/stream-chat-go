package stream_chat

import (
	"net/http"

	"github.com/francoispqt/gojay"
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

type devices []Device

func (d *devices) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var dev Device
	if err := dec.Object(&dev); err != nil {
		return err
	}
	*d = append(*d, dev)
	return nil
}

type devicesResponse struct {
	Devices devices `json:"devices"`
}

func (d *devicesResponse) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	if key == "devices" {
		return dec.Array(&d.Devices)
	}
	return nil
}

func (d *devicesResponse) NKeys() int {
	return 1
}

// Get list of devices for user
func (c *Client) GetDevices(userId string) (devices []Device, err error) {
	params := map[string][]string{
		"user_id": {userId},
	}

	var resp devicesResponse

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
