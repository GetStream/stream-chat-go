package stream_chat //nolint: golint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Mute represents a user mute.
type Mute struct {
	User      User       `json:"user"`
	Target    User       `json:"target"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Expires   *time.Time `json:"expires"`
}

// ChannelMute represents a channel mute.
type ChannelMute struct {
	User      User       `json:"user"`
	Channel   Channel    `json:"channel"`
	Expires   *time.Time `json:"expires"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name,omitempty"`
	Image string   `json:"image,omitempty"`
	Role  string   `json:"role,omitempty"`
	Teams []string `json:"teams,omitempty"`

	Online    bool `json:"online,omitempty"`
	Invisible bool `json:"invisible,omitempty"`

	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	LastActive *time.Time `json:"last_active,omitempty"`

	Mutes                    []*Mute                `json:"mutes,omitempty"`
	ChannelMutes             []*ChannelMute         `json:"channel_mutes,omitempty"`
	ExtraData                map[string]interface{} `json:"-"`
	RevokeTokensIssuedBefore *time.Time             `json:"revoke_tokens_issued_before,omitempty"`
}

type userForJSON User

// UnmarshalJSON implements json.Unmarshaler.
func (u *User) UnmarshalJSON(data []byte) error {
	var u2 userForJSON
	if err := json.Unmarshal(data, &u2); err != nil {
		return err
	}
	*u = User(u2)

	if err := json.Unmarshal(data, &u.ExtraData); err != nil {
		return err
	}

	removeFromMap(u.ExtraData, *u)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (u User) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(u.ExtraData, userForJSON(u))
}

type muteOptions struct {
	Expiration int `json:"timeout,omitempty"`

	TargetID  string   `json:"target_id"`
	TargetIDs []string `json:"target_ids"`
	UserID    string   `json:"user_id"`
}

type MuteOption func(*muteOptions)

func MuteWithExpiration(expiration int) func(*muteOptions) {
	return func(opt *muteOptions) {
		opt.Expiration = expiration
	}
}

// MuteUser mutes targetID.
func (c *Client) MuteUser(ctx context.Context, targetID, mutedBy string, options ...MuteOption) (*Response, error) {
	switch {
	case targetID == "":
		return nil, errors.New("targetID should not be empty")
	case mutedBy == "":
		return nil, errors.New("mutedBy should not be empty")
	}

	opts := &muteOptions{
		TargetID: targetID,
		UserID:   mutedBy,
	}

	for _, fn := range options {
		fn(opts)
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/mute", nil, opts, &resp)
	return &resp, err
}

// MuteUsers mutes all users in targetIDs.
func (c *Client) MuteUsers(ctx context.Context, targetIDs []string, mutedBy string, options ...MuteOption) (*Response, error) {
	switch {
	case len(targetIDs) == 0:
		return nil, errors.New("targetIDs should not be empty")
	case mutedBy == "":
		return nil, errors.New("mutedBy should not be empty")
	}

	opts := &muteOptions{
		TargetIDs: targetIDs,
		UserID:    mutedBy,
	}

	for _, fn := range options {
		fn(opts)
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/mute", nil, opts, &resp)
	return &resp, err
}

// UnmuteUser unmute targetID.
func (c *Client) UnmuteUser(ctx context.Context, targetID, unmutedBy string) (*Response, error) {
	switch {
	case targetID == "":
		return nil, errors.New("targetID should not be empty")
	case unmutedBy == "":
		return nil, errors.New("unmutedBy should not be empty")
	}

	opts := &muteOptions{
		TargetID: targetID,
		UserID:   unmutedBy,
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/unmute", nil, opts, &resp)
	return &resp, err
}

// UnmuteUsers unmute all users in targetIDs
func (c *Client) UnmuteUsers(ctx context.Context, targetIDs []string, unmutedBy string) (*Response, error) {
	switch {
	case len(targetIDs) == 0:
		return nil, errors.New("target IDs is empty")
	case unmutedBy == "":
		return nil, errors.New("user ID is empty")
	}

	opts := &muteOptions{
		TargetIDs: targetIDs,
		UserID:    unmutedBy,
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/unmute", nil, opts, &resp)
	return &resp, err
}

func (c *Client) FlagUser(ctx context.Context, targetID, flaggedBy string) (*Response, error) {
	switch {
	case targetID == "":
		return nil, errors.New("targetID should not be empty")
	case flaggedBy == "":
		return nil, errors.New("flaggedBy should not be empty")
	}

	options := map[string]string{
		"target_user_id": targetID,
		"user_id":        flaggedBy,
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/flag", nil, options, &resp)
	return &resp, err
}

type ExportUserResponse struct {
	*User
	Response
}

func (c *Client) ExportUser(ctx context.Context, targetID string) (*ExportUserResponse, error) {
	if targetID == "" {
		return nil, errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "export")

	var resp ExportUserResponse
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &resp)
	return &resp, err
}

func (c *Client) DeactivateUser(ctx context.Context, targetID string, options map[string]interface{}) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "deactivate")

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, p, nil, options, &resp)
	return &resp, err
}

func (c *Client) ReactivateUser(ctx context.Context, targetID string, options map[string]interface{}) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("target ID is empty")
	}

	p := path.Join("users", url.PathEscape(targetID), "reactivate")

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, p, nil, options, &resp)
	return &resp, err
}

type deleteUserOptions struct {
	MarkMessagesDeleted string
	HardDelete          string
	DeleteConversations string
}

type DeleteUserOption func(*deleteUserOptions)

func DeleteUserWithHardDelete() func(*deleteUserOptions) {
	return func(opt *deleteUserOptions) {
		opt.HardDelete = "true"
	}
}

func DeleteUserWithMarkMessagesDeleted() func(*deleteUserOptions) {
	return func(opt *deleteUserOptions) {
		opt.MarkMessagesDeleted = "true"
	}
}

func DeleteUserWithDeleteConversations() func(*deleteUserOptions) {
	return func(opt *deleteUserOptions) {
		opt.DeleteConversations = "true"
	}
}

func (c *Client) DeleteUser(ctx context.Context, targetID string, options ...DeleteUserOption) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("targetID should not be empty")
	}

	option := &deleteUserOptions{}
	for _, fn := range options {
		fn(option)
	}

	params := url.Values{}
	params.Set("mark_messages_deleted", option.MarkMessagesDeleted)
	params.Set("hard_delete", option.HardDelete)
	params.Set("delete_conversation_channels", option.DeleteConversations)

	p := path.Join("users", url.PathEscape(targetID))

	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, p, params, nil, &resp)
	return &resp, err
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

type UpsertUserResponse struct {
	User *User
	Response
}

// UpsertUser is a single user version of UpsertUsers for convenience.
func (c *Client) UpsertUser(ctx context.Context, user *User) (*UpsertUserResponse, error) {
	resp, err := c.UpsertUsers(ctx, user)
	return &UpsertUserResponse{
		User:     resp.Users[user.ID],
		Response: resp.Response,
	}, err
}

type UsersResponse struct {
	Users map[string]*User `json:"users"`
	Response
}

// UpsertUsers creates the given users. If a user doesn't exist, it will be created.
// Otherwise, custom data will be extended or updated. Missing data is never removed.
func (c *Client) UpsertUsers(ctx context.Context, users ...*User) (*UsersResponse, error) {
	if len(users) == 0 {
		return nil, errors.New("users are not set")
	}

	req := usersRequest{Users: make(map[string]userRequest, len(users))}
	for _, u := range users {
		req.Users[u.ID] = userRequest{User: u}
	}

	var resp UsersResponse
	err := c.makeRequest(ctx, http.MethodPost, "users", nil, req, &resp)
	return &resp, err
}

// PartialUserUpdate request; Set and Unset fields can be set at same time, but should not be same field,
// for example you cannot set 'field.path.name' and unset 'field.path' at the same time.
// Field path should not contain spaces or dots (dot is path separator).
type PartialUserUpdate struct {
	ID    string                 `json:"id"`              // User ID, required
	Set   map[string]interface{} `json:"set,omitempty"`   // map of field.name => value; optional
	Unset []string               `json:"unset,omitempty"` // list of field names to unset
}

// PartialUpdateUser makes partial update for single user.
func (c *Client) PartialUpdateUser(ctx context.Context, update PartialUserUpdate) (*User, error) {
	resp, err := c.PartialUpdateUsers(ctx, []PartialUserUpdate{update})
	if err != nil {
		return nil, err
	}

	if user, ok := resp.Users[update.ID]; ok {
		return user, nil
	}

	return nil, fmt.Errorf("response error: no user with such ID in response: %s", update.ID)
}

type partialUserUpdateReq struct {
	Users []PartialUserUpdate `json:"users"`
}

// PartialUpdateUsers makes partial update for users.
func (c *Client) PartialUpdateUsers(ctx context.Context, updates []PartialUserUpdate) (*UsersResponse, error) {
	var resp UsersResponse

	err := c.makeRequest(ctx, http.MethodPatch, "users", nil, partialUserUpdateReq{Users: updates}, &resp)
	return &resp, err
}

// RevokeUserToken revoke token for a user issued before given time.
func (c *Client) RevokeUserToken(ctx context.Context, userID string, before *time.Time) (*Response, error) {
	return c.RevokeUsersTokens(ctx, []string{userID}, before)
}

// RevokeUsersTokens revoke tokens for users issued before given time.
func (c *Client) RevokeUsersTokens(ctx context.Context, userIDs []string, before *time.Time) (*Response, error) {
	userUpdates := make([]PartialUserUpdate, 0)
	for _, userID := range userIDs {
		userUpdate := PartialUserUpdate{
			ID:  userID,
			Set: make(map[string]interface{}),
		}
		if before == nil {
			userUpdate.Set["revoke_tokens_issued_before"] = nil
		} else {
			userUpdate.Set["revoke_tokens_issued_before"] = before.Format(time.RFC3339)
		}
		userUpdates = append(userUpdates, userUpdate)
	}

	resp, err := c.PartialUpdateUsers(ctx, userUpdates)
	return &resp.Response, err
}
