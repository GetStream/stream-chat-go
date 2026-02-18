package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShadowBanUser(t *testing.T) {
	c := initClient(t)
	userA := randomUser(t, c)
	userB := randomUser(t, c)
	userC := randomUser(t, c)
	ctx := context.Background()

	ch := initChannel(t, c, userA.ID, userB.ID, userC.ID)
	resp, err := c.CreateChannel(ctx, ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	ch = resp.Channel

	// shadow ban userB globally
	_, err = c.ShadowBan(ctx, userB.ID, userA.ID)
	require.NoError(t, err)

	// shadow ban userC on channel
	_, err = ch.ShadowBan(ctx, userC.ID, userA.ID)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	messageResp, err := ch.SendMessage(ctx, msg, userB.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.False(t, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.True(t, messageResp.Message.Shadowed)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(ctx, msg, userC.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.False(t, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.True(t, messageResp.Message.Shadowed)

	_, err = c.UnBanUser(ctx, userB.ID)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(ctx, msg, userB.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.False(t, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.False(t, messageResp.Message.Shadowed)

	_, err = ch.UnBanUser(ctx, userC.ID)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(ctx, msg, userC.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.False(t, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.False(t, messageResp.Message.Shadowed)
}

func TestBanUnbanUser(t *testing.T) {
	c := initClient(t)
	target := randomUser(t, c)
	user := randomUser(t, c)
	ctx := context.Background()

	_, err := c.BanUser(ctx, target.ID, user.ID, BanWithReason("spammer"), BanWithExpiration(60))
	require.NoError(t, err)

	resp, err := c.QueryBannedUsers(ctx, &QueryBannedUsersOptions{
		QueryOption: &QueryOption{Filter: map[string]interface{}{
			"user_id": map[string]string{"$eq": target.ID},
		}},
	})
	require.NoError(t, err)
	require.Equal(t, "spammer", resp.Bans[0].Reason)
	require.NotZero(t, resp.Bans[0].Expires)

	_, err = c.UnBanUser(ctx, target.ID)
	require.NoError(t, err)

	resp, err = c.QueryBannedUsers(ctx, &QueryBannedUsersOptions{
		QueryOption: &QueryOption{Filter: map[string]interface{}{
			"user_id": map[string]string{"$eq": target.ID},
		}},
	})
	require.NoError(t, err)
	require.Empty(t, resp.Bans)
}

func TestChannelBanUnban(t *testing.T) {
	c := initClient(t)
	target := randomUser(t, c)
	user := randomUser(t, c)
	ch := initChannel(t, c, user.ID, target.ID)
	ctx := context.Background()

	_, err := ch.BanUser(ctx, target.ID, user.ID, BanWithReason("spammer"), BanWithExpiration(60))
	require.NoError(t, err)

	_, err = ch.UnBanUser(ctx, target.ID)
	require.NoError(t, err)

	resp, err := c.QueryBannedUsers(ctx, &QueryBannedUsersOptions{
		QueryOption: &QueryOption{Filter: map[string]interface{}{
			"channel_cid": map[string]string{"$eq": ch.CID},
		}},
	})
	require.NoError(t, err)
	require.Empty(t, resp.Bans)
}

func TestQueryFutureChannelBans(t *testing.T) {
	c := initClient(t)
	creator := randomUser(t, c)
	target1 := randomUser(t, c)
	target2 := randomUser(t, c)
	ctx := context.Background()

	// Create a channel to use for future channel bans
	ch := initChannel(t, c, creator.ID)

	// Ban both targets from future channels created by creator
	_, err := c.BanUser(ctx, target1.ID, creator.ID, BanWithBanFromFutureChannels(), BanWithChannel(ch.Type, ch.ID), BanWithReason("test ban 1"))
	require.NoError(t, err)

	_, err = c.BanUser(ctx, target2.ID, creator.ID, BanWithBanFromFutureChannels(), BanWithChannel(ch.Type, ch.ID), BanWithReason("test ban 2"))
	require.NoError(t, err)

	// Query all future channel bans by creator
	resp, err := c.QueryFutureChannelBans(ctx, &QueryFutureChannelBansOptions{
		UserID: creator.ID,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(resp.Bans), 2)

	// Query with target_user_id filter - should only return the specific target
	// Note: When filtering by target_user_id, the API doesn't return the user object
	// since it's already known from the filter
	resp, err = c.QueryFutureChannelBans(ctx, &QueryFutureChannelBansOptions{
		UserID:       creator.ID,
		TargetUserID: target1.ID,
	})
	require.NoError(t, err)
	require.Len(t, resp.Bans, 1)
	require.Equal(t, "test ban 1", resp.Bans[0].Reason)

	// Query for the other target
	resp, err = c.QueryFutureChannelBans(ctx, &QueryFutureChannelBansOptions{
		UserID:       creator.ID,
		TargetUserID: target2.ID,
	})
	require.NoError(t, err)
	require.Len(t, resp.Bans, 1)
	require.Equal(t, "test ban 2", resp.Bans[0].Reason)

	// Cleanup - unban both users with RemoveFutureChannelsBan
	_, err = c.UnBanUser(ctx, target1.ID, UnbanWithRemoveFutureChannelsBan(), UnbanWithCreatedBy(creator.ID))
	require.NoError(t, err)
	_, err = c.UnBanUser(ctx, target2.ID, UnbanWithRemoveFutureChannelsBan(), UnbanWithCreatedBy(creator.ID))
	require.NoError(t, err)
}

func ExampleClient_BanUser() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	// ban a user for 60 minutes from all channel
	_, _ = client.BanUser(ctx, "eviluser", "modUser", BanWithExpiration(60), BanWithReason("Banned for one hour"))

	// ban a user from the livestream:fortnite channel
	channel := client.Channel("livestream", "fortnite")
	_, _ = channel.BanUser(ctx, "eviluser", "modUser", BanWithReason("Profanity is not allowed here"))

	// remove ban from channel
	channel = client.Channel("livestream", "fortnite")
	_, _ = channel.UnBanUser(ctx, "eviluser")

	// remove global ban
	_, _ = client.UnBanUser(ctx, "eviluser")
}
