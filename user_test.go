package stream_chat // nolint: golint

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ShadowBanUser(t *testing.T) {
	c := initClient(t)
	userA := randomUser(t, c)
	userB := randomUser(t, c)
	userC := randomUser(t, c)

	ch := initChannel(t, c, userA.ID, userB.ID, userC.ID)
	resp, err := c.CreateChannel(context.Background(), ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	ch = resp.Channel

	// shadow ban userB globally
	_, err = c.ShadowBan(context.Background(), userB.ID, userA.ID, nil)
	require.NoError(t, err)

	// shadow ban userC on channel
	_, err = ch.ShadowBan(context.Background(), userC.ID, userA.ID, nil)
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

	_, err = c.RemoveShadowBan(context.Background(), userB.ID, nil)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	messageResp, err = ch.SendMessage(context.Background(), msg, userB.ID)
	require.NoError(t, err)

	msg = messageResp.Message
	require.Equal(t, false, msg.Shadowed)

	messageResp, err = c.GetMessage(context.Background(), msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, messageResp.Message.Shadowed)

	_, err = ch.RemoveShadowBan(context.Background(), userC.ID)
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

func TestClient_BanUser(t *testing.T) {
}

func TestClient_DeactivateUser(t *testing.T) {
}

func TestClient_DeleteUser(t *testing.T) {
}

func TestClient_ExportUser(t *testing.T) {}

func TestClient_FlagUser(t *testing.T) {
}

func TestClient_MuteUser(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	err := c.MuteUser(context.Background(), randomUser(t, c).ID, user.ID, nil)
	require.NoError(t, err, "MuteUser should not return an error")

	resp, err := c.QueryUsers(context.Background(), &QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": user.ID},
		},
	})

	users := resp.Users
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.NotEmptyf(t, users[0].Mutes, "user should have Mutes: %+v", users[0])

	mute := users[0].Mutes[0]
	assert.NotEmpty(t, mute.User, "mute should have a User")
	assert.NotEmpty(t, mute.Target, "mute should have a Target")
	assert.Empty(t, mute.Expires, "mute should have no Expires")

	user = randomUser(t, c)
	// when timeout is given, expiration field should be set on mute
	err = c.MuteUser(context.Background(), randomUser(t, c).ID, user.ID, map[string]interface{}{"timeout": 60})
	require.NoError(t, err, "MuteUser should not return an error")

	resp, err = c.QueryUsers(context.Background(), &QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": user.ID},
		},
	})

	users = resp.Users
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.NotEmptyf(t, users[0].Mutes, "user should have Mutes: %+v", users[0])

	mute = users[0].Mutes[0]
	assert.NotEmpty(t, mute.User, "mute should have a User")
	assert.NotEmpty(t, mute.Target, "mute should have a Target")
	assert.NotEmpty(t, mute.Expires, "mute should have Expires")
}

func TestClient_MuteUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	targetIDs := randomUsersID(t, c, 2)

	err := c.MuteUsers(context.Background(), targetIDs, user.ID, map[string]interface{}{"timeout": 60})
	require.NoError(t, err, "MuteUsers should not return an error")

	resp, err := c.QueryUsers(context.Background(), &QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": user.ID},
		},
	})

	users := resp.Users
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.NotEmptyf(t, users[0].Mutes, "user should have Mutes: %+v", users[0])

	for _, mute := range users[0].Mutes {
		assert.NotEmpty(t, mute.Expires, "mute should have Expires")
	}
}

func TestClient_UnBanUser(t *testing.T) {
}

func TestClient_UnFlagUser(t *testing.T) {
}

func TestClient_UnmuteUser(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	mutedUser := randomUser(t, c)
	err := c.MuteUser(context.Background(), mutedUser.ID, user.ID, nil)
	require.NoError(t, err, "MuteUser should not return an error")

	err = c.UnmuteUser(context.Background(), mutedUser.ID, user.ID)
	assert.NoError(t, err)
}

func TestClient_UnmuteUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	targetIDs := []string{randomUser(t, c).ID, randomUser(t, c).ID}
	err := c.MuteUsers(context.Background(), targetIDs, user.ID, nil)
	require.NoError(t, err, "MuteUsers should not return an error")

	err = c.UnmuteUsers(context.Background(), targetIDs, user.ID)
	assert.NoError(t, err, "unmute users")
}

func TestClient_UpsertUsers(t *testing.T) {
	c := initClient(t)

	user := &User{ID: randomString(10)}

	resp, err := c.UpsertUsers(context.Background(), user)
	require.NoError(t, err, "update users")

	assert.Contains(t, resp, user.ID)
	assert.NotEmpty(t, resp[user.ID].CreatedAt)
	assert.NotEmpty(t, resp[user.ID].UpdatedAt)
}

func TestClient_PartialUpdateUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)

	update := PartialUserUpdate{
		ID: user.ID,
		Set: map[string]interface{}{
			"test": map[string]interface{}{
				"passed": true,
			},
		},
	}

	got, err := c.PartialUpdateUsers(context.Background(), []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test",
		"extra data contains: %v", got[user.ID].ExtraData)
	assert.Equal(t, got[user.ID].ExtraData["test"], map[string]interface{}{
		"passed": true,
	})

	update = PartialUserUpdate{
		ID:    user.ID,
		Unset: []string{"test.passed"},
	}

	got, err = c.PartialUpdateUsers(context.Background(), []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test", "extra data contains", got[user.ID].ExtraData)
	assert.Empty(t, got[user.ID].ExtraData["test"], "extra data field removed")
}

func ExampleClient_UpsertUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_, err := client.UpsertUser(context.Background(), &User{
		ID:   "tommaso",
		Name: "Tommaso",
		Role: "Admin",
	})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}

func ExampleClient_ExportUser() {
	client, _ := NewClient("XXXX", "XXXX")

	user, _ := client.ExportUser(context.Background(), "userID", nil)
	log.Printf("%#v", user)
}

func ExampleClient_DeactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.DeactivateUser(context.Background(), "userID", nil)
}

func ExampleClient_ReactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.ReactivateUser(context.Background(), "userID", nil)
}

func ExampleClient_DeleteUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.DeleteUser(context.Background(), "userID", nil)
}

func ExampleClient_DeleteUser_hard() {
	client, _ := NewClient("XXXX", "XXXX")

	options := map[string][]string{
		"mark_messages_deleted": {"true"},
		"hard_delete":           {"true"},
	}

	_ = client.DeleteUser(context.Background(), "userID", options)
}

func ExampleClient_BanUser() {
	client, _ := NewClient("XXXX", "XXXX")

	// ban a user for 60 minutes from all channel
	_, _ = client.BanUser(context.Background(), "eviluser", "modUser",
		map[string]interface{}{"timeout": 60, "reason": "Banned for one hour"})

	// ban a user from the livestream:fortnite channel
	channel := client.Channel("livestream", "fortnite")
	_, _ = channel.BanUser(context.Background(), "eviluser", "modUser",
		map[string]interface{}{"reason": "Profanity is not allowed here"})

	// remove ban from channel
	channel = client.Channel("livestream", "fortnite")
	_, _ = channel.UnBanUser(context.Background(), "eviluser", nil)

	// remove global ban
	_, _ = client.UnBanUser(context.Background(), "eviluser", nil)
}
