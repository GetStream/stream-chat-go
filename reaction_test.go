package stream_chat // nolint: golint

import (
	"context"
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
	_, err := channel.SendReaction(context.Background(), reaction, msgID, userID)
	if err != nil {
		log.Fatalf("Found Error: %v", err)
	}
}

func TestChannel_SendReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer func() {
		_, err := ch.Delete(context.Background())
		require.NoError(t, err, "delete channel")
	}()

	user := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		User: user,
	}
	resp, err := ch.SendMessage(context.Background(), msg, user.ID)
	require.NoError(t, err, "send message")

	reaction := Reaction{Type: "love"}
	reactionResp, err := ch.SendReaction(context.Background(), &reaction, resp.Message.ID, user.ID)
	require.NoError(t, err, "send reaction")

	assert.Equal(t, 1, reactionResp.Message.ReactionCounts[reaction.Type], "reaction count", reaction)

	assert.Condition(t, reactionExistsCondition(reactionResp.Message.LatestReactions, reaction.Type), "latest reaction exists")
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
		_, err := ch.Delete(context.Background())
		require.NoError(t, err, "delete channel")
	}()

	user := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		User: user,
	}
	resp, err := ch.SendMessage(context.Background(), msg, user.ID)
	require.NoError(t, err, "send message")

	reaction := Reaction{Type: "love"}
	reactionResp, err := ch.SendReaction(context.Background(), &reaction, resp.Message.ID, user.ID)
	require.NoError(t, err, "send reaction")

	reactionResp, err = ch.DeleteReaction(context.Background(), reactionResp.Message.ID, reaction.Type, user.ID)
	require.NoError(t, err, "delete reaction")

	assert.Equal(t, 0, reactionResp.Message.ReactionCounts[reaction.Type], "reaction count")
	assert.Empty(t, reactionResp.Message.LatestReactions, "latest reactions empty")
}

func TestChannel_GetReactions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer func() {
		_, err := ch.Delete(context.Background())
		require.NoError(t, err, "delete channel")
	}()

	user := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		User: user,
	}
	resp, err := ch.SendMessage(context.Background(), msg, user.ID)
	require.NoError(t, err, "send message")
	msg = resp.Message

	reactions, err := ch.GetReactions(context.Background(), msg.ID, nil)
	require.NoError(t, err, "get reactions")
	assert.Empty(t, reactions, "reactions empty")

	reaction := Reaction{Type: "love"}

	reactionResp, err := ch.SendReaction(context.Background(), &reaction, msg.ID, user.ID)
	require.NoError(t, err, "send reaction")

	reactionsResp, err := ch.GetReactions(context.Background(), reactionResp.Message.ID, nil)
	require.NoError(t, err, "get reactions")

	assert.Condition(t, reactionExistsCondition(reactionsResp.Reactions, reaction.Type), "reaction exists")
}
