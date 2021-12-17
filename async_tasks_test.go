package stream_chat //nolint: golint

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_DeleteChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user := randomUser(t, c)
	msg := &Message{Text: "test message"}

	_, err := ch.SendMessage(context.Background(), msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without CIDs in parameter
	_, err = c.DeleteChannels(context.Background(), []string{}, true)
	require.Error(t, err)

	resp1, err := c.DeleteChannels(context.Background(), []string{ch.CID}, true)
	require.NoError(t, err)
	require.NotEmpty(t, resp1.TaskID)

	for i := 0; i < 10; i++ {
		resp2, err := c.GetTask(context.Background(), resp1.TaskID)
		require.NoError(t, err)
		require.Equal(t, resp1.TaskID, resp2.TaskID)

		if resp2.Status == TaskStatusCompleted {
			require.Equal(t, resp2.Result[ch.CID], map[string]interface{}{"status": "ok"})
			return
		}

		time.Sleep(time.Second)
	}
}

func TestClient_DeleteUsers(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user := randomUser(t, c)

	msg := &Message{Text: "test message"}

	_, err := ch.SendMessage(context.Background(), msg, user.ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	// should fail without userIDs in parameter
	_, err = c.DeleteUsers(context.Background(), []string{}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.Error(t, err)

	resp1, err := c.DeleteUsers(context.Background(), []string{user.ID}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.TaskID)

	for i := 0; i < 10; i++ {
		resp2, err := c.GetTask(context.Background(), resp1.TaskID)
		require.NoError(t, err)
		require.Equal(t, resp1.TaskID, resp2.TaskID)

		if resp2.Status == TaskStatusCompleted {
			require.Equal(t, resp2.Result[user.ID], map[string]interface{}{"status": "ok"})
			return
		}

		time.Sleep(time.Second)
	}

	require.True(t, false, "task did not succeed")
}

func TestClient_ExportChannels(t *testing.T) {
	c := initClient(t)
	ch1 := initChannel(t, c)
	ch2 := initChannel(t, c)

	chMembers := ch1.Members
	chMembers = append(chMembers, ch2.Members...)

	defer func() {
		for _, u := range chMembers {
			_, _ = c.DeleteUser(context.Background(), u.UserID,
				DeleteUserWithDeleteConversations(),
				DeleteUserWithHardDelete(),
				DeleteUserWithMarkMessagesDeleted())
		}
	}()

	t.Run("Return error if there are 0 channels", func(t *testing.T) {
		_, err := c.ExportChannels(context.Background(), nil, nil, nil)
		require.Error(t, err)
	})

	t.Run("Return error if exportable channel structs are incorrect", func(t *testing.T) {
		expChannels := []*ExportableChannel{
			{Type: "", ID: ch1.ID},
		}
		_, err := c.ExportChannels(context.Background(), expChannels, nil, nil)
		require.Error(t, err)
	})

	t.Run("Export channels with no error", func(t *testing.T) {
		expChannels := []*ExportableChannel{
			{Type: ch1.Type, ID: ch1.ID},
			{Type: ch2.Type, ID: ch2.ID},
		}

		resp1, err := c.ExportChannels(context.Background(), expChannels, nil, nil)
		require.NoError(t, err)
		require.NotEmpty(t, resp1.TaskID)

		for i := 0; i < 10; i++ {
			task, err := c.GetExportChannelsTask(context.Background(), resp1.TaskID)
			require.NoError(t, err)
			require.Equal(t, resp1.TaskID, task.TaskID)
			require.NotEmpty(t, task.Status)

			if task.Status == TaskStatusCompleted {
				break
			}

			time.Sleep(time.Second)
		}
	})
}
