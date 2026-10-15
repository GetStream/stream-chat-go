package stream_chat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// taskWaitBudget is how long waitForTask polls before failing. The shared CI
// app sometimes takes well over 10s to pick up async tasks.
const taskWaitBudget = 60 * time.Second

// waitForTask polls get until the task completes and returns it. It fails the
// test right away if the task reports failed, and on timeout, in both cases
// including the last seen status and result so a slow task can be told apart
// from a failed one.
func waitForTask(t *testing.T, get func(context.Context, string) (*TaskResponse, error), taskID string) *TaskResponse {
	t.Helper()
	ctx := context.Background()

	var last *TaskResponse
	for deadline := time.Now().Add(taskWaitBudget); time.Now().Before(deadline); time.Sleep(time.Second) {
		task, err := get(ctx, taskID)
		require.NoError(t, err)
		require.Equal(t, taskID, task.TaskID)
		last = task

		switch task.Status {
		case TaskStatusCompleted:
			return task
		case TaskStatusFailed:
			require.FailNowf(t, "task failed", "task %s: status=%q result=%v", taskID, task.Status, task.Result)
		}
	}

	require.FailNowf(t, "task did not complete", "task %s after %s: last status=%q result=%v", taskID, taskWaitBudget, last.Status, last.Result)
	return nil
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

	task := waitForTask(t, c.GetTask, resp1.TaskID)
	require.Equal(t, map[string]interface{}{"status": "ok"}, task.Result[ch.CID])
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

	task := waitForTask(t, c.GetTask, resp1.TaskID)
	require.Equal(t, map[string]interface{}{"status": "ok"}, task.Result[user.ID])
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

		waitForTask(t, c.GetExportChannelsTask, resp1.TaskID)
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

		task := waitForTask(t, c.GetTask, resp.TaskID)
		require.Contains(t, task.Result["url"], "/exports/users/")
	})
}
