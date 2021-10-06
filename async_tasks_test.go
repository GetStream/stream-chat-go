package stream_chat //nolint: golint

import (
	"testing"

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
