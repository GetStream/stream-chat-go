package stream_chat //nolint: golint

import (
	"testing"
)

func TestClient_GetApp(t *testing.T) {
	c := initClient(t)
	_, err := c.GetAppConfig()
	mustNoError(t, err)
}

func TestClient_UpdateAppSettings(t *testing.T) {
	c := initClient(t)

	settings := NewAppSettings().
		SetDisableAuth(true).
		SetDisablePermissions(true)

	err := c.UpdateAppSettings(settings)
	mustNoError(t, err)
}
