package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Devices(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	devices := []*Device{
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderFirebase},
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderAPNS},
	}

	for _, dev := range devices {
		require.NoError(t, c.AddDevice(dev), "add device")
		defer func(dev *Device) {
			require.NoError(t, c.DeleteDevice(user.ID, dev.ID), "delete device")
		}(dev)

		resp, err := c.GetDevices(user.ID)
		require.NoError(t, err, "get devices")

		assert.True(t, deviceIDExists(resp, dev.ID), "device with ID %s was created", dev.ID)
	}
}

func deviceIDExists(dev []*Device, id string) bool {
	for _, d := range dev {
		if d.ID == id {
			return true
		}
	}
	return false
}

func ExampleClient_AddDevice() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.AddDevice(&Device{
		ID:           "2ffca4ad6599adc9b5202d15a5286d33c19547d472cd09de44219cda5ac30207",
		UserID:       "elon",
		PushProvider: PushProviderAPNS,
	})
}

func ExampleClient_DeleteDevice() {
	client, _ := NewClient("XXXX", "XXXX")

	deviceID := "2ffca4ad6599adc9b5202d15a5286d33c19547d472cd09de44219cda5ac30207"
	userID := "elon"
	_ = client.DeleteDevice(userID, deviceID)
}
