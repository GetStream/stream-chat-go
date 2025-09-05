package stream_chat

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChannel_MarkDelivered(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel with members
	membersID := randomUsersID(t, c, 2)
	ch := initChannel(t, c, membersID...)

	// Send a message to the channel
	msg, err := ch.SendMessage(ctx, &Message{Text: "test message for delivery"}, ch.CreatedBy.ID)
	require.NoError(t, err)

	t.Run("successful mark delivered with full options", func(t *testing.T) {
		userID := membersID[0]
		clientID := "test-client-123"
		connectionID := "test-connection-456"

		options := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{ch.cid(): msg.Message.ID},
			UserID:                  userID,
			ClientID:                &clientID,
			ConnectionID:            &connectionID,
		}

		resp, err := ch.MarkDelivered(ctx, options)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("successful mark delivered with minimal options", func(t *testing.T) {
		userID := membersID[1]

		options := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{ch.cid(): msg.Message.ID},
			UserID:                  userID,
		}

		resp, err := ch.MarkDelivered(ctx, options)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("successful mark delivered with user object", func(t *testing.T) {
		user := &User{ID: membersID[0]}

		options := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{ch.cid(): msg.Message.ID},
			User:                    user,
		}

		resp, err := ch.MarkDelivered(ctx, options)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("error when options is nil", func(t *testing.T) {
		resp, err := ch.MarkDelivered(ctx, nil)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "options must not be nil")
	})

	t.Run("error when channel_delivered_message is empty", func(t *testing.T) {
		userID := membersID[0]

		options := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{},
			UserID:                  userID,
		}

		resp, err := ch.MarkDelivered(ctx, options)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "channel_delivered_message must not be empty")
	})

	t.Run("mark delivered for multiple channels", func(t *testing.T) {
		// Create another channel
		ch2 := initChannel(t, c, membersID...)
		msg2, err := ch2.SendMessage(ctx, &Message{Text: "test message 2"}, ch2.CreatedBy.ID)
		require.NoError(t, err)

		userID := membersID[0]

		options := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{
				ch.cid():  msg.Message.ID,
				ch2.cid(): msg2.Message.ID,
			},
			UserID: userID,
		}

		resp, err := ch.MarkDelivered(ctx, options)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

func TestChannel_MarkDeliveredSimple(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel with members
	membersID := randomUsersID(t, c, 2)
	ch := initChannel(t, c, membersID...)

	// Send a message to the channel
	msg, err := ch.SendMessage(ctx, &Message{Text: "test message for simple delivery"}, ch.CreatedBy.ID)
	require.NoError(t, err)

	t.Run("successful mark delivered simple", func(t *testing.T) {
		userID := membersID[0]
		messageID := msg.Message.ID

		resp, err := ch.MarkDeliveredSimple(ctx, userID, messageID)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("error when userID is empty", func(t *testing.T) {
		resp, err := ch.MarkDeliveredSimple(ctx, "", msg.Message.ID)
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "user ID must not be empty")
	})

	t.Run("error when messageID is empty", func(t *testing.T) {
		userID := membersID[0]

		resp, err := ch.MarkDeliveredSimple(ctx, userID, "")
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "message ID must not be empty")
	})

	t.Run("error when both userID and messageID are empty", func(t *testing.T) {
		resp, err := ch.MarkDeliveredSimple(ctx, "", "")
		require.Error(t, err)
		require.Nil(t, resp)
		require.Contains(t, err.Error(), "user ID must not be empty")
	})
}

func TestChannel_MarkDelivered_Integration(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel with members
	membersID := randomUsersID(t, c, 3)
	ch := initChannel(t, c, membersID...)

	// Send multiple messages
	msg1, err := ch.SendMessage(ctx, &Message{Text: "message 1"}, ch.CreatedBy.ID)
	require.NoError(t, err)

	msg2, err := ch.SendMessage(ctx, &Message{Text: "message 2"}, ch.CreatedBy.ID)
	require.NoError(t, err)

	t.Run("mark different messages as delivered for different users", func(t *testing.T) {
		// Mark message 1 as delivered for user 1
		resp1, err := ch.MarkDeliveredSimple(ctx, membersID[0], msg1.Message.ID)
		require.NoError(t, err)
		require.NotNil(t, resp1)

		// Mark message 2 as delivered for user 2
		resp2, err := ch.MarkDeliveredSimple(ctx, membersID[1], msg2.Message.ID)
		require.NoError(t, err)
		require.NotNil(t, resp2)

		// Mark both messages as delivered for user 3
		options := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{
				ch.cid(): msg1.Message.ID,
			},
			UserID: membersID[2],
		}
		resp3, err := ch.MarkDelivered(ctx, options)
		require.NoError(t, err)
		require.NotNil(t, resp3)

		// Mark message 2 as delivered for user 3 as well
		resp4, err := ch.MarkDeliveredSimple(ctx, membersID[2], msg2.Message.ID)
		require.NoError(t, err)
		require.NotNil(t, resp4)
	})
}

func TestMarkDeliveredOptions_JSON(t *testing.T) {
	t.Run("marshal and unmarshal MarkDeliveredOptions", func(t *testing.T) {
		userID := "test-user-123"
		clientID := "test-client-456"
		connectionID := "test-connection-789"
		user := &User{ID: userID, Name: "Test User"}

		original := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{
				"messaging:general": "msg-123",
				"messaging:private": "msg-456",
			},
			ClientID:     &clientID,
			ConnectionID: &connectionID,
			User:         user,
			UserID:       userID,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled MarkDeliveredOptions
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		// Verify the data
		require.Equal(t, original.ChannelDeliveredMessage, unmarshaled.ChannelDeliveredMessage)
		require.Equal(t, *original.ClientID, *unmarshaled.ClientID)
		require.Equal(t, *original.ConnectionID, *unmarshaled.ConnectionID)
		require.Equal(t, original.UserID, unmarshaled.UserID)
		require.Equal(t, original.User.ID, unmarshaled.User.ID)
		require.Equal(t, original.User.Name, unmarshaled.User.Name)
	})

	t.Run("marshal with minimal options", func(t *testing.T) {
		userID := "test-user-123"

		original := &MarkDeliveredOptions{
			ChannelDeliveredMessage: map[string]string{
				"messaging:general": "msg-123",
			},
			UserID: userID,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled MarkDeliveredOptions
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		// Verify the data
		require.Equal(t, original.ChannelDeliveredMessage, unmarshaled.ChannelDeliveredMessage)
		require.Equal(t, original.UserID, unmarshaled.UserID)
		require.Nil(t, unmarshaled.ClientID)
		require.Nil(t, unmarshaled.ConnectionID)
		require.Nil(t, unmarshaled.User)
	})
}

func TestDeliveryReceipts_JSON(t *testing.T) {
	t.Run("marshal and unmarshal DeliveryReceipts", func(t *testing.T) {
		original := &DeliveryReceipts{
			Enabled: true,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled DeliveryReceipts
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		// Verify the data
		require.Equal(t, original.Enabled, unmarshaled.Enabled)
	})

	t.Run("marshal disabled DeliveryReceipts", func(t *testing.T) {
		original := &DeliveryReceipts{
			Enabled: false,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled DeliveryReceipts
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		// Verify the data
		require.Equal(t, original.Enabled, unmarshaled.Enabled)
	})
}

func TestChannelRead_WithDeliveryFields(t *testing.T) {
	t.Run("marshal and unmarshal ChannelRead with delivery fields", func(t *testing.T) {
		now := time.Now()
		messageID := "msg-123"

		original := &ChannelRead{
			User:                   &User{ID: "user-123", Name: "Test User"},
			LastRead:               now,
			UnreadMessages:         5,
			LastDeliveredAt:        &now,
			LastDeliveredMessageID: &messageID,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled ChannelRead
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		// Verify the data
		require.Equal(t, original.User.ID, unmarshaled.User.ID)
		require.Equal(t, original.User.Name, unmarshaled.User.Name)
		require.Equal(t, original.LastRead.Unix(), unmarshaled.LastRead.Unix())
		require.Equal(t, original.UnreadMessages, unmarshaled.UnreadMessages)
		require.Equal(t, original.LastDeliveredAt.Unix(), unmarshaled.LastDeliveredAt.Unix())
		require.Equal(t, *original.LastDeliveredMessageID, *unmarshaled.LastDeliveredMessageID)
	})

	t.Run("marshal ChannelRead without delivery fields", func(t *testing.T) {
		now := time.Now()

		original := &ChannelRead{
			User:           &User{ID: "user-123", Name: "Test User"},
			LastRead:       now,
			UnreadMessages: 5,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(original)
		require.NoError(t, err)

		// Unmarshal back
		var unmarshaled ChannelRead
		err = json.Unmarshal(jsonData, &unmarshaled)
		require.NoError(t, err)

		// Verify the data
		require.Equal(t, original.User.ID, unmarshaled.User.ID)
		require.Equal(t, original.User.Name, unmarshaled.User.Name)
		require.Equal(t, original.LastRead.Unix(), unmarshaled.LastRead.Unix())
		require.Equal(t, original.UnreadMessages, unmarshaled.UnreadMessages)
		require.Nil(t, unmarshaled.LastDeliveredAt)
		require.Nil(t, unmarshaled.LastDeliveredMessageID)
	})
}
