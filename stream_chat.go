//go:generate go run github.com/getstream/easyjson/easyjson -pkg -all

// package stream_chat //nolint:golint provides "chat as an API" via stream
package stream_chat //nolint:golint

import (
	"time"
)

// for interfaces type matching
var (
	_ StreamClient  = (*Client)(nil)
	_ StreamChannel = (*Channel)(nil)
)

// StreamClient is a client for chat
type StreamClient interface { //nolint:golint
	GetAppConfig() (*AppConfig, error)
	UpdateAppSettings(settings *AppSettings) error

	AddDevice(device *Device) error
	DeleteDevice(userID string, deviceID string) error
	GetDevices(userID string) (devices []*Device, err error)

	CreateChannel(chanType string, chanID string, userID string, data map[string]interface{}) (*Channel, error)

	CreateChannelType(chType *ChannelType) (*ChannelType, error)
	DeleteChannelType(chType string) error
	GetChannelType(chanType string) (ct *ChannelType, err error)
	ListChannelTypes() (map[string]*ChannelType, error)
	UpdateChannelType(name string, data map[string]interface{}) error

	CreateToken(userID string, expire time.Time) ([]byte, error)
	VerifyWebhook(body []byte, signature []byte) (valid bool)

	DeleteMessage(msgID string) error
	GetMessage(msgID string) (*Message, error)
	MarkAllRead(userID string) error
	UpdateMessage(msg *Message, msgID string) (*Message, error)
	FlagMessage(msgID string) error
	UnflagMessage(msgID string) error

	QueryUsers(q *QueryOption, sort ...*SortOption) ([]*User, error)
	QueryChannels(q *QueryOption, sort ...*SortOption) ([]*Channel, error)
	Search(request SearchRequest) ([]*Message, error)

	BanUser(targetID string, userID string, options ...Option) error
	DeactivateUser(targetID string, options ...Option) error
	ReactivateUser(targetID string, options ...Option) error
	DeleteUser(targetID string, options ...Option) error
	ExportUser(targetID string, options ...Option) (user *User, err error)
	FlagUser(targetID string, options ...Option) error
	MuteUser(targetID string, userID string) error
	MuteUsers(targetIDs []string, userID string) error
	UnBanUser(targetID string, options ...Option) error
	UnFlagUser(targetID string, options ...Option) error
	UnmuteUser(targetID string, userID string) error
	UnmuteUsers(targetIDs []string, userID string) error
	UpdateUser(user *User) (*User, error)
	UpdateUsers(users ...*User) (map[string]*User, error)
	PartialUpdateUser(update PartialUserUpdate) (*User, error)
	PartialUpdateUsers(updates []PartialUserUpdate) (map[string]*User, error)
}

// StreamChannel is a channel of communication
type StreamChannel interface { //nolint:golint

	AddMembers(userIDs []string, message *Message) error
	AddModerators(userIDs ...string) error
	AddModeratorsWithMessage(userIDs []string, msg *Message) error
	BanUser(targetID string, userID string, options ...Option) error
	Delete() error
	DemoteModerators(userIDs ...string) error
	DemoteModeratorsWithMessage(userIDs []string, msg *Message) error
	MarkRead(userID string, options ...Option) error
	RemoveMembers(userIDs []string, message *Message) error
	Truncate() error
	UnBanUser(targetID string, options ...Option) error
	Update(data map[string]interface{}, message *Message) error
	Query(data map[string]interface{}) error
	Show(userID string) error
	Hide(userID string) error
	HideWithHistoryClear(userID string) error
	InviteMembers(userIDs ...string) error
	InviteMembersWithMessage(userIDs []string, msg *Message) error
	SendFile(request SendFileRequest) (url string, err error)
	SendImage(request SendFileRequest) (url string, err error)
	DeleteFile(location string) error
	DeleteImage(location string) error
	AcceptInvite(userID string, message *Message) error
	RejectInvite(userID string, message *Message) error

	SendEvent(event *Event, userID string) error

	SendMessage(message *Message, userID string) (*Message, error)
	GetReplies(parentID string, options ...Option) (replies []*Message, err error)
	SendAction(msgID string, formData map[string]string) (*Message, error)

	DeleteReaction(messageID string, reactionType string, userID string) (*Message, error)
	GetReactions(messageID string, options ...Option) ([]*Reaction, error)
	SendReaction(reaction *Reaction, messageID string, userID string) (*Message, error)
}
