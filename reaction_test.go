package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannel_SendReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	mustNoError(t, err, "send message")

	reaction := Reaction{Type: "love"}

	msg, err = ch.SendReaction(&reaction, msg.ID, serverUser.ID)
	mustNoError(t, err, "send reaction")

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
	defer ch.Delete()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	mustNoError(t, err, "send message")

	reaction := Reaction{Type: "love"}

	msg, err = ch.SendReaction(&reaction, msg.ID, serverUser.ID)
	mustNoError(t, err, "send reaction")

	msg, err = ch.DeleteReaction(msg.ID, reaction.Type, serverUser.ID)
	mustNoError(t, err, "delete reaction")

	assert.Equal(t, 0, msg.ReactionCounts[reaction.Type], "reaction count")
	assert.Empty(t, msg.LatestReactions, "latest reactions empty")
}

func TestChannel_GetReactions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	mustNoError(t, err, "send message")

	reactions, err := ch.GetReactions(msg.ID, nil)
	mustNoError(t, err, "get reactions")
	assert.Empty(t, reactions, "reactions empty")

	reaction := Reaction{Type: "love"}

	msg, err = ch.SendReaction(&reaction, msg.ID, serverUser.ID)
	mustNoError(t, err, "send reaction")

	reactions, err = ch.GetReactions(msg.ID, nil)
	mustNoError(t, err, "get reactions")

	assert.Condition(t, reactionExistsCondition(reactions, reaction.Type), "reaction exists")
}
