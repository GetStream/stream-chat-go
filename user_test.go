package stream_chat // nolint: golint

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	initChannel(t, c)

	require.NotEmptyf(t, testUsers, "there should be at least 1 test user: %+v", testUsers)
	user := testUsers[0]

	err := c.MuteUser(user.ID, serverUser.ID, nil)
	require.NoError(t, err, "MuteUser should not return an error")

	users, err := c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": serverUser.ID},
		},
	})
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.NotEmptyf(t, users[0].Mutes, "user should have Mutes: %+v", users[0])

	mute := users[0].Mutes[0]
	assert.NotEmpty(t, mute.User, "mute should have a User")
	assert.NotEmpty(t, mute.Target, "mute should have a Target")
	assert.Empty(t, mute.Expires, "mute should have no Expires")

	// when timeout is given, expiration field should be set on mute
	err = c.MuteUser(user.ID, serverUser.ID, map[string]interface{}{"timeout": 60})
	require.NoError(t, err, "MuteUser should not return an error")

	users, err = c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": serverUser.ID},
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
	initChannel(t, c)

	require.GreaterOrEqualf(t, len(testUsers), 2, "there should be at least 2 test users: %+v", testUsers)
	usernames := []string{testUsers[0].ID, testUsers[1].ID}

	err := c.MuteUsers(usernames, serverUser.ID, map[string]interface{}{"timeout": 60})
	require.NoError(t, err, "MuteUsers should not return an error")

	users, err := c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": serverUser.ID},
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
	err := c.UnmuteUser(randomUser().ID, serverUser.ID)
	assert.NoError(t, err)
}

func TestClient_UnmuteUsers(t *testing.T) {
	c := initClient(t)

	users := []string{randomUser().ID, randomUser().ID}

	err := c.UnmuteUsers(users, serverUser.ID)
	assert.NoError(t, err, "unmute users")
}

func TestClient_UpdateUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	resp, err := c.UpdateUsers(user)
	require.NoError(t, err, "update users")

	assert.Contains(t, resp, user.ID)
	assert.NotEmpty(t, resp[user.ID].CreatedAt)
	assert.NotEmpty(t, resp[user.ID].UpdatedAt)
}

func TestClient_PartialUpdateUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser()

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

func ExampleClient_UpdateUser() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

	_, err := client.UpdateUser(&User{
		ID:   "tommaso",
		Name: "Tommaso",
		Role: "Admin",
	})
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
}

func ExampleClient_ExportUser() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

	user, _ := client.ExportUser("userID", nil)
	log.Printf("%#v", user)
}

func ExampleClient_DeactivateUser() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

	_ = client.DeactivateUser("userID", nil)
}

func ExampleClient_ReactivateUser() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

	_ = client.ReactivateUser("userID", nil)
}

func ExampleClient_DeleteUser() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

	_ = client.DeleteUser("userID", nil)
}

func ExampleClient_DeleteUser_hard() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

	options := map[string][]string{
		"mark_messages_deleted": {"true"},
		"hard_delete":           {"true"},
	}

	_ = client.DeleteUser("userID", options)
}

func ExampleClient_BanUser() {
	client, _ := NewClient("XXXX", []byte("XXXX"))

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
