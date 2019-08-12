package stream_chat

import (
	"net/http"
	"net/url"
	"path"
	"time"
)

type EventType string

const (
	EventUserPresenceChanged            EventType = "user.presence.changed"
	EventUserWatchingStart              EventType = "user.watching.start"
	EventUserWatchingStop               EventType = "user.watching.stop"
	EventUserUpdated                    EventType = "user.updated"
	EventTypingStart                    EventType = "typing.start"
	EventTypingStop                     EventType = "typing.stop"
	EventMessageNew                     EventType = "message.new"
	EventMessageUpdated                 EventType = "message.updated"
	EventMessageDeleted                 EventType = "message.deleted"
	EventMessageRead                    EventType = "message.read"
	EventReactionNew                    EventType = "reaction.new"
	EventReactionDeleted                EventType = "reaction.deleted"
	EventMemberAdded                    EventType = "member.added"
	EventMemberUpdated                  EventType = "member.updated"
	EventMemberRemoved                  EventType = "member.removed"
	EventChannelUpdated                 EventType = "channel.updated"
	EventChannelDeleted                 EventType = "channel.deleted"
	EventHealthCheck                    EventType = "health.check"
	EventNotificationNewMessage         EventType = "notification.message_new"
	EventNotificationMarkRead           EventType = "notification.mark_read"
	EventNotificationInvited            EventType = "notification.invited"
	EventNotificationInviteAccepted     EventType = "notification.invite_accepted"
	EventNotificationAddedToChannel     EventType = "notification.added_to_channel"
	EventNotificationRemovedFromChannel EventType = "notification.removed_from_channel"
	EventNotificationMutesUpdated       EventType = "notification.mutes_updated"
)

type Event struct {
	CID          string         `json:"cid,omitempty"` // Channel ID
	Type         EventType      `json:"type"`          // Event type, one of Event* constants
	Message      *Message       `json:"message,omitempty"`
	Reaction     *Reaction      `json:"reaction,omitempty"`
	Channel      *Channel       `json:"channel,omitempty"`
	Member       *ChannelMember `json:"member,omitempty"`
	User         *User          `json:"user,omitempty"`
	UserID       string         `json:"user_id,omitempty"`
	OwnUser      *User          `json:"me,omitempty"`
	WatcherCount int            `json:"watcher_count,omitempty"`

	ExtraData map[string]interface{} `json:"-"`

	CreatedAt time.Time `json:"created_at,omitempty"`
}

type eventRequest struct {
	Event Event `json:"event"`
}

// SendEvent sends an event on this channel
//
// event: event data, ie {type: 'message.read'}
// userID: the ID of the user sending the event
func (ch *Channel) SendEvent(event Event, userID string) error {
	if event.User == nil {
		event.User = &User{ID: userID}
	}

	req := eventRequest{Event: event}
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "event")

	return ch.client.makeRequest(http.MethodPost, p, nil, req, nil)
}
