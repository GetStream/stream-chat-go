package stream_chat

import (
	"context"
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
	ID  string `json:"id"`
	CID string `json:"cid"`

	Text string `json:"text"`
	HTML string `json:"html"`

	Type   MessageType `json:"type,omitempty"` // one of MessageType* constants
	Silent bool        `json:"silent,omitempty"`

	User            *User          `json:"user"`
	UserID          string         `json:"user_id"`
	Attachments     []*Attachment  `json:"attachments"`
	LatestReactions []*Reaction    `json:"latest_reactions"` // last reactions
	OwnReactions    []*Reaction    `json:"own_reactions"`
	ReactionCounts  map[string]int `json:"reaction_counts"`
	ReactionScores  map[string]int `json:"reaction_scores"`

	ParentID           string  `json:"parent_id,omitempty"`       // id of parent message if it's reply
	ShowInChannel      bool    `json:"show_in_channel,omitempty"` // show reply message also in channel
	ThreadParticipants []*User `json:"thread_participants,omitempty"`

	ReplyCount      int       `json:"reply_count,omitempty"`
	QuotedMessage   *Message  `json:"quoted_message,omitempty"`
	QuotedMessageID string    `json:"quoted_message_id,omitempty"`
	MentionedUsers  []*User   `json:"mentioned_users"`

	Command string `json:"command,omitempty"`

	Shadowed   bool       `json:"shadowed,omitempty"`
	Pinned     bool       `json:"pinned,omitempty"`
	PinnedAt   *time.Time `json:"pinned_at,omitempty"`
	PinnedBy   *User      `json:"pinned_by,omitempty"`
	PinExpires *time.Time `json:"pin_expires,omitempty"`

	ImageModerationLabels map[string][]string `json:"image_labels,omitempty"`

	MML  string            `json:"mml,omitempty"`
	I18n map[string]string `json:"i18n,omitempty"`

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
		Text:            m.Text,
		Type:            m.Type,
		Attachments:     m.Attachments,
		UserID:          m.UserID,
		ExtraData:       m.ExtraData,
		Pinned:          m.Pinned,
		ParentID:        m.ParentID,
		ShowInChannel:   m.ShowInChannel,
		Silent:          m.Silent,
		QuotedMessageID: m.QuotedMessageID,
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
	Message                messageRequestMessage `json:"message"`
	SkipPush               bool                  `json:"skip_push,omitempty"`
	SkipEnrichURL          bool                  `json:"skip_enrich_url,omitempty"`
	Pending                bool                  `json:"pending,omitempty"`
	IsPendingMessage       bool                  `json:"is_pending_message,omitempty"`
	PendingMessageMetadata map[string]string     `json:"pending_message_metadata,omitempty"`
	KeepChannelHidden      bool                  `json:"keep_channel_hidden,omitempty"`
}

type messageRequestMessage struct {
	Text            string                 `json:"text"`
	Type            MessageType            `json:"type" validate:"omitempty,oneof=system"`
	Attachments     []*Attachment          `json:"attachments"`
	UserID          string                 `json:"user_id"`
	MentionedUsers  []string               `json:"mentioned_users"`
	ParentID        string                 `json:"parent_id"`
	ShowInChannel   bool                   `json:"show_in_channel"`
	Silent          bool                   `json:"silent"`
	QuotedMessageID string                 `json:"quoted_message_id"`
	Pinned          bool                   `json:"pinned"`
	ExtraData       map[string]interface{} `json:"-"`
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

// MessageSkipEnrichURL is a flag that disables enrichment of the URLs in the message
func MessageSkipEnrichURL(r *messageRequest) {
	if r != nil {
		r.SkipEnrichURL = true
	}
}

// MessagePending is a flag that makes this a pending message
func MessagePending(r *messageRequest) {
	if r != nil {
		r.Pending = true
	}
}

// MessagePendingMessageMetadata saves metadata to the pending message
func MessagePendingMessageMetadata(metadata map[string]string) SendMessageOption {
	return func(r *messageRequest) {
		if r != nil {
			r.PendingMessageMetadata = metadata
		}
	}
}

func KeepChannelHidden(r *messageRequest) {
	if r != nil {
		r.KeepChannelHidden = true
	}
}

type MessageResponse struct {
	Message                *Message          `json:"message"`
	PendingMessageMetadata map[string]string `json:"pending_message_metadata,omitempty"`
	Response
}

// SendMessage sends a message to the channel. Returns full message details from server.
func (ch *Channel) SendMessage(ctx context.Context, message *Message, userID string, options ...SendMessageOption) (*MessageResponse, error) {
	switch {
	case message == nil:
		return nil, errors.New("message is nil")
	case userID == "":
		return nil, errors.New("user ID must be not empty")
	}

	message.User = &User{ID: userID}
	message.UserID = userID
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "message")

	req := message.toRequest()
	for _, op := range options {
		op(&req)
	}

	var resp MessageResponse
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, req, &resp)
	return &resp, err
}

// MarkAllRead marks all messages as read for userID.
func (c *Client) MarkAllRead(ctx context.Context, userID string) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"user": map[string]string{
			"id": userID,
		},
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "channels/read", nil, data, &resp)
	return &resp, err
}

// GetMessage returns message by ID.
func (c *Client) GetMessage(ctx context.Context, msgID string) (*MessageResponse, error) {
	if msgID == "" {
		return nil, errors.New("message ID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(msgID))

	var resp MessageResponse
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &resp)
	return &resp, err
}

// UpdateMessage updates message with given msgID.
func (c *Client) UpdateMessage(ctx context.Context, msg *Message, msgID string) (*MessageResponse, error) {
	switch {
	case msg == nil:
		return nil, errors.New("message is nil")
	case msgID == "":
		return nil, errors.New("message ID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(msgID))

	var resp MessageResponse
	err := c.makeRequest(ctx, http.MethodPost, p, nil, msg.toRequest(), &resp)
	return &resp, err
}

type MessagePartialUpdateRequest struct {
	PartialUpdate
	UserID        string `json:"user_id"`
	SkipEnrichURL bool   `json:"skip_enrich_url"`
}

// PartialUpdateMessage partially updates message with given msgID.
func (c *Client) PartialUpdateMessage(ctx context.Context, messageID string, updates *MessagePartialUpdateRequest) (*MessageResponse, error) {
	switch {
	case len(updates.Set) == 0 && len(updates.Unset) == 0:
		return nil, errors.New("set or unset should not be empty")
	case messageID == "":
		return nil, errors.New("messageID should not be empty")
	}

	p := path.Join("messages", url.PathEscape(messageID))

	var resp MessageResponse
	err := c.makeRequest(ctx, http.MethodPut, p, nil, updates, &resp)
	return &resp, err
}

// PinMessage pins the message with given msgID.
func (c *Client) PinMessage(ctx context.Context, msgID, pinnedByID string, expiration *time.Time) (*MessageResponse, error) {
	updates := PartialUpdate{
		Set: map[string]interface{}{
			"pinned": true,
		},
	}
	if expiration != nil {
		updates.Set["pin_expires"] = expiration
	}

	request := MessagePartialUpdateRequest{
		PartialUpdate: updates,
		UserID:        pinnedByID,
	}

	return c.PartialUpdateMessage(ctx, msgID, &request)
}

// UnPinMessage unpins the message with given msgID.
func (c *Client) UnPinMessage(ctx context.Context, msgID, userID string) (*MessageResponse, error) {
	updates := PartialUpdate{
		Set: map[string]interface{}{
			"pinned": false,
		},
	}

	request := MessagePartialUpdateRequest{
		PartialUpdate: updates,
		UserID:        userID,
	}

	return c.PartialUpdateMessage(ctx, msgID, &request)
}

func (c *Client) CommitMessage(ctx context.Context, msgID string) (*Response, error) {
	if msgID == "" {
		return nil, errors.New("message ID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(msgID), "commit")
	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, p, nil, nil, &resp)
	return &resp, err

}

// DeleteMessage soft deletes the message with given msgID.
func (c *Client) DeleteMessage(ctx context.Context, msgID string) (*Response, error) {
	return c.deleteMessage(ctx, msgID, false)
}

// HardDeleteMessage deletes the message with given msgID. This is permanent.
func (c *Client) HardDeleteMessage(ctx context.Context, msgID string) (*Response, error) {
	return c.deleteMessage(ctx, msgID, true)
}

func (c *Client) deleteMessage(ctx context.Context, msgID string, hard bool) (*Response, error) {
	if msgID == "" {
		return nil, errors.New("message ID must be not empty")
	}
	p := path.Join("messages", url.PathEscape(msgID))

	params := map[string][]string{}
	if hard {
		params["hard"] = []string{"true"}
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, p, params, nil, &resp)
	return &resp, err
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

// FlagMessage flags the message with given msgID.
func (c *Client) FlagMessage(ctx context.Context, msgID, userID string) (*Response, error) {
	if msgID == "" {
		return nil, errors.New("message ID is empty")
	}

	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	options := map[string]interface{}{
		"target_message_id": msgID,
		"user_id":           userID,
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/flag", nil, options, &resp)
	return &resp, err
}

type RepliesResponse struct {
	Messages []*Message `json:"messages"`
	Response
}

// GetReplies returns list of the message replies for a parent message.
// options: Pagination params, ie {limit:10, idlte: 10}
func (ch *Channel) GetReplies(ctx context.Context, parentID string, options map[string][]string) (*RepliesResponse, error) {
	if parentID == "" {
		return nil, errors.New("parent ID is empty")
	}

	p := path.Join("messages", url.PathEscape(parentID), "replies")

	var resp RepliesResponse
	err := ch.client.makeRequest(ctx, http.MethodGet, p, options, nil, &resp)
	return &resp, err
}

type sendActionRequest struct {
	MessageID string            `json:"message_id"`
	FormData  map[string]string `json:"form_data"`
}

// SendAction for a message.
func (ch *Channel) SendAction(ctx context.Context, msgID string, formData map[string]string) (*MessageResponse, error) {
	switch {
	case msgID == "":
		return nil, errors.New("message ID is empty")
	case len(formData) == 0:
		return nil, errors.New("form data is empty")
	}

	p := path.Join("messages", url.PathEscape(msgID), "action")

	data := sendActionRequest{MessageID: msgID, FormData: formData}

	var resp MessageResponse
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

type TranslationResponse struct {
	Message *Message `json:"message"`
	Response
}

// TranslateMessage translates the message with given msgID to the given language.
func (c *Client) TranslateMessage(ctx context.Context, msgID, language string) (*TranslationResponse, error) {
	p := "messages/" + url.PathEscape(msgID) + "/translate"
	var resp TranslationResponse
	err := c.makeRequest(ctx, http.MethodPost, p, nil, map[string]string{"language": language}, &resp)

	return &resp, err
}
