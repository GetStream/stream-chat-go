package stream_chat //nolint: golint

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_GetApp(t *testing.T) {
	c := initClient(t)
	_, err := c.GetAppConfig(context.Background())
	require.NoError(t, err)
}

func TestClient_UpdateAppSettings(t *testing.T) {
	c := initClient(t)

	settings := NewAppSettings().
		SetDisableAuth(true).
		SetDisablePermissions(true)

	_, err := c.UpdateAppSettings(context.Background(), settings)
	require.NoError(t, err)
}

// See https://getstream.io/chat/docs/app_settings_auth/ for
// more details.
func ExampleClient_UpdateAppSettings_disable_auth() {
	client, err := NewClient("XXXXXXXXXXXX", "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// disable auth checks, allows dev token usage
	settings := NewAppSettings().SetDisableAuth(true)
	_, err = client.UpdateAppSettings(context.Background(), settings)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// re-enable auth checks
	_, err = client.UpdateAppSettings(context.Background(), NewAppSettings().SetDisableAuth(false))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}

func ExampleClient_UpdateAppSettings_disable_permission() {
	client, err := NewClient("XXXX", "XXXX")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// disable permission checkse
	settings := NewAppSettings().SetDisablePermissions(true)
	_, err = client.UpdateAppSettings(context.Background(), settings)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// re-enable permission checks
	_, err = client.UpdateAppSettings(context.Background(), NewAppSettings().SetDisablePermissions(false))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}
