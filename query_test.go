package stream_chat

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
