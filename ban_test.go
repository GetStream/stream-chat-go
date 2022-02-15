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
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.Equal(t, true, messageResp.Message.Shadowed)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(ctx, msg, userC.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.Equal(t, true, messageResp.Message.Shadowed)

	_, err = c.UnBanUser(ctx, userB.ID)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(ctx, msg, userB.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, messageResp.Message.Shadowed)

	_, err = ch.UnBanUser(ctx, userC.ID)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(ctx, msg, userC.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(ctx, msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, messageResp.Message.Shadowed)
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
	require.Equal(t, resp.Bans[0].Reason, "spammer")
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
