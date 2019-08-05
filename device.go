package stream_chat

import "net/http"

const (
	PushProviderAPNS     = pushProvider("apns")
	PushPrivoderFirebase = pushProvider("firebase")
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
	//The user ID for this device.
	UserID string
	//The device ID.
	DeviceID string
	//The push provider for this device. One of constants PushProvider*
	Provider pushProvider
}

func (d *Device) fromHash(hash map[string]string) {
	d.UserID = hash["user_id"]
	d.DeviceID = hash["id"]
	d.Provider = pushProvider(hash["provider"])
}

func (d *Device) toHash() map[string]interface{} {
	return map[string]interface{}{
		"user_id":  d.UserID,
		"id":       d.DeviceID,
		"provider": d.Provider,
	}
}

// Get list of devices for user
func (c *client) GetDevices(userId string) (devices []Device, err error) {
	params := map[string][]string{
		"user_id": {userId},
	}

	var resp map[string][]map[string]string
	err = c.makeRequest(http.MethodGet, "devices", params, nil, &resp)

	if err != nil {
		return nil, err
	}

	if devs := resp["devices"]; devs != nil {
		devices = make([]Device, len(devs))
		for i := range devices {
			devices[i].fromHash(devs[i])
		}
	}

	return devices, err
}

// Add device to a user. Provider should be one of PushProvider* constant
func (c *client) AddDevice(device Device) error {
	return c.makeRequest(http.MethodPost, "devices", nil, device, nil)
}

// Delete a device for a user
func (c *client) DeleteDevice(userId string, deviceID string) error {
	params := map[string][]string{
		"id":      {deviceID},
		"user_id": {userId},
	}

	return c.makeRequest(http.MethodDelete, "devices", params, nil, nil)
}
