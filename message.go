package stream_chat // nolint: golint

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
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

	Type   MessageType `json:"type,omitempty"` // one of MessageType* constants
	Silent bool        `json:"silent,omitempty"`

	User            *User          `json:"user"`
	Attachments     []*Attachment  `json:"attachments"`
	LatestReactions []*Reaction    `json:"latest_reactions"` // last reactions
	OwnReactions    []*Reaction    `json:"own_reactions"`
	ReactionCounts  map[string]int `json:"reaction_counts"`

	ParentID      string `json:"parent_id"`       // id of parent message if it's reply
	ShowInChannel bool   `json:"show_in_channel"` // show reply message also in channel

	ReplyCount int `json:"reply_count,omitempty"`

	MentionedUsers []*User `json:"mentioned_users"`

	Shadowed bool       `json:"shadowed,omitempty"`
	PinnedAt *time.Time `json:"pinned_at,omitempty"`
	PinnedBy *User      `json:"pinned_by,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	ExtraData map[string]interface{} `json:"-"`
}

type messageForJSON Message

// UnmarshalJSON implements json.Unmarshaler.
func (m *Message) UnmarshalJSON(data []byte) error {
	var m2 messageForJSON
	if err := json.Unmarshal(data, &m2); err != nil {
		return err
	}
	*m = Message(m2)

	if err := json.Unmarshal(data, &m.ExtraData); err != nil {
		return err
	}
	removeFromMap(m.ExtraData, *m)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m Message) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(m.ExtraData, messageForJSON(m))
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
		Silent:        m.Silent,
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
	Message  messageRequestMessage `json:"message"`
	SkipPush bool                  `json:"skip_push,omitempty"`
}

type messageRequestMessage struct {
	Text           string             `json:"text"`
	Attachments    []*Attachment      `json:"attachments"`
	User           messageRequestUser `json:"user"`
	MentionedUsers []string           `json:"mentioned_users"`
	ParentID       string             `json:"parent_id"`
	ShowInChannel  bool               `json:"show_in_channel"`
	Silent         bool               `json:"silent"`

	ExtraData map[string]interface{} `json:"-"`
}

type messageRequestForJSON messageRequestMessage

func (s *messageRequestMessage) UnmarshalJSON(data []byte) error {
	var s2 messageRequestForJSON
	if err := json.Unmarshal(data, &s2); err != nil {
		return err
	}
	*s = messageRequestMessage(s2)

	if err := json.Unmarshal(data, &s.ExtraData); err != nil {
		return err
	}

	removeFromMap(s.ExtraData, *s)
	return nil
}

func (s messageRequestMessage) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(s.ExtraData, messageRequestForJSON(s))
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

type attachmentForJSON Attachment

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attachment) UnmarshalJSON(data []byte) error {
	var a2 attachmentForJSON
	if err := json.Unmarshal(data, &a2); err != nil {
		return err
	}
	*a = Attachment(a2)

	if err := json.Unmarshal(data, &a.ExtraData); err != nil {
		return err
	}

	removeFromMap(a.ExtraData, *a)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (a Attachment) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(a.ExtraData, attachmentForJSON(a))
}

// SendMessageOption is an option that modifies behavior of send message request.
type SendMessageOption func(*messageRequest)

// MessageSkipPush is a flag that be given to SendMessage if you don't want to generate
// any push notifications.
func MessageSkipPush(r *messageRequest) {
	if r != nil {
		r.SkipPush = true
	}
}

// SendMessage sends a message to the channel. Returns full message details from server.
func (ch *Channel) SendMessage(message *Message, userID string, options ...SendMessageOption) (*Message, error) {
	switch {
	case message == nil:
		return nil, errors.New("message is nil")
	case userID == "":
		return nil, errors.New("user ID must be not empty")
	}

	message.User = &User{ID: userID}
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "message")

	req := message.toRequest()
	for _, op := range options {
		op(&req)
	}

	var resp messageResponse
	err := ch.client.makeRequest(http.MethodPost, p, nil, req, &resp)
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

// PartialUpdateMessage partially updates message with given msgID.
func (c *Client) PartialUpdateMessage(msgID string, updates PartialUpdate, options map[string]interface{}) (*Message, error) {
	switch {
	case len(updates.Set) == 0 && len(updates.Unset) == 0:
		return nil, errors.New("updates is empty")
	case options == nil:
		options = map[string]interface{}{}
	case msgID == "":
		return nil, errors.New("message ID must be not empty")
	}

	var resp messageResponse

	p := path.Join("messages", url.PathEscape(msgID))

	data := map[string]interface{}{
		"Set":   updates.Set,
		"Unset": updates.Unset,
	}
	for k, v := range options {
		data[k] = v
	}

	err := c.makeRequest(http.MethodPut, p, nil, data, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Message, nil
}

// PinMessage pins the message with given msgID.
func (c *Client) PinMessage(msgID, pinnedByID string, expiration *time.Duration) (*Message, error) {
	updates := PartialUpdate{
		Set: map[string]interface{}{
			"pinned": true,
		},
	}
	if expiration != nil {
		updates.Set["pin_expires"] = int(expiration.Milliseconds())
	}

	options := map[string]interface{}{
		"user_id": pinnedByID,
	}

	return c.PartialUpdateMessage(msgID, updates, options)
}

// UnPinMessage unpins the message with given msgID.
func (c *Client) UnPinMessage(msgID, userID string) (*Message, error) {
	updates := PartialUpdate{
		Set: map[string]interface{}{
			"pinned": false,
		},
	}

	options := map[string]interface{}{
		"user_id": userID,
	}

	return c.PartialUpdateMessage(msgID, updates, options)
}

func (c *Client) DeleteMessage(msgID string) error {
	return c.deleteMessage(msgID, false)
}

func (c *Client) HardDeleteMessage(msgID string) error {
	return c.deleteMessage(msgID, true)
}

func (c *Client) deleteMessage(msgID string, hard bool) error {
	if msgID == "" {
		return errors.New("message ID must be not empty")
	}
	p := path.Join("messages", url.PathEscape(msgID))

	params := map[string][]string{}
	if hard {
		params["hard"] = []string{"true"}
	}
	return c.makeRequest(http.MethodDelete, p, params, nil, nil)
}

type MessageFlag struct {
	CreatedByAutomod bool `json:"created_by_automod"`
	ModerationResult *struct {
		MessageID            string `json:"message_id"`
		Action               string `json:"action"`
		ModeratedBy          string `json:"moderated_by"`
		BlockedWord          string `json:"blocked_word"`
		BlocklistName        string `json:"blocklist_name"`
		ModerationThresholds *struct {
			Explicit *struct {
				Flag  float32 `json:"flag"`
				Block float32 `json:"block"`
			} `json:"explicit"`
			Spam *struct {
				Flag  float32 `json:"flag"`
				Block float32 `json:"block"`
			} `json:"spam"`
			Toxic *struct {
				Flag  float32 `json:"flag"`
				Block float32 `json:"block"`
			} `json:"toxic"`
		} `json:"moderation_thresholds"`
		AIModerationResponse *struct {
			Toxic    float32 `json:"toxic"`
			Explicit float32 `json:"explicit"`
			Spam     float32 `json:"spam"`
		} `json:"ai_moderation_response"`
		UserKarma    float64   `json:"user_karma"`
		UserBadKarma bool      `json:"user_bad_karma"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"moderation_result"`
	User    *User    `json:"user"`
	Message *Message `json:"message"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ReviewedAt time.Time `json:"reviewed_at"`
	ReviewedBy *User     `json:"reviewed_by"`
	ApprovedAt time.Time `json:"approved_at"`
	RejectedAt time.Time `json:"rejected_at"`
}

func (c *Client) FlagMessage(msgID, userID string) error {
	if msgID == "" {
		return errors.New("message ID is empty")
	}

	if userID == "" {
		return errors.New("user ID is empty")
	}

	options := map[string]interface{}{
		"target_message_id": msgID,
		"user_id":           userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/flag", nil, options, nil)
}

func (c *Client) UnflagMessage(msgID, userID string) error {
	if msgID == "" {
		return errors.New("message ID is empty")
	}

	if userID == "" {
		return errors.New("user ID is empty")
	}

	options := map[string]interface{}{
		"target_message_id": msgID,
		"user_id":           userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/unflag", nil, options, nil)
}

type repliesResponse struct {
	Messages []*Message `json:"messages"`
}

// GetReplies returns list of the message replies for a parent message.
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
