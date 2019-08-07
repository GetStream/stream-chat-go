package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Devices(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	devices := []Device{
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderFirebase},
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderAPNS},
	}

	for _, dev := range devices {
		err := c.AddDevice(dev)
		defer c.DeleteDevice(user.ID, dev.ID)
		mustNoError(t, err)

		resp, err := c.GetDevices(user.ID)
		mustNoError(t, err)

		assert.True(t, deviceIDExists(resp, dev.ID), "device with ID %s was created", dev.ID)
	}
}

func deviceIDExists(dev []Device, id string) bool {
	for _, d := range dev {
		if d.ID == id {
			return true
		}
	}
	return false
}
