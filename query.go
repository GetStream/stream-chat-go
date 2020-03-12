package stream

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/getstream/easyjson"
)

// QueryOption helps to build queries see
// https://getstream.io/chat/docs/#query_syntax
// for full details of the syntax.
type QueryOption struct {
	Filter map[string]interface{} `json:"-,extra"` //nolint: staticcheck

	Limit  int `json:"limit,omitempty"`  // pagination option: limit number of results
	Offset int `json:"offset,omitempty"` // pagination option: offset to return items from
}

type SortDirection = int

const (
	// SortAscending sets the specified field to sort in an ascending manner.
	SortAscending SortDirection = 1
	// SortDescending sets the specified field to sort in a descending manner.
	SortDescending SortDirection = -1
)

// SortField designates a string as being a sortable field. The default sort fields
// are all declared as constants.
type SortField = string

const (
	// SortFieldUserID lets you sort by the user's ID.
	SortFieldUserID SortField = "user_id"

	// SortFieldLastActive lets you sort by when the user was last active.
	SortFieldLastActive SortField = "last_active"
)

// SortOption is used to configure a sort field.
type SortOption struct {
	// Field is the naem of the field to sort by. This should be snake case.
	Field SortField `json:"field"`

	// Direction is the direction to sort the field by. This can be either
	// SortAscending (i.e. 1) or SortDescending (i.e. -1).
	Direction SortDirection `json:"direction"`
}

type queryUsersRequest struct {
	FilterConditions *QueryOption  `json:"filter_conditions,omitempty"`
	Sort             []*SortOption `json:"sort,omitempty"`
}

type queryUsersResponse struct {
	Users []*User `json:"users"`
}

// QueryUsers returns list of users that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in oder of sort options.
func (c *Client) QueryUsers(q *QueryOption, sort ...*SortOption) ([]*User, error) {
	qp := queryUsersRequest{
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

	return resp.Users, err
}

type queryChannelRequest struct {
	Watch    bool `json:"watch"`
	State    bool `json:"state"`
	Presence bool `json:"presence"`

	FilterConditions *QueryOption  `json:"filter_conditions,omitempty"`
	Sort             []*SortOption `json:"sort,omitempty"`
}

type queryChannelResponse struct {
	Channels []queryChannelResponseData `json:"channels"`
}

type queryChannelResponseData struct {
	Channel  *Channel         `json:"channel"`
	Messages []*Message       `json:"messages"`
	Read     []*ChannelRead   `json:"read"`
	Members  []*ChannelMember `json:"members"`
}

// QueryChannels returns list of channels with members and messages, that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in oder of sort options.
func (c *Client) QueryChannels(q *QueryOption, sort ...*SortOption) ([]*Channel, error) {
	qp := queryChannelRequest{
		State:            true,
		FilterConditions: q,
		Sort:             sort,
	}

	data, err := easyjson.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := make(url.Values)
	values.Set("payload", string(data))

	var resp queryChannelResponse
	err = c.makeRequest(http.MethodGet, "channels", values, nil, &resp)

	result := make([]*Channel, len(resp.Channels))
	for i, data := range resp.Channels {
		result[i] = data.Channel
		result[i].Members = data.Members
		result[i].Messages = data.Messages
		result[i].Read = data.Read
		result[i].client = c
	}

	return result, err
}

type SearchRequest struct {
	// Required
	Query   string                 `json:"query"`
	Filters map[string]interface{} `json:"filter_conditions"`

	// Pagination, optional
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type searchResponse struct {
	Results []searchMessageResponse `json:"results"`
}

type searchMessageResponse struct {
	Message *Message `json:"message"`
}

// Search returns channels matching for given keyword;
func (c *Client) Search(request SearchRequest) ([]*Message, error) {
	var buf strings.Builder

	_, err := easyjson.MarshalToWriter(request, &buf)
	if err != nil {
		return nil, err
	}

	var values = url.Values{}
	values.Set("payload", buf.String())

	var result searchResponse
	err = c.makeRequest(http.MethodGet, "search", values, nil, &result)

	messages := make([]*Message, 0, len(result.Results))
	for _, res := range result.Results {
		messages = append(messages, res.Message)
	}

	return messages, err
}
