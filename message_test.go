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

func TestClient_SendMessage(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)

	ctx := context.Background()

	ch := initChannel(t, c, user.ID)
	resp1, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	msg := &Message{ID: randomString(10), Text: "test message", MML: "test mml", HTML: "test HTML"}
	messageResp, err := resp1.Channel.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err)
	require.Equal(t, ch.CID, messageResp.Message.CID)
	require.Equal(t, user.ID, messageResp.Message.User.ID)
	require.Equal(t, msg.ID, messageResp.Message.ID)
	require.Equal(t, msg.Text, messageResp.Message.Text)
	require.Equal(t, msg.MML, messageResp.Message.MML)
	require.Equal(t, msg.HTML, messageResp.Message.HTML)
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

func TestClient_SendMessage_WithPendingFalse(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)

	ctx := context.Background()

	ch := initChannel(t, c, user.ID)
	resp1, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "message with WithPending(false) - non-pending message"}
	messageResp, err := resp1.Channel.SendMessage(ctx, msg, user.ID, WithPending(false))
	require.NoError(t, err)

	// Get the message to verify it's not in pending state
	gotMsg, err := c.GetMessage(ctx, messageResp.Message.ID)
	require.NoError(t, err)

	// No need to commit the message as it's already in non-pending state
	// The message should be immediately available without requiring a commit
	require.NotNil(t, gotMsg.Message)
	require.Equal(t, msg.Text, gotMsg.Message.Text)
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

func TestMessage_ChannelRoleInMember(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create two users: one default member and one with a custom channel role
	userMember := randomUser(t, c)
	userCustom := randomUser(t, c)

	// Create a channel and assign the custom role to the second user
	chanID := randomString(12)
	chResp, err := c.CreateChannel(ctx, "messaging", chanID, userMember.ID, &ChannelRequest{
		ChannelMembers: []*ChannelMember{
			{UserID: userMember.ID, ChannelRole: "channel_member"},
			{UserID: userCustom.ID, ChannelRole: "custom_role"},
		},
	})
	require.NoError(t, err, "create channel")
	ch := chResp.Channel

	// Send a message as the default member
	msgMember := &Message{Text: "message from channel_member"}
	respMember, err := ch.SendMessage(ctx, msgMember, userMember.ID)
	require.NoError(t, err, "send message member")
	require.NotNil(t, respMember.Message.Member)
	assert.Equal(t, "channel_member", respMember.Message.Member.ChannelRole)

	// Send a message as the custom-role member
	msgCustom := &Message{Text: "message from custom_role"}
	respCustom, err := ch.SendMessage(ctx, msgCustom, userCustom.ID)
	require.NoError(t, err, "send message custom role")
	require.NotNil(t, respCustom.Message.Member)
	assert.Equal(t, "custom_role", respCustom.Message.Member.ChannelRole)

	// Fetch channel state and verify both messages retain the correct channel_role
	queryResp, err := c.QueryChannels(ctx, &QueryOption{
		Filter: map[string]interface{}{"cid": ch.CID},
		UserID: userMember.ID,
	})
	require.NoError(t, err, "query channel")
	require.Len(t, queryResp.Channels, 1, "one channel should match filter")

	roles := map[string]string{
		userMember.ID: "channel_member",
		userCustom.ID: "custom_role",
	}
	for _, m := range queryResp.Channels[0].Messages {
		expectedRole, ok := roles[m.User.ID]
		if !ok {
			continue // skip system messages or others
		}
		require.NotNil(t, m.Member)
		assert.Equal(t, expectedRole, m.Member.ChannelRole,
			"user %s should have role %s", m.User.ID, expectedRole)
	}

func TestClient_DeleteMessageWithOptions_DeleteForMe(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)
	ctx := context.Background()

	ch := initChannel(t, c, user.ID)
	resp1, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	// Send a message to delete
	msg := &Message{Text: "test message to delete for me"}
	messageResp, err := resp1.Channel.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err)

	// Test delete for me only
	_, err = c.DeleteMessageWithOptions(ctx, messageResp.Message.ID, DeleteMessageWithDeleteForMe(user.ID))
	require.NoError(t, err)
}
