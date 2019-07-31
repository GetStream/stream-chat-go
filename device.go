package stream_chat

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
	GetDevices(userId UserID) ([]Device, error)
	// Add device to a user. Provider should be one of PushProvider* constant
	AddDevice(userId UserID, deviceID DeviceID, provider pushProvider) error
	// Delete a device for a user
	DeleteDevice(userId UserID, deviceID DeviceID) error
}

type Device struct {
}

// Get list of devices for user
func (c *client) GetDevices(userId UserID) ([]Device, error) {
	panic("implement me")
}

// Add device to a user. Provider should be one of PushProvider* constant
func (c *client) AddDevice(userId UserID, deviceID DeviceID, provider pushProvider) error {
	panic("implement me")
}

// Delete a device for a user
func (c *client) DeleteDevice(userId UserID, deviceID DeviceID) error {
	panic("implement me")
}
