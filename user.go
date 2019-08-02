package stream_chat

import "time"

type UserID string

type User struct {
	ID    UserID `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Role  string `json:"role"`

	Online    bool `json:"online"`
	Invisible bool `json:"invisible"`

	LastActive time.Time `json:"last_active"`

	Mutes []*Mute `json:"mutes"`

	ExtraData map[string]interface{}

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) toHash() map[string]interface{} {
	return nil
}

func (u *User) MarshalJSON() (data []byte, err error) {
	return
}

type UserAPI interface {
	MuteUser(userID UserID, targetID UserID) error
	UnmuteUser(userID UserID, targetID UserID) error
	FlagUser(userID UserID, options ...interface{}) error
	UnFlagUser(userID UserID, options ...interface{}) error
	BanUser(id UserID, options map[string]interface{}) error
	UnBanUser(id UserID) error
	ExportUser(id UserID, options ...interface{}) (user interface{}, err error)
	DeactivateUser(id UserID, options ...interface{}) error
	DeleteUser(id UserID, options ...interface{}) error
	UpdateUser(id UserID, options ...interface{}) error
	// TODO: QueryUsers()
}

func (*client) MuteUser(userID UserID, targetID UserID) error {
	panic("implement me")
}

func (*client) UnmuteUser(userID UserID, targetID UserID) error {
	panic("implement me")
}

func (*client) FlagUser(userID UserID, options ...interface{}) error {
	panic("implement me")
}

func (*client) UnFlagUser(userID UserID, options ...interface{}) error {
	panic("implement me")
}

func (*client) BanUser(id UserID, options map[string]interface{}) error {
	panic("implement me")
}

func (*client) UnBanUser(id UserID) error {
	panic("implement me")
}

func (*client) ExportUser(id UserID, options ...interface{}) (user interface{}, err error) {
	panic("implement me")
}

func (*client) DeactivateUser(id UserID, options ...interface{}) error {
	panic("implement me")
}

func (*client) DeleteUser(id UserID, options ...interface{}) error {
	panic("implement me")
}

func (*client) UpdateUser(id UserID, options ...interface{}) error {
	panic("implement me")
}
