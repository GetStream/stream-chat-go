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

func TestClient_SendMessage_Pending(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)

	ctx := context.Background()

	ch := initChannel(t, c, user.ID)
	resp1, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test pending message"}
	metadata := map[string]string{"my": "metadata"}
	messageResp, err := resp1.Channel.SendMessage(ctx, msg, user.ID, MessagePending, MessagePendingMessageMetadata(metadata))
	require.NoError(t, err)
	require.Equal(t, metadata, messageResp.PendingMessageMetadata)

	gotMsg, err := c.GetMessage(ctx, messageResp.Message.ID)
	require.NoError(t, err)
	require.Equal(t, metadata, gotMsg.PendingMessageMetadata)

	_, err = c.CommitMessage(ctx, messageResp.Message.ID)
	require.NoError(t, err)
}

func TestClient_SendMessage_SkipEnrichURL(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)

	ctx := context.Background()

	ch := initChannel(t, c, user.ID)
	resp1, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test message with link to https://getstream.io"}
	messageResp, err := resp1.Channel.SendMessage(ctx, msg, user.ID, MessageSkipEnrichURL)
	require.NoError(t, err)
	require.Len(t, messageResp.Message.Attachments, 0)

	time.Sleep(3 * time.Second)
	gotMsg, err := c.GetMessage(ctx, messageResp.Message.ID)
	require.NoError(t, err)
	require.Len(t, gotMsg.Message.Attachments, 0)
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

func TestClient_SendMessage_KeepChannelHidden(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)

	ctx := context.Background()

	ch := initChannel(t, c, user.ID)
	resp, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	_, err = resp.Channel.Hide(ctx, user.ID)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	_, err = resp.Channel.SendMessage(ctx, msg, user.ID, KeepChannelHidden)
	require.NoError(t, err)

	result, err := c.QueryChannels(ctx, &QueryOption{
		Filter: map[string]interface{}{"cid": resp.Channel.CID},
		UserID: user.ID,
	})
	require.NoError(t, err)
	require.Len(t, result.Channels, 0)
}
