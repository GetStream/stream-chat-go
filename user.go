package stream_chat

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/getstream/easyjson"
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

	Mutes []*Mute `json:"mutes"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastActive time.Time `json:"last_active"`

	ExtraData map[string]interface{} `json:"-,extra"`
}

// Create a mute
// targetID: the user getting muted
// userID: the user muting the target
func (c *Client) MuteUser(targetID string, userID string) error {
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case userID == "":
		return errors.New("user ID is empty")
	}

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
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case userID == "":
		return errors.New("user ID is empty")
	}

	data := map[string]interface{}{
		"target_id": targetID,
		"user_id":   userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/unmute", nil, data, nil)
}

func (c *Client) FlagUser(targetID string, options map[string]interface{}) error {
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case len(options) == 0:
		return errors.New("flag user: options must be not empty")
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/flag", nil, options, nil)
}

func (c *Client) UnFlagUser(targetID string, options map[string]interface{}) error {
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case options == nil:
		options = map[string]interface{}{}
	}

	options["target_user_id"] = targetID

	return c.makeRequest(http.MethodPost, "moderation/unflag", nil, options, nil)
}

func (c *Client) BanUser(targetID string, userID string, options map[string]interface{}) error {
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case userID == "":
		return errors.New("user ID is empty")
	case options == nil:
		options = map[string]interface{}{}
	}

	options["target_user_id"] = targetID
	options["user_id"] = userID

	return c.makeRequest(http.MethodPost, "moderation/ban", nil, options, nil)
}

func (c *Client) UnBanUser(targetID string, options map[string]string) error {
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case options == nil:
		options = map[string]string{}
	}

	var params = map[string][]string{}

	for k, v := range options {
		params[k] = []string{v}
	}

	params["target_user_id"] = []string{targetID}

	return c.makeRequest(http.MethodDelete, "moderation/ban", params, nil, nil)
}

func (c *Client) ExportUser(targetID string, options map[string][]string) (user *User, err error) {
	if targetID == "" {
		return user, errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "export")

	err = c.makeRequest(http.MethodGet, p, options, nil, user)

	return user, err
}

func (c *Client) DeactivateUser(targetID string, options map[string]interface{}) error {
	if targetID == "" {
		return errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "deactivate")

	return c.makeRequest(http.MethodPost, p, nil, options, nil)
}

func (c *Client) DeleteUser(targetID string, options map[string][]string) error {
	if targetID == "" {
		return errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID))

	return c.makeRequest(http.MethodDelete, p, options, nil, nil)
}

type usersResponse struct {
	Users map[string]*User `json:"users"`
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

// UpdateUsers send update users request, returns updated user info
func (c *Client) UpdateUsers(users ...*User) (map[string]*User, error) {
	if len(users) == 0 {
		return nil, errors.New("users are not set")
	}

	req := usersRequest{Users: make(map[string]userRequest, len(users))}
	for _, u := range users {
		req.Users[u.ID] = userRequest{User: u}
	}

	var resp usersResponse

	err := c.makeRequest(http.MethodPost, "users", nil, req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Users, err
}

type QueryOption struct {
	Query map[string]interface{} `json:"-,extra"` // https://getstream.io/chat/docs/#query_syntax

	PaginationOption
}

type PaginationOption struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type SortOption struct {
	Field     string `json:"field"`
	Direction int    `json:"direction"` // [-1, 1]
}

type queryParams struct {
	Watch    bool `json:"watch,omitempty"`
	State    bool `json:"state,omitempty"`
	Presence bool `json:"presence,omitempty"`

	FilterConditions *QueryOption  `json:"filter_conditions,omitempty"`
	Sort             []*SortOption `json:"sort,omitempty"`
}

type queryUsersResponse struct {
	Users []*User `json:"users"`
}

func (c *Client) QueryUsers(q *QueryOption, sort ...*SortOption) (map[string]*User, error) {
	qp := queryParams{
		State:            true,
		Watch:            false,
		Presence:         false,
		FilterConditions: q,
		Sort:             sort,
	}

	data, err := easyjson.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := make(url.Values)
	values.Set("payload", string(data))

	var resp queryUsersResponse
	err = c.makeRequest(http.MethodGet, "users", values, nil, &resp)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*User, len(resp.Users))
	for _, u := range resp.Users {
		result[u.ID] = u
	}

	return result, nil
}
