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

func TestClient_UpdateAppSettingsWithFileUploadConfig(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Save original settings for cleanup
	originalSettings, err := c.GetAppSettings(ctx)
	require.NoError(t, err)

	// Cleanup: restore original settings after test
	defer func() {
		cleanupSettings := NewAppSettings()
		cleanupSettings.FileUploadConfig = originalSettings.App.FileUploadConfig
		_, err := c.UpdateAppSettings(ctx, cleanupSettings)
		require.NoError(t, err)
	}()

	// Test updating app settings with file upload config including size limit
	sizeLimit := 10485760 // 10MB
	fileUploadConfig := &FileUploadConfig{
		AllowedFileExtensions: []string{".pdf", ".doc", ".txt"},
		AllowedMimeTypes:      []string{"application/pdf", "text/plain"},
		SizeLimit:             &sizeLimit,
	}

	settings := NewAppSettings()
	settings.FileUploadConfig = fileUploadConfig

	_, err = c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)
}

func TestClient_GetAppSettingsWithFileUploadConfig(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Save original settings for cleanup
	originalSettings, err := c.GetAppSettings(ctx)
	require.NoError(t, err)

	// Cleanup: restore original settings after test
	defer func() {
		cleanupSettings := NewAppSettings()
		cleanupSettings.FileUploadConfig = originalSettings.App.FileUploadConfig
		_, err := c.UpdateAppSettings(ctx, cleanupSettings)
		require.NoError(t, err)
	}()

	// First, set up file upload config with size limit
	sizeLimit := 5242880 // 5MB
	fileUploadConfig := &FileUploadConfig{
		AllowedFileExtensions: []string{".jpg", ".png", ".txt"},
		AllowedMimeTypes:      []string{"image/jpeg", "image/png", "text/plain"},
		SizeLimit:             &sizeLimit,
	}

	settings := NewAppSettings()
	settings.FileUploadConfig = fileUploadConfig
	_, err = c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)

	resp, err := c.GetAppSettings(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp.App.FileUploadConfig)

	// Verify all fields are present and correct
	require.Equal(t, []string{".jpg", ".png", ".txt"}, resp.App.FileUploadConfig.AllowedFileExtensions)
	require.Equal(t, []string{"image/jpeg", "image/png", "text/plain"}, resp.App.FileUploadConfig.AllowedMimeTypes)
	require.NotNil(t, resp.App.FileUploadConfig.SizeLimit)
	require.Equal(t, sizeLimit, *resp.App.FileUploadConfig.SizeLimit)
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

func TestClientUpdateEventHooks(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	t.Run("webhook event hooks", func(t *testing.T) {
		eventHooks := []EventHook{
			{
				HookType:   WebhookHook,
				Enabled:    true,
				EventTypes: []string{"message.new"},
				WebhookURL: "http://google.com",
			},
			{
				HookType:   WebhookHook,
				Enabled:    true,
				EventTypes: []string{"message.new"},
				WebhookURL: "http://google.nl",
				Product:    ProductChat,
			},
		}

		settings := NewAppSettings().SetEventHooks(eventHooks)
		_, err := c.UpdateAppSettings(ctx, settings)
		require.NoError(t, err)
	})

	t.Run("pending message async moderation config", func(t *testing.T) {
		eventHooks := []EventHook{
			{
				HookType:   PendingMessage,
				Enabled:    true,
				WebhookURL: "http://google.com",
				TimeoutMs:  10000,
				Callback: &Callback{
					Mode: CallbackModeREST,
				},
			},
		}

		settings := NewAppSettings().SetEventHooks(eventHooks)
		_, err := c.UpdateAppSettings(ctx, settings)
		require.NoError(t, err)
	})

	t.Run("SQS event hooks", func(t *testing.T) {
		eventHooks := []EventHook{
			{
				HookType:    SQSHook,
				Enabled:     true,
				EventTypes:  []string{"message.new"},
				SQSQueueURL: "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue",
				SQSRegion:   "us-east-1",
				SQSAuthType: "keys",
				SQSKey:      "some key",
				SQSSecret:   "some secret",
			},
		}

		settings := NewAppSettings().SetEventHooks(eventHooks)
		_, err := c.UpdateAppSettings(ctx, settings)
		require.NoError(t, err)
	})

	t.Run("SNS event hooks", func(t *testing.T) {
		eventHooks := []EventHook{
			{
				HookType:    SNSHook,
				Enabled:     true,
				EventTypes:  []string{"message.new"},
				SNSTopicARN: "arn:aws:sns:us-east-1:123456789012:my-topic",
				SNSRegion:   "us-east-1",
				SNSAuthType: "keys",
				SNSKey:      "some key",
				SNSSecret:   "some secret",
			},
		}

		settings := NewAppSettings().SetEventHooks(eventHooks)
		_, err := c.UpdateAppSettings(ctx, settings)
		require.NoError(t, err)
	})

	t.Run("clear event hooks", func(t *testing.T) {
		// First set some event hooks
		eventHooks := []EventHook{
			{
				HookType:   WebhookHook,
				Enabled:    true,
				EventTypes: []string{"message.new"},
				WebhookURL: "http://example.com",
			},
		}
		settings := NewAppSettings().SetEventHooks(eventHooks)
		_, err := c.UpdateAppSettings(ctx, settings)
		require.NoError(t, err)

		// Now clear them with empty array
		emptySettings := NewAppSettings().SetEventHooks([]EventHook{})
		_, err = c.UpdateAppSettings(ctx, emptySettings)
		require.NoError(t, err)

		// Verify they are cleared
		_, err = c.GetAppSettings(ctx)
		require.NoError(t, err)
	})
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
