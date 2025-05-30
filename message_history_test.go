package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageHistory(t *testing.T) {
	client := initClient(t)
	users := randomUsers(t, client, 2)
	user1 := users[0]
	user2 := users[1]

	ch := initChannel(t, client, user1.ID)

	ctx := context.Background()
	initialText := "initial text"
	customField := "custom_field"
	initialCustomFieldValue := "custom value"
	// send a message with initial text (user1)
	response, err := ch.SendMessage(ctx, &Message{Text: initialText, ExtraData: map[string]interface{}{customField: initialCustomFieldValue}}, user1.ID)
	require.NoError(t, err)
	message := response.Message

	updatedText1 := "updated text"
	updatedCustomFieldValue := "updated custom value"
	// update the message by user1
	_, err = client.UpdateMessage(ctx, &Message{Text: updatedText1, ExtraData: map[string]interface{}{customField: updatedCustomFieldValue}, UserID: user1.ID}, message.ID)
	assert.NoError(t, err)

	updatedText2 := "updated text 2"
	// update the message by user2
	_, err = client.UpdateMessage(ctx, &Message{Text: updatedText2, UserID: user2.ID}, message.ID)
	assert.NoError(t, err)

	t.Run("test query", func(t *testing.T) {
		req := QueryMessageHistoryRequest{
			Filter: map[string]interface{}{
				"message_id": message.ID,
			},
		}
		messageHistoryResponse, err := client.QueryMessageHistory(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, messageHistoryResponse)

		history := messageHistoryResponse.MessageHistory
		assert.Equal(t, 2, len(history))

		firstUpdate := history[1]
		assert.Equal(t, initialText, firstUpdate.Text)
		assert.Equal(t, user1.ID, firstUpdate.MessageUpdatedByID)
		assert.Equal(t, initialCustomFieldValue, firstUpdate.ExtraData[customField].(string))

		secondUpdate := history[0]
		assert.Equal(t, updatedText1, secondUpdate.Text)
		assert.Equal(t, user1.ID, secondUpdate.MessageUpdatedByID)
		assert.Equal(t, updatedCustomFieldValue, secondUpdate.ExtraData[customField].(string))
	})

	t.Run("test sorting", func(t *testing.T) {
		sortedHistoryQueryRequest := QueryMessageHistoryRequest{
			Filter: map[string]interface{}{
				"message_id": message.ID,
			},
			Sort: []*SortOption{
				{
					Field:     "message_updated_at",
					Direction: 1,
				},
			},
		}
		sortedHistoryResponse, err := client.QueryMessageHistory(ctx, sortedHistoryQueryRequest)
		assert.NoError(t, err)
		assert.NotNil(t, sortedHistoryResponse)

		sortedHistory := sortedHistoryResponse.MessageHistory
		assert.Equal(t, 2, len(sortedHistory))

		firstUpdate := sortedHistory[0]
		assert.Equal(t, initialText, firstUpdate.Text)
		assert.Equal(t, user1.ID, firstUpdate.MessageUpdatedByID)

		secondUpdate := sortedHistory[1]
		assert.Equal(t, updatedText1, secondUpdate.Text)
		assert.Equal(t, user1.ID, secondUpdate.MessageUpdatedByID)
	})
}
