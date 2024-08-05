package stream_chat

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func ExampleClient_SendReaction() {
	client := &Client{}
	msgID := "123"
	userID := "bob-1"
	ctx := context.Background()

	reaction := &Reaction{
		Type:      "love",
		ExtraData: map[string]interface{}{"my_custom_field": 123},
	}
	_, err := client.SendReaction(ctx, reaction, msgID, userID)
	if err != nil {
		log.Fatalf("Found Error: %v", err)
	}
}

func TestChannel_SendReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	user := randomUser(t, c)
	ctx := context.Background()
	msg := &Message{
		Text: "test message",
		User: user,
	}

	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err, "send message")

	reaction := Reaction{Type: "love"}
	reactionResp, err := c.SendReaction(ctx, &reaction, resp.Message.ID, user.ID)
	require.NoError(t, err, "send reaction")

	require.Equal(t, 1, reactionResp.Message.ReactionCounts[reaction.Type], "reaction count", reaction)

	require.Condition(t, reactionExistsCondition(reactionResp.Message.LatestReactions, reaction.Type), "latest reaction exists")
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

func TestClient_DeleteReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	user := randomUser(t, c)
	ctx := context.Background()
	msg := &Message{
		Text: "test message",
		User: user,
	}

	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err, "send message")

	reaction := Reaction{Type: "love"}
	reactionResp, err := c.SendReaction(ctx, &reaction, resp.Message.ID, user.ID)
	require.NoError(t, err, "send reaction")

	reactionResp, err = c.DeleteReaction(ctx, reactionResp.Message.ID, reaction.Type, user.ID)
	require.NoError(t, err, "delete reaction")

	require.Equal(t, 0, reactionResp.Message.ReactionCounts[reaction.Type], "reaction count")
	require.Empty(t, reactionResp.Message.LatestReactions, "latest reactions empty")
}

func TestClient_GetReactions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	user := randomUser(t, c)
	ctx := context.Background()
	msg := &Message{
		Text: "test message",
		User: user,
	}

	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err, "send message")
	msg = resp.Message

	reactionsResp, err := c.GetReactions(ctx, msg.ID, nil)
	require.NoError(t, err, "get reactions")
	require.Empty(t, reactionsResp.Reactions, "reactions empty")

	reaction := Reaction{Type: "love"}

	reactionResp, err := c.SendReaction(ctx, &reaction, msg.ID, user.ID)
	require.NoError(t, err, "send reaction")

	reactionsResp, err = c.GetReactions(ctx, reactionResp.Message.ID, nil)
	require.NoError(t, err, "get reactions")

	require.Condition(t, reactionExistsCondition(reactionsResp.Reactions, reaction.Type), "reaction exists")
}
