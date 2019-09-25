// Package stream_chat provides chat via stream api
//nolint: golint
package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Devices(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	devices := []*Device{
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderFirebase},
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderAPNS},
	}

	for _, dev := range devices {
		mustNoError(t, c.AddDevice(dev), "add device")
		defer func(dev *Device) {
			mustNoError(t, c.DeleteDevice(user.ID, dev.ID), "delete device")
		}(dev)

		resp, err := c.GetDevices(user.ID)
		mustNoError(t, err, "get devices")

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
