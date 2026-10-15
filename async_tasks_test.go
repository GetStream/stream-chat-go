package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// requireTaskReadable checks the parts of an async task the SDK owns: the
// task-status endpoint can be called with the returned ID and answers with a
// known status. It does not wait for the task to finish, because when that
// happens depends on the backend task queue, not on the SDK.
func requireTaskReadable(t *testing.T, get func(context.Context, string) (*TaskResponse, error), taskID string) {
	t.Helper()

	task, err := get(context.Background(), taskID)
	require.NoError(t, err)
	require.Equal(t, taskID, task.TaskID)
	require.Contains(t, []TaskStatus{
		TaskStatusWaiting,
		TaskStatusPending,
		TaskStatusRunning,
		TaskStatusCompleted,
		TaskStatusFailed,
	}, task.Status)
}

func TestClient_DeleteChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	user := randomUser(t, c)
	msg := &Message{Text: "test message"}

	_, err := ch.SendMessage(ctx, msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without CIDs in parameter
	_, err = c.DeleteChannels(ctx, []string{}, true)
	require.Error(t, err)

	resp1, err := c.DeleteChannels(ctx, []string{ch.CID}, true)
	require.NoError(t, err)
	require.NotEmpty(t, resp1.TaskID)

	requireTaskReadable(t, c.GetTask, resp1.TaskID)
}

func TestClient_DeleteUsers(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	user := randomUser(t, c)

	msg := &Message{Text: "test message"}

	_, err := ch.SendMessage(ctx, msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without userIDs in parameter
	_, err = c.DeleteUsers(ctx, []string{}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.Error(t, err)

	resp1, err := c.DeleteUsers(ctx, []string{user.ID}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.TaskID)

	requireTaskReadable(t, c.GetTask, resp1.TaskID)
}

func TestClient_ExportChannels(t *testing.T) {
	c := initClient(t)
	ch1 := initChannel(t, c)
	ch2 := initChannel(t, c)
	ctx := context.Background()

	t.Run("Return error if there are 0 channels", func(t *testing.T) {
		_, err := c.ExportChannels(ctx, nil, nil)
		require.Error(t, err)
	})

	t.Run("Return error if exportable channel structs are incorrect", func(t *testing.T) {
		expChannels := []*ExportableChannel{
			{Type: "", ID: ch1.ID},
		}
		_, err := c.ExportChannels(ctx, expChannels, nil)
		require.Error(t, err)
	})

	t.Run("Export channels with no error", func(t *testing.T) {
		expChannels := []*ExportableChannel{
			{Type: ch1.Type, ID: ch1.ID},
			{Type: ch2.Type, ID: ch2.ID},
		}

		resp1, err := c.ExportChannels(ctx, expChannels, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp1.TaskID)

		requireTaskReadable(t, c.GetExportChannelsTask, resp1.TaskID)
	})
}

func TestClient_ExportUsers(t *testing.T) {
	c := initClient(t)
	ch1 := initChannel(t, c)
	ctx := context.Background()

	t.Run("Return error if there are 0 user IDs", func(t *testing.T) {
		_, err := c.ExportUsers(ctx, nil)
		require.Error(t, err)
	})

	t.Run("Export users with no error", func(t *testing.T) {
		resp, err := c.ExportUsers(ctx, []string{ch1.CreatedBy.ID})
		require.NoError(t, err)
		require.NotEmpty(t, resp.TaskID)

		requireTaskReadable(t, c.GetTask, resp.TaskID)
	})
}
