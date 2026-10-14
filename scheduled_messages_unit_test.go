package stream_chat

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const scheduledMessageJSON = `{
	"id": "sm-1",
	"cid": "messaging:chan",
	"user": {"id": "user-1"},
	"message": {
		"text": "hello later",
		"mentioned_users": ["user-2"],
		"custom": {"mood": "happy"}
	},
	"scheduled_at": "2030-01-02T03:04:05Z",
	"status": "scheduled",
	"created_at": "2029-12-31T00:00:00Z",
	"updated_at": "2029-12-31T00:00:00Z"
}`

type recordedRequest struct {
	method string
	path   string
	body   map[string]interface{}
}

func newScheduledMessageMockClient(t *testing.T, response string) (*Client, *recordedRequest) {
	t.Helper()

	rec := &recordedRequest{}
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		rec.method = r.Method
		rec.path = r.URL.EscapedPath()
		rec.body = nil
		if b, err := io.ReadAll(r.Body); err == nil && len(b) > 0 {
			require.NoError(t, json.Unmarshal(b, &rec.body))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	})
	return c, rec
}

func TestChannel_SendMessageWithScheduledAt(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{
		"duration": "0.01",
		"message": {"id": "sm-1", "text": "hello later"},
		"scheduled_message": `+scheduledMessageJSON+`
	}`)
	ch := &Channel{Type: "messaging", ID: "chan", client: c}

	at := time.Date(2030, 1, 2, 4, 4, 5, 0, time.FixedZone("CET", 3600))
	resp, err := ch.SendMessage(ctx, &Message{Text: "hello later"}, "user-1", WithScheduledAt(at))
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, rec.method)
	assert.Equal(t, "/channels/messaging/chan/message", rec.path)
	assert.Equal(t, "2030-01-02T03:04:05Z", rec.body["scheduled_at"])

	require.NotNil(t, resp.ScheduledMessage)
	sm := resp.ScheduledMessage
	assert.Equal(t, "sm-1", sm.ID)
	assert.Equal(t, "messaging:chan", sm.CID)
	assert.Equal(t, "user-1", sm.User.ID)
	assert.Equal(t, ScheduledMessageStatusScheduled, sm.Status)
	assert.Equal(t, time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC), sm.ScheduledAt.UTC())
	assert.Equal(t, "hello later", sm.Message.Text)
	assert.Equal(t, []string{"user-2"}, sm.Message.MentionedUsers)
	assert.Equal(t, "happy", sm.Message.Custom["mood"])
	assert.Equal(t, "sm-1", resp.Message.ID)
}

func TestChannel_SendMessageWithoutScheduledAt(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{"duration": "0.01", "message": {"id": "m-1"}}`)
	ch := &Channel{Type: "messaging", ID: "chan", client: c}

	resp, err := ch.SendMessage(ctx, &Message{Text: "now"}, "user-1")
	require.NoError(t, err)
	assert.NotContains(t, rec.body, "scheduled_at")
	assert.Nil(t, resp.ScheduledMessage)
}

func TestClient_UpdateScheduledMessage(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{"duration": "0.01", "scheduled_message": `+scheduledMessageJSON+`}`)

	at := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	resp, err := c.UpdateScheduledMessage(ctx, "sm-1", &Message{
		ID:        "ignored",
		Text:      "edited",
		ExtraData: map[string]interface{}{"mood": "happy"},
	}, &at)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPatch, rec.method)
	assert.Equal(t, "/scheduled_messages/sm-1", rec.path)
	assert.Equal(t, "2030-01-02T03:04:05Z", rec.body["scheduled_at"])
	msg, ok := rec.body["message"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "edited", msg["text"])
	assert.Equal(t, "happy", msg["mood"])
	assert.NotContains(t, msg, "id")
	assert.Equal(t, "sm-1", resp.ScheduledMessage.ID)
}

func TestClient_UpdateScheduledMessageOnlyScheduledAt(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{"duration": "0.01", "scheduled_message": `+scheduledMessageJSON+`}`)

	at := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	_, err := c.UpdateScheduledMessage(ctx, "sm-1", nil, &at)
	require.NoError(t, err)
	assert.NotContains(t, rec.body, "message")
	assert.Equal(t, "2030-01-02T03:04:05Z", rec.body["scheduled_at"])
}

func TestClient_UpdateScheduledMessageValidation(t *testing.T) {
	ctx := context.Background()
	c, _ := newScheduledMessageMockClient(t, `{}`)

	_, err := c.UpdateScheduledMessage(ctx, "", &Message{Text: "x"}, nil)
	require.Error(t, err)

	_, err = c.UpdateScheduledMessage(ctx, "sm-1", nil, nil)
	require.Error(t, err)
}

func TestClient_DeleteScheduledMessage(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{"duration": "0.01", "scheduled_message": {"id": "sm-1", "status": "canceled"}}`)

	resp, err := c.DeleteScheduledMessage(ctx, "sm-1")
	require.NoError(t, err)
	assert.Equal(t, http.MethodDelete, rec.method)
	assert.Equal(t, "/scheduled_messages/sm-1", rec.path)
	assert.Equal(t, ScheduledMessageStatusCanceled, resp.ScheduledMessage.Status)

	_, err = c.DeleteScheduledMessage(ctx, "")
	require.Error(t, err)
}

func TestClient_SendScheduledMessageNow(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{"duration": "0.01", "scheduled_message": `+scheduledMessageJSON+`}`)

	resp, err := c.SendScheduledMessageNow(ctx, "sm/1")
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, rec.method)
	assert.Equal(t, "/scheduled_messages/sm%2F1/send_now", rec.path)
	assert.Equal(t, "sm-1", resp.ScheduledMessage.ID)

	_, err = c.SendScheduledMessageNow(ctx, "")
	require.Error(t, err)
}

func TestClient_QueryScheduledMessages(t *testing.T) {
	ctx := context.Background()
	c, rec := newScheduledMessageMockClient(t, `{
		"duration": "0.01",
		"scheduled_messages": [
			{"id": "sm-1", "status": "scheduled", "channel": {"cid": "messaging:chan", "type": "messaging", "id": "chan", "name": "General"}},
			{"id": "sm-2", "status": "failed", "failure_reason": "channel_frozen"}
		],
		"next": "cursor-2"
	}`)

	resp, err := c.QueryScheduledMessages(ctx, &QueryScheduledMessagesOptions{
		UserID: "user-1",
		Filter: &ScheduledMessagesFilter{
			CID:    "messaging:chan",
			Status: []ScheduledMessageStatus{ScheduledMessageStatusScheduled, ScheduledMessageStatusFailed},
		},
		Sort:  []*SortOption{{Field: "scheduled_at", Direction: 1}},
		Limit: 100,
		Next:  "cursor-1",
	})
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, rec.method)
	assert.Equal(t, "/scheduled_messages/query", rec.path)
	assert.Equal(t, "user-1", rec.body["user_id"])
	assert.Equal(t, map[string]interface{}{
		"cid":    "messaging:chan",
		"status": []interface{}{"scheduled", "failed"},
	}, rec.body["filter"])
	assert.Equal(t, []interface{}{map[string]interface{}{"field": "scheduled_at", "direction": float64(1)}}, rec.body["sort"])
	assert.Equal(t, float64(100), rec.body["limit"])
	assert.Equal(t, "cursor-1", rec.body["next"])

	require.Len(t, resp.ScheduledMessages, 2)
	assert.Equal(t, "General", resp.ScheduledMessages[0].Channel.Name)
	assert.Equal(t, ScheduledMessageFailureChannelFrozen, resp.ScheduledMessages[1].FailureReason)
	assert.Equal(t, "cursor-2", resp.Next)

	_, err = c.QueryScheduledMessages(ctx, &QueryScheduledMessagesOptions{})
	require.Error(t, err)
	_, err = c.QueryScheduledMessages(ctx, nil)
	require.Error(t, err)
}

func TestEvent_ScheduledMessageSent(t *testing.T) {
	var e Event
	require.NoError(t, json.Unmarshal([]byte(`{
		"type": "scheduled_message.sent",
		"cid": "messaging:chan",
		"message_id": "final-id",
		"scheduled_message": {"id": "sm-1", "status": "sent"}
	}`), &e))

	assert.Equal(t, EventScheduledMessageSent, e.Type)
	require.NotNil(t, e.ScheduledMessage)
	assert.Equal(t, ScheduledMessageStatusSent, e.ScheduledMessage.Status)
	assert.Equal(t, "final-id", e.ExtraData["message_id"])
	assert.NotContains(t, e.ExtraData, "scheduled_message")
}
