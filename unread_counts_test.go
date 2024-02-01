package stream_chat

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnreadCounts(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)
	ch := initChannel(t, c, user.ID)

	ctx := context.Background()
	msg := &Message{Text: "test message"}
	randSender := randomString(5)
	for i := 0; i < 5; i++ {
		_, err := ch.SendMessage(ctx, msg, randSender)
		require.NoError(t, err)
	}

	resp, err := c.UnreadCounts(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, 5, resp.TotalUnreadCount)
	require.Equal(t, 1, len(resp.Channels))
	require.Equal(t, ch.CID, resp.Channels[0].ChannelID)
	require.Equal(t, 5, resp.Channels[0].UnreadCount)
	require.Equal(t, 1, len(resp.ChannelType))
	require.Equal(t, strings.Split(ch.CID, ":")[0], resp.ChannelType[0].ChannelType)
	require.Equal(t, 5, resp.ChannelType[0].UnreadCount)
}
func TestUnreadCountsBatch(t *testing.T) {
	t.Skip()
	c := initClient(t)
	user1 := randomUser(t, c)
	user2 := randomUser(t, c)
	ch := initChannel(t, c, user1.ID, user2.ID)

	ctx := context.Background()
	msg := &Message{Text: "test message"}
	randSender := randomString(5)
	for i := 0; i < 5; i++ {
		_, err := ch.SendMessage(ctx, msg, randSender)
		require.NoError(t, err)
	}

	nonexistant := randomString(5)
	resp, err := c.UnreadCountsBatch(ctx, []string{user1.ID, user2.ID, nonexistant})
	require.NoError(t, err)
	require.Equal(t, 2, len(resp.CountsByUser))
	require.Contains(t, resp.CountsByUser, user1.ID)
	require.Contains(t, resp.CountsByUser, user2.ID)

	// user 1 counts
	require.Equal(t, 5, resp.CountsByUser[user1.ID].TotalUnreadCount)
	require.Equal(t, 1, len(resp.CountsByUser[user1.ID].Channels))
	require.Equal(t, ch.CID, resp.CountsByUser[user1.ID].Channels[0].ChannelID)
	require.Equal(t, 5, resp.CountsByUser[user1.ID].Channels[0].UnreadCount)
	require.Equal(t, 1, len(resp.CountsByUser[user1.ID].ChannelType))
	require.Equal(t, strings.Split(ch.CID, ":")[0], resp.CountsByUser[user1.ID].ChannelType[0].ChannelType)
	require.Equal(t, 5, resp.CountsByUser[user1.ID].ChannelType[0].UnreadCount)

	// user 2 counts
	require.Equal(t, 5, resp.CountsByUser[user2.ID].TotalUnreadCount)
	require.Equal(t, 1, len(resp.CountsByUser[user2.ID].Channels))
	require.Equal(t, ch.CID, resp.CountsByUser[user2.ID].Channels[0].ChannelID)
	require.Equal(t, 5, resp.CountsByUser[user2.ID].Channels[0].UnreadCount)
	require.Equal(t, 1, len(resp.CountsByUser[user2.ID].ChannelType))
	require.Equal(t, strings.Split(ch.CID, ":")[0], resp.CountsByUser[user2.ID].ChannelType[0].ChannelType)
	require.Equal(t, 5, resp.CountsByUser[user2.ID].ChannelType[0].UnreadCount)
}
