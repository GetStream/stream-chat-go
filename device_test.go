package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_Devices(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	user := randomUser(t, c)

	devices := []*Device{
		{UserID: user.ID, ID: "xxxx", PushProvider: PushProviderFirebase},
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderAPNS},
	}

	for _, dev := range devices {
		_, err := c.AddDevice(ctx, dev)
		require.NoError(t, err, "add device")

		resp, err := c.GetDevices(ctx, user.ID)
		require.NoError(t, err, "get devices")

		require.True(t, deviceIDExists(resp.Devices, dev.ID), "device with ID %s was created", dev.ID)
		_, err = c.DeleteDevice(ctx, user.ID, dev.ID)
		require.NoError(t, err, "delete device")
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
	ctx := context.Background()

	_, _ = client.AddDevice(ctx, &Device{
		ID:           "2ffca4ad6599adc9b5202d15a5286d33c19547d472cd09de44219cda5ac30207",
		UserID:       "elon",
		PushProvider: PushProviderAPNS,
	})
}

func ExampleClient_DeleteDevice() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	deviceID := "2ffca4ad6599adc9b5202d15a5286d33c19547d472cd09de44219cda5ac30207"
	userID := "elon"
	_, _ = client.DeleteDevice(ctx, userID, deviceID)
}
