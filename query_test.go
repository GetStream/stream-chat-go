package stream_chat

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_QueryUsers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	const n = 5
	ids := make([]string, n)
	t.Cleanup(func() {
		for _, id := range ids {
			if id != "" {
				_, _ = c.DeleteUser(ctx, id)
			}
		}
	})

	for i := n - 1; i > -1; i-- {
		u := &User{ID: randomString(30), ExtraData: map[string]interface{}{"order": n - i - 1}}
		_, err := c.UpsertUser(ctx, u)
		require.NoError(t, err)
		ids[i] = u.ID
		time.Sleep(200 * time.Millisecond)
	}

	_, err := c.DeactivateUser(ctx, ids[n-1])
	require.NoError(t, err)

	t.Parallel()
	t.Run("Query all", func(tt *testing.T) {
		results, err := c.QueryUsers(ctx, &QueryUsersOptions{
			QueryOption: QueryOption{
				Filter: map[string]interface{}{
					"id": map[string]interface{}{
						"$in": ids,
					},
				},
			},
		})

		require.NoError(tt, err)
		require.Len(tt, results.Users, len(ids)-1)
	})

	t.Run("Query with offset/limit", func(tt *testing.T) {
		offset := 1

		results, err := c.QueryUsers(ctx, &QueryUsersOptions{
			QueryOption: QueryOption{
				Filter: map[string]interface{}{
					"id": map[string]interface{}{
						"$in": ids,
					},
				},
				Offset: offset,
				Limit:  2,
			},
		})

		require.NoError(tt, err)
		require.Len(tt, results.Users, 2)

		require.Equal(tt, results.Users[0].ID, ids[offset])
		require.Equal(tt, results.Users[1].ID, ids[offset+1])
	})

	t.Run("Query with deactivated", func(tt *testing.T) {
		results, err := c.QueryUsers(ctx, &QueryUsersOptions{
			QueryOption: QueryOption{
				Filter: map[string]interface{}{
					"id": map[string]interface{}{
						"$in": ids,
					},
				},
			},
			IncludeDeactivatedUsers: true,
		})

		require.NoError(tt, err)
		require.Len(tt, results.Users, len(ids))
	})
}

func TestClient_QueryChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	_, err := ch.SendMessage(ctx, &Message{Text: "abc"}, "some")
	require.NoError(t, err)
	_, err = ch.SendMessage(ctx, &Message{Text: "abc"}, "some")
	require.NoError(t, err)

	messageLimit := 1
	resp, err := c.QueryChannels(ctx, &QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]interface{}{
				"$eq": ch.ID,
			},
		},
		MessageLimit: &messageLimit,
	})

	require.NoError(t, err, "query channels error")
	require.Equal(t, ch.ID, resp.Channels[0].ID, "received channel ID")
	require.Len(t, resp.Channels[0].Messages, messageLimit)
}

func TestClient_Search(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	user1, user2 := randomUser(t, c), randomUser(t, c)

	ch := initChannel(t, c, user1.ID, user2.ID)

	text := randomString(10)

	_, err := ch.SendMessage(ctx, &Message{Text: text + " " + randomString(25)}, user1.ID)
	require.NoError(t, err)

	_, err = ch.SendMessage(ctx, &Message{Text: text + " " + randomString(25)}, user2.ID)
	require.NoError(t, err)

	t.Run("Query", func(tt *testing.T) {
		resp, err := c.Search(ctx, SearchRequest{Query: text, Filters: map[string]interface{}{
			"members": map[string][]string{
				"$in": {user1.ID, user2.ID},
			},
		}})

		require.NoError(tt, err)

		assert.Len(tt, resp.Messages, 2)
	})
	t.Run("Message filters", func(tt *testing.T) {
		resp, err := c.Search(ctx, SearchRequest{
			Filters: map[string]interface{}{
				"members": map[string][]string{
					"$in": {user1.ID, user2.ID},
				},
			},
			MessageFilters: map[string]interface{}{
				"text": map[string]interface{}{
					"$q": text,
				},
			},
		})
		require.NoError(tt, err)

		assert.Len(tt, resp.Messages, 2)
	})
	t.Run("Query and message filters error", func(tt *testing.T) {
		_, err := c.Search(ctx, SearchRequest{
			Filters: map[string]interface{}{
				"members": map[string][]string{
					"$in": {user1.ID, user2.ID},
				},
			},
			MessageFilters: map[string]interface{}{
				"text": map[string]interface{}{
					"$q": text,
				},
			},
			Query: text,
		})
		require.Error(tt, err)
	})
	t.Run("Offset and sort error", func(tt *testing.T) {
		_, err := c.Search(ctx, SearchRequest{
			Filters: map[string]interface{}{
				"members": map[string][]string{
					"$in": {user1.ID, user2.ID},
				},
			},
			Offset: 1,
			Query:  text,
			Sort: []SortOption{{
				Field:     "created_at",
				Direction: -1,
			}},
		})
		require.Error(tt, err)
	})
	t.Run("Offset and next error", func(tt *testing.T) {
		_, err := c.Search(ctx, SearchRequest{
			Filters: map[string]interface{}{
				"members": map[string][]string{
					"$in": {user1.ID, user2.ID},
				},
			},
			Offset: 1,
			Query:  text,
			Next:   randomString(5),
		})
		require.Error(tt, err)
	})
}

func TestClient_SearchWithFullResponse(t *testing.T) {
	t.Skip()
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	user1, user2 := randomUser(t, c), randomUser(t, c)

	text := randomString(10)

	messageIDs := make([]string, 6)
	for i := 0; i < 6; i++ {
		userID := user1.ID
		if i%2 == 0 {
			userID = user2.ID
		}
		messageID := fmt.Sprintf("%d-%s", i, text)
		_, err := ch.SendMessage(ctx, &Message{
			ID:   messageID,
			Text: text + " " + randomString(25),
		}, userID)
		require.NoError(t, err)

		messageIDs[6-i] = messageID
	}

	got, err := c.SearchWithFullResponse(ctx, SearchRequest{
		Query: text,
		Filters: map[string]interface{}{
			"members": map[string][]string{
				"$in": {user1.ID, user2.ID},
			},
		},
		Sort: []SortOption{
			{Field: "created_at", Direction: -1},
		},
		Limit: 3,
	})

	gotMessageIDs := make([]string, 0, 6)
	require.NoError(t, err)
	assert.NotEmpty(t, got.Next)
	assert.Len(t, got.Results, 3)
	for _, result := range got.Results {
		gotMessageIDs = append(gotMessageIDs, result.Message.ID)
	}
	got, err = c.SearchWithFullResponse(ctx, SearchRequest{
		Query: text,
		Filters: map[string]interface{}{
			"members": map[string][]string{
				"$in": {user1.ID, user2.ID},
			},
		},
		Next:  got.Next,
		Limit: 3,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, got.Previous)
	assert.Empty(t, got.Next)
	assert.Len(t, got.Results, 3)
	for _, result := range got.Results {
		gotMessageIDs = append(gotMessageIDs, result.Message.ID)
	}
	assert.Equal(t, messageIDs, gotMessageIDs)
}

func TestClient_QueryMessageFlags(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	user1, user2 := randomUser(t, c), randomUser(t, c)
	for user1.ID == user2.ID {
		user2 = randomUser(t, c)
	}

	// send 2 messages
	text := randomString(10)
	resp, err := ch.SendMessage(ctx, &Message{Text: text + " " + randomString(25)}, user1.ID)
	require.NoError(t, err)
	msg1 := resp.Message

	resp, err = ch.SendMessage(ctx, &Message{Text: text + " " + randomString(25)}, user2.ID)
	require.NoError(t, err)
	msg2 := resp.Message

	// flag 2 messages
	_, err = c.FlagMessage(ctx, msg2.ID, user1.ID)
	require.NoError(t, err)

	_, err = c.FlagMessage(ctx, msg1.ID, user2.ID)
	require.NoError(t, err)

	// both flags show up in this query by channel_cid
	got, err := c.QueryMessageFlags(ctx, &QueryOption{
		Filter: map[string]interface{}{
			"channel_cid": map[string][]string{
				"$in": {ch.cid()},
			},
		},
	})
	require.NoError(t, err)
	assert.Len(t, got.Flags, 2)

	// one flag shows up in this query by user_id
	got, err = c.QueryMessageFlags(ctx, &QueryOption{
		Filter: map[string]interface{}{
			"user_id": user1.ID,
		},
	})
	require.NoError(t, err)
	assert.Len(t, got.Flags, 1)
}

func TestQueryChannelsResponse_ParsedPredefinedFilter(t *testing.T) {
	// Test that ParsedPredefinedFilterResponse is correctly unmarshaled from JSON
	jsonData := `{
		"channels": [],
		"predefined_filter": {
			"name": "user_messaging",
			"filter": {"type": "messaging", "members": {"$in": ["user123"]}},
			"sort": [{"field": "last_message_at", "direction": -1}]
		},
		"duration": "0.01s"
	}`

	var resp queryChannelResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	require.NotNil(t, resp.PredefinedFilter)
	require.Equal(t, "user_messaging", resp.PredefinedFilter.Name)
	require.NotNil(t, resp.PredefinedFilter.Filter)
	require.Equal(t, "messaging", resp.PredefinedFilter.Filter["type"])
	require.Len(t, resp.PredefinedFilter.Sort, 1)
	require.Equal(t, "last_message_at", resp.PredefinedFilter.Sort[0].Field)
	require.Equal(t, -1, resp.PredefinedFilter.Sort[0].Direction)
}

func TestQueryChannelsResponse_NoPredefinedFilter(t *testing.T) {
	// Test that response without predefined_filter has nil PredefinedFilter
	jsonData := `{
		"channels": [],
		"duration": "0.01s"
	}`

	var resp queryChannelResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	require.Nil(t, resp.PredefinedFilter)
}

func TestClient_QueryFlagReportsAndReview(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	user1, user2 := randomUser(t, c), randomUser(t, c)
	msg, err := ch.SendMessage(ctx, &Message{Text: randomString(25)}, user1.ID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = ch.Delete(ctx)
		_, _ = c.DeleteUser(ctx, user1.ID, DeleteUserWithHardDelete())
		_, _ = c.DeleteUser(ctx, user2.ID, DeleteUserWithHardDelete())
	})

	_, err = c.FlagMessage(ctx, msg.Message.ID, user1.ID)
	require.NoError(t, err)

	resp, err := c.QueryFlagReports(ctx, &QueryFlagReportsRequest{
		FilterConditions: map[string]interface{}{"message_id": msg.Message.ID},
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.FlagReports)

	flagResp, err := c.ReviewFlagReport(ctx, resp.FlagReports[0].ID, &ReviewFlagReportRequest{
		ReviewResult: "reviewed",
		UserID:       user2.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, flagResp.FlagReport)
}
