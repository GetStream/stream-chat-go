package stream_chat

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	_, err := c.MuteUser(context.Background(), randomUser(t, c).ID, user.ID)
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
	_, err = c.MuteUser(context.Background(), randomUser(t, c).ID, user.ID, MuteWithExpiration(60))
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

	_, err := c.MuteUsers(context.Background(), targetIDs, user.ID, MuteWithExpiration(60))
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

func TestClient_UnmuteUser(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	mutedUser := randomUser(t, c)
	_, err := c.MuteUser(context.Background(), mutedUser.ID, user.ID)
	require.NoError(t, err, "MuteUser should not return an error")

	_, err = c.UnmuteUser(context.Background(), mutedUser.ID, user.ID)
	assert.NoError(t, err)
}

func TestClient_CreateGuestUser(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	u := &User{ID: randomString(10)}
	resp, err := c.CreateGuestUser(ctx, u)
	if err != nil {
		// Sometimes the guest user access is disabled on app level
		// so let's ignore errors here
		return
	}
	require.NotNil(t, resp.AccessToken)
	require.NotNil(t, resp.User)
}

func TestClient_UnmuteUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser(t, c)
	targetIDs := []string{randomUser(t, c).ID, randomUser(t, c).ID}
	_, err := c.MuteUsers(context.Background(), targetIDs, user.ID)
	require.NoError(t, err, "MuteUsers should not return an error")

	_, err = c.UnmuteUsers(context.Background(), targetIDs, user.ID)
	assert.NoError(t, err, "unmute users")
}

func TestClient_UpsertUsers(t *testing.T) {
	c := initClient(t)

	user := &User{ID: randomString(10)}

	resp, err := c.UpsertUsers(context.Background(), user)
	require.NoError(t, err, "update users")

	assert.Contains(t, resp.Users, user.ID)
	assert.NotEmpty(t, resp.Users[user.ID].CreatedAt)
	assert.NotEmpty(t, resp.Users[user.ID].UpdatedAt)
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

	resp, err := c.PartialUpdateUsers(context.Background(), []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	got := resp.Users
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

	resp, err = c.PartialUpdateUsers(context.Background(), []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	got = resp.Users
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

	user, _ := client.ExportUser(context.Background(), "userID")
	log.Printf("%#v", user)
}

func ExampleClient_DeactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_, _ = client.DeactivateUser(context.Background(), "userID")
}

func ExampleClient_ReactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_, _ = client.ReactivateUser(context.Background(), "userID")
}

func ExampleClient_DeleteUser() {
	client, _ := NewClient("XXXX", "XXXX")

	_, _ = client.DeleteUser(context.Background(), "userID")
}

func ExampleClient_DeleteUser_hard() {
	client, _ := NewClient("XXXX", "XXXX")

	_, _ = client.DeleteUser(context.Background(), "userID",
		DeleteUserWithHardDelete(),
		DeleteUserWithMarkMessagesDeleted(),
	)
}
