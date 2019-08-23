package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestClient_DeleteMessage(t *testing.T) {
	// TODO: add test cases.
}

func TestClient_MarkAllRead(t *testing.T) {
	// TODO: Add test cases.
}

func TestClient_UpdateMessage(t *testing.T) {
	// TODO: Add test cases.
}
