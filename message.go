package stream_chat

import (
	"errors"
	"net/http"
	"net/url"
	"path"
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

	MentionedUsers []User `json:"mentioned_users"`
}

func (m Message) toRequest() messageRequest {
	var req messageRequest

	req.Message = messageRequestMessage{
		Text:        m.Text,
		Attachments: m.Attachments,
		User:        messageRequestUser{ID: m.User.ID},
		ExtraData:   m.ExtraData,
	}

	if len(m.MentionedUsers) > 0 {
		req.Message.MentionedUsers = make([]string, 0, len(m.MentionedUsers))
		for _, u := range m.MentionedUsers {
			req.Message.MentionedUsers = append(req.Message.MentionedUsers, u.ID)
		}
	}

	return req
}

type Attachment struct {
}

type messageResponse struct {
	Message  Message `json:"message"`
	duration string
}

type messageRequest struct {
	Message messageRequestMessage `json:"message"`
}

type messageRequestUser struct {
	ID string `json:"id"`
}

type messageRequestMessage struct {
	Text           string                 `json:"text"`
	Attachments    []Attachment           `json:"attachments"`
	User           messageRequestUser     `json:"user"`
	MentionedUsers []string               `json:"mentioned_users"`
	ExtraData      map[string]interface{} `json:"-,extra"`
}

// SendMessage sends a message to the channel.
// *Message will be updated from response body
func (ch *Channel) SendMessage(message *Message, userID string) error {
	var resp messageResponse

	message.User = &User{ID: userID}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "message")

	err := ch.client.makeRequest(http.MethodPost, p, nil, message.toRequest(), &resp)
	if err != nil {
		return err
	}

	*message = resp.Message

	return nil
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

func (c *Client) UpdateMessage(msg Message, msgID string) error {
	if msgID == "" {
		return errors.New("message ID must be not empty")
	}

	return c.makeRequest(http.MethodPost, "messages/"+msgID, nil, msg.toRequest(), nil)
}

func (c *Client) DeleteMessage(msgID string) error {
	return c.makeRequest(http.MethodDelete, "messages/"+msgID, nil, nil, nil)
}
