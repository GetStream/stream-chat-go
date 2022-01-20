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

func TestClient_CheckSqs(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	req := &CheckSQSRequest{SqsURL: "https://foo.com/bar", SqsKey: "key", SqsSecret: "secret"}
	resp, err := c.CheckSqs(ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, resp.Error)
	require.Equal(t, "error", resp.Status)
	require.NotNil(t, resp.Data)
}

func TestClient_CheckPush(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	user := randomUser(t, c)
	msgResp, _ := ch.SendMessage(ctx, &Message{Text: "text"}, user.ID)
	skipDevices := true

	req := &CheckPushRequest{MessageID: msgResp.Message.ID, SkipDevices: &skipDevices, UserID: user.ID}
	resp, err := c.CheckPush(ctx, req)

	require.NoError(t, err)
	require.Equal(t, msgResp.Message.ID, resp.RenderedMessage["message_id"])
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
