//go:generate go run github.com/getstream/easyjson/easyjson -pkg -all
package stream_chat

import (
	"time"
)

// for interfaces type matching
var (
	_ StreamClient  = (*Client)(nil)
	_ StreamChannel = (*Channel)(nil)
)

type StreamClient interface {
	AddDevice(device *Device) error
	BanUser(targetID string, userID string, options map[string]interface{}) error
	CreateChannelType(chType *ChannelType) (*ChannelType, error)
	CreateToken(userID string, expire time.Time) ([]byte, error)
	DeactivateUser(targetID string, options map[string]interface{}) error
	DeleteChannelType(chType string) error
	DeleteDevice(userID string, deviceID string) error
	DeleteMessage(msgID string) error
	DeleteUser(targetID string, options map[string][]string) error
	ExportUser(targetID string, options map[string][]string) (user *User, err error)
	FlagUser(targetID string, options map[string]interface{}) error
	GetChannelType(chanType string) (ct *ChannelType, err error)
	GetDevices(userId string) (devices []*Device, err error)
	ListChannelTypes() (map[string]*ChannelType, error)
	MarkAllRead(userID string) error
	MuteUser(targetID string, userID string) error
	UnBanUser(targetID string, options map[string]string) error
	UnFlagUser(targetID string, options map[string]interface{}) error
	UnmuteUser(targetID string, userID string) error
	UpdateMessage(msg *Message, msgID string) (*Message, error)
	UpdateUsers(users ...*User) (map[string]*User, error)
}

type StreamChannel interface {
	AddMembers(userIDs ...string) error
	AddModerators(userIDs ...string) error
	BanUser(targetID string, userID string, options map[string]interface{}) error
	Delete() error
	DeleteReaction(messageID string, reactionType string, userID string) (*Message, error)
	DemoteModerators(userIDs ...string) error
	GetReactions(messageID string, options map[string][]string) ([]*Reaction, error)
	GetReplies(parentID string, options map[string][]string) (replies []*Message, err error)
	MarkRead(userID string, options map[string]interface{}) error
	RemoveMembers(userIDs ...string) error
	SendEvent(event *Event, userID string) error
	SendMessage(message *Message, userID string) (*Message, error)
	SendReaction(reaction *Reaction, messageID string, userID string) (*Message, error)
	Truncate() error
	UnBanUser(targetID string, options map[string]string) error
	Update(options map[string]interface{}, message string) error
}
