package stream_chat

type UserID string

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
