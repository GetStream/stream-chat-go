package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_QueryUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	users, err := c.QueryUsers(&QueryOption{Query: map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": user.ID,
		},
	}})

	mustNoError(t, err)
	assert.Equal(t, user.ID, users[0].ID)
}

func TestClient_QueryChannels(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	got, err := c.QueryChannels(&QueryOption{Query: map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": ch.ID,
		},
	}})

	mustNoError(t, err)
	assert.Equal(t, ch.ID, got[0].ID)
}
