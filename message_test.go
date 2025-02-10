package stream_chat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	require.Empty(t, messageResp.Message.Attachments)

	time.Sleep(3 * time.Second)
	gotMsg, err := c.GetMessage(ctx, messageResp.Message.ID)
	require.NoError(t, err)
	require.Empty(t, gotMsg.Message.Attachments)
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
	require.Empty(t, result.Channels)
}

func TestClient_UpdateRestrictedVisibilityMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	adminUser := randomUserWithRole(t, c, "admin")
	user1 := randomUser(t, c)
	user2 := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		RestrictedVisibility: []string{
			user1.ID,
		},
	}

	resp, err := ch.SendMessage(ctx, msg, adminUser.ID)
	require.NoError(t, err, "send message")

	msg = resp.Message
	msg.RestrictedVisibility = []string{user2.ID}
	msg.UserID = adminUser.ID
	resp, err = c.UpdateMessage(ctx, msg, msg.ID)
	require.NoError(t, err, "send message")
	assert.Equal(t, []string{user2.ID}, resp.Message.RestrictedVisibility)
}

func TestClient_PartialUpdateRestrictedVisibilityMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	adminUser := randomUserWithRole(t, c, "admin")
	user1 := randomUser(t, c)
	user2 := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		RestrictedVisibility: []string{
			user1.ID,
		},
	}

	messageResponse, err := ch.SendMessage(ctx, msg, adminUser.ID)
	require.NoError(t, err, "send message")

	resp, err := c.PartialUpdateMessage(ctx, messageResponse.Message.ID, &MessagePartialUpdateRequest{
		UserID: adminUser.ID,
		PartialUpdate: PartialUpdate{
			Set: map[string]interface{}{
				"restricted_visibility": []string{user2.ID},
			},
		},
	})
	require.NoError(t, err, "send message")
	assert.Equal(t, []string{user2.ID}, resp.Message.RestrictedVisibility)
}
