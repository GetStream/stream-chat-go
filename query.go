package stream_chat // nolint: golint

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

type Map map[string]interface{}

type QueryOption struct {
	// https://getstream.io/chat/docs/#query_syntax
	Filter Map           `json:"filter_conditions,omitempty"`
	Sort   []*SortOption `json:"sort,omitempty"`

	UserID string `json:"user_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`  // pagination option: limit number of results
	Offset int    `json:"offset,omitempty"` // pagination option: offset to return items from
}

// UnmarshalUnknown implements the `easyjson.UnknownsUnmarshaler` interface.
func (q *QueryOption) UnmarshalUnknown(in *jlexer.Lexer, key string) {
	if q.Filter == nil {
		q.Filter = make(map[string]interface{}, 1)
	}
	q.Filter[key] = in.Interface()
}

// MarshalUnknowns implements the `easyjson.UnknownsMarshaler` interface.
func (q QueryOption) MarshalUnknowns(out *jwriter.Writer, first bool) {
	for key, val := range q.Filter {
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

type SortOption struct {
	Field     string `json:"field"`     // field name to sort by,from json tags(in camel case), for example created_at
	Direction int    `json:"direction"` // [-1, 1]
}

type queryRequest struct {
	Watch    bool `json:"watch"`
	State    bool `json:"state"`
	Presence bool `json:"presence"`

	UserID string `json:"user_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`

	FilterConditions map[string]interface{} `json:"filter_conditions,omitempty"`
	Sort             []*SortOption          `json:"sort,omitempty"`
}

type queryUsersResponse struct {
	Users []*User `json:"users"`
}

// QueryUsers returns list of users that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in oder of sort options.
func (c *Client) QueryUsers(q *QueryOption, sort ...*SortOption) ([]*User, error) {
	qp := queryRequest{
		FilterConditions: q.Filter,
		Sort:             sort,
	}

	data, err := json.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := make(url.Values)
	values.Set("payload", string(data))

	var resp queryUsersResponse
	err = c.makeRequest(http.MethodGet, "users", values, nil, &resp)

	return resp.Users, err
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
	qp := queryRequest{
		State:            true,
		FilterConditions: q.Filter,
		Sort:             sort,
		UserID:           q.UserID,
		Limit:            q.Limit,
		Offset:           q.Offset,
	}

	data, err := json.Marshal(&qp)
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

// PaginationOptions are passed to the PageQuery* functions. They are optional
// and allow you to configure some aspects of it.
type PaginationOptions struct {
	Limit          int
	StartingOffset int
}

// getOptions compresses multiple options into one, and makes sure defaults are
// set.
func getPageOptions(options []PaginationOptions) PaginationOptions {
	result := PaginationOptions{}

	for _, opt := range options {
		if opt.StartingOffset != 0 {
			result.StartingOffset = opt.StartingOffset
		}
		if opt.Limit != 0 {
			result.Limit = opt.Limit
		}
	}

	if result.Limit == 0 {
		result.Limit = 30
	}

	return result
}

// ChannelPaginationFunc is a function that is given a list of channels. It
// returns true if it wants another page of channels, false if you wish to
// stop.
type ChannelPaginationFunc func([]*Channel) bool

// UserPaginationFunc is a function that is given a list of users. It
// returns true if it wants another page of users, false if you wish to
// stop.
type UserPaginationFunc func([]*User) bool

// PageQueryChannels allows you to paginate through a query results. It takes a
// paginationFunc and an optional set of PaginationOptions.
func (c *Client) PageQueryChannels(q *QueryOption, paginationFunc ChannelPaginationFunc, options ...PaginationOptions) error {
	if paginationFunc == nil {
		return errors.New("must pass a pagination function")
	}

	opt := getPageOptions(options)

	// If we are given a set of queryOptions then make a copy of it so we don't
	// mess with upstreams.
	var newQ *QueryOption
	if q != nil {
		newQ = &QueryOption{}
		*newQ = *q
	}

	for i := (opt.StartingOffset / opt.Limit); ; i++ {
		newQ.Limit = opt.Limit
		newQ.Offset = opt.Limit * i

		res, err := c.QueryChannels(newQ)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			return nil
		}

		// If the paginationFunc indicates we should not continue then we stop.
		if ok := paginationFunc(res); !ok {
			return nil
		}
	}
}

// PageQueryUsers allows you to paginate through a query results. It takes a
// paginationFunc and an optional set of PaginationOptions.
func (c *Client) PageQueryUsers(q *QueryOption, paginationFunc UserPaginationFunc, options ...PaginationOptions) error {
	if paginationFunc == nil {
		return errors.New("must pass a pagination function")
	}

	opt := getPageOptions(options)

	// If we are given a set of queryOptions then make a copy of it so we don't
	// mess with upstreams.
	var newQ *QueryOption
	if q != nil {
		newQ = &QueryOption{}
		*newQ = *q
	}

	for i := (opt.StartingOffset / opt.Limit); ; i++ {
		newQ.Limit = opt.Limit
		newQ.Offset = opt.Limit * i

		res, err := c.QueryUsers(newQ)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			return nil
		}

		// If the paginationFunc indicates we should not continue then we stop.
		if ok := paginationFunc(res); !ok {
			return nil
		}
	}
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
