package stream_chat

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_GetApp(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	_, err := c.GetAppSettings(ctx)
	require.NoError(t, err)
}

func TestClient_UpdateAppSettings(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	settings := NewAppSettings().
		SetDisableAuth(true).
		SetDisablePermissions(true)

	_, err := c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)
}

func TestClient_CheckAsyncModeConfig(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	settings := NewAppSettings().
		SetAsyncModerationConfig(
			AsyncModerationConfiguration{
				Callback: &AsyncModerationCallback{
					Mode:      "CALLBACK_MODE_REST",
					ServerURL: "https://example.com/gosdk",
				},
				Timeout: 10000,
			},
		)

	_, err := c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)
}

func TestClient_UpdateAppSettingsClearing(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	sqsURL := "https://example.com"
	sqsKey := "some key"
	sqsSecret := "some secret"

	settings := NewAppSettings()
	settings.SqsURL = &sqsURL
	settings.SqsKey = &sqsKey
	settings.SqsSecret = &sqsSecret

	_, err := c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)

	sqsURL = ""
	settings.SqsURL = &sqsURL
	_, err = c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)

	s, err := c.GetAppSettings(ctx)
	require.NoError(t, err)
	require.Equal(t, *settings.SqsURL, *s.App.SqsURL)
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

func TestClient_CheckSns(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	req := &CheckSNSRequest{SnsTopicARN: "arn:aws:sns:us-east-1:123456789012:sns-topic", SnsKey: "key", SnsSecret: "secret"}
	resp, err := c.CheckSns(ctx, req)

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
	ctx := context.Background()

	// disable auth checks, allows dev token usage
	settings := NewAppSettings().SetDisableAuth(true)
	_, err = client.UpdateAppSettings(ctx, settings)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// re-enable auth checks
	_, err = client.UpdateAppSettings(ctx, NewAppSettings().SetDisableAuth(false))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}

func ExampleClient_UpdateAppSettings_disable_permission() {
	client, err := NewClient("XXXX", "XXXX")
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	ctx := context.Background()

	// disable permission checkse
	settings := NewAppSettings().SetDisablePermissions(true)
	_, err = client.UpdateAppSettings(ctx, settings)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	// re-enable permission checks
	_, err = client.UpdateAppSettings(ctx, NewAppSettings().SetDisablePermissions(false))
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}
