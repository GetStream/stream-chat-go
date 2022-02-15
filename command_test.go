package stream_chat

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
	ctx := context.Background()

	resp, err := c.CreateCommand(ctx, cmd)
	require.NoError(t, err, "create command")

	t.Cleanup(func() {
		_, _ = c.DeleteCommand(ctx, cmd.Name)
	})

	return resp.Command
}

func TestClient_GetCommand(t *testing.T) {
	c := initClient(t)
	cmd := prepareCommand(t, c)
	ctx := context.Background()

	resp, err := c.GetCommand(ctx, cmd.Name)
	require.NoError(t, err, "get command")

	assert.Equal(t, cmd.Name, resp.Command.Name)
	assert.Equal(t, cmd.Description, resp.Command.Description)
}

func TestClient_ListCommands(t *testing.T) {
	c := initClient(t)
	cmd := prepareCommand(t, c)
	ctx := context.Background()

	resp, err := c.ListCommands(ctx)
	require.NoError(t, err, "list commands")

	assert.Contains(t, resp.Commands, cmd)
}

func TestClient_UpdateCommand(t *testing.T) {
	c := initClient(t)
	cmd := prepareCommand(t, c)
	ctx := context.Background()

	update := Command{Description: "new description"}
	resp, err := c.UpdateCommand(ctx, cmd.Name, &update)
	require.NoError(t, err, "update command")

	assert.Equal(t, cmd.Name, resp.Command.Name)
	assert.Equal(t, "new description", resp.Command.Description)
}

// See https://getstream.io/chat/docs/custom_commands/ for more details.
func ExampleClient_CreateCommand() {
	client := &Client{}
	ctx := context.Background()

	newCommand := &Command{
		Name:        "my-command",
		Description: "my command",
		Args:        "[@username]",
		Set:         "custom_cmd_set",
	}

	_, _ = client.CreateCommand(ctx, newCommand)
}

func ExampleClient_ListCommands() {
	client := &Client{}
	ctx := context.Background()
	_, _ = client.ListCommands(ctx)
}

func ExampleClient_GetCommand() {
	client := &Client{}
	ctx := context.Background()
	_, _ = client.GetCommand(ctx, "my-command")
}

func ExampleClient_UpdateCommand() {
	client := &Client{}
	ctx := context.Background()

	update := Command{Description: "updated description"}
	_, _ = client.UpdateCommand(ctx, "my-command", &update)
}

func ExampleClient_DeleteCommand() {
	client := &Client{}
	ctx := context.Background()

	_, _ = client.DeleteCommand(ctx, "my-command")
}
