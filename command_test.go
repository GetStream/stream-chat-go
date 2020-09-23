package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareCommand(t *testing.T, c *Client) *Command {
	cmd := &Command{
		Name:        randomString(10),
		Description: "test command",
	}

	cmd, err := c.CreateCommand(cmd)
	require.NoError(t, err, "create command")

	return cmd
}

func TestClient_GetCommand(t *testing.T) {
	c := initClient(t)

	cmd := prepareCommand(t, c)
	defer func() {
		_ = c.DeleteCommand(cmd.Name)
	}()

	got, err := c.GetCommand(cmd.Name)
	require.NoError(t, err, "get command")

	assert.Equal(t, cmd.Name, got.Name)
	assert.Equal(t, cmd.Description, got.Description)
}

func TestClient_ListCommands(t *testing.T) {
	c := initClient(t)

	cmd := prepareCommand(t, c)
	defer func() {
		_ = c.DeleteCommand(cmd.Name)
	}()

	got, err := c.ListCommands()
	require.NoError(t, err, "list commands")

	assert.Contains(t, got, cmd)
}

func TestClient_UpdateCommand(t *testing.T) {
	c := initClient(t)

	cmd := prepareCommand(t, c)
	defer func() {
		_ = c.DeleteCommand(cmd.Name)
	}()

	got, err := c.UpdateCommand(cmd.Name, map[string]interface{}{
		"description": "new description",
	})
	require.NoError(t, err, "update command")

	assert.Equal(t, cmd.Name, got.Name)
	assert.Equal(t, "new description", got.Description)
}

// See https://getstream.io/chat/docs/custom_commands/ for more details.
func ExampleClient_CreateCommand() {
	client := &Client{}

	newCommand := &Command{
		Name:        "my-command",
		Description: "my command",
		Args:        "[@username]",
		Set:         "custom_cmd_set",
	}

	_, _ = client.CreateCommand(newCommand)
}

func ExampleClient_ListCommands() {
	client := &Client{}
	_, _ = client.ListCommands()
}

func ExampleClient_GetCommand() {
	client := &Client{}
	_, _ = client.GetCommand("my-command")
}

func ExampleClient_UpdateCommand() {
	client := &Client{}

	_, _ = client.UpdateCommand("my-command", map[string]interface{}{
		"description": "updated description",
	})
}

func ExampleClient_DeleteCommand() {
	client := &Client{}

	_ = client.DeleteCommand("my-command")
}
