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
	ch := initChannel(t, c, user)

	ctx := context.Background()
	msg := &Message{Text: "test message"}
	randSender := randomString(5)
	var messageID string
	for i := 0; i < 5; i++ {
		resp, err := ch.SendMessage(ctx, msg, randSender)
		require.NoError(t, err)
		messageID = resp.Message.ID
	}

	resp, err := c.UnreadCounts(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, 5, resp.TotalUnreadCount)
	require.Len(t, resp.Channels, 1)
	require.Equal(t, ch.CID, resp.Channels[0].ChannelID)
	require.Equal(t, 5, resp.Channels[0].UnreadCount)
	require.Len(t, resp.ChannelType, 1)
	require.Equal(t, strings.Split(ch.CID, ":")[0], resp.ChannelType[0].ChannelType)
	require.Equal(t, 5, resp.ChannelType[0].UnreadCount)

	// test unread threads
	threadMsg := &Message{Text: "test thread", ParentID: messageID}
	_, err = ch.SendMessage(ctx, threadMsg, user.ID)
	require.NoError(t, err)
	_, err = ch.SendMessage(ctx, threadMsg, randSender)
	require.NoError(t, err)

	resp, err = c.UnreadCounts(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, 1, resp.TotalUnreadThreadsCount)
	require.Len(t, resp.Threads, 1)
	require.Equal(t, messageID, resp.Threads[0].ParentMessageID)
}

func TestUnreadCountsBatch(t *testing.T) {
	c := initClient(t)
	user1 := randomUser(t, c)
	user2 := randomUser(t, c)
	ch := initChannel(t, c, user1, user2)

	ctx := context.Background()
	msg := &Message{Text: "test message"}
	randSender := randomString(5)
	var messageID string
	for i := 0; i < 5; i++ {
		resp, err := ch.SendMessage(ctx, msg, randSender)
		require.NoError(t, err)
		messageID = resp.Message.ID
	}

	nonexistant := randomString(5)
	resp, err := c.UnreadCountsBatch(ctx, []string{user1.ID, user2.ID, nonexistant})
	require.NoError(t, err)
	require.Len(t, resp.CountsByUser, 2)
	require.Contains(t, resp.CountsByUser, user1.ID)
	require.Contains(t, resp.CountsByUser, user2.ID)
	require.NotContains(t, resp.CountsByUser, nonexistant)

	// user 1 counts
	require.Equal(t, 5, resp.CountsByUser[user1.ID].TotalUnreadCount)
	require.Len(t, resp.CountsByUser[user1.ID].Channels, 1)
	require.Equal(t, ch.CID, resp.CountsByUser[user1.ID].Channels[0].ChannelID)
	require.Equal(t, 5, resp.CountsByUser[user1.ID].Channels[0].UnreadCount)
	require.Len(t, resp.CountsByUser[user1.ID].ChannelType, 1)
	require.Equal(t, strings.Split(ch.CID, ":")[0], resp.CountsByUser[user1.ID].ChannelType[0].ChannelType)
	require.Equal(t, 5, resp.CountsByUser[user1.ID].ChannelType[0].UnreadCount)

	// user 2 counts
	require.Equal(t, 5, resp.CountsByUser[user2.ID].TotalUnreadCount)
	require.Len(t, resp.CountsByUser[user2.ID].Channels, 1)
	require.Equal(t, ch.CID, resp.CountsByUser[user2.ID].Channels[0].ChannelID)
	require.Equal(t, 5, resp.CountsByUser[user2.ID].Channels[0].UnreadCount)
	require.Len(t, resp.CountsByUser[user2.ID].ChannelType, 1)
	require.Equal(t, strings.Split(ch.CID, ":")[0], resp.CountsByUser[user2.ID].ChannelType[0].ChannelType)
	require.Equal(t, 5, resp.CountsByUser[user2.ID].ChannelType[0].UnreadCount)

	// test unread threads
	threadMsg := &Message{Text: "test thread", ParentID: messageID}
	_, err = ch.SendMessage(ctx, threadMsg, user1.ID)
	require.NoError(t, err)
	_, err = ch.SendMessage(ctx, threadMsg, user2.ID)
	require.NoError(t, err)
	_, err = ch.SendMessage(ctx, threadMsg, randSender)
	require.NoError(t, err)

	resp, err = c.UnreadCountsBatch(ctx, []string{user1.ID, user2.ID, nonexistant})
	require.NoError(t, err)

	// user 1 thread counts
	require.Equal(t, 1, resp.CountsByUser[user1.ID].TotalUnreadThreadsCount)
	require.Len(t, resp.CountsByUser[user1.ID].Threads, 1)
	require.Equal(t, messageID, resp.CountsByUser[user1.ID].Threads[0].ParentMessageID)

	// user 2 thread counts
	require.Equal(t, resp.CountsByUser[user2.ID].TotalUnreadThreadsCount, 1)
	require.Len(t, resp.CountsByUser[user2.ID].Threads, 1)
	require.Equal(t, messageID, resp.CountsByUser[user2.ID].Threads[0].ParentMessageID)
}
