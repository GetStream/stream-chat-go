package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannel_CreateDraft(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	channel := initChannel(t, c)
	user := randomUser(t, c)

	// Create a draft message
	message := &messageRequestMessage{
		Text: "This is a draft message",
	}

	resp, err := channel.CreateDraft(ctx, user.ID, message)
	require.NoError(t, err)
	assert.Equal(t, "This is a draft message", resp.Draft.Message.Text)
	assert.Equal(t, channel.CID, resp.Draft.ChannelCID)
	assert.NotNil(t, resp.Draft.Channel)
	assert.Equal(t, channel.CID, resp.Draft.Channel.CID)
}

func TestChannel_GetDraft(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	channel := initChannel(t, c)
	user := randomUser(t, c)

	// Create a draft message
	message := &messageRequestMessage{
		Text: "This is a draft message",
	}

	createResp, err := channel.CreateDraft(ctx, user.ID, message)
	require.NoError(t, err)

	// Get the draft
	resp, err := channel.GetDraft(ctx, nil, user.ID)
	require.NoError(t, err)
	assert.Equal(t, createResp.Draft.Message.ID, resp.Draft.Message.ID)
	assert.Equal(t, "This is a draft message", resp.Draft.Message.Text)
}

func TestChannel_DeleteDraft(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	channel := initChannel(t, c)
	user := randomUser(t, c)

	// Create a draft message
	message := &messageRequestMessage{
		Text:   "This is a draft message",
		UserID: user.ID,
	}

	_, err := channel.CreateDraft(ctx, user.ID, message)
	require.NoError(t, err)

	// Delete the draft
	resp, err := channel.DeleteDraft(ctx, user.ID, nil)
	require.NoError(t, err)
	// Just verify the response is received, not specific fields
	assert.NotNil(t, resp)

	// Verify the draft is deleted
	_, err = channel.GetDraft(ctx, nil, user.ID)
	require.Error(t, err)
}

func TestChannel_CreateDraftInThread(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	channel := initChannel(t, c)

	// Create a parent message
	userID := randomUser(t, c).ID
	parentMsg, err := channel.SendMessage(ctx, &Message{
		Text: "Parent message",
		User: &User{ID: userID},
	}, userID)
	require.NoError(t, err)

	// Create a draft message in thread
	parentID := parentMsg.Message.ID
	message := &messageRequestMessage{
		Text:     "This is a draft reply",
		ParentID: parentID,
	}

	resp, err := channel.CreateDraft(ctx, userID, message)
	require.NoError(t, err)
	assert.Equal(t, "This is a draft reply", resp.Draft.Message.Text)
	assert.Equal(t, parentID, *resp.Draft.Message.ParentID)

	// Get the draft in thread
	getResp, err := channel.GetDraft(ctx, &parentID, userID)
	require.NoError(t, err)
	assert.Equal(t, resp.Draft.Message.ID, getResp.Draft.Message.ID)
	assert.Equal(t, parentID, *getResp.Draft.Message.ParentID)
}

func TestClient_QueryDrafts(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	channel := initChannel(t, c)
	user := randomUser(t, c)

	// Create a draft message
	message1 := &messageRequestMessage{
		Text: "Draft 1",
	}
	_, err := channel.CreateDraft(ctx, user.ID, message1)
	require.NoError(t, err)

	// Create a second channel
	channel2 := initChannel(t, c)

	// Create a draft in the second channel
	message2 := &messageRequestMessage{
		Text: "Draft 2",
	}
	_, err = channel2.CreateDraft(ctx, user.ID, message2)
	require.NoError(t, err)

	// Query all drafts
	resp, err := c.QueryDrafts(ctx, &QueryDraftsOptions{UserID: user.ID, Limit: 10})
	require.NoError(t, err)

	// Verify we have at least 2 drafts
	assert.GreaterOrEqual(t, len(resp.Drafts), 2)

	// Check if we can find our drafts
	foundDraft1 := false
	foundDraft2 := false
	for _, draft := range resp.Drafts {
		if draft.Message.Text == "Draft 1" {
			foundDraft1 = true
		} else if draft.Message.Text == "Draft 2" {
			foundDraft2 = true
		}
	}
	assert.True(t, foundDraft1, "First draft not found")
	assert.True(t, foundDraft2, "Second draft not found")
}

func TestClient_QueryDraftsWithFilters(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	user := randomUser(t, c)
	channel1 := initChannel(t, c, user)

	// Create a draft message
	draft1 := &messageRequestMessage{
		Text: "Draft in channel 1",
	}
	_, err := channel1.CreateDraft(ctx, user.ID, draft1)
	require.NoError(t, err)

	// Create a second channel
	channel2 := initChannel(t, c, user)

	// Create a draft in the second channel
	draft2 := &messageRequestMessage{
		Text: "Draft in channel 2",
	}
	_, err = channel2.CreateDraft(ctx, user.ID, draft2)
	require.NoError(t, err)

	// Query all drafts for the user
	resp, err := c.QueryDrafts(ctx, &QueryDraftsOptions{UserID: user.ID})
	require.NoError(t, err)
	assert.Equal(t, 2, len(resp.Drafts))

	// Query drafts for a specific channel
	resp, err = c.QueryDrafts(ctx, &QueryDraftsOptions{
		UserID: user.ID,
		Filter: map[string]interface{}{
			"channel_cid": channel2.CID,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(resp.Drafts))
	assert.Equal(t, channel2.CID, resp.Drafts[0].ChannelCID)
	assert.Equal(t, "Draft in channel 2", resp.Drafts[0].Message.Text)

	// Query drafts with sort
	resp, err = c.QueryDrafts(ctx, &QueryDraftsOptions{
		UserID: user.ID,
		Sort: []*SortOption{
			{Field: "created_at", Direction: 1},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, 2, len(resp.Drafts))
	assert.Equal(t, channel1.CID, resp.Drafts[0].ChannelCID)
	assert.Equal(t, channel2.CID, resp.Drafts[1].ChannelCID)

	// Query drafts with pagination
	resp, err = c.QueryDrafts(ctx, &QueryDraftsOptions{
		UserID: user.ID,
		Limit:  1,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(resp.Drafts))
	assert.Equal(t, channel2.CID, resp.Drafts[0].ChannelCID)
	assert.NotNil(t, resp.Next)

	// Query drafts with pagination using next token
	resp, err = c.QueryDrafts(ctx, &QueryDraftsOptions{
		UserID: user.ID,
		Limit:  1,
		Next:   *resp.Next,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, len(resp.Drafts))
	assert.Equal(t, channel1.CID, resp.Drafts[0].ChannelCID)
}
