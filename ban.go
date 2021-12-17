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
func (c *Client) UnBanUser(ctx context.Context, targetID string) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("targetID should not be empty")
	}

	params := url.Values{}
	params.Set("target_user_id", targetID)

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
	options = append(options, banFromChannel(ch.ID, ch.Type))
	return ch.client.BanUser(ctx, targetID, bannedBy, options...)
}

// UnBanUser removes the ban for targetID from the channel ch.
func (ch *Channel) UnBanUser(ctx context.Context, targetID string) (*Response, error) {
	if targetID == "" {
		return nil, errors.New("targetID should not be empty")
	}

	params := url.Values{}
	params.Set("target_user_id", targetID)
	params.Set("id", ch.ID)
	params.Set("type", ch.Type)

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodDelete, "moderation/ban", params, nil, &resp)
	return &resp, err
}

// ShadowBan shadow bans targetID on the channel ch.
func (ch *Channel) ShadowBan(ctx context.Context, targetID, bannedByID string, options ...BanOption) (*Response, error) {
	options = append(options, banWithShadow(), banFromChannel(ch.ID, ch.Type))
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

	values := make(url.Values)
	values.Set("payload", string(data))

	var resp QueryBannedUsersResponse
	err = c.makeRequest(ctx, http.MethodGet, "query_banned_users", values, nil, &resp)
	return &resp, err
}

type banOptions struct {
	Reason     string `json:"reason,omitempty"`
	Expiration int    `json:"timeout,omitempty"`

	TargetUserID string `json:"target_user_id"`
	BannedBy     string `json:"user_id"`
	Shadow       bool   `json:"shadow"`

	// ID and Type of the channel when acting on a channel member.
	ID   string `json:"id"`
	Type string `json:"type"`
}

type BanOption func(*banOptions)

func BanWithReason(reason string) func(*banOptions) {
	return func(opt *banOptions) {
		opt.Reason = reason
	}
}

// BanWithExpiration set when the ban will expire. Should be in minutes.
// eg. to ban during one hour: BanWithExpiration(60)
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

func banFromChannel(id, _type string) func(*banOptions) {
	return func(opt *banOptions) {
		opt.ID = id
		opt.Type = _type
	}
}
