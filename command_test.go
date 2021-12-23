package stream_chat // nolint: golint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareCommand(t *testing.T, c *Client) *Command {
	cmd := &Command{
		Name:        randomString(10),
		Description: "test command",
	}

	resp, err := c.CreateCommand(context.Background(), cmd)
	require.NoError(t, err, "create command")

	return resp.Command
}

func TestClient_GetCommand(t *testing.T) {
	c := initClient(t)

	cmd := prepareCommand(t, c)
	defer func() {
		_, _ = c.DeleteCommand(context.Background(), cmd.Name)
	}()

	resp, err := c.GetCommand(context.Background(), cmd.Name)
	require.NoError(t, err, "get command")

	assert.Equal(t, cmd.Name, resp.Command.Name)
	assert.Equal(t, cmd.Description, resp.Command.Description)
}

func TestClient_ListCommands(t *testing.T) {
	c := initClient(t)

	cmd := prepareCommand(t, c)
	defer func() {
		_, _ = c.DeleteCommand(context.Background(), cmd.Name)
	}()

	resp, err := c.ListCommands(context.Background())
	require.NoError(t, err, "list commands")

	assert.Contains(t, resp.Commands, cmd)
}

func TestClient_UpdateCommand(t *testing.T) {
	c := initClient(t)

	cmd := prepareCommand(t, c)
	defer func() {
		_, _ = c.DeleteCommand(context.Background(), cmd.Name)
	}()

	resp, err := c.UpdateCommand(context.Background(), cmd.Name, map[string]interface{}{
		"description": "new description",
	})
	require.NoError(t, err, "update command")

	assert.Equal(t, cmd.Name, resp.Command.Name)
	assert.Equal(t, "new description", resp.Command.Description)
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

	_, _ = client.CreateCommand(context.Background(), newCommand)
}

func ExampleClient_ListCommands() {
	client := &Client{}
	_, _ = client.ListCommands(context.Background())
}

func ExampleClient_GetCommand() {
	client := &Client{}
	_, _ = client.GetCommand(context.Background(), "my-command")
}

func ExampleClient_UpdateCommand() {
	client := &Client{}

	_, _ = client.UpdateCommand(context.Background(), "my-command", map[string]interface{}{
		"description": "updated description",
	})
}

func ExampleClient_DeleteCommand() {
	client := &Client{}

	_, _ = client.DeleteCommand(context.Background(), "my-command")
}
