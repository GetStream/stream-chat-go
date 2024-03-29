package stream_chat

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	AutoModDisabled modType = "disabled"
	AutoModSimple   modType = "simple"
	AutoModAI       modType = "AI"

	ModBehaviourFlag  modBehaviour = "flag"
	ModBehaviourBlock modBehaviour = "block"

	defaultMessageLength = 5000

	MessageRetentionForever = "infinite"
)

var defaultChannelTypes = []string{
	"messaging",
	"team",
	"livestream",
	"commerce",
	"gaming",
}

type (
	modType      string
	modBehaviour string
)

type ChannelTypePermission struct {
	Name   string `json:"name"`   // required
	Action string `json:"action"` // one of: Deny Allow

	Resources []string `json:"resources"` // required
	Roles     []string `json:"roles"`
	Owner     bool     `json:"owner"`
	Priority  int      `json:"priority"` // required
}

type ChannelType struct {
	ChannelConfig

	Commands []*Command `json:"commands"`
	// Deprecated: Use Permissions V2 API instead,
	// that can be found in permission_client.go.
	// See https://getstream.io/chat/docs/go-golang/migrating_from_legacy/?language=go
	Permissions []*ChannelTypePermission `json:"permissions"`
	Grants      map[string][]string      `json:"grants"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ct *ChannelType) toRequest() channelTypeRequest {
	req := channelTypeRequest{ChannelType: ct}

	for _, cmd := range ct.Commands {
		req.Commands = append(req.Commands, cmd.Name)
	}

	if len(req.Commands) == 0 {
		req.Commands = []string{"all"}
	}

	return req
}

// NewChannelType returns initialized ChannelType with default values.
func NewChannelType(name string) *ChannelType {
	ct := &ChannelType{ChannelConfig: DefaultChannelConfig}
	ct.Name = name

	return ct
}

type channelTypeRequest struct {
	*ChannelType

	Commands []string `json:"commands"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type ChannelTypeResponse struct {
	*ChannelType

	Commands []string `json:"commands"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Response
}

// CreateChannelType adds new channel type.
func (c *Client) CreateChannelType(ctx context.Context, chType *ChannelType) (*ChannelTypeResponse, error) {
	if chType == nil {
		return nil, errors.New("channel type is nil")
	}

	var resp ChannelTypeResponse

	err := c.makeRequest(ctx, http.MethodPost, "channeltypes", nil, chType.toRequest(), &resp)
	if err != nil {
		return nil, err
	}
	if resp.ChannelType == nil {
		return nil, errors.New("unexpected error: channel type response is nil")
	}

	for _, cmd := range resp.Commands {
		resp.ChannelType.Commands = append(resp.ChannelType.Commands, &Command{Name: cmd})
	}

	return &resp, nil
}

type GetChannelTypeResponse struct {
	*ChannelType
	Response
}

// GetChannelType returns information about channel type.
func (c *Client) GetChannelType(ctx context.Context, chanType string) (*GetChannelTypeResponse, error) {
	if chanType == "" {
		return nil, errors.New("channel type is empty")
	}

	p := path.Join("channeltypes", url.PathEscape(chanType))

	var resp GetChannelTypeResponse
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &resp)
	return &resp, err
}

type ChannelTypesResponse struct {
	ChannelTypes map[string]*ChannelType `json:"channel_types"`
	Response
}

// ListChannelTypes returns all channel types.
func (c *Client) ListChannelTypes(ctx context.Context) (*ChannelTypesResponse, error) {
	var resp ChannelTypesResponse
	err := c.makeRequest(ctx, http.MethodGet, "channeltypes", nil, nil, &resp)
	return &resp, err
}

// UpdateChannelType updates channel type.
func (c *Client) UpdateChannelType(ctx context.Context, name string, options map[string]interface{}) (*Response, error) {
	switch {
	case name == "":
		return nil, errors.New("channel type name is empty")
	case len(options) == 0:
		return nil, errors.New("options are empty")
	}

	p := path.Join("channeltypes", url.PathEscape(name))
	var resp Response
	err := c.makeRequest(ctx, http.MethodPut, p, nil, options, &resp)
	return &resp, err
}

// DeleteChannelType deletes channel type.
func (c *Client) DeleteChannelType(ctx context.Context, name string) (*Response, error) {
	if name == "" {
		return nil, errors.New("channel type name is empty")
	}

	p := path.Join("channeltypes", url.PathEscape(name))
	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, p, nil, nil, &resp)
	return &resp, err
}
