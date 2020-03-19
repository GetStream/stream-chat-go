package stream_chat //nolint: golint

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// Mute represents an user mute
type Mute struct {
	User      User      `json:"user"`
	Target    User      `json:"target"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChannelMute represents a channel mute
type ChannelMute struct {
	User      User       `json:"user"`
	Target    Channel    `json:"channel"`
	Expires   *time.Time `json:"expires"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Role  string `json:"role,omitempty"`

	Online    bool `json:"online,omitempty"`
	Invisible bool `json:"invisible,omitempty"`

	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	LastActive *time.Time `json:"last_active,omitempty"`

	Mutes       []*Mute                `json:"mutes,omitempty"`
	ChannelMute []*ChannelMute         `json:"channel_mutes,omitempty"`
	ExtraData   map[string]interface{} `json:"-,extra"` //nolint: staticcheck
}

// UnmarshalUnknown implements the `easyjson.UnknownsUnmarshaler` interface.
func (u *User) UnmarshalUnknown(in *jlexer.Lexer, key string) {
	if u.ExtraData == nil {
		u.ExtraData = make(map[string]interface{}, 1)
	}
	u.ExtraData[key] = in.Interface()
}

// MarshalUnknowns implements the `easyjson.UnknownsMarshaler` interface.
func (u User) MarshalUnknowns(out *jwriter.Writer, first bool) {
	for key, val := range u.ExtraData {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.String(key)
		out.RawByte(':')
		out.Raw(json.Marshal(val))
	}
}

// MuteUser create a mute
// targetID: the user getting muted
// userID: the user muting the target
func (c *Client) MuteUser(targetID, userID string) error {
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

// MuteUsers creates a mute
// targetID: the user getting muted
// userID: the user muting the target
func (c *Client) MuteUsers(targetIDs []string, userID string) error {
	switch {
	case len(targetIDs) == 0:
		return errors.New("target IDs are empty")
	case userID == "":
		return errors.New("user ID is empty")
	}

	data := map[string]interface{}{
		"target_ids": targetIDs,
		"user_id":    userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/mute", nil, data, nil)
}

// UnmuteUser removes a mute
// targetID: the user getting un-muted
// userID: the user muting the target
func (c *Client) UnmuteUser(targetID, userID string) error {
	switch {
	case targetID == "":
		return errors.New("target IDs is empty")
	case userID == "":
		return errors.New("user ID is empty")
	}

	data := map[string]interface{}{
		"target_id": targetID,
		"user_id":   userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/unmute", nil, data, nil)
}

// UnmuteUsers removes a mute
// targetID: the users getting un-muted
// userID: the user muting the target
func (c *Client) UnmuteUsers(targetIDs []string, userID string) error {
	switch {
	case len(targetIDs) == 0:
		return errors.New("target IDs is empty")
	case userID == "":
		return errors.New("user ID is empty")
	}

	data := map[string]interface{}{
		"target_ids": targetIDs,
		"user_id":    userID,
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

func (c *Client) BanUser(targetID, userID string, options map[string]interface{}) error {
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

	params := url.Values{}

	for k, v := range options {
		params.Add(k, v)
	}
	params.Set("target_user_id", targetID)

	return c.makeRequest(http.MethodDelete, "moderation/ban", params, nil, nil)
}

func (c *Client) ExportUser(targetID string, options map[string][]string) (user *User, err error) {
	if targetID == "" {
		return user, errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "export")
	user = &User{}

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

func (c *Client) ReactivateUser(targetID string, options map[string]interface{}) error {
	if targetID == "" {
		return errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "reactivate")

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
	Users map[string]userRequest `json:"users"`
}

type userRequest struct {
	*User

	// readonly fields
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	LastActive time.Time `json:"-"`
}

// UpdateUser sending update users request, returns updated user info
func (c *Client) UpdateUser(user *User) (*User, error) {
	users, err := c.UpdateUsers(user)
	return users[user.ID], err
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

// PartialUserUpdate request; Set and Unset fields can be set at same time, but should not be same field,
// for example you cannot set 'field.path.name' and unset 'field.path' at the same time.
// Field path should not contain spaces or dots (dot is path separator)
type PartialUserUpdate struct {
	ID    string                 `json:"id"`              // User ID, required
	Set   map[string]interface{} `json:"set,omitempty"`   // map of field.name => value; optional
	Unset []string               `json:"unset,omitempty"` // list of field names to unset
}

// PartialUpdateUser makes partial update for single user
func (c *Client) PartialUpdateUser(update PartialUserUpdate) (*User, error) {
	res, err := c.PartialUpdateUsers([]PartialUserUpdate{update})
	if err != nil {
		return nil, err
	}

	if user, ok := res[update.ID]; ok {
		return user, nil
	}

	return nil, fmt.Errorf("response error: no user with such ID in response: %s", update.ID)
}

type partialUserUpdateReq struct {
	Users []PartialUserUpdate `json:"users"`
}

// PartialUpdateUsers makes partial update for users
func (c *Client) PartialUpdateUsers(updates []PartialUserUpdate) (map[string]*User, error) {
	var resp usersResponse

	err := c.makeRequest(http.MethodPatch, "users", nil, partialUserUpdateReq{Users: updates}, &resp)

	return resp.Users, err
}
