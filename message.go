package stream_chat

import (
	"encoding/json"
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
