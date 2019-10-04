package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	user := randomUser()

	err := c.MuteUser(user.ID, serverUser.ID)
	mustNoError(t, err, "mute user")

	users, err := c.QueryUsers(&QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]string{"$eq": serverUser.ID},
		}})

	mustNoError(t, err, "query users")

	assert.Lenf(t, users[0].Mutes, 1, "user mutes exists: %+v", users[0])

	mute := users[0].Mutes[0]
	assert.NotEmpty(t, mute.User, "mute has user")
	assert.NotEmpty(t, mute.Target, "mute has target")
}

func TestClient_UnBanUser(t *testing.T) {
}

func TestClient_UnFlagUser(t *testing.T) {
}

func TestClient_UnmuteUser(t *testing.T) {
}

func TestClient_UpdateUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	resp, err := c.UpdateUsers(user)
	mustNoError(t, err, "update users")

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
	mustNoError(t, err, "partial update user")

	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test", "extra data contains", got[user.ID].ExtraData)
	assert.Equal(t, got[user.ID].ExtraData["test"], map[string]interface{}{
		"passed": true,
	})

	update = PartialUserUpdate{
		ID:    user.ID,
		Unset: []string{"test.passed"},
	}

	got, err = c.PartialUpdateUsers([]PartialUserUpdate{update})
	mustNoError(t, err, "partial update user")

	assert.Contains(t, got, user.ID)
	assert.Contains(t, got[user.ID].ExtraData, "test", "extra data contains", got[user.ID].ExtraData)
	assert.Empty(t, got[user.ID].ExtraData["test"], "extra data field removed")
}

func TestClient_GetOrCreateUsers(t *testing.T) {
	c := initClient(t)

	user := randomUser()

	user, err := c.UpdateUser(user)
	mustNoError(t, err)

	newUser := &User{ID: randomString(10), ExtraData: map[string]interface{}{
		"test": true,
	}}

	got, err := c.GetOrCreateUsers([]*User{user, newUser})
	mustNoError(t, err)

	assert.Contains(t, got, user.ID, "response contains user")
	assert.Contains(t, got, newUser.ID, "response contains new user")
	assert.Equal(t, got[newUser.ID].ExtraData["test"], true, "new user has extra field")
}
