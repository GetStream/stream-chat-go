package stream_chat // nolint: golint

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"

	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

type MessageType string

const (
	MessageTypeRegular   MessageType = "regular"
	MessageTypeError     MessageType = "error"
	MessageTypeReply     MessageType = "reply"
	MessageTypeSystem    MessageType = "system"
	MessageTypeEphemeral MessageType = "ephemeral"
)

type Message struct {
	ID string `json:"id"`

	Text string `json:"text"`
	HTML string `json:"html"`

	Type MessageType `json:"type,omitempty"` // one of MessageType* constants

	User            *User          `json:"user"`
	Attachments     []*Attachment  `json:"attachments"`
	LatestReactions []*Reaction    `json:"latest_reactions"` // last reactions
	OwnReactions    []*Reaction    `json:"own_reactions"`
	ReactionCounts  map[string]int `json:"reaction_counts"`

	ParentID      string `json:"parent_id"`       // id of parent message if it's reply
	ShowInChannel bool   `json:"show_in_channel"` // show reply message also in channel

	ReplyCount int `json:"reply_count,omitempty"`

	MentionedUsers []*User `json:"mentioned_users"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// any other fields the user wants to attach a message
	ExtraData map[string]interface{} `json:"-"`
}

// UnmarshalUnknown implements the `easyjson.UnknownsUnmarshaler` interface.
func (m *Message) UnmarshalUnknown(in *jlexer.Lexer, key string) {
	if m.ExtraData == nil {
		m.ExtraData = make(map[string]interface{}, 1)
	}
	m.ExtraData[key] = in.Interface()
}

// MarshalUnknowns implements the `easyjson.UnknownsMarshaler` interface.
func (m Message) MarshalUnknowns(out *jwriter.Writer, first bool) {
	for key, val := range m.ExtraData {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.String(key)
		out.RawByte(':')
		out.Raw(json.Marshal(val))
	}
}

func (m *Message) toRequest() messageRequest {
	var req messageRequest

	req.Message = messageRequestMessage{
		Text:          m.Text,
		Attachments:   m.Attachments,
		User:          messageRequestUser{ID: m.User.ID},
		ExtraData:     m.ExtraData,
		ParentID:      m.ParentID,
		ShowInChannel: m.ShowInChannel,
	}

	if len(m.MentionedUsers) > 0 {
		req.Message.MentionedUsers = make([]string, 0, len(m.MentionedUsers))
		for _, u := range m.MentionedUsers {
			req.Message.MentionedUsers = append(req.Message.MentionedUsers, u.ID)
		}
	}

	return req
}

type messageRequest struct {
	Message messageRequestMessage `json:"message"`
}

type messageRequestMessage struct {
	Text           string             `json:"text"`
	Attachments    []*Attachment      `json:"attachments"`
	User           messageRequestUser `json:"user"`
	MentionedUsers []string           `json:"mentioned_users"`
	ParentID       string             `json:"parent_id"`
	ShowInChannel  bool               `json:"show_in_channel"`

	ExtraData map[string]interface{} `json:"-"`
}

func (s *messageRequestMessage) UnmarshalUnknown(in *jlexer.Lexer, key string) {
	if s.ExtraData == nil {
		s.ExtraData = make(map[string]interface{}, 1)
	}
	s.ExtraData[key] = in.Interface()
}

func (s messageRequestMessage) MarshalUnknowns(out *jwriter.Writer, first bool) {
	for key, val := range s.ExtraData {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.String(key)
		out.RawByte(':')
		out.Raw(json.Marshal(val))
	}
}

type messageRequestUser struct {
	ID string `json:"id"`
}

type messageResponse struct {
	Message *Message `json:"message"`
}

type Attachment struct {
	Type string `json:"type,omitempty"` // text, image, audio, video

	AuthorName string `json:"author_name,omitempty"`
	Title      string `json:"title,omitempty"`
	TitleLink  string `json:"title_link,omitempty"`
	Text       string `json:"text,omitempty"`

	ImageURL    string `json:"image_url,omitempty"`
	ThumbURL    string `json:"thumb_url,omitempty"`
	AssetURL    string `json:"asset_url,omitempty"`
	OGScrapeURL string `json:"og_scrape_url,omitempty"`

	ExtraData map[string]interface{} `json:"-"`
}

// UnmarshalUnknown implements the `easyjson.UnknownsUnmarshaler` interface.
func (a *Attachment) UnmarshalUnknown(in *jlexer.Lexer, key string) {
	if a.ExtraData == nil {
		a.ExtraData = make(map[string]interface{}, 1)
	}
	a.ExtraData[key] = in.Interface()
}

// MarshalUnknowns implements the `easyjson.UnknownsMarshaler` interface.
func (a Attachment) MarshalUnknowns(out *jwriter.Writer, first bool) {
	for key, val := range a.ExtraData {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.String(key)
		out.RawByte(':')
		out.Raw(json.Marshal(val))
	}
}

// SendMessage sends a message to the channel. Returns full message details from server.
func (ch *Channel) SendMessage(message *Message, userID string) (*Message, error) {
	switch {
	case message == nil:
		return nil, errors.New("message is nil")
	case userID == "":
		return nil, errors.New("user ID must be not empty")
	}

	var resp messageResponse

	message.User = &User{ID: userID}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "message")

	err := ch.client.makeRequest(http.MethodPost, p, nil, message.toRequest(), &resp)
	if err != nil {
		return nil, err
	}

	return resp.Message, nil
}

// MarkAllRead marks all messages as read for userID.
func (c *Client) MarkAllRead(userID string) error {
	if userID == "" {
		return errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"user": map[string]string{
			"id": userID,
		},
	}

	return c.makeRequest(http.MethodPost, "channels/read", nil, data, nil)
}

// GetMessage returns message by ID.
func (c *Client) GetMessage(msgID string) (*Message, error) {
	if msgID == "" {
		return nil, errors.New("message ID must be not empty")
	}

	var resp messageResponse

	p := path.Join("messages", url.PathEscape(msgID))

	err := c.makeRequest(http.MethodGet, p, nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Message, nil
}

// UpdateMessage updates message with given msgID.
func (c *Client) UpdateMessage(msg *Message, msgID string) (*Message, error) {
	switch {
	case msg == nil:
		return nil, errors.New("message is nil")
	case msgID == "":
		return nil, errors.New("message ID must be not empty")
	}

	var resp messageResponse

	p := path.Join("messages", url.PathEscape(msgID))

	err := c.makeRequest(http.MethodPost, p, nil, msg.toRequest(), &resp)
	if err != nil {
		return nil, err
	}

	return resp.Message, nil
}

func (c *Client) DeleteMessage(msgID string) error {
	if msgID == "" {
		return errors.New("message ID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(msgID))

	return c.makeRequest(http.MethodDelete, p, nil, nil, nil)
}

func (c *Client) FlagMessage(msgID string) error {
	if msgID == "" {
		return errors.New("message ID is empty")
	}

	options := map[string]interface{}{
		"target_message_id": msgID,
	}

	return c.makeRequest(http.MethodPost, "moderation/flag", nil, options, nil)
}

func (c *Client) UnflagMessage(msgID string) error {
	if msgID == "" {
		return errors.New("message ID is empty")
	}

	options := map[string]interface{}{
		"target_message_id": msgID,
	}

	return c.makeRequest(http.MethodPost, "moderation/unflag", nil, options, nil)
}

type repliesResponse struct {
	Messages []*Message `json:"messages"`
}

// GetReplies returns list of the message replies for a parent message
// options: Pagination params, ie {limit:10, idlte: 10}
func (ch *Channel) GetReplies(parentID string, options map[string][]string) ([]*Message, error) {
	if parentID == "" {
		return nil, errors.New("parent ID is empty")
	}

	p := path.Join("messages", url.PathEscape(parentID), "replies")

	var resp repliesResponse

	err := ch.client.makeRequest(http.MethodGet, p, options, nil, &resp)

	return resp.Messages, err
}

type sendActionRequest struct {
	MessageID string            `json:"message_id"`
	FormData  map[string]string `json:"form_data"`
}

// SendAction for a message.
func (ch *Channel) SendAction(msgID string, formData map[string]string) (*Message, error) {
	switch {
	case msgID == "":
		return nil, errors.New("message ID is empty")
	case len(formData) == 0:
		return nil, errors.New("form data is empty")
	}

	p := path.Join("messages", url.PathEscape(msgID), "action")

	data := sendActionRequest{MessageID: msgID, FormData: formData}

	var resp messageResponse

	err := ch.client.makeRequest(http.MethodPost, p, nil, data, &resp)
	return resp.Message, err
}
