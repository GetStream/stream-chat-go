package stream_chat

import (
	"net/http"
	"net/url"

	"github.com/getstream/easyjson"
)

type QueryOption struct {
	Filter map[string]interface{} `json:"-,extra"` // https://getstream.io/chat/docs/#query_syntax

	Limit  int `json:"limit,omitempty"`  // pagination option: limit number of results
	Offset int `json:"offset,omitempty"` // pagination option: offset to return items from
}

type SortOption struct {
	Field     string `json:"field"`     // field name to sort by,from json tags(in camel case), for example created_at
	Direction int    `json:"direction"` // [-1, 1]
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
	Channels []channelResponse `json:"channels"`
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

	channels := make([]*Channel, len(resp.Channels))
	for i, data := range resp.Channels {
		ch := c.newChannel()
		ch.update(data)
		channels[i] = ch
	}

	return channels, err
}
