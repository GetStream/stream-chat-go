package stream

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Mute struct {
	User      User      `json:"user"`
	Target    User      `json:"target"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	ExtraData map[string]interface{} `json:"-,extra"` //nolint: staticcheck

	Mutes []*Mute `json:"mutes,omitempty"`
}

// MuteUser create a mute
// targetID: the user getting muted
// userID: the user muting the target
func (c *Client) MuteUser(targetID, userID string) error {
	switch {
	case targetID == "":
		return ErrorMissingTargetID
	case userID == "":
		return ErrorMissingUserID
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
		return ErrorMissingUserID
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
		return ErrorMissingTargetID
	case userID == "":
		return ErrorMissingUserID
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
		return ErrorEmptyTargetID
	case userID == "":
		return ErrorMissingUserID
	}

	data := map[string]interface{}{
		"target_ids": targetIDs,
		"user_id":    userID,
	}

	return c.makeRequest(http.MethodPost, "moderation/unmute", nil, data, nil)
}

func (c *Client) FlagUser(targetID string, options ...Option) error {
	switch {
	case targetID == "":
		return ErrorMissingTargetID
	case len(options) == 0:
		return errors.New("flag user: options must be not empty")
	}

	options = append(options, NewOption(optionKeyTargetUserID, targetID))

	return c.makeRequestWithOptions(http.MethodPost, "moderation/flag", nil, options, nil)
}

func (c *Client) UnFlagUser(targetID string, options ...Option) error {
	switch {
	case targetID == "":
		return ErrorMissingTargetID
	}

	options = append(options, NewOption(optionKeyTargetUserID, targetID))

	return c.makeRequestWithOptions(http.MethodPost, "moderation/unflag", nil, options, nil)
}

func (c *Client) BanUser(targetID, userID string, options ...Option) error {
	return c.banUser(&banUserInput{
		TargetID: targetID,
		UserID:   userID,
	}, options...)
}

type banUserInput struct {
	// TargetID is the ID of the person to be banned.
	TargetID string `json:"target_user_id,omitempty"`

	// UserID is the ID of the person doing the banning.
	UserID string `json:"user_id,omitempty"`

	Reason  *string `json:"reason,omitempty"`
	Timeout *int    `json:"timeout,omitempty"`

	// ChannelType and ChannelID are used if the ban is channel specific.
	ChannelType *string `json:"type,omitempty"`
	ChannelID   *string `json:"id,omitempty"`

	Extra map[string]interface{} `json:"-"`
}

func (b *banUserInput) AddOptions(options []Option) {
	for _, opt := range options {
		switch opt.Key() {
		case optionKeyTimeout:
			seconds := opt.Value().(int)
			b.Timeout = &seconds
		case optionKeyReason:
			reason := optionAsString(opt)
			b.Reason = &reason
		}
	}
}

func (b banUserInput) validate() error {
	switch {
	case b.TargetID == "":
		return ErrorMissingTargetID
	case b.UserID == "":
		return ErrorMissingUserID
	}

	return nil
}

func (c *Client) banUser(input *banUserInput, options ...Option) error {
	if err := input.validate(); err != nil {
		return err
	}

	input.AddOptions(options)

	return c.makeRequestWithOptions(http.MethodPost, "moderation/ban", nil, input, nil)
}

func (c *Client) UnBanUser(targetID string, options ...Option) error {
	switch {
	case targetID == "":
		return ErrorMissingTargetID
	}

	options = append(options, makeTargetID(targetID))

	return c.makeRequestWithOptions(http.MethodDelete, "moderation/ban", options, nil, nil)
}

func (c *Client) ExportUser(targetID string, options ...Option) (*User, error) {
	if targetID == "" {
		return nil, ErrorMissingTargetID
	}

	p := path.Join("users", url.PathEscape(targetID), "export")

	user := &User{}
	err := c.makeRequestWithOptions(http.MethodGet, p, options, nil, user)

	return user, err
}

func (c *Client) DeactivateUser(targetID string, options ...Option) error {
	if targetID == "" {
		return ErrorMissingTargetID
	}

	p := path.Join("users", url.PathEscape(targetID), "deactivate")

	return c.makeRequestWithOptions(http.MethodPost, p, nil, options, nil)
}

// ReactivateUser reactivates targetID.
func (c *Client) ReactivateUser(targetID string, options ...Option) error {
	if targetID == "" {
		return ErrorMissingTargetID
	}

	p := path.Join("users", url.PathEscape(targetID), "reactivate")

	return c.makeRequestWithOptions(http.MethodPost, p, nil, options, nil)
}

// DeleteUser deletes the targetID. See the Option documentation for more
// details on the options that can be set.
func (c *Client) DeleteUser(targetID string, options ...Option) error {
	if targetID == "" {
		return ErrorMissingTargetID
	}

	p := path.Join("users", url.PathEscape(targetID))

	return c.makeRequestWithOptions(http.MethodDelete, p, options, nil, nil)
}

type usersResponse struct {
	Users map[string]*User `json:"users"`
}

type usersRequest struct {
	Users map[string]userRequest `json:"users"`
}

type userRequest struct {
	*User

	// extra data doesn't work for embedded structs
	ExtraData map[string]interface{} `json:"-,extra"` //nolint: staticcheck

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
		req.Users[u.ID] = userRequest{User: u, ExtraData: u.ExtraData}
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
