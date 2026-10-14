package stream_chat_test

import (
	"context"
	"errors"
	"testing"
	"time"

	stream_chat "github.com/GetStream/stream-chat-go/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// scheduleTestMessage schedules a message in a new channel, skipping the test
// when the app does not have scheduled messages enabled.
func scheduleTestMessage(t *testing.T, c *stream_chat.Client) (*stream_chat.Channel, *stream_chat.User, *stream_chat.ScheduledMessage) {
	t.Helper()

	ctx := context.Background()
	ch, user := createTestChannel(t, c)

	resp, err := ch.SendMessage(ctx, &stream_chat.Message{Text: "scheduled hello"}, user.ID,
		stream_chat.WithScheduledAt(time.Now().Add(time.Hour)))
	var apiErr stream_chat.Error
	if errors.As(err, &apiErr) && apiErr.ExceptionFields["reason"] == "scheduled_messages_disabled" {
		t.Skip("scheduled messages are not enabled for this app")
	}
	require.NoError(t, err)
	if resp.ScheduledMessage == nil {
		t.Skip("the API does not support scheduled messages yet; the message was sent immediately")
	}

	return ch, user, resp.ScheduledMessage
}

func TestClient_ScheduledMessagesLifecycle(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	ch, user, sm := scheduleTestMessage(t, c)

	assert.NotEmpty(t, sm.ID)
	assert.Equal(t, ch.CID, sm.CID)
	assert.Equal(t, user.ID, sm.User.ID)
	assert.Equal(t, stream_chat.ScheduledMessageStatusScheduled, sm.Status)
	assert.Equal(t, "scheduled hello", sm.Message.Text)

	t.Run("Query", func(t *testing.T) {
		resp, err := c.QueryScheduledMessages(ctx, &stream_chat.QueryScheduledMessagesOptions{
			UserID: user.ID,
			Filter: &stream_chat.ScheduledMessagesFilter{CID: ch.CID},
		})
		require.NoError(t, err)
		require.Len(t, resp.ScheduledMessages, 1)
		assert.Equal(t, sm.ID, resp.ScheduledMessages[0].ID)
		require.NotNil(t, resp.ScheduledMessages[0].Channel)
		assert.Equal(t, ch.CID, resp.ScheduledMessages[0].Channel.CID)
	})

	t.Run("Update", func(t *testing.T) {
		at := time.Now().Add(2 * time.Hour).Truncate(time.Second)
		resp, err := c.UpdateScheduledMessage(ctx, sm.ID, &stream_chat.Message{Text: "edited"}, &at)
		require.NoError(t, err)
		assert.Equal(t, "edited", resp.ScheduledMessage.Message.Text)
		assert.True(t, at.Equal(resp.ScheduledMessage.ScheduledAt))
	})

	t.Run("Delete", func(t *testing.T) {
		resp, err := c.DeleteScheduledMessage(ctx, sm.ID)
		require.NoError(t, err)
		assert.Equal(t, sm.ID, resp.ScheduledMessage.ID)

		query, err := c.QueryScheduledMessages(ctx, &stream_chat.QueryScheduledMessagesOptions{
			UserID: user.ID,
			Filter: &stream_chat.ScheduledMessagesFilter{CID: ch.CID},
		})
		require.NoError(t, err)
		assert.Empty(t, query.ScheduledMessages)
	})
}

func TestClient_ScheduledMessagesSendNow(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	_, _, sm := scheduleTestMessage(t, c)

	resp, err := c.SendScheduledMessageNow(ctx, sm.ID)
	require.NoError(t, err)
	assert.Equal(t, sm.ID, resp.ScheduledMessage.ID)
	assert.False(t, resp.ScheduledMessage.ScheduledAt.After(time.Now().Add(time.Minute)))
}
