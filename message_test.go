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
	ch, err := c.CreateChannel(context.Background(), ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	resp, err := ch.SendMessage(context.Background(), msg, userB.ID)
	require.NoError(t, err)

	msg, err = c.PinMessage(context.Background(), resp.Message.ID, userA.ID, nil)
	require.NoError(t, err)
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	msg, err = c.UnPinMessage(context.Background(), msg.ID, userA.ID)
	require.NoError(t, err)
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)

	expireAt := time.Now().Add(3 * time.Second)
	msg, err = c.PinMessage(context.Background(), msg.ID, userA.ID, &expireAt)
	require.NoError(t, err)
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	time.Sleep(3 * time.Second)
	msg, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)
}
