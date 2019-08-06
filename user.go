package stream_chat

import (
	"errors"
	"net/http"
	"time"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Role  string `json:"role"`

	Online    bool `json:"online"`
	Invisible bool `json:"invisible"`

	LastActive time.Time `json:"last_active"`

	Mutes []Mute `json:"mutes"`

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
	MuteUser(userID string, targetID string) error
	UnmuteUser(userID string, targetID string) error
	FlagUser(userID string, options ...interface{}) error
	UnFlagUser(userID string, options ...interface{}) error
	BanUser(id string, options map[string]interface{}) error
	UnBanUser(id string) error
	ExportUser(id string, options ...interface{}) (user interface{}, err error)
	DeactivateUser(id string, options ...interface{}) error
	DeleteUser(id string, options ...interface{}) error
	UpdateUser(id string, options ...interface{}) error
	// TODO: QueryUsers()
}

// Create a mute
// targetID: the user getting muted
// userID: the user muting the target
func (c *Client) MuteUser(targetID string, userID string) error {
	data := map[string]interface{}{
		"target_id": targetID,
		"user_id":   userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/mute", nil, data, nil)
}

// Removes a mute
// targetID: the user getting un-muted
// userID: the user muting the target
func (c *Client) UnmuteUser(targetID string, userID string) error {
	data := map[string]interface{}{
		"target_id": targetID,
		"user_id":   userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/unmute", nil, data, nil)
}

func (c *Client) FlagUser(targetID string, options map[string]interface{}) error {
	if options == nil || len(options) == 0 {
		return errors.New("flag user: options must be not empty")
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/flag", nil, options, nil)
}

func (c *Client) UnFlagUser(targetID string, options map[string]interface{}) error {
	if options == nil {
		return errors.New("flag user: options are nil")
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/unflag", nil, options, nil)
}

func (c *Client) BanUser(targetID string, options map[string]interface{}) error {
	if options == nil {
		options = map[string]interface{}{}
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/ban", nil, options, nil)
}

func (c *Client) UnBanUser(targetID string, options map[string]string) error {
	var params = map[string][]string{}

	for k, v := range options {
		params[k] = []string{v}
	}

	params["target_user_id"] = []string{targetID}

	return c.makeRequest(http.MethodDelete, "moderation/ban", params, nil, nil)
}

func (c *Client) ExportUser(targetID string, options map[string][]string) (user User, err error) {
	path := "users/" + targetID + "/export"

	err = c.makeRequest(http.MethodGet, path, options, nil, &user)

	return user, err
}

func (c *Client) DeactivateUser(targetID string, options map[string]interface{}) error {
	path := "users/" + targetID + "/deactivate"

	return c.makeRequest(http.MethodPost, path, nil, options, nil)
}

func (c *Client) DeleteUser(targetID string, options map[string][]string) error {
	path := "users/" + targetID

	return c.makeRequest(http.MethodDelete, path, options, nil, nil)
}

func (c *Client) UpdateUsers(users ...User) error {
	if len(users) == 0 {
		return errors.New("users are not set")
	}

	usersMap := make(map[string]User, len(users))
	for _, u := range users {
		usersMap[u.ID] = u
	}

	data := map[string]interface{}{
		"users": usersMap,
	}

	return c.makeRequest(http.MethodPost, "users", nil, data, nil)
}
