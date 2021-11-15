package stream_chat // nolint: golint

import (
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
	ch, err := c.CreateChannel(ch.Type, ch.ID, userA.ID, nil)
	require.NoError(t, err)

	// shadow ban userB globally
	err = c.ShadowBan(userB.ID, userA.ID, nil)
	require.NoError(t, err)

	// shadow ban userC on channel
	err = ch.ShadowBan(userC.ID, userA.ID, nil)
	require.NoError(t, err)

	msg := &Message{Text: "test message"}
	msg, err = ch.SendMessage(msg, userB.ID)
	require.NoError(t, err)
	require.Equal(t, false, msg.Shadowed)

	msg, err = c.GetMessage(msg.ID)
	require.NoError(t, err)
	require.Equal(t, true, msg.Shadowed)

	msg = &Message{Text: "test message"}
	msg, err = ch.SendMessage(msg, userC.ID)
	require.NoError(t, err)
	require.Equal(t, false, msg.Shadowed)

	msg, err = c.GetMessage(msg.ID)
	require.NoError(t, err)
	require.Equal(t, true, msg.Shadowed)

	err = c.RemoveShadowBan(userB.ID, nil)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	msg, err = ch.SendMessage(msg, userB.ID)
	require.NoError(t, err)
	require.Equal(t, false, msg.Shadowed)

	msg, err = c.GetMessage(msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, msg.Shadowed)

	err = ch.RemoveShadowBan(userC.ID, nil)
	require.NoError(t, err)

	msg = &Message{Text: "test message"}
	msg, err = ch.SendMessage(msg, userC.ID)
	require.NoError(t, err)
	require.Equal(t, false, msg.Shadowed)

	msg, err = c.GetMessage(msg.ID)
	require.NoError(t, err)
	require.Equal(t, false, msg.Shadowed)
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
	err := c.MuteUser(randomUser(t, c).ID, user.ID, nil)
	require.NoError(t, err, "MuteUser should not return an error")

	users, err := c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": user.ID},
		},
	})
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.NotEmptyf(t, users[0].Mutes, "user should have Mutes: %+v", users[0])

	mute := users[0].Mutes[0]
	assert.NotEmpty(t, mute.User, "mute should have a User")
	assert.NotEmpty(t, mute.Target, "mute should have a Target")
	assert.Empty(t, mute.Expires, "mute should have no Expires")

	user = randomUser(t, c)
	// when timeout is given, expiration field should be set on mute
	err = c.MuteUser(randomUser(t, c).ID, user.ID, map[string]interface{}{"timeout": 60})
	require.NoError(t, err, "MuteUser should not return an error")

	users, err = c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": user.ID},
		},
	})
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

	err := c.MuteUsers(targetIDs, user.ID, map[string]interface{}{"timeout": 60})
	require.NoError(t, err, "MuteUsers should not return an error")

	users, err := c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": user.ID},
		},
	})
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
	err := c.MuteUser(mutedUser.ID, user.ID, nil)
	require.NoError(t, err, "MuteUser should not return an error")

	err = c.UnmuteUser(mutedUser.ID, user.ID)
	assert.NoError(t, err)
}

func TestClient_UnmuteUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	targetIDs := []string{randomUser(t, c).ID, randomUser(t, c).ID}
	err := c.MuteUsers(targetIDs, user.ID, nil)
	require.NoError(t, err, "MuteUsers should not return an error")

	err = c.UnmuteUsers(targetIDs, user.ID)
	assert.NoError(t, err, "unmute users")
}

func TestClient_UpsertUsers(t *testing.T) {
	c := initClient(t)

	user := &User{ID: randomString(10)}

	resp, err := c.UpsertUsers(user)
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

	got, err := c.PartialUpdateUsers([]PartialUserUpdate{update})
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

	got, err = c.PartialUpdateUsers([]PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test", "extra data contains", got[user.ID].ExtraData)
	assert.Empty(t, got[user.ID].ExtraData["test"], "extra data field removed")
}

func ExampleClient_UpsertUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_, err := client.UpsertUser(&User{
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

	user, _ := client.ExportUser("userID", nil)
	log.Printf("%#v", user)
}

func ExampleClient_DeactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.DeactivateUser("userID", nil)
}

func ExampleClient_ReactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.ReactivateUser("userID", nil)
}

func ExampleClient_DeleteUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_ = client.DeleteUser("userID", nil)
}

func ExampleClient_DeleteUser_hard() {
	client, _ := NewClient("XXXX", "XXXX")

	options := map[string][]string{
		"mark_messages_deleted": {"true"},
		"hard_delete":           {"true"},
	}

	_ = client.DeleteUser("userID", options)
}

func ExampleClient_BanUser() {
	client, _ := NewClient("XXXX", "XXXX")

	// ban a user for 60 minutes from all channel
	_ = client.BanUser("eviluser", "modUser",
		map[string]interface{}{"timeout": 60, "reason": "Banned for one hour"})

	// ban a user from the livestream:fortnite channel
	channel := client.Channel("livestream", "fortnite")
	_ = channel.BanUser("eviluser", "modUser",
		map[string]interface{}{"reason": "Profanity is not allowed here"})

	// remove ban from channel
	channel = client.Channel("livestream", "fortnite")
	_ = channel.UnBanUser("eviluser", nil)

	// remove global ban
	_ = client.UnBanUser("eviluser", nil)
}
