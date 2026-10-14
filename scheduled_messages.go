package stream_chat

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

// ScheduledMessageStatus is the status of a scheduled message.
type ScheduledMessageStatus string

const (
	// ScheduledMessageStatusScheduled means the message waits for its send time.
	ScheduledMessageStatusScheduled ScheduledMessageStatus = "scheduled"
	// ScheduledMessageStatusSending means the message is being sent. It can no
	// longer be updated or deleted.
	ScheduledMessageStatusSending ScheduledMessageStatus = "sending"
	// ScheduledMessageStatusFailed means the message could not be sent. See
	// ScheduledMessage.FailureReason. A failed message can only be deleted.
	ScheduledMessageStatusFailed ScheduledMessageStatus = "failed"
	// ScheduledMessageStatusCanceled is reported when a scheduled message is deleted.
	ScheduledMessageStatusCanceled ScheduledMessageStatus = "canceled"
	// ScheduledMessageStatusSent is reported only by the scheduled_message.sent event.
	ScheduledMessageStatusSent ScheduledMessageStatus = "sent"
)

// ScheduledMessageFailureReason says why a scheduled message failed to send.
type ScheduledMessageFailureReason string

const (
	ScheduledMessageFailureBanned              ScheduledMessageFailureReason = "banned"
	ScheduledMessageFailureNoPermission        ScheduledMessageFailureReason = "no_permission"
	ScheduledMessageFailureChannelFrozen       ScheduledMessageFailureReason = "channel_frozen"
	ScheduledMessageFailureChannelDeleted      ScheduledMessageFailureReason = "channel_deleted"
	ScheduledMessageFailureUserDeleted         ScheduledMessageFailureReason = "user_deleted"
	ScheduledMessageFailureModerationBlocked   ScheduledMessageFailureReason = "moderation_blocked"
	ScheduledMessageFailureHookRejected        ScheduledMessageFailureReason = "hook_rejected"
	ScheduledMessageFailureExpired             ScheduledMessageFailureReason = "expired"
	ScheduledMessageFailureInternalError       ScheduledMessageFailureReason = "internal_error"
	ScheduledMessageFailureInvalidMessage      ScheduledMessageFailureReason = "invalid_message"
	ScheduledMessageFailurePendingNotSupported ScheduledMessageFailureReason = "pending_not_supported"
)

// ScheduledMessageContent is the message a scheduled message will send. It has
// no ID: the ID of the sent message is only known once it is sent (see the
// scheduled_message.sent event).
type ScheduledMessageContent struct {
	Text                 string                 `json:"text"`
	MML                  string                 `json:"mml,omitempty"`
	HTML                 string                 `json:"html,omitempty"`
	Type                 MessageType            `json:"type,omitempty"`
	Attachments          []*Attachment          `json:"attachments,omitempty"`
	ParentID             string                 `json:"parent_id,omitempty"`
	ShowInChannel        bool                   `json:"show_in_channel,omitempty"`
	QuotedMessageID      string                 `json:"quoted_message_id,omitempty"`
	MentionedUsers       []string               `json:"mentioned_users,omitempty"` // IDs of the mentioned users
	MentionedRoles       []string               `json:"mentioned_roles,omitempty"`
	MentionedGroupIDs    []string               `json:"mentioned_group_ids,omitempty"`
	MentionedChannel     bool                   `json:"mentioned_channel,omitempty"`
	MentionedHere        bool                   `json:"mentioned_here,omitempty"`
	RestrictedVisibility []string               `json:"restricted_visibility,omitempty"`
	Silent               bool                   `json:"silent,omitempty"`
	Pinned               bool                   `json:"pinned,omitempty"`
	PinExpires           *time.Time             `json:"pin_expires,omitempty"`
	PollID               string                 `json:"poll_id,omitempty"`
	SharedLocation       *SharedLocation        `json:"shared_location,omitempty"`
	Custom               map[string]interface{} `json:"custom,omitempty"`
}

// ScheduledMessageChannel identifies the channel of a scheduled message. It is
// only set in QueryScheduledMessages results.
type ScheduledMessageChannel struct {
	CID   string `json:"cid"`
	Type  string `json:"type"`
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}

// ScheduledMessage is a message stored to be sent later.
type ScheduledMessage struct {
	ID            string                        `json:"id"`
	CID           string                        `json:"cid"`
	User          *User                         `json:"user"`
	Message       *ScheduledMessageContent      `json:"message"`
	ScheduledAt   time.Time                     `json:"scheduled_at"`
	Status        ScheduledMessageStatus        `json:"status"`
	FailureReason ScheduledMessageFailureReason `json:"failure_reason,omitempty"`
	Channel       *ScheduledMessageChannel      `json:"channel,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WithScheduledAt schedules the message to be sent at the given time instead of
// sending it now. The response then carries ScheduledMessage, and Message holds
// the content as it will be sent; its ID is the scheduled message ID, not the ID
// of the sent message.
//
// The send time must be at least one minute and at most the app's maximum lead
// time (120 days by default) from now. Scheduled messages must be enabled for
// the app.
func WithScheduledAt(scheduledAt time.Time) SendMessageOption {
	return func(r *messageRequest) {
		if r != nil {
			t := scheduledAt.UTC()
			r.ScheduledAt = &t
		}
	}
}

// ScheduledMessageResponse is the response of the scheduled message endpoints.
type ScheduledMessageResponse struct {
	ScheduledMessage *ScheduledMessage `json:"scheduled_message"`
	Response
}

type updateScheduledMessageRequest struct {
	Message     *messageRequestMessage `json:"message,omitempty"`
	ScheduledAt *time.Time             `json:"scheduled_at,omitempty"`
}

// UpdateScheduledMessage replaces the message and/or the send time of a
// scheduled message. At least one of message and scheduledAt must be set.
// message, when set, replaces the scheduled message as a whole; its ID and user
// are ignored, the sender stays the user who scheduled it.
func (c *Client) UpdateScheduledMessage(ctx context.Context, scheduledMessageID string, message *Message, scheduledAt *time.Time) (*ScheduledMessageResponse, error) {
	if scheduledMessageID == "" {
		return nil, errors.New("scheduled message ID is empty")
	}
	if message == nil && scheduledAt == nil {
		return nil, errors.New("message or scheduledAt is required")
	}

	var req updateScheduledMessageRequest
	if message != nil {
		m := message.toRequest().Message
		m.ID = ""
		req.Message = &m
	}
	if scheduledAt != nil {
		t := scheduledAt.UTC()
		req.ScheduledAt = &t
	}

	p := path.Join("scheduled_messages", url.PathEscape(scheduledMessageID))

	var resp ScheduledMessageResponse
	err := c.makeRequest(ctx, http.MethodPatch, p, nil, req, &resp)
	return &resp, err
}

// DeleteScheduledMessage cancels a scheduled message so it is never sent, or
// removes a failed one.
func (c *Client) DeleteScheduledMessage(ctx context.Context, scheduledMessageID string) (*ScheduledMessageResponse, error) {
	if scheduledMessageID == "" {
		return nil, errors.New("scheduled message ID is empty")
	}

	p := path.Join("scheduled_messages", url.PathEscape(scheduledMessageID))

	var resp ScheduledMessageResponse
	err := c.makeRequest(ctx, http.MethodDelete, p, nil, nil, &resp)
	return &resp, err
}

// SendScheduledMessageNow moves the send time of a scheduled message to now and
// queues the send.
func (c *Client) SendScheduledMessageNow(ctx context.Context, scheduledMessageID string) (*ScheduledMessageResponse, error) {
	if scheduledMessageID == "" {
		return nil, errors.New("scheduled message ID is empty")
	}

	p := path.Join("scheduled_messages", url.PathEscape(scheduledMessageID), "send_now")

	var resp ScheduledMessageResponse
	err := c.makeRequest(ctx, http.MethodPost, p, nil, map[string]interface{}{}, &resp)
	return &resp, err
}

// ScheduledMessagesFilter narrows QueryScheduledMessages. Empty fields mean no
// restriction.
type ScheduledMessagesFilter struct {
	// CID returns only the scheduled messages of this channel.
	CID string `json:"cid,omitempty"`
	// Status returns only these statuses: scheduled, sending and/or failed.
	Status []ScheduledMessageStatus `json:"status,omitempty"`
}

// QueryScheduledMessagesOptions represents the options for the QueryScheduledMessages request.
type QueryScheduledMessagesOptions struct {
	// UserID is the user whose scheduled messages are listed. Required.
	UserID string `json:"user_id"`

	Filter *ScheduledMessagesFilter `json:"filter,omitempty"`

	// Sort supports only scheduled_at ascending.
	Sort []*SortOption `json:"sort,omitempty"`

	// Limit is the page size, 1 to 100. Defaults to 20.
	Limit int `json:"limit,omitempty"`
	// Next is the 'next' value of a previous response, to continue from that point.
	Next string `json:"next,omitempty"`
}

// QueryScheduledMessagesResponse is one page of scheduled messages.
type QueryScheduledMessagesResponse struct {
	// ScheduledMessages lists pending messages first, then failed ones, each by scheduled_at.
	ScheduledMessages []*ScheduledMessage `json:"scheduled_messages"`
	// Next is to be used as the 'next' option to get the next page. Empty on the last page.
	Next string `json:"next,omitempty"`
	Response
}

// QueryScheduledMessages lists the scheduled and failed messages of one user
// across channels.
func (c *Client) QueryScheduledMessages(ctx context.Context, options *QueryScheduledMessagesOptions) (*QueryScheduledMessagesResponse, error) {
	if options == nil || options.UserID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	var resp QueryScheduledMessagesResponse
	err := c.makeRequest(ctx, http.MethodPost, "scheduled_messages/query", nil, options, &resp)
	return &resp, err
}
