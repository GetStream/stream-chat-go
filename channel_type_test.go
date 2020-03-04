package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareChannelType(t *testing.T, c *Client) *ChannelType {
	ct := NewChannelType(randomString(10))

	ct, err := c.CreateChannelType(ct)
	mustNoError(t, err, "create channel type")

	return ct
}

func TestClient_GetChannelType(t *testing.T) {
	c := initClient(t)

	ct := prepareChannelType(t, c)
	defer func() {
		mustNoError(t, c.DeleteChannelType(ct.Name), "delete channel type")
	}()

	got, err := c.GetChannelType(ct.Name)
	mustNoError(t, err, "get channel type")

	assert.Equal(t, ct.Name, got.Name)
	assert.Equal(t, len(ct.Commands), len(got.Commands))
	assert.Equal(t, ct.Permissions, got.Permissions)
}

func TestClient_ListChannelTypes(t *testing.T) {
	c := initClient(t)

	ct := prepareChannelType(t, c)
	defer func() {
		mustNoError(t, c.DeleteChannelType(ct.Name), "delete channel type")
	}()

	got, err := c.ListChannelTypes()
	mustNoError(t, err, "list channel types")

	assert.Contains(t, got, ct.Name)
}

// See https://getstream.io/chat/docs/channel_features/ for more details.
func ExampleClient_CreateChannelType() {
	client := &Client{}

	newChannelType := &ChannelType{
		// Copy the default settings.
		ChannelConfig: DefaultChannelConfig,
	}

	newChannelType.Name = "public"
	newChannelType.Mutes = false
	newChannelType.Reactions = false
	newChannelType.Permissions = append(newChannelType.Permissions,
		&Permission{
			Name:      "Allow reads for all",
			Priority:  999,
			Resources: []string{"ReadChannel", "CreateMessage"},
			Action:    "Allow",
		},
		&Permission{
			Name:      "Deny all",
			Priority:  1,
			Resources: []string{"*"},
			Action:    "Deny",
		},
	)

	client.CreateChannelType(newChannelType)
}

func ExampleClient_ListChannelTypes() {
	client := &Client{}
	client.ListChannelTypes()
}

func ExampleClient_GetChannelType() {
	client := &Client{}
	client.GetChannelType("public")
}

func ExampleClient_UpdateChannelType() {
	client := &Client{}

	client.UpdateChannelType("public", map[string]interface{}{
		"permissions": []map[string]interface{}{
			map[string]interface{}{
				"name":      "Allow reads for all",
				"priority":  999,
				"resources": []string{"ReadChannel", "CreateMessage"},
				"role":      "*",
				"action":    "Allow",
			},
			map[string]interface{}{
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

	client.UpdateChannelType("public", map[string]interface{}{
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

	client.UpdateChannelType(
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

	client.UpdateChannelType(
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

	client.DeleteChannelType("public")
}
