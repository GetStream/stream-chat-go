package stream_chat

import "time"

type EventType int

const (
	UserPresenceChanged EventType = iota + 1
	UserWatchingStart
	UserWatchingStop
	UserUpdated
	TypingStart
	TypingStop
	MessageNew
	MessageUpdated
	MessageDeleted
	MessageRead
	ReactionNew
	ReactionDeleted
	MemberAdded
	MemberUpdated
	MemberRemoved
	ChannelUpdated
	HealthCheck
	NotificationNewMessage
	NotificationMarkRead
	NotificationInvited
	NotificationInviteAccepted
	NotificationAddedToChannel
	NotificationRemovedFromChannel
	NotificationMutesUpdated
	ChannelDeleted
	// local events
	ConnectionChanged
	ConnectionRecovered
)

func (t EventType) String() string {
	return EventTypeMap[t]
}

var EventTypeMap = map[EventType]string{
	UserPresenceChanged:            "user.presence.changed",
	UserWatchingStart:              "user.watching.start",
	UserWatchingStop:               "user.watching.stop",
	UserUpdated:                    "user.updated",
	TypingStart:                    "typing.start",
	TypingStop:                     "typing.stop",
	MessageNew:                     "message.new",
	MessageUpdated:                 "message.updated",
	MessageDeleted:                 "message.deleted",
	MessageRead:                    "message.read",
	ReactionNew:                    "reaction.new",
	ReactionDeleted:                "reaction.deleted",
	MemberAdded:                    "member.added",
	MemberUpdated:                  "member.updated",
	MemberRemoved:                  "member.removed",
	ChannelUpdated:                 "channel.updated",
	ChannelDeleted:                 "channel.deleted",
	HealthCheck:                    "health.check",
	NotificationNewMessage:         "notification.message_new",
	NotificationMarkRead:           "notification.mark_read",
	NotificationInvited:            "notification.invited",
	NotificationInviteAccepted:     "notification.invite_accepted",
	NotificationAddedToChannel:     "notification.added_to_channel",
	NotificationRemovedFromChannel: "notification.removed_from_channel",
	NotificationMutesUpdated:       "notification.mutes_updated",

	// local events
	ConnectionChanged:   "connection.changed",
	ConnectionRecovered: "connection.recovered",
}

//var StringToEventType map[string]EventType
//
//func init() {
//	StringToEventType = make(map[string]EventType, len(EventTypeMap))
//	for k, v := range EventTypeMap {
//		StringToEventType[v] = k
//	}
//}

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
