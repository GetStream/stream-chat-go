package stream_chat // nolint: golint

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

// EventType marks which of the various sub-types of a webhook event you are
// receiving or sending.
type EventType string

const (
	// EventMessageNew is fired when a new message is added.
	EventMessageNew EventType = "message.new"
	// EventMessageUpdated is fired when a message is updated.
	EventMessageUpdated EventType = "message.updated"
	// EventMessageDeleted is fired when a message is deleted.
	EventMessageDeleted EventType = "message.deleted"
	// EventMessageRead is fired when a user calls mark as read.
	EventMessageRead EventType = "message.read"

	// EventReactionNew is fired when a message reaction is added.
	EventReactionNew EventType = "reaction.new"
	// EventReactionDeleted is fired when a message reaction deleted.
	EventReactionDeleted EventType = "reaction.deleted"

	// EventMemberAdded is fired when a member is added to a channel.
	EventMemberAdded EventType = "member.added"
	// EventMemberUpdated is fired when a member is updated.
	EventMemberUpdated EventType = "member.updated"
	// EventMemberRemoved is fired when a member is removed from a channel.
	EventMemberRemoved EventType = "member.removed"

	// EventChannelUpdated is fired when a channel is updated.
	EventChannelUpdated EventType = "channel.updated"
	// EventChannelDeleted is fired when a channel is deleted.
	EventChannelDeleted EventType = "channel.deleted"
	// EventChannelTruncated is fired when a channel is truncated.
	EventChannelTruncated EventType = "channel.truncated"

	// EventHealthCheck is fired when a user is updated.
	EventHealthCheck EventType = "health.check"

	// EventNotificationNewMessage and family are fired when a notification is
	// created, marked read, invited to a channel, and so on.
	EventNotificationNewMessage         EventType = "notification.message_new"
	EventNotificationMarkRead           EventType = "notification.mark_read"
	EventNotificationInvited            EventType = "notification.invited"
	EventNotificationInviteAccepted     EventType = "notification.invite_accepted"
	EventNotificationAddedToChannel     EventType = "notification.added_to_channel"
	EventNotificationRemovedFromChannel EventType = "notification.removed_from_channel"
	EventNotificationMutesUpdated       EventType = "notification.mutes_updated"

	// EventTypingStart and EventTypingStop are fired when a user starts or stops typing.
	EventTypingStart EventType = "typing.start"
	EventTypingStop  EventType = "typing.stop"

	// EventUserMuted is fired when a user is muted.
	EventUserMuted EventType = "user.muted"
	// EventUserUnmuted is fired when a user is unmuted.
	EventUserUnmuted         EventType = "user.unmuted"
	EventUserPresenceChanged EventType = "user.presence.changed"
	EventUserWatchingStart   EventType = "user.watching.start"
	EventUserWatchingStop    EventType = "user.watching.stop"
	EventUserUpdated         EventType = "user.updated"
)

// Event is received from a webhook, or sent with the SendEvent function.
type Event struct {
	CID          string           `json:"cid,omitempty"` // Channel ID
	Type         EventType        `json:"type"`          // Event type, one of Event* constants
	Message      *Message         `json:"message,omitempty"`
	Reaction     *Reaction        `json:"reaction,omitempty"`
	Channel      *Channel         `json:"channel,omitempty"`
	Member       *ChannelMember   `json:"member,omitempty"`
	Members      []*ChannelMember `json:"members,omitempty"`
	User         *User            `json:"user,omitempty"`
	UserID       string           `json:"user_id,omitempty"`
	OwnUser      *User            `json:"me,omitempty"`
	WatcherCount int              `json:"watcher_count,omitempty"`

	ExtraData map[string]interface{} `json:"-"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

type eventForJSON Event

func (e *Event) UnmarshalJSON(data []byte) error {
	var e2 eventForJSON
	if err := json.Unmarshal(data, &e2); err != nil {
		return err
	}
	*e = Event(e2)

	if err := json.Unmarshal(data, &e.ExtraData); err != nil {
		return err
	}

	removeFromMap(e.ExtraData, *e)
	return nil
}

func (e Event) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(e.ExtraData, eventForJSON(e))
}

type eventRequest struct {
	Event *Event `json:"event"`
}

// SendEvent sends an event on this channel.
func (ch *Channel) SendEvent(event *Event, userID string) error {
	if event == nil {
		return errors.New("event is nil")
	}

	event.User = &User{ID: userID}

	req := eventRequest{Event: event}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "event")

	return ch.client.makeRequest(http.MethodPost, p, nil, req, nil)
}

// SendUserCustomEvent sends a custom event to all connected clients for the user userID.
func (c *Client) SendUserCustomEvent(event *Event, userID string) error {
	if event == nil {
		return errors.New("event is nil")
	}
	if userID == "" {
		return errors.New("userID should not be empty")
	}

	req := eventRequest{Event: event}

	p := path.Join("users", url.PathEscape(userID), "event")

	return c.makeRequest(http.MethodPost, p, nil, req, nil)
}
