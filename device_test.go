package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Devices(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	resp, err := c.GetDevices(user.ID)
	mustNoError(t, err)

	if !assert.Len(t, resp, 0, "devices are not empty in response") {
		t.FailNow()
	}

	devices := []Device{
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderFirebase},
		{UserID: user.ID, ID: randomString(12), PushProvider: PushProviderAPNS},
	}

	for _, dev := range devices {
		err = c.AddDevice(dev)
		mustNoError(t, err)

		resp, err = c.GetDevices(user.ID)
		mustNoError(t, err)

		if !assert.Len(t, resp, 1) {
			t.FailNow()
		}

		err = c.DeleteDevice(user.ID, dev.ID)
		mustNoError(t, err)
	}
}
