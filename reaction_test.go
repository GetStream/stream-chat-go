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
	msg := Message{
		Text: "test message",
		User: user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	reaction := Reaction{Type: "love"}

	err = ch.SendReaction(&msg, &reaction, serverUser.ID)
	mustNoError(t, err)

	assert.Equal(t, 1, msg.ReactionCounts[reaction.Type], "reaction count", reaction)
	assert.Contains(t, msg.LatestReactions, reaction, "latest reactions exists")
}

func TestChannel_DeleteReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	reaction := Reaction{Type: "love"}

	err = ch.SendReaction(&msg, &reaction, serverUser.ID)
	mustNoError(t, err)

	err = ch.DeleteReaction(&msg, reaction.Type, serverUser.ID)
	mustNoError(t, err)

	assert.Equal(t, 0, msg.ReactionCounts[reaction.Type], "reaction count")
	assert.Empty(t, msg.LatestReactions, "latest reactions empty")
}

func TestChannel_GetReactions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	reactions, err := ch.GetReactions(msg.ID, nil)
	mustNoError(t, err)
	assert.Empty(t, reactions, "reactions empty")

	reaction := Reaction{Type: "love"}

	err = ch.SendReaction(&msg, &reaction, serverUser.ID)
	mustNoError(t, err)

	reactions, err = ch.GetReactions(msg.ID, nil)
	mustNoError(t, err)

	assert.Contains(t, reactions, reaction, "reaction exists")
}
