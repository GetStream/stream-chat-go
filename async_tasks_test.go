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

	msg, err := ch.SendMessage(msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without CIDs in parameter
	_, err = c.DeleteChannels([]string{}, true)
	require.Error(t, err)

	taskID, err := c.DeleteChannels([]string{ch.CID}, true)
	require.NoError(t, err)
	require.NotEmpty(t, taskID)

	resp, err := c.GetTask(taskID)
	require.NoError(t, err)
	require.Equal(t, taskID, resp.TaskID)
	require.NotEmpty(t, resp.Status)
}

func TestClient_ExportChannels(t *testing.T) {
	c := initClient(t)
	ch1 := initChannel(t, c)
	ch2 := initChannel(t, c)

	chMembers := append(ch1.Members, ch2.Members...)

	defer func() {
		options := map[string][]string{
			"delete_conversation_channels": {"true"},
			"mark_messages_deleted":        {"true"},
			"hard_delete":                  {"true"},
		}

		for _, u := range chMembers {
			_ = c.DeleteUser(u.UserID, options)
		}
	}()

	t.Run("Return error if there are 0 channels", func(t *testing.T) {
		_, err := c.ExportChannels(nil, nil, nil)
		require.Error(t, err)
	})

	t.Run("Return error if there are >25 channels", func(t *testing.T) {
		var expChannels []*ExportableChannel
		for i := 0; i <= 26; i++ {
			expChannels = append(expChannels,
				&ExportableChannel{Type: fmt.Sprint("test_", i), ID: fmt.Sprint("test_", i)})
		}
		_, err := c.ExportChannels(expChannels, nil, nil)
		require.Error(t, err)
	})

	t.Run("Return error if exportable channel structs are incorrect", func(t *testing.T) {
		expChannels := []*ExportableChannel{
			{Type: "", ID: ch1.ID},
		}
		_, err := c.ExportChannels(expChannels, nil, nil)
		require.Error(t, err)
	})

	t.Run("Export channels with no error", func(t *testing.T) {
		expChannels := []*ExportableChannel{
			{Type: ch1.Type, ID: ch1.ID},
			{Type: ch2.Type, ID: ch2.ID},
		}

		taskID, err := c.ExportChannels(expChannels, nil, nil)
		require.NoError(t, err)
		require.NotEmpty(t, taskID)

		for i := 0; i < 10; i++ {
			task, err := c.GetExportChannelsTask(taskID)
			require.NoError(t, err)
			require.Equal(t, taskID, task.TaskID)
			require.NotEmpty(t, task.Status)

			if task.Status == "completed" {
				break
			}

			time.Sleep(time.Second)
		}
	})
}
