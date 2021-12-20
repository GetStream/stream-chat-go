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

	ch := initChannel(t, c, userA.ID, userB.ID, userC.ID)
	resp, err := c.CreateChannel(context.Background(), ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	ch = resp.Channel

	// shadow ban userB globally
	_, err = c.ShadowBan(context.Background(), userB.ID, userA.ID)
	require.NoError(t, err)

	// shadow ban userC on channel
	_, err = ch.ShadowBan(context.Background(), userC.ID, userA.ID)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	messageResp, err := ch.SendMessage(context.Background(), msg, userB.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)
	require.Equal(t, true, messageResp.Message.Shadowed)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(context.Background(), msg, userC.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)
	require.Equal(t, true, messageResp.Message.Shadowed)

	_, err = c.UnBanUser(context.Background(), userB.ID)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(context.Background(), msg, userB.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, messageResp.Message.Shadowed)

	_, err = ch.UnBanUser(context.Background(), userC.ID)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(context.Background(), msg, userC.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, messageResp.Message.Shadowed)
}

func TestBanUnbanUser(t *testing.T) {
	c := initClient(t)
	target := randomUser(t, c)
	user := randomUser(t, c)

	_, err := c.BanUser(context.Background(), target.ID, user.ID, BanWithReason("spammer"), BanWithExpiration(60))
	require.NoError(t, err)

	resp, err := c.QueryBannedUsers(context.Background(), &QueryBannedUsersOptions{
		QueryOption: &QueryOption{Filter: map[string]interface{}{
			"user_id": map[string]string{"$eq": target.ID},
		}},
	})
	require.NoError(t, err)
	require.Equal(t, resp.Bans[0].Reason, "spammer")
	require.NotZero(t, resp.Bans[0].Expires)

	_, err = c.UnBanUser(context.Background(), target.ID)
	require.NoError(t, err)

	resp, err = c.QueryBannedUsers(context.Background(), &QueryBannedUsersOptions{
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

	_, err := ch.BanUser(context.Background(), target.ID, user.ID, BanWithReason("spammer"), BanWithExpiration(60))
	require.NoError(t, err)

	_, err = ch.UnBanUser(context.Background(), target.ID)
	require.NoError(t, err)

	resp, err := c.QueryBannedUsers(context.Background(), &QueryBannedUsersOptions{
		QueryOption: &QueryOption{Filter: map[string]interface{}{
			"channel_cid": map[string]string{"$eq": ch.CID},
		}},
	})
	require.NoError(t, err)
	require.Empty(t, resp.Bans)
}

func ExampleClient_BanUser() {
	client, _ := NewClient("XXXX", "XXXX")

	// ban a user for 60 minutes from all channel
	_, _ = client.BanUser(context.Background(), "eviluser", "modUser", BanWithExpiration(60), BanWithReason("Banned for one hour"))

	// ban a user from the livestream:fortnite channel
	channel := client.Channel("livestream", "fortnite")
	_, _ = channel.BanUser(context.Background(), "eviluser", "modUser", BanWithReason("Profanity is not allowed here"))

	// remove ban from channel
	channel = client.Channel("livestream", "fortnite")
	_, _ = channel.UnBanUser(context.Background(), "eviluser")

	// remove global ban
	_, _ = client.UnBanUser(context.Background(), "eviluser")
}
