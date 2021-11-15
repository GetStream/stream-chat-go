package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_PinMessage(t *testing.T) {
	c := initClient(t)
	userA := randomUser(t, c)
	userB := randomUser(t, c)

	ch := initChannel(t, c, userA.ID, userB.ID)
	ch, err := c.CreateChannel(ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	msg, err = ch.SendMessage(msg, userB.ID)
	require.NoError(t, err)

	msg, err = c.PinMessage(msg.ID, userA.ID, nil)
	require.NoError(t, err)
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	msg, err = c.UnPinMessage(msg.ID, userA.ID)
	require.NoError(t, err)
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)
}
