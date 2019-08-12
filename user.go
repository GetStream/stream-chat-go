package stream_chat

import (
	"errors"
	"net/http"
	"time"
)

type Mute struct {
	User      User
	Target    User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Role  string `json:"role"`

	Online    bool `json:"online"`
	Invisible bool `json:"invisible"`

	Mutes []Mute `json:"mutes"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastActive time.Time `json:"last_active"`

	ExtraData map[string]interface{} `json:"-,extra"`
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

type usersResponse struct {
	Users map[string]User `json:"users"`
}

type usersRequest struct {
	Users map[string]userRequest `json:"Users"`
}

type userRequest struct {
	*User
	// readonly fields
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	LastActive time.Time `json:"-"`
}

// UpdateUsers send update users request; each user will be updated from response
func (c *Client) UpdateUsers(users ...*User) error {
	if len(users) == 0 {
		return errors.New("users are not set")
	}

	req := usersRequest{Users: make(map[string]userRequest, len(users))}
	for _, u := range users {
		req.Users[u.ID] = userRequest{User: u}
	}

	var resp usersResponse

	err := c.makeRequest(http.MethodPost, "users", nil, req, &resp)
	if err != nil {
		return err
	}

	for _, usr := range users {
		if u, ok := resp.Users[usr.ID]; ok {
			*usr = u
		}
	}

	return err
}
