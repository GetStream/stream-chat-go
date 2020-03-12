package stream

import (
	"testing"

	"log"
	"os"
	"time"

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

	_, err := ch.SendMessage(&Message{Text: text + " " + randomString(25)}, user1.ID)
	mustNoError(t, err)

	_, err = ch.SendMessage(&Message{Text: text + " " + randomString(25)}, user2.ID)
	mustNoError(t, err)

	got, err := c.Search(SearchRequest{Query: text, Filters: map[string]interface{}{
		"members": map[string][]string{
			"$in": {user1.ID, user2.ID},
		},
	}})

	mustNoError(t, err)

	assert.Len(t, got, 2)
}

func ExampleQueryOption() {
	client, err := NewClient(
		os.Getenv("STREAM_CHAT_API_KEY"),
		os.Getenv("STREAM_CHAT_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	prefix := time.Now().UTC().Format("2006-01-02-15-04-")
	chName := prefix + "awesome_sauce"

	ch, err := client.CreateChannel(
		ChannelTypeLabelTeam,
		chName,
		"calvin",
		nil,
	)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	ch.SendMessage(&Message{
		Text: "Hello Calvin",
		User: &User{
			ID: "hobbes",
		},
	}, "hobbes")

	ch.SendMessage(&Message{
		Text: "Hello Hobbes",
		User: &User{
			ID: "calvin",
		},
	}, "calvin")

	sort := []*SortOption{
		{Field: SortFieldUserID, Direction: SortAscending},
	}

	users, err := client.QueryUsers(&QueryOption{}, sort...)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	for _, user := range users {
		log.Printf("%v", user.ID)
	}
}
