package stream_chat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateChannel(t *testing.T) {
	c := initClient(t)

	t.Run("get existing channel", func(t *testing.T) {
		ch := initChannel(t, c)
		got, err := CreateChannel(c, ch.Type, ch.ID, serverUser.ID, nil)
		mustNoError(t, err, "create channel", ch)

		assert.Equal(t, c, got.client, "client link")
		assert.Equal(t, ch.Type, got.Type, "channel type")
		assert.Equal(t, ch.ID, got.ID, "channel id")
		assert.Equal(t, got.MemberCount, ch.MemberCount, "member count")
		assert.Len(t, got.Members, got.MemberCount, "members length")
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
			got, err := CreateChannel(c, tt._type, tt.id, tt.userID, tt.data)
			if tt.wantErr {
				mustError(t, err, "create channel", tt)
			} else {
				mustNoError(t, err, "create channel", tt)
			}

			assert.Equal(t, tt._type, got.Type, "channel type")
			assert.Equal(t, tt.id, got.ID, "channel id")
			assert.Equal(t, tt.userID, got.CreatedBy.ID, "channel created by")
		})
	}
}

func TestChannel_AddMembers(t *testing.T) {
	c := initClient(t)

	chanID := randomString(12)

	ch, err := CreateChannel(c, "messaging", chanID, serverUser.ID, nil)
	mustNoError(t, err, "create channel")
	defer ch.Delete()

	assert.Empty(t, ch.Members, "members are empty")

	user := randomUser()

	err = ch.AddMembers(user.ID)
	mustNoError(t, err, "add members")

	// refresh channel state
	mustNoError(t, ch.refresh(), "refresh channel")

	assert.Equal(t, user.ID, ch.Members[0].User.ID, "members contain user id")
}

func TestChannel_Moderation(t *testing.T) {
	c := initClient(t)

	// init random channel
	chanID := randomString(12)
	ch, err := CreateChannel(c, "messaging", chanID, serverUser.ID, nil)
	mustNoError(t, err, "create channel")
	defer ch.Delete()

	assert.Empty(t, ch.Members, "members are empty")

	user := randomUser()

	err = ch.AddModerators(user.ID)
	mustNoError(t, err, "add moderators")

	// refresh channel state
	mustNoError(t, ch.refresh(), "refresh channel")

	assert.Equal(t, user.ID, ch.Members[0].User.ID, "user exists")
	assert.Equal(t, "moderator", ch.Members[0].Role, "user role is moderator")

	err = ch.DemoteModerators(user.ID)
	mustNoError(t, err, "demote moderators")

	// refresh channel state
	mustNoError(t, ch.refresh(), "refresh channel")

	assert.Equal(t, user.ID, ch.Members[0].User.ID, "user exists")
	assert.Equal(t, "member", ch.Members[0].Role, "user role is member")
}

func TestChannel_BanUser(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()

	err := ch.BanUser(user.ID, serverUser.ID, nil)
	mustNoError(t, err, "ban user")

	err = ch.BanUser(user.ID, serverUser.ID, map[string]interface{}{
		"timeout": 3600,
		"reason":  "offensive language is not allowed here",
	})
	mustNoError(t, err, "ban user")

	err = ch.UnBanUser(user.ID, nil)
	mustNoError(t, err, "unban user")
}

func TestChannel_Delete(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)

	err := ch.Delete()
	mustNoError(t, err, "delete channel")
}

func TestChannel_GetReplies(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()

	msg := &Message{Text: "test message"}

	msg, err := ch.SendMessage(msg, user.ID)
	mustNoError(t, err, "send message")

	reply := &Message{Text: "test reply", ParentID: msg.ID, Type: MessageTypeReply}
	reply, err = ch.SendMessage(reply, serverUser.ID)
	mustNoError(t, err, "send reply")

	replies, err := ch.GetReplies(msg.ID, nil)
	mustNoError(t, err, "get replies")
	assert.Len(t, replies, 1)
}

func TestChannel_MarkRead(t *testing.T) {

}

func TestChannel_RemoveMembers(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	err := ch.RemoveMembers(user.ID)

	mustNoError(t, err, "remove members")

	for _, member := range ch.Members {
		assert.NotEqual(t, member.User.ID, user.ID, "member is not present")
	}
}

func TestChannel_SendEvent(t *testing.T) {

}

func TestChannel_SendMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}

	msg, err := ch.SendMessage(msg, serverUser.ID)
	mustNoError(t, err, "send message")
	// check that message was updated
	assert.NotEmpty(t, msg.ID, "message has ID")
	assert.NotEmpty(t, msg.HTML, "message has HTML body")
}

func TestChannel_Truncate(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	defer ch.Delete()

	user := randomUser()
	msg := &Message{
		Text: "test message",
		User: user,
	}
	msg, err := ch.SendMessage(msg, serverUser.ID)
	mustNoError(t, err, "send message")

	// refresh channel state
	mustNoError(t, ch.refresh(), "refresh channel")

	assert.Equal(t, ch.Messages[0].ID, msg.ID, "message exists")

	err = ch.Truncate()
	mustNoError(t, err, "truncate channel")

	// refresh channel state
	mustNoError(t, ch.refresh(), "refresh channel")

	assert.Empty(t, ch.Messages, "message not exists")
}

func TestChannel_Update(t *testing.T) {

}
