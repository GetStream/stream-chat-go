package stream_chat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateChannel(t *testing.T) {
	c := initClient(t)

	t.Run("get existing channel", func(t *testing.T) {
		ch := initChannel(t, c)
		got, err := c.CreateChannel(ch.Type, ch.ID, serverUser.ID, nil)
		mustNoError(t, err)

		assert.Equal(t, c, got.client)
		assert.Equal(t, ch.Type, got.Type)
		assert.Equal(t, ch.ID, got.ID)
		assert.Equal(t, got.MemberCount, ch.MemberCount)
		assert.Len(t, got.Members, got.MemberCount)
	})

	tests := []struct {
		_type   string
		id      string
		userID  string
		data    map[string]interface{}
		wantErr bool
	}{
		{"messaging", randomString(12), serverUser.ID, nil, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("create new channel %s:%s", tt._type, tt.id), func(t *testing.T) {
			got, err := c.CreateChannel(tt._type, tt.id, tt.userID, tt.data)
			if tt.wantErr {
				mustError(t, err)
			} else {
				mustNoError(t, err)
			}

			assert.Equal(t, tt._type, got.Type)
			assert.Equal(t, tt.id, got.ID)
			assert.Equal(t, tt.userID, got.CreatedBy.ID)
		})
	}
}

func TestChannel_AddMembers(t *testing.T) {
	c := initClient(t)

	chanID := randomString(12)

	ch, err := c.CreateChannel("messaging", chanID, serverUser.ID, nil)
	mustNoError(t, err)
	defer ch.Delete()

	assert.Empty(t, ch.Members)

	user := randomUser()

	err = ch.AddMembers([]string{user.ID})
	mustNoError(t, err)

	// refresh channel state
	mustNoError(t, ch.refresh())

	assert.Equal(t, user.ID, ch.Members[0].User.ID)
}

func TestChannel_Moderation(t *testing.T) {
	c := initClient(t)

	// init random channel
	chanID := randomString(12)
	ch, err := c.CreateChannel("messaging", chanID, serverUser.ID, nil)
	mustNoError(t, err)
	defer ch.Delete()

	assert.Empty(t, ch.Members)

	user := randomUser()

	err = ch.AddModerators([]string{user.ID})
	mustNoError(t, err)

	// refresh channel state
	mustNoError(t, ch.refresh())

	assert.Equal(t, user.ID, ch.Members[0].User.ID)
	assert.Equal(t, "moderator", ch.Members[0].Role)

	err = ch.DemoteModerators([]string{user.ID})
	// refresh channel state
	mustNoError(t, ch.refresh())

	assert.Equal(t, user.ID, ch.Members[0].User.ID)
	assert.Equal(t, "member", ch.Members[0].Role)
}

func TestChannel_BanUser(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()

	err := ch.BanUser(user.ID, serverUser.ID, nil)
	mustNoError(t, err)

	err = ch.BanUser(user.ID, serverUser.ID, map[string]interface{}{
		"timeout": 3600,
		"reason":  "offensive language is not allowed here",
	})
	mustNoError(t, err)

	err = ch.UnBanUser(user.ID, nil)
	mustNoError(t, err)
}

func TestChannel_Delete(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	err := ch.Delete()
	mustNoError(t, err)
}

func TestChannel_GetReplies(t *testing.T) {

}

func TestChannel_MarkRead(t *testing.T) {

}

func TestChannel_RemoveMembers(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	err := ch.RemoveMembers([]string{user.ID})

	mustNoError(t, err)

	for _, member := range ch.Members {
		assert.NotEqual(t, member.User.ID, user.ID)
	}
}

func TestChannel_SendEvent(t *testing.T) {

}

func TestChannel_SendMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: &user,
	}

	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)
	// check that message was updated
	assert.NotEmpty(t, msg.ID)
	assert.NotEmpty(t, msg.HTML)
}

func TestChannel_SendReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: &user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	reaction := Reaction{Type: "love"}

	err = ch.SendReaction(&msg, &reaction, serverUser.ID)
	mustNoError(t, err)

	assert.Equal(t, 1, msg.ReactionCounts[reaction.Type])
	assert.Contains(t, msg.LatestReactions, reaction)
}

func TestChannel_DeleteReaction(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: &user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	reaction := Reaction{Type: "love"}

	err = ch.SendReaction(&msg, &reaction, serverUser.ID)
	mustNoError(t, err)

	err = ch.DeleteReaction(&msg, reaction.Type, serverUser.ID)
	mustNoError(t, err)

	assert.Equal(t, 0, msg.ReactionCounts[reaction.Type])
	assert.Empty(t, msg.LatestReactions)
}

func TestChannel_GetReactions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: &user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	reactions, err := ch.GetReactions(msg.ID, nil)
	mustNoError(t, err)
	assert.Empty(t, reactions)

	reaction := Reaction{Type: "love"}

	err = ch.SendReaction(&msg, &reaction, serverUser.ID)
	mustNoError(t, err)

	reactions, err = ch.GetReactions(msg.ID, nil)

	assert.Contains(t, reactions, reaction)
}

func TestChannel_Truncate(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := Message{
		Text: "test message",
		User: &user,
	}
	err := ch.SendMessage(&msg, serverUser.ID)
	mustNoError(t, err)

	// refresh channel state
	mustNoError(t, ch.refresh())

	assert.Equal(t, ch.Messages[0].ID, msg.ID)

	err = ch.Truncate()
	mustNoError(t, err)

	// refresh channel state
	mustNoError(t, ch.refresh())

	assert.Empty(t, ch.Messages)
}

func TestChannel_Update(t *testing.T) {

}

func Test_addUserID(t *testing.T) {
	id := "someid"

	params := map[string]interface{}{
		"test": 1,
	}

	addUserID(params, id)

	assert.Equal(t, map[string]interface{}{"id": id}, params["user"])
}
