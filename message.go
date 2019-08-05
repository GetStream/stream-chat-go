package stream_chat

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Message struct {
	ID string `json:"id"`

	Text string `json:"text"`
	HTML string `json:"html"`

	Type messageType `json:"type"`

	User            *User          `json:"user"`
	Attachments     []Attachment   `json:"attachments"`
	LatestReactions []Reaction     `json:"latest_reactions"` // last reactions
	OwnReactions    []Reaction     `json:"own_reactions"`
	ReactionCounts  map[string]int `json:"reaction_counts"`

	ReplyCount int `json:"reply_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// any other fields the user wants to attach a message
	ExtraData map[string]interface{}

	MentionedUsers []User `json:"mentioned_users"`
}

func (m *Message) MarshallJSON() ([]byte, error) {
	return json.Marshal(m.toHash())
}

func (m *Message) toHash() map[string]interface{} {
	var data = map[string]interface{}{}
	for k, v := range m.ExtraData {
		data[k] = v
	}

	data["text"] = m.Text

	if len(m.Attachments) > 0 {
		data["attachments"] = m.Attachments
	}

	if m.User != nil {
		data["user"] = m.User
	}

	if len(m.MentionedUsers) > 0 {
		data["mentioned_users"] = m.MentionedUsers
	}

	return data
}

type Attachment struct {
}

// MarkAllRead marks all messages as read for userID
func (c *client) MarkAllRead(userID string) error {
	data := map[string]interface{}{
		"user": map[string]string{
			"id": userID,
		},
	}

	return c.makeRequest(http.MethodPost, "channels/read", nil, data, nil)
}

func (c *client) UpdateMessage(msg Message) error {
	if msg.ID == "" {
		return errors.New("message ID must be not empty")
	}

	data := map[string]interface{}{
		"message": msg.toHash(),
	}

	return c.makeRequest(http.MethodPost, "messages/"+msg.ID, nil, data, nil)
}

func (c *client) DeleteMessage(msgID string) error {
	return c.makeRequest(http.MethodDelete, "messages/"+msgID, nil, nil, nil)
}
