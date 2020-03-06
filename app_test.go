package stream

import (
	"log"
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

// See https://getstream.io/chat/docs/app_settings_auth/ for
// more details.
func ExampleClient_UpdateAppSettings_disable_auth() {
	client, err := NewClient("XXXXXXXXXXXX", []byte("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// disable auth checks, allows dev token usage
	settings := NewAppSettings().SetDisableAuth(true)
	err = client.UpdateAppSettings(settings)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// re-enable auth checks
	err = client.UpdateAppSettings(NewAppSettings().SetDisableAuth(false))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}

func ExampleClient_UpdateAppSettings_disable_permission() {
	client, err := NewClient("XXXX", []byte("XXXX"))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// disable permission checkse
	settings := NewAppSettings().SetDisablePermissions(true)
	err = client.UpdateAppSettings(settings)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// re-enable permission checks
	err = client.UpdateAppSettings(NewAppSettings().SetDisablePermissions(false))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}
