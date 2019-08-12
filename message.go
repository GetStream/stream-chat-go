package stream_chat

import (
	"errors"
	"net/http"
	"time"
)

type attachments []Attachment

type Message struct {
	ID string `json:"id"`

	Text string `json:"text"`
	HTML string `json:"html"`

	Type messageType `json:"type"`

	User            *User          `json:"user"`
	Attachments     attachments    `json:"attachments"`
	LatestReactions []Reaction     `json:"latest_reactions"` // last reactions
	OwnReactions    []Reaction     `json:"own_reactions"`
	ReactionCounts  map[string]int `json:"reaction_counts"`

	ReplyCount int `json:"reply_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// any other fields the user wants to attach a message
	ExtraData map[string]interface{}

	MentionedUsers users `json:"mentioned_users"`
}

type Attachment struct {
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

	req := messageRequest{Message: msg}

	return c.makeRequest(http.MethodPost, "messages/"+msg.ID, nil, req, nil)
}

func (c *Client) DeleteMessage(msgID string) error {
	return c.makeRequest(http.MethodDelete, "messages/"+msgID, nil, nil, nil)
}
