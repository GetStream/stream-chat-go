package stream_chat

import (
	"errors"
	"net/http"
	"time"
)

type User struct {
	ID    string
	Name  string
	Image string
	Role  string

	Online    bool
	Invisible bool

	Mutes Mutes

	ExtraData map[string]interface{}

	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastActive time.Time
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
	if len(options) == 0 {
		return errors.New("flag user: options must be not empty")
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/flag", nil, options, nil)
}

func (c *Client) UnFlagUser(targetID string, options map[string]interface{}) error {
	if options == nil {
		options = map[string]interface{}{}
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/unflag", nil, options, nil)
}

func (c *Client) BanUser(targetID string, userID string, options map[string]interface{}) error {
	if options == nil {
		options = map[string]interface{}{}
	}

	options["target_user_id"] = targetID
	options["user_id"] = userID

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

// UpdateUsers send update users request; each user will be updated from response
func (c *Client) UpdateUsers(users ...*User) error {
	if len(users) == 0 {
		return errors.New("users are not set")
	}

	// users search table for unmarshal
	usersMap := map[string]*User{}

	payload := map[string]map[string]interface{}{
		"users": {},
	}

	// marshal users
	for _, u := range users {
		usersMap[u.ID] = u
		payload["users"][u.ID] = u.marshalMap()
	}

	var resp struct{ Users map[string]interface{} }

	err := c.makeRequest(http.MethodPost, "users", nil, payload, &resp)
	if err != nil {
		return err
	}

	for k, v := range resp.Users {
		switch val := v.(type) {
		case map[string]interface{}:
			usersMap[k].unmarshalMap(val)
		default:
			// TODO: logging
		}
	}

	return err
}
