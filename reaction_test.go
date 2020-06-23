package stream_chat // nolint: golint

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleChannel_SendReaction() {
	channel := &Channel{}
	msgID := "123"
	userID := "bob-1"

	reaction := &Reaction{
		Type:      "love",
		ExtraData: map[string]interface{}{"my_custom_field": 123},
	}
	_, err := channel.SendReaction(reaction, msgID, userID)
	if err != nil {
		log.Fatalf("Found Error: %v", err)
	}
}

func TestChannel_SendReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer func() {
		require.NoError(t, ch.Delete(), "delete channel")
	}()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	require.NoError(t, err, "send message")

	reaction := Reaction{Type: "love"}

	msg, err = ch.SendReaction(&reaction, msg.ID, serverUser.ID)
	require.NoError(t, err, "send reaction")

	assert.Equal(t, 1, msg.ReactionCounts[reaction.Type], "reaction count", reaction)

	assert.Condition(t, reactionExistsCondition(msg.LatestReactions, reaction.Type), "latest reaction exists")
}

func reactionExistsCondition(reactions []*Reaction, searchType string) func() bool {
	return func() bool {
		for _, r := range reactions {
			if r.Type == searchType {
				return true
			}
		}
		return false
	}
}

func TestChannel_DeleteReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer func() {
		require.NoError(t, ch.Delete(), "delete channel")
	}()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	require.NoError(t, err, "send message")

	reaction := Reaction{Type: "love"}

	msg, err = ch.SendReaction(&reaction, msg.ID, serverUser.ID)
	require.NoError(t, err, "send reaction")

	msg, err = ch.DeleteReaction(msg.ID, reaction.Type, serverUser.ID)
	require.NoError(t, err, "delete reaction")

	assert.Equal(t, 0, msg.ReactionCounts[reaction.Type], "reaction count")
	assert.Empty(t, msg.LatestReactions, "latest reactions empty")
}

func TestChannel_GetReactions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer func() {
		require.NoError(t, ch.Delete(), "delete channel")
	}()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	require.NoError(t, err, "send message")

	reactions, err := ch.GetReactions(msg.ID, nil)
	require.NoError(t, err, "get reactions")
	assert.Empty(t, reactions, "reactions empty")

	reaction := Reaction{Type: "love"}

	msg, err = ch.SendReaction(&reaction, msg.ID, serverUser.ID)
	require.NoError(t, err, "send reaction")

	reactions, err = ch.GetReactions(msg.ID, nil)
	require.NoError(t, err, "get reactions")

	assert.Condition(t, reactionExistsCondition(reactions, reaction.Type), "reaction exists")
}
