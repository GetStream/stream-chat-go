package stream_chat

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/francoispqt/gojay"
)

type attachments []Attachment

func (a *attachments) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var at Attachment
	if err := dec.Object(&at); err != nil {
		return err
	}
	*a = append(*a, at)
	return nil
}

type reactions []Reaction

func (r *reactions) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var react Reaction
	if err := dec.Object(&react); err != nil {
		return err
	}
	*r = append(*r, react)
	return nil
}

type reactionsCount map[string]int

func (r *reactionsCount) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	var i int
	if err := dec.Int(&i); err != nil {
		return err
	}
	(*r)[key] = i
	return nil
}

func (r *reactionsCount) NKeys() int {
	return 0
}

type Message struct {
	ID string `json:"id"`

	Text string `json:"text"`
	HTML string `json:"html"`

	Type messageType `json:"type"`

	User            *User          `json:"user"`
	Attachments     attachments    `json:"attachments"`
	LatestReactions reactions      `json:"latest_reactions"` // last reactions
	OwnReactions    reactions      `json:"own_reactions"`
	ReactionCounts  reactionsCount `json:"reaction_counts"`

	ReplyCount int `json:"reply_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// any other fields the user wants to attach a message
	ExtraData map[string]interface{}

	MentionedUsers users `json:"mentioned_users"`
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

	data["user"] = map[string]interface{}{}

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

func (a *Attachment) UnmarshalJSONObject(*gojay.Decoder, string) error {
	return nil
}

func (a *Attachment) NKeys() int {
	return 0
}

// MarkAllRead marks all messages as read for userID
func (c *Client) MarkAllRead(userID string) error {
	data := map[string]interface{}{
		"user": map[string]string{
			"id": userID,
		},
	}

	return c.makeRequest(http.MethodPost, "channels/read", nil, data, nil)
}

func (c *Client) UpdateMessage(msg Message) error {
	if msg.ID == "" {
		return errors.New("message ID must be not empty")
	}

	data := map[string]interface{}{
		"message": msg.toHash(),
	}

	return c.makeRequest(http.MethodPost, "messages/"+msg.ID, nil, data, nil)
}

func (c *Client) DeleteMessage(msgID string) error {
	return c.makeRequest(http.MethodDelete, "messages/"+msgID, nil, nil, nil)
}
