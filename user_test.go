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

	resp, err := c.UpdateUsers(user)
	mustNoError(t, err, "update users")

	assert.NotEmpty(t, resp[user.ID].CreatedAt)
	assert.NotEmpty(t, resp[user.ID].UpdatedAt)
}
