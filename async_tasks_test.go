package stream_chat //nolint: golint

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_DeleteChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user := randomUser()

	msg := &Message{Text: "test message"}

	_, err := ch.SendMessage(msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without CIDs in parameter
	_, err = c.DeleteChannels([]string{}, true)
	require.Error(t, err)

	nonExistingCID := "non-existing-cid"
	resp, err := c.DeleteChannels([]string{ch.CID, nonExistingCID}, true)
	require.NoError(t, err)
	require.NotEmpty(t, resp.TaskID)
	require.NotEmpty(t, resp.Result)

	require.Equal(t, "ok", resp.Result[ch.CID].Status)
	require.Equal(t, "error", resp.Result[nonExistingCID].Status)
	require.Equal(t, fmt.Sprintf("channel %q does not exist", nonExistingCID), resp.Result[nonExistingCID].Error)

	for i := 0; i < 10; i++ {
		resp, err := c.GetTask(resp.TaskID)
		require.NoError(t, err)
		require.Equal(t, resp.TaskID, resp.TaskID)

		if resp.Status == "completed" {
			require.Equal(t, resp.Result[ch.CID], map[string]interface{}{"status": "ok"})
			return
		}

		time.Sleep(time.Second)
	}
}

func TestClient_DeleteUsers(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user := randomUser()

	msg := &Message{Text: "test message"}

	_, err := ch.SendMessage(msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without userIDs in parameter
	_, err = c.DeleteUsers([]string{}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.Error(t, err)

	taskID, err := c.DeleteUsers([]string{user.ID}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.NoError(t, err)
	require.NotEmpty(t, taskID)

	for i := 0; i < 10; i++ {
		resp, err := c.GetTask(taskID)
		require.NoError(t, err)
		require.Equal(t, taskID, resp.TaskID)

		if resp.Status == "completed" {
			require.Equal(t, resp.Result[user.ID], map[string]interface{}{"status": "ok"})
			return
		}

		time.Sleep(time.Second)
	}

	require.True(t, false, "task did not succeed")
}
