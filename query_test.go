package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_QueryUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	users, err := c.QueryUsers(&QueryOption{Filter: map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": user.ID,
		}},
	})

	mustNoError(t, err, "query users error")

	if assert.NotEmpty(t, users, "query users exists") {
		assert.Equal(t, user.ID, users[0].ID, "received user ID")
	}
}

func TestClient_QueryChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	got, err := c.QueryChannels(&QueryOption{Filter: map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": ch.ID,
		},
	}})

	mustNoError(t, err, "query channels error")

	if assert.NotEmpty(t, got, "query channels exists") {
		assert.Equal(t, ch.ID, got[0].ID, "received channel ID")
	}
}

func TestClient_Search(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	user1, user2 := randomUser(), randomUser()

	text := randomString(10)

	_, err := ch.SendMessage(&Message{
		Text:      text + " " + randomString(25),
		ExtraData: map[string]interface{}{"color": "green"},
	}, user1.ID)
	mustNoError(t, err)

	_, err = ch.SendMessage(&Message{
		Text:      text + " " + randomString(25),
		ExtraData: map[string]interface{}{"color": "red"},
	}, user2.ID)
	mustNoError(t, err)

	t.Run("search message text", func(t *testing.T) {
		got, err := c.Search(SearchRequest{Query: text, Filters: map[string]interface{}{
			"members": map[string][]string{
				"$in": {user1.ID, user2.ID},
			},
		}})

		mustNoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("search text and custom props via message filters", func(t *testing.T) {
		got, err := c.Search(SearchRequest{
			Filters: map[string]interface{}{
				"members": map[string][]string{
					"$in": {user1.ID, user2.ID},
				}},
			MessageFilters: map[string]interface{}{
				"text":  map[string]interface{}{"$q": text},
				"color": "green",
			},
		})

		mustNoError(t, err)
		assert.Len(t, got, 1)
	})
}
