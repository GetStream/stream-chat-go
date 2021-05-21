package stream_chat // nolint: golint

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_QueryUsers(t *testing.T) {
	c := initClient(t)

	const n = 4
	ids := make([]string, n)
	defer func() {
		for _, id := range ids {
			if id != "" {
				_ = c.DeleteUser(id, nil)
			}
		}
	}()

	for i := n - 1; i > -1; i-- {
		u := &User{ID: randomString(30), ExtraData: map[string]interface{}{"order": n - i - 1}}
		_, err := c.UpsertUser(u)
		require.NoError(t, err)
		ids[i] = u.ID
		time.Sleep(200 * time.Millisecond)
	}

	t.Parallel()
	t.Run("Query all", func(tt *testing.T) {
		results, err := c.QueryUsers(&QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]interface{}{
					"$in": ids,
				},
			},
		})

		require.NoError(tt, err)
		require.Len(tt, results, len(ids))
	})

	t.Run("Query with offset/limit", func(tt *testing.T) {
		offset := 1

		results, err := c.QueryUsers(
			&QueryOption{
				Filter: map[string]interface{}{
					"id": map[string]interface{}{
						"$in": ids,
					},
				},
				Offset: offset,
				Limit:  2,
			},
		)

		require.NoError(tt, err)
		require.Len(tt, results, 2)

		require.Equal(tt, results[0].ID, ids[offset])
		require.Equal(tt, results[1].ID, ids[offset+1])
	})
}

func TestClient_QueryChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	_, err := ch.SendMessage(&Message{Text: "abc"}, "some")
	require.NoError(t, err)
	_, err = ch.SendMessage(&Message{Text: "abc"}, "some")
	require.NoError(t, err)

	messageLimit := 1
	got, err := c.QueryChannels(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]interface{}{
				"$eq": ch.ID,
			},
		},
		MessageLimit: &messageLimit,
	})

	require.NoError(t, err, "query channels error")
	require.Equal(t, ch.ID, got[0].ID, "received channel ID")
	require.Len(t, got[0].Messages, messageLimit)
}

func TestClient_Search(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user1, user2 := randomUser(), randomUser()

	text := randomString(10)

	_, err := ch.SendMessage(&Message{Text: text + " " + randomString(25)}, user1.ID)
	require.NoError(t, err)

	_, err = ch.SendMessage(&Message{Text: text + " " + randomString(25)}, user2.ID)
	require.NoError(t, err)

	got, err := c.Search(SearchRequest{Query: text, Filters: map[string]interface{}{
		"members": map[string][]string{
			"$in": {user1.ID, user2.ID},
		},
	}})

	require.NoError(t, err)

	assert.Len(t, got, 2)
}

func TestClient_QueryMessageFlags(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user1, user2 := randomUser(), randomUser()
	for user1.ID == user2.ID {
		user2 = randomUser()
	}

	// send 2 messages
	text := randomString(10)
	msg1, err := ch.SendMessage(&Message{Text: text + " " + randomString(25)}, user1.ID)
	require.NoError(t, err)
	msg2, err := ch.SendMessage(&Message{Text: text + " " + randomString(25)}, user2.ID)
	require.NoError(t, err)

	// flag 2 messages
	err = c.FlagMessage(msg2.ID, user1.ID)
	require.NoError(t, err)

	err = c.FlagMessage(msg1.ID, user2.ID)
	require.NoError(t, err)

	// both flags show up in this query by channel_cid
	got, err := c.QueryMessageFlags(&QueryOption{
		Filter: map[string]interface{}{
			"channel_cid": map[string][]string{
				"$in": {ch.cid()},
			},
		},
	})
	require.NoError(t, err)
	assert.Len(t, got, 2)

	// one flag shows up in this query by user_id
	got, err = c.QueryMessageFlags(&QueryOption{
		Filter: map[string]interface{}{
			"user_id": user1.ID,
		},
	})
	require.NoError(t, err)
	assert.Len(t, got, 1)

	// unflag these 2 messages
	err = c.UnflagMessage(msg1.ID, user2.ID)
	require.NoError(t, err)
	err = c.UnflagMessage(msg2.ID, user1.ID)
	require.NoError(t, err)

	// none should show up
	got, err = c.QueryMessageFlags(&QueryOption{
		Filter: map[string]interface{}{"channel_cid": ch.cid()},
	})
	require.NoError(t, err)
	assert.Len(t, got, 0)
}
