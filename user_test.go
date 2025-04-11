package stream_chat

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_MuteUser(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	user := randomUser(t, c)
	_, err := c.MuteUser(ctx, randomUser(t, c).ID, user.ID)
	require.NoError(t, err, "MuteUser should not return an error")

	resp, err := c.QueryUsers(ctx, &QueryUsersOptions{
		QueryOption: QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]string{"$eq": user.ID},
			},
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
	_, err = c.MuteUser(ctx, randomUser(t, c).ID, user.ID, MuteWithExpiration(60))
	require.NoError(t, err, "MuteUser should not return an error")

	resp, err = c.QueryUsers(ctx, &QueryUsersOptions{
		QueryOption: QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]string{"$eq": user.ID},
			},
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
	ctx := context.Background()

	user := randomUser(t, c)
	targetIDs := randomUsersID(t, c, 2)

	_, err := c.MuteUsers(ctx, targetIDs, user.ID, MuteWithExpiration(60))
	require.NoError(t, err, "MuteUsers should not return an error")

	resp, err := c.QueryUsers(ctx, &QueryUsersOptions{
		QueryOption: QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]string{"$eq": user.ID},
			},
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

func TestClient_BlockUsers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	blockingUser := randomUser(t, c)
	blockedUser := randomUser(t, c)

	_, err := c.BlockUser(ctx, blockedUser.ID, blockingUser.ID)
	require.NoError(t, err, "BlockUser should not return an error")

	resp, err := c.QueryUsers(ctx, &QueryUsersOptions{
		QueryOption: QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]string{"$eq": blockingUser.ID},
			},
		},
	})

	users := resp.Users
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.Equal(t, len(users[0].BlockedUserIDs), 1)

	require.Equal(t, users[0].BlockedUserIDs[0], blockedUser.ID)
}

func TestClient_UnblockUsersGetBlockedUsers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	blockingUser := randomUser(t, c)
	blockedUser := randomUser(t, c)

	_, err := c.BlockUser(ctx, blockedUser.ID, blockingUser.ID)
	require.NoError(t, err, "BlockUser should not return an error")

	resp, err := c.QueryUsers(ctx, &QueryUsersOptions{
		QueryOption: QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]string{"$eq": blockingUser.ID},
			},
		},
	})

	users := resp.Users
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.Equal(t, len(users[0].BlockedUserIDs), 1)
	require.Equal(t, users[0].BlockedUserIDs[0], blockedUser.ID)

	getRes, err := c.GetBlockedUser(ctx, blockingUser.ID)
	require.Equal(t, 1, len(getRes.BlockedUsers))
	require.Equal(t, blockedUser.ID, getRes.BlockedUsers[0].BlockedUserID)

	_, err = c.UnblockUser(ctx, blockedUser.ID, blockingUser.ID)
	require.NoError(t, err, "UnblockUser should not return an error")

	resp, err = c.QueryUsers(ctx, &QueryUsersOptions{
		QueryOption: QueryOption{
			Filter: map[string]interface{}{
				"id": map[string]string{"$eq": blockingUser.ID},
			},
		},
	})

	users = resp.Users
	require.NoError(t, err, "QueryUsers should not return an error")
	require.NotEmptyf(t, users, "QueryUsers should return a user: %+v", users)
	require.Equal(t, 0, len(users[0].BlockedUserIDs))
}

func TestClient_UnmuteUser(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	user := randomUser(t, c)
	mutedUser := randomUser(t, c)

	_, err := c.MuteUser(ctx, mutedUser.ID, user.ID)
	require.NoError(t, err, "MuteUser should not return an error")

	_, err = c.UnmuteUser(ctx, mutedUser.ID, user.ID)
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
	ctx := context.Background()
	user := randomUser(t, c)

	targetIDs := []string{randomUser(t, c).ID, randomUser(t, c).ID}
	_, err := c.MuteUsers(ctx, targetIDs, user.ID)
	require.NoError(t, err, "MuteUsers should not return an error")

	_, err = c.UnmuteUsers(ctx, targetIDs, user.ID)
	assert.NoError(t, err, "unmute users")
}

func TestClient_UpsertUsers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	user := &User{ID: randomString(10)}

	resp, err := c.UpsertUsers(ctx, user)
	require.NoError(t, err, "update users")

	assert.Contains(t, resp.Users, user.ID)
	assert.NotEmpty(t, resp.Users[user.ID].CreatedAt)
	assert.NotEmpty(t, resp.Users[user.ID].UpdatedAt)
}

func TestClient_UpsertUsersWithRoleAndTeamsRole(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	user := &User{
		ID:        randomString(10),
		Role:      "admin",
		Teams:     []string{"blue"},
		TeamsRole: map[string]string{"blue": "admin"},
	}

	resp, err := c.UpsertUsers(ctx, user)
	require.NoError(t, err, "update users with role and teams_role")

	assert.Contains(t, resp.Users, user.ID)
	assert.Equal(t, "admin", resp.Users[user.ID].Role)
	assert.Equal(t, []string{"blue"}, resp.Users[user.ID].Teams)
	assert.Equal(t, map[string]string{"blue": "admin"}, resp.Users[user.ID].TeamsRole)
	assert.NotEmpty(t, resp.Users[user.ID].CreatedAt)
	assert.NotEmpty(t, resp.Users[user.ID].UpdatedAt)
}

func TestClient_UpdatePrivacySettings(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	user := &User{ID: randomString(10)}

	resp, err := c.UpsertUser(ctx, user)
	require.NoError(t, err, "update users")

	require.Equal(t, resp.User.ID, user.ID)
	require.Nil(t, resp.User.PrivacySettings)

	user = resp.User
	user.PrivacySettings = &PrivacySettings{
		TypingIndicators: &TypingIndicators{
			Enabled: false,
		},
	}
	resp, err = c.UpsertUser(ctx, user)
	require.NoError(t, err, "update users")

	require.Equal(t, resp.User.ID, user.ID)
	require.NotNil(t, resp.User.PrivacySettings)
	require.False(t, resp.User.PrivacySettings.TypingIndicators.Enabled)
	require.Nil(t, resp.User.PrivacySettings.ReadReceipts)

	user = resp.User
	user.PrivacySettings = &PrivacySettings{
		TypingIndicators: &TypingIndicators{
			Enabled: true,
		},
		ReadReceipts: &ReadReceipts{
			Enabled: false,
		},
	}
	resp, err = c.UpsertUser(ctx, user)
	require.NoError(t, err, "update users")

	require.Equal(t, resp.User.ID, user.ID)
	require.NotNil(t, resp.User.PrivacySettings)
	require.True(t, resp.User.PrivacySettings.TypingIndicators.Enabled)
	require.False(t, resp.User.PrivacySettings.ReadReceipts.Enabled)
}

func TestClient_PartialUpdateUsers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	user := randomUser(t, c)

	update := PartialUserUpdate{
		ID: user.ID,
		Set: map[string]interface{}{
			"test": map[string]interface{}{
				"passed": true,
			},
		},
	}

	resp, err := c.PartialUpdateUsers(ctx, []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	got := resp.Users
	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test", "extra data contains: %v", got[user.ID].ExtraData)
	assert.Equal(t, map[string]interface{}{"passed": true}, got[user.ID].ExtraData["test"])

	update = PartialUserUpdate{
		ID:    user.ID,
		Unset: []string{"test.passed"},
	}

	resp, err = c.PartialUpdateUsers(ctx, []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	got = resp.Users
	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test", "extra data contains", got[user.ID].ExtraData)
	assert.Empty(t, got[user.ID].ExtraData["test"], "extra data field removed")
}

func TestClient_PartialUpdatePrivacySettings(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	user := &User{ID: randomString(10)}

	upsertResponse, err := c.UpsertUser(ctx, user)
	require.NoError(t, err, "update users")

	require.Equal(t, upsertResponse.User.ID, user.ID)
	require.Nil(t, upsertResponse.User.PrivacySettings)

	update := PartialUserUpdate{
		ID: user.ID,
		Set: map[string]interface{}{
			"privacy_settings": map[string]interface{}{
				"typing_indicators": map[string]interface{}{
					"enabled": true,
				},
			},
		},
	}

	partialUpdateResponse, err := c.PartialUpdateUsers(ctx, []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")

	require.True(t, partialUpdateResponse.Users[user.ID].PrivacySettings.TypingIndicators.Enabled)
	require.Nil(t, partialUpdateResponse.Users[user.ID].PrivacySettings.ReadReceipts)

	update = PartialUserUpdate{
		ID: user.ID,
		Set: map[string]interface{}{
			"privacy_settings": map[string]interface{}{
				"read_receipts": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}
	partialUpdateResponse, err = c.PartialUpdateUsers(ctx, []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user")
	require.True(t, partialUpdateResponse.Users[user.ID].PrivacySettings.TypingIndicators.Enabled)
	require.False(t, partialUpdateResponse.Users[user.ID].PrivacySettings.ReadReceipts.Enabled)
}

func TestClient_PartialUpdateUserWithTeam(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// First create a basic user
	user := &User{ID: randomString(10)}
	upsertResp, err := c.UpsertUser(ctx, user)
	require.NoError(t, err, "create user")
	assert.Equal(t, upsertResp.User.ID, user.ID)

	// Partially update the user with team and teams_role
	update := PartialUserUpdate{
		ID: user.ID,
		Set: map[string]interface{}{
			"teams":      []string{"blue"},
			"teams_role": map[string]string{"blue": "admin"},
		},
	}

	partialResp, err := c.PartialUpdateUsers(ctx, []PartialUserUpdate{update})
	require.NoError(t, err, "partial update user with team")

	// Verify the changes
	assert.Contains(t, partialResp.Users, user.ID)
	assert.Equal(t, []string{"blue"}, partialResp.Users[user.ID].Teams)
	assert.Equal(t, map[string]string{"blue": "admin"}, partialResp.Users[user.ID].TeamsRole)
}

func TestClient_RestoreUsers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	userId := randomString(10)
	users := []*User{
		{
			ID: userId,
		},
	}
	// create users
	_, err := c.UpsertUsers(ctx, users...)
	require.NoError(t, err, "UpsertUsers should not return an error")

	_, err = c.DeleteUser(ctx, userId)
	require.NoError(t, err, "DeactivateUsers should not return an error")

	// Test error case: empty userIDs
	t.Run("Empty userIDs", func(t *testing.T) {
		_, err := c.RestoreUsers(ctx, []string{})
		require.Error(t, err, "RestoreUsers should return an error when userIDs is empty")
		require.Equal(t, "userIDs are empty", err.Error(), "Error message should match")
	})

	// Test successful case
	t.Run("Restore deactivated users", func(t *testing.T) {
		// Get the users to verify they are deactivated
		resp, err := c.QueryUsers(ctx, &QueryUsersOptions{
			QueryOption: QueryOption{
				Filter: map[string]interface{}{
					"id": map[string]interface{}{
						"$in": []string{userId},
					},
				},
			},
		})

		require.NoError(t, err, "QueryUsers should not return an error")
		require.Empty(t, resp.Users, "Response users should be empty")

		for _, user := range resp.Users {
			require.Contains(t, userId, user.ID, "User should be in the list of deactivated users")
		}

		// Restore the users
		restoreResp, err := c.RestoreUsers(ctx, []string{userId})
		require.NoError(t, err, "RestoreUsers should not return an error")
		require.NotNil(t, restoreResp, "Response should not be nil")

		// Verify users are restored by querying them without IncludeDeactivatedUsers
		verifyResp, err := c.QueryUsers(ctx, &QueryUsersOptions{
			QueryOption: QueryOption{
				Filter: map[string]interface{}{
					"id": map[string]interface{}{
						"$in": []string{userId},
					},
				},
			},
		})
		require.NoError(t, err, "QueryUsers should not return an error")
		for _, user := range verifyResp.Users {
			require.Contains(t, []string{userId}, user.ID, "User should be in the list of restored users")
		}
	})
}
func ExampleClient_UpsertUser() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	_, err := client.UpsertUser(ctx, &User{
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
	ctx := context.Background()

	user, _ := client.ExportUser(ctx, "userID")
	log.Printf("%#v", user)
}

func ExampleClient_DeactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	_, _ = client.DeactivateUser(ctx, "userID")
}

func ExampleClient_ReactivateUser() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	_, _ = client.ReactivateUser(ctx, "userID")
}

func ExampleClient_DeleteUser() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	_, _ = client.DeleteUser(ctx, "userID")
}

func ExampleClient_DeleteUser_hard() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	_, _ = client.DeleteUser(ctx, "userID",
		DeleteUserWithHardDelete(),
		DeleteUserWithMarkMessagesDeleted(),
	)
}
