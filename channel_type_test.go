package stream_chat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareChannelType(t *testing.T, c *Client) *ChannelType {
	ct := NewChannelType(randomString(10))
	ctx := context.Background()

	resp, err := c.CreateChannelType(ctx, ct)
	require.NoError(t, err, "create channel type")
	time.Sleep(6 * time.Second)

	t.Cleanup(func() {
		for i := 0; i < 5; i++ {
			_, err = c.DeleteChannelType(ctx, ct.Name)
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
	})

	return resp.ChannelType
}

func TestClient_GetChannelType(t *testing.T) {
	c := initClient(t)
	ct := prepareChannelType(t, c)
	ctx := context.Background()

	resp, err := c.GetChannelType(ctx, ct.Name)
	require.NoError(t, err, "get channel type")

	assert.Equal(t, ct.Name, resp.ChannelType.Name)
	assert.Equal(t, len(ct.Commands), len(resp.ChannelType.Commands))
	assert.Equal(t, ct.Permissions, resp.ChannelType.Permissions)
	assert.NotEmpty(t, resp.Grants)
}

func TestClient_ListChannelTypes(t *testing.T) {
	c := initClient(t)
	ct := prepareChannelType(t, c)
	ctx := context.Background()

	resp, err := c.ListChannelTypes(ctx)
	require.NoError(t, err, "list channel types")

	assert.Contains(t, resp.ChannelTypes, ct.Name)
}

func TestClient_UpdateChannelTypePushNotifications(t *testing.T) {
	c := initClient(t)
	ct := prepareChannelType(t, c)
	ctx := context.Background()

	// default is on
	require.True(t, ct.PushNotifications)

	_, err := c.UpdateChannelType(ctx, ct.Name, map[string]interface{}{"push_notifications": false})
	require.NoError(t, err)

	resp, err := c.GetChannelType(ctx, ct.Name)
	require.NoError(t, err)
	require.False(t, resp.ChannelType.PushNotifications)
}

// See https://getstream.io/chat/docs/channel_features/ for more details.
func ExampleClient_CreateChannelType() {
	client := &Client{}
	ctx := context.Background()

	newChannelType := &ChannelType{
		// Copy the default settings.
		ChannelConfig: DefaultChannelConfig,
	}

	newChannelType.Name = "public"
	newChannelType.Mutes = false
	newChannelType.Reactions = false
	newChannelType.Permissions = append(newChannelType.Permissions,
		&ChannelTypePermission{
			Name:      "Allow reads for all",
			Priority:  999,
			Resources: []string{"ReadChannel", "CreateMessage"},
			Action:    "Allow",
		},
		&ChannelTypePermission{
			Name:      "Deny all",
			Priority:  1,
			Resources: []string{"*"},
			Action:    "Deny",
		},
	)

	_, _ = client.CreateChannelType(ctx, newChannelType)
}

func ExampleClient_ListChannelTypes() {
	client := &Client{}
	ctx := context.Background()
	_, _ = client.ListChannelTypes(ctx)
}

func ExampleClient_GetChannelType() {
	client := &Client{}
	ctx := context.Background()
	_, _ = client.GetChannelType(ctx, "public")
}

func ExampleClient_UpdateChannelType() {
	client := &Client{}
	ctx := context.Background()

	_, _ = client.UpdateChannelType(ctx, "public", map[string]interface{}{
		"permissions": []map[string]interface{}{
			{
				"name":      "Allow reads for all",
				"priority":  999,
				"resources": []string{"ReadChannel", "CreateMessage"},
				"role":      "*",
				"action":    "Allow",
			},
			{
				"name":      "Deny all",
				"priority":  1,
				"resources": []string{"*"},
				"role":      "*",
				"action":    "Deny",
			},
		},
		"replies":  false,
		"commands": []string{"all"},
	})
}

func ExampleClient_UpdateChannelType_bool() {
	client := &Client{}
	ctx := context.Background()

	_, _ = client.UpdateChannelType(ctx, "public", map[string]interface{}{
		"typing_events":  false,
		"read_events":    true,
		"connect_events": true,
		"search":         false,
		"reactions":      true,
		"replies":        false,
		"mutes":          true,
	})
}

func ExampleClient_UpdateChannelType_other() {
	client := &Client{}
	ctx := context.Background()

	_, _ = client.UpdateChannelType(ctx,
		"public",
		map[string]interface{}{
			"automod":            "disabled",
			"message_retention":  "7",
			"max_message_length": 140,
			"commands":           []interface{}{"ban", "unban"},
		},
	)
}

func ExampleClient_UpdateChannelType_permissions() {
	client := &Client{}
	ctx := context.Background()

	_, _ = client.UpdateChannelType(ctx,
		"public",
		map[string]interface{}{
			"permissions": []map[string]interface{}{
				{
					"name":      "Allow reads for all",
					"priority":  999,
					"resources": []string{"ReadChannel", "CreateMessage"},
					"role":      "*",
					"action":    "Allow",
				},
				{
					"name":      "Deny all",
					"priority":  1,
					"resources": []string{"*"},
					"action":    "Deny",
				},
			},
		},
	)
}

func ExampleClient_DeleteChannelType() {
	client := &Client{}
	ctx := context.Background()

	_, _ = client.DeleteChannelType(ctx, "public")
}
