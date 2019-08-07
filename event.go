package stream_chat

import "time"

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

	// local events
	EventConnectionChanged   EventType = "connection.changed"
	EventConnectionRecovered EventType = "connection.recovered"
)

type Event struct {
	// Channel ID
	CID          string         `json:"cid,omitempty"`
	Type         EventType      `json:"type"`
	Message      *Message       `json:"message,omitempty"`
	Reaction     *Reaction      `json:"reaction,omitempty"`
	Channel      *Channel       `json:"channel,omitempty"`
	Member       *ChannelMember `json:"member,omitempty"`
	User         *User          `json:"user,omitempty"`
	UserID       string         `json:"user_id,omitempty"`
	OwnUser      *User          `json:"me,omitempty"`
	WatcherCount int            `json:"watcher_count,omitempty"`

	ExtraData map[string]interface{} `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}

func (ev Event) toHash() map[string]interface{} {
	return nil
}
