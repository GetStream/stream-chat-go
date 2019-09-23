// Package stream_chat provides chat via stream api
//nolint: golint
package stream_chat

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

	t.Run("update users", func(t *testing.T) {
		resp, err := c.UpdateUsers(user)
		mustNoError(t, err, "update users")

		assert.Contains(t, resp, user.ID)
		assert.NotEmpty(t, resp[user.ID].CreatedAt)
		assert.NotEmpty(t, resp[user.ID].UpdatedAt)
	})

	t.Run("partial update", func(t *testing.T) {
		extra := map[string]interface{}{
			"test":   true,
			"random": randomString(12),
		}

		resp, err := c.UpdateUsers(&User{ID: user.ID, ExtraData: extra})
		mustNoError(t, err, "update users")

		assert.Contains(t, resp, user.ID)
		assert.Contains(t, resp[user.ID].ExtraData, "test", "extra data contains", resp[user.ID].ExtraData)
		assert.Contains(t, resp[user.ID].ExtraData, "random", "extra data contains", resp[user.ID].ExtraData)
		assert.Equal(t, extra["test"], resp[user.ID].ExtraData["test"], "extra data equal", resp[user.ID].ExtraData)
		assert.Equal(t, extra["random"], resp[user.ID].ExtraData["random"], "extra data equal", resp[user.ID].ExtraData)
	})

	t.Run("remove custom field", func(t *testing.T) {
		extra := map[string]interface{}{
			"test":   true,
			"random": randomString(12),
		}

		_, err := c.UpdateUsers(&User{ID: user.ID, ExtraData: extra})
		mustNoError(t, err, "update users")

		resp, err := c.UpdateUsers(&User{ID: user.ID, ExtraData: map[string]interface{}{"test": nil, "1": 1}})
		assert.Contains(t, resp, user.ID)
		assert.NotContains(t, resp[user.ID].ExtraData, "test")
		assert.Contains(t, resp[user.ID].ExtraData, "random")
	})
}
