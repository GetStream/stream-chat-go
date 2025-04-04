package stream_chat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type QueryOption struct {
	// https://getstream.io/chat/docs/#query_syntax
	Filter map[string]interface{} `json:"filter_conditions,omitempty"`
	Sort   []*SortOption          `json:"sort,omitempty"`

	UserID       string `json:"user_id,omitempty"`
	Limit        int    `json:"limit,omitempty"`  // pagination option: limit number of results
	Offset       int    `json:"offset,omitempty"` // pagination option: offset to return items from
	MessageLimit *int   `json:"message_limit,omitempty"`
	MemberLimit  *int   `json:"member_limit,omitempty"`
}

type SortOption struct {
	Field     string `json:"field"`     // field name to sort by,from json tags(in camel case), for example created_at
	Direction int    `json:"direction"` // [-1, 1]
}

type queryRequest struct {
	Watch    bool `json:"watch"`
	State    bool `json:"state"`
	Presence bool `json:"presence"`

	UserID                  string `json:"user_id,omitempty"`
	Limit                   int    `json:"limit,omitempty"`
	Offset                  int    `json:"offset,omitempty"`
	MemberLimit             *int   `json:"member_limit,omitempty"`
	MessageLimit            *int   `json:"message_limit,omitempty"`
	IncludeDeactivatedUsers bool   `json:"include_deactivated_users,omitempty"`

	FilterConditions map[string]interface{} `json:"filter_conditions,omitempty"`
	Sort             []*SortOption          `json:"sort,omitempty"`
}

type QueryFlagReportsRequest struct {
	FilterConditions map[string]interface{} `json:"filter_conditions,omitempty"`
	Limit            int                    `json:"limit,omitempty"`
	Offset           int                    `json:"offset,omitempty"`
}

type FlagReport struct {
	ID            string    `json:"id"`
	Message       *Message  `json:"message"`
	FlagsCount    int       `json:"flags_count"`
	MessageUserID string    `json:"message_user_id"`
	ChannelCid    string    `json:"channel_cid"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type QueryUsersOptions struct {
	QueryOption

	IncludeDeactivatedUsers bool `json:"include_deactivated_users"`
}

type QueryUsersResponse struct {
	Users []*User `json:"users"`
	Response
}

// QueryUsers returns list of users that match QueryUsersOptions.
// If any number of SortOption are set, result will be sorted by field and direction in the order of sort options.
func (c *Client) QueryUsers(ctx context.Context, q *QueryUsersOptions, sorters ...*SortOption) (*QueryUsersResponse, error) {
	qp := queryRequest{
		FilterConditions:        q.Filter,
		Limit:                   q.Limit,
		Offset:                  q.Offset,
		IncludeDeactivatedUsers: q.IncludeDeactivatedUsers,
		Sort:                    sorters,
	}

	data, err := json.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("payload", string(data))

	var resp QueryUsersResponse
	err = c.makeRequest(ctx, http.MethodGet, "users", values, nil, &resp)
	return &resp, err
}

type queryChannelResponse struct {
	Channels []queryChannelResponseData `json:"channels"`
	Response
}

type queryChannelResponseData struct {
	Channel         *Channel         `json:"channel"`
	Messages        []*Message       `json:"messages"`
	Read            []*ChannelRead   `json:"read"`
	Members         []*ChannelMember `json:"members"`
	PendingMessages []*Message       `json:"pending_messages"`
	PinnedMessages  []*Message       `json:"pinned_messages"`
}

type QueryChannelsResponse struct {
	Channels []*Channel
	Response
}

// QueryChannels returns list of channels with members and messages, that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in oder of sort options.
func (c *Client) QueryChannels(ctx context.Context, q *QueryOption, sort ...*SortOption) (*QueryChannelsResponse, error) {
	qp := queryRequest{
		State:            true,
		FilterConditions: q.Filter,
		Sort:             sort,
		UserID:           q.UserID,
		Limit:            q.Limit,
		Offset:           q.Offset,
		MemberLimit:      q.MemberLimit,
		MessageLimit:     q.MessageLimit,
	}

	var resp queryChannelResponse
	if err := c.makeRequest(ctx, http.MethodPost, "channels", nil, qp, &resp); err != nil {
		return nil, err
	}

	result := make([]*Channel, len(resp.Channels))
	for i, data := range resp.Channels {
		result[i] = data.Channel
		result[i].Members = data.Members
		result[i].Messages = data.Messages
		result[i].PendingMessages = data.PendingMessages
		result[i].PinnedMessages = data.PinnedMessages
		result[i].Read = data.Read
		result[i].client = c
	}

	return &QueryChannelsResponse{Channels: result, Response: resp.Response}, nil
}

type SearchRequest struct {
	// Required
	Query          string                 `json:"query"`
	Filters        map[string]interface{} `json:"filter_conditions"`
	MessageFilters map[string]interface{} `json:"message_filter_conditions"`

	// Pagination, optional
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Next   string `json:"next,omitempty"`

	// Sort, optional
	Sort []SortOption `json:"sort,omitempty"`
}

type SearchFullResponse struct {
	Results  []SearchMessageResponse `json:"results"`
	Next     string                  `json:"next,omitempty"`
	Previous string                  `json:"previous,omitempty"`
	Response
}

type SearchMessageResponse struct {
	Message *Message `json:"message"`
}

type SearchResponse struct {
	Messages []*Message
	Response
}

// Search returns messages matching for given keyword.
func (c *Client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	result, err := c.SearchWithFullResponse(ctx, request)
	if err != nil {
		return nil, err
	}
	messages := make([]*Message, 0, len(result.Results))
	for _, res := range result.Results {
		messages = append(messages, res.Message)
	}

	resp := SearchResponse{
		Messages: messages,
		Response: result.Response,
	}
	return &resp, nil
}

// SearchWithFullResponse performs a search and returns the full results.
func (c *Client) SearchWithFullResponse(ctx context.Context, request SearchRequest) (*SearchFullResponse, error) {
	if request.Offset != 0 {
		if len(request.Sort) > 0 || request.Next != "" {
			return nil, errors.New("cannot use Offset with Next or Sort parameters")
		}
	}
	if request.Query != "" && len(request.MessageFilters) != 0 {
		return nil, errors.New("can only specify Query or MessageFilters, not both")
	}
	var buf strings.Builder

	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("payload", buf.String())

	var result SearchFullResponse
	if err := c.makeRequest(ctx, http.MethodGet, "search", values, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type QueryMessageFlagsResponse struct {
	Flags []*MessageFlag `json:"flags"`
	Response
}

// QueryMessageFlags returns list of message flags that match QueryOption.
func (c *Client) QueryMessageFlags(ctx context.Context, q *QueryOption) (*QueryMessageFlagsResponse, error) {
	qp := queryRequest{
		FilterConditions: q.Filter,
		Limit:            q.Limit,
		Offset:           q.Offset,
	}

	data, err := json.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("payload", string(data))

	var resp QueryMessageFlagsResponse
	err = c.makeRequest(ctx, http.MethodGet, "moderation/flags/message", values, nil, &resp)
	return &resp, err
}

type QueryFlagReportsResponse struct {
	Response
	FlagReports []*FlagReport `json:"flag_reports"`
}

// QueryFlagReports returns list of flag reports that match the filter.
func (c *Client) QueryFlagReports(ctx context.Context, q *QueryFlagReportsRequest) (*QueryFlagReportsResponse, error) {
	resp := &QueryFlagReportsResponse{}
	err := c.makeRequest(ctx, http.MethodPost, "moderation/reports", nil, q, &resp)
	return resp, err
}
