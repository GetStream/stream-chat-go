package stream_chat // nolint: golint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Devices(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)

	devices := []*Device{
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderFirebase},
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderAPNS},
	}

	for _, dev := range devices {
		_, err := c.AddDevice(context.Background(), dev)
		require.NoError(t, err, "add device")
		defer func(dev *Device) {
			_, err := c.DeleteDevice(context.Background(), user.ID, dev.ID)
			require.NoError(t, err, "delete device")
		}(dev)

		resp, err := c.GetDevices(context.Background(), user.ID)
		require.NoError(t, err, "get devices")

		assert.True(t, deviceIDExists(resp.Devices, dev.ID), "device with ID %s was created", dev.ID)
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

	_, _ = client.AddDevice(context.Background(), &Device{
		ID:           "2ffca4ad6599adc9b5202d15a5286d33c19547d472cd09de44219cda5ac30207",
		UserID:       "elon",
		PushProvider: PushProviderAPNS,
	})
}

func ExampleClient_DeleteDevice() {
	client, _ := NewClient("XXXX", "XXXX")

	deviceID := "2ffca4ad6599adc9b5202d15a5286d33c19547d472cd09de44219cda5ac30207"
	userID := "elon"
	_, _ = client.DeleteDevice(context.Background(), userID, deviceID)
}
