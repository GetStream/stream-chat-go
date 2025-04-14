package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_QueryThreads(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	t.Run("basic query", func(t *testing.T) {
		membersID, ch, parentMsg, replyMsg := testThreadSetup(t, c, 3)

		query := &QueryThreadsRequest{
			Filter: map[string]any{
				"channel_cid": map[string]any{
					"$eq": ch.CID,
				},
			},
			Sort: &SortParamRequestList{
				{
					Field:     "created_at",
					Direction: -1,
				},
			},
			PagerRequest: PagerRequest{
				Limit: intPtr(10),
			},
			UserID: membersID[0],
		}

		resp, err := c.QueryThreads(ctx, query)
		require.NoError(t, err)
		require.NotNil(t, resp, "response should not be nil")
		require.NotEmpty(t, resp.Threads, "threads should not be empty")

		thread := resp.Threads[0]
		assertThreadData(t, thread, ch, parentMsg, replyMsg)
		assertThreadParticipants(t, thread, ch.CreatedBy.ID)

		assert.Empty(t, resp.PagerResponse)
	})

	t.Run("with pagination", func(t *testing.T) {
		membersID, ch, parentMsg1, replyMsg1 := testThreadSetup(t, c, 3)

		// Create a second thread
		parentMsg2, err := ch.SendMessage(ctx, &Message{Text: "Parent message for thread 2"}, ch.CreatedBy.ID)
		require.NoError(t, err, "send second parent message")

		replyMsg2, err := ch.SendMessage(ctx, &Message{
			Text:     "Reply message 2",
			ParentID: parentMsg2.Message.ID,
		}, ch.CreatedBy.ID)
		require.NoError(t, err, "send second reply message")

		// First page query
		query := &QueryThreadsRequest{
			Filter: map[string]any{
				"channel_cid": map[string]any{
					"$eq": ch.CID,
				},
			},
			Sort: &SortParamRequestList{
				{
					Field:     "created_at",
					Direction: 1,
				},
			},
			PagerRequest: PagerRequest{
				Limit: intPtr(1),
			},
			UserID: membersID[0],
		}

		resp, err := c.QueryThreads(ctx, query)
		require.NoError(t, err)
		require.NotNil(t, resp, "response should not be nil")
		require.NotEmpty(t, resp.Threads, "threads should not be empty")

		thread := resp.Threads[0]
		assertThreadData(t, thread, ch, parentMsg1, replyMsg1)
		assertThreadParticipants(t, thread, ch.CreatedBy.ID)

		// Second page query
		query2 := &QueryThreadsRequest{
			Filter: map[string]any{
				"channel_cid": map[string]any{
					"$eq": ch.CID,
				},
			},
			Sort: &SortParamRequestList{
				{
					Field:     "created_at",
					Direction: -1,
				},
			},
			PagerRequest: PagerRequest{
				Limit: intPtr(1),
				Next:  resp.Next,
			},
			UserID: membersID[0],
		}

		resp, err = c.QueryThreads(ctx, query2)
		require.NoError(t, err)
		require.NotNil(t, resp, "response should not be nil")
		require.NotEmpty(t, resp.Threads, "threads should not be empty")

		thread = resp.Threads[0]
		assertThreadData(t, thread, ch, parentMsg2, replyMsg2)
		assertThreadParticipants(t, thread, ch.CreatedBy.ID)
	})
}

// testThreadSetup creates a channel with members and returns necessary test data
func testThreadSetup(t *testing.T, c *Client, numMembers int) ([]string, *Channel, *MessageResponse, *MessageResponse) {
	membersID := randomUsersID(t, c, numMembers)
	ch := initChannel(t, c, membersID...)

	// Create a parent message
	parentMsg, err := ch.SendMessage(context.Background(), &Message{Text: "Parent message for thread"}, ch.CreatedBy.ID)
	require.NoError(t, err, "send parent message")

	// Create a thread by sending a reply
	replyMsg, err := ch.SendMessage(context.Background(), &Message{
		Text:     "Reply message",
		ParentID: parentMsg.Message.ID,
	}, ch.CreatedBy.ID)
	require.NoError(t, err, "send reply message")

	return membersID, ch, parentMsg, replyMsg
}

// assertThreadData validates common thread data fields
func assertThreadData(t *testing.T, thread ThreadResponse, ch *Channel, parentMsg, replyMsg *MessageResponse) {
	assert.Equal(t, ch.CID, thread.ChannelCID, "channel CID should match")
	assert.Equal(t, parentMsg.Message.ID, thread.ParentMessageID, "parent message ID should match")
	assert.Equal(t, ch.CreatedBy.ID, thread.CreatedByUserID, "created by user ID should match")
	assert.Equal(t, 1, thread.ReplyCount, "reply count should be 1")
	assert.Equal(t, 1, thread.ParticipantCount, "participant count should be 1")
	assert.Equal(t, parentMsg.Message.Text, thread.Title, "title should not be empty")
	assert.Equal(t, replyMsg.Message.CreatedAt, thread.CreatedAt, "created at should not be zero")
	assert.Equal(t, replyMsg.Message.UpdatedAt, thread.UpdatedAt, "updated at should not be zero")
	assert.Nil(t, thread.DeletedAt, "deleted at should be nil")
}

// assertThreadParticipants validates thread participant data
func assertThreadParticipants(t *testing.T, thread ThreadResponse, createdByID string) {
	require.Len(t, thread.Participants, 1, "should have one participant")
	assert.Equal(t, createdByID, thread.Participants[0].UserID, "participant user ID should match")
	assert.NotZero(t, thread.Participants[0].CreatedAt, "participant created at should not be zero")
	assert.NotZero(t, thread.Participants[0].LastReadAt, "participant last read at should not be zero")
}

// Helper function to create a pointer to an int
func intPtr(i int) *int {
	return &i
}
