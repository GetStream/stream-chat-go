package stream_chat // nolint: golint

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_PinMessage(t *testing.T) {
	c := initClient(t)
	userA := randomUser(t, c)
	userB := randomUser(t, c)

	ch := initChannel(t, c, userA.ID, userB.ID)
	resp1, err := c.CreateChannel(context.Background(), ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	messageResp, err := resp1.Channel.SendMessage(context.Background(), msg, userB.ID)
	require.NoError(t, err)

	messageResp, err = c.PinMessage(context.Background(), messageResp.Message.ID, userA.ID, nil)
	require.NoError(t, err)

	msg = messageResp.Message
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	messageResp, err = c.UnPinMessage(context.Background(), msg.ID, userA.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)

	expireAt := time.Now().Add(3 * time.Second)
	messageResp, err = c.PinMessage(context.Background(), msg.ID, userA.ID, &expireAt)
	require.NoError(t, err)

	msg = messageResp.Message
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	time.Sleep(3 * time.Second)
	messageResp, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)
}
