package stream_chat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_TranslateMessage(t *testing.T) {
	c := initClient(t)
	u := randomUser(t, c)
	ch := initChannel(t, c, u.ID)
	ctx := context.Background()

	msg := &Message{Text: "test message"}
	messageResp, err := ch.SendMessage(ctx, msg, u.ID)
	require.NoError(t, err)

	translated, err := c.TranslateMessage(ctx, messageResp.Message.ID, "es")
	require.NoError(t, err)
	require.Equal(t, "mensaje de prueba", translated.Message.I18n["es_text"])
}

func TestClient_PinMessage(t *testing.T) {
	c := initClient(t)
	userA := randomUser(t, c)
	userB := randomUser(t, c)
	ctx := context.Background()

	ch := initChannel(t, c, userA.ID, userB.ID)
	resp1, err := c.CreateChannel(ctx, ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	messageResp, err := resp1.Channel.SendMessage(ctx, msg, userB.ID)
	require.NoError(t, err)

	msgWithOptions := &Message{Text: "test message"}
	quotedMsgResp, err := resp1.Channel.SendMessage(ctx, msgWithOptions, userB.ID, func(msg *messageRequest) {
		msg.Message.QuotedMessageID = messageResp.Message.ID
	})
	require.NoError(t, err)
	require.Equal(t, messageResp.Message.ID, quotedMsgResp.Message.QuotedMessageID)

	messageResp, err = c.PinMessage(ctx, messageResp.Message.ID, userA.ID, nil)
	require.NoError(t, err)

	msg = messageResp.Message
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	messageResp, err = c.UnPinMessage(ctx, msg.ID, userA.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)

	expireAt := time.Now().Add(3 * time.Second)
	messageResp, err = c.PinMessage(ctx, msg.ID, userA.ID, &expireAt)
	require.NoError(t, err)

	msg = messageResp.Message
	require.NotZero(t, msg.PinnedAt)
	require.NotZero(t, msg.PinnedBy)
	require.Equal(t, userA.ID, msg.PinnedBy.ID)

	time.Sleep(3 * time.Second)
	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Zero(t, msg.PinnedAt)
	require.Zero(t, msg.PinnedBy)
}
