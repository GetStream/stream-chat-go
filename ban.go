package stream_chat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

// BanUser bans targetID.
func (c *Client) BanUser(ctx context.Context, targetID, bannedBy string, options ...BanOption) (*Response, error) {
	switch {
	case targetID == "":
		return nil, errors.New("targetID should not be empty")
	case bannedBy == "":
		return nil, errors.New("bannedBy should not be empty")
	}

	opts := &banOptions{
		TargetUserID: targetID,
		BannedBy:     bannedBy,
	}

	for _, fn := range options {
		fn(opts)
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "moderation/ban", nil, opts, &resp)
	return &resp, err
}

// UnBanUser removes the ban for targetID.
func (c *Client) UnBanUser(ctx context.Context, targetID string, options ...UnbanOption) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("targetID should not be empty")
	}

	opts := &unbanOptions{}
	for _, fn := range options {
		fn(opts)
	}

	params := url.Values{}
	params.Set("target_user_id", targetID)
	if opts.RemoveFutureChannelsBan {
		params.Set("remove_future_channels_ban", "true")
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, "moderation/ban", params, nil, &resp)
	return &resp, err
}

// ShadowBan shadow bans targetID.
func (c *Client) ShadowBan(ctx context.Context, targetID, bannedByID string, options ...BanOption) (*Response, error) {
	options = append(options, banWithShadow())
	return c.BanUser(ctx, targetID, bannedByID, options...)
}

// BanUser bans targetID on the channel ch.
func (ch *Channel) BanUser(ctx context.Context, targetID, bannedBy string, options ...BanOption) (*Response, error) {
	options = append(options, banFromChannel(ch.Type, ch.ID))
	return ch.client.BanUser(ctx, targetID, bannedBy, options...)
}

// UnBanUser removes the ban for targetID from the channel ch.
func (ch *Channel) UnBanUser(ctx context.Context, targetID string, options ...UnbanOption) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("targetID should not be empty")
	}

	opts := &unbanOptions{}
	for _, fn := range options {
		fn(opts)
	}

	params := url.Values{}
	params.Set("target_user_id", targetID)
	params.Set("id", ch.ID)
	params.Set("type", ch.Type)
	if opts.RemoveFutureChannelsBan {
		params.Set("remove_future_channels_ban", "true")
	}

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodDelete, "moderation/ban", params, nil, &resp)
	return &resp, err
}

// ShadowBan shadow bans targetID on the channel ch.
func (ch *Channel) ShadowBan(ctx context.Context, targetID, bannedByID string, options ...BanOption) (*Response, error) {
	options = append(options, banWithShadow(), banFromChannel(ch.Type, ch.ID))
	return ch.client.ShadowBan(ctx, targetID, bannedByID, options...)
}

type QueryBannedUsersOptions struct {
	*QueryOption
}

type QueryBannedUsersResponse struct {
	Bans []*Ban `json:"bans"`
	Response
}

type Ban struct {
	Channel   *Channel   `json:"channel,omitempty"`
	User      *User      `json:"user"`
	Expires   *time.Time `json:"expires,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	Shadow    bool       `json:"shadow,omitempty"`
	BannedBy  *User      `json:"banned_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// QueryBannedUsers filters and returns a list of banned users.
// Banned users can be retrieved in different ways:
// 1) Using the dedicated query bans endpoint
// 2) User Search: you can add the banned:true condition to your search. Please note that
// this will only return users that were banned at the app-level and not the ones
// that were banned only on channels.
func (c *Client) QueryBannedUsers(ctx context.Context, q *QueryBannedUsersOptions, sorters ...*SortOption) (*QueryBannedUsersResponse, error) {
	qp := queryRequest{
		FilterConditions: q.Filter,
		Limit:            q.Limit,
		Offset:           q.Offset,
		Sort:             sorters,
	}

	data, err := json.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("payload", string(data))

	var resp QueryBannedUsersResponse
	err = c.makeRequest(ctx, http.MethodGet, "query_banned_users", values, nil, &resp)
	return &resp, err
}

type banOptions struct {
	Reason     string `json:"reason,omitempty"`
	Expiration int    `json:"timeout,omitempty"`

	TargetUserID          string `json:"target_user_id"`
	BannedBy              string `json:"user_id"`
	Shadow                bool   `json:"shadow"`
	BanFromFutureChannels bool   `json:"ban_from_future_channels,omitempty"`

	// ID and Type of the channel when acting on a channel member.
	ID   string `json:"id"`
	Type string `json:"type"`
}

type unbanOptions struct {
	RemoveFutureChannelsBan bool `json:"remove_future_channels_ban,omitempty"`
}

type BanOption func(*banOptions)

func BanWithReason(reason string) func(*banOptions) {
	return func(opt *banOptions) {
		opt.Reason = reason
	}
}

// BanWithExpiration set when the ban will expire. Should be in minutes.
// eg. to ban during one hour: BanWithExpiration(60).
func BanWithExpiration(expiration int) func(*banOptions) {
	return func(opt *banOptions) {
		opt.Expiration = expiration
	}
}

func banWithShadow() func(*banOptions) {
	return func(opt *banOptions) {
		opt.Shadow = true
	}
}

func banFromChannel(_type, id string) func(*banOptions) {
	return func(opt *banOptions) {
		opt.Type = _type
		opt.ID = id
	}
}

// BanWithBanFromFutureChannels when set to true, the user will be automatically
// banned from all future channels created by the user who issued the ban.
func BanWithBanFromFutureChannels() func(*banOptions) {
	return func(opt *banOptions) {
		opt.BanFromFutureChannels = true
	}
}

type UnbanOption func(*unbanOptions)

// UnbanWithRemoveFutureChannelsBan when set to true, also removes the future
// channel ban, so the user will no longer be auto-banned in new channels.
func UnbanWithRemoveFutureChannelsBan() func(*unbanOptions) {
	return func(opt *unbanOptions) {
		opt.RemoveFutureChannelsBan = true
	}
}

// FutureChannelBan represents a future channel ban entry.
type FutureChannelBan struct {
	User      *User      `json:"user"`
	Expires   *time.Time `json:"expires,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	Shadow    bool       `json:"shadow,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// QueryFutureChannelBansOptions contains options for querying future channel bans.
type QueryFutureChannelBansOptions struct {
	UserID             string `json:"user_id,omitempty"`
	ExcludeExpiredBans bool   `json:"exclude_expired_bans,omitempty"`
	Limit              int    `json:"limit,omitempty"`
	Offset             int    `json:"offset,omitempty"`
}

// QueryFutureChannelBansResponse is the response from QueryFutureChannelBans.
type QueryFutureChannelBansResponse struct {
	Bans []*FutureChannelBan `json:"bans"`
	Response
}

// QueryFutureChannelBans queries future channel bans.
// Future channel bans are automatically applied when a user creates a new channel
// or adds a member to an existing channel.
func (c *Client) QueryFutureChannelBans(ctx context.Context, opts *QueryFutureChannelBansOptions) (*QueryFutureChannelBansResponse, error) {
	data, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("payload", string(data))

	var resp QueryFutureChannelBansResponse
	err = c.makeRequest(ctx, http.MethodGet, "query_future_channel_bans", values, nil, &resp)
	return &resp, err
}
