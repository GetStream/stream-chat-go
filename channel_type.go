package stream

import (
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

type modType string
type modBehaviour string

type Permission struct {
	Name   string `json:"name"`   // required
	Action string `json:"action"` // one of: Deny Allow

	Resources []string `json:"resources"` // required
	Roles     []string `json:"roles"`
	Owner     bool     `json:"owner"`
	Priority  int      `json:"priority"` // required
}

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Args        string `json:"args"`
	Set         string `json:"set"`
}

type ChannelType struct {
	ChannelConfig

	Commands    []*Command    `json:"commands"`
	Permissions []*Permission `json:"permissions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ct *ChannelType) toRequest() channelTypeRequest {
	req := channelTypeRequest{ChannelType: ct}

	if len(req.Commands) == 0 {
		req.Commands = []string{"all"}
	}

	return req
}

// NewChannelType returns initialized ChannelType with default values
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

type channelTypeResponse struct {
	ChannelTypes map[string]*ChannelType `json:"channel_types"`
}

// CreateChannelType adds new channel type
func (c *Client) CreateChannelType(chType *ChannelType) (*ChannelType, error) {
	if chType == nil {
		return nil, errors.New("channel type is nil")
	}

	var resp channelTypeRequest

	err := c.makeRequest(http.MethodPost, "channeltypes", nil, chType.toRequest(), &resp)
	if err != nil {
		return nil, err
	}
	if resp.ChannelType == nil {
		return nil, errors.New("unexpected error: channel type response is nil")
	}

	for _, cmd := range resp.Commands {
		resp.ChannelType.Commands = append(resp.ChannelType.Commands, &Command{Name: cmd})
	}

	return resp.ChannelType, nil
}

// GetChannelType returns information about channel type
func (c *Client) GetChannelType(chanType string) (*ChannelType, error) {
	if chanType == "" {
		return nil, errors.New("channel type is empty")
	}

	p := path.Join("channeltypes", url.PathEscape(chanType))

	ct := ChannelType{}

	err := c.makeRequest(http.MethodGet, p, nil, nil, &ct)

	return &ct, err
}

// ListChannelTypes returns all channel types
func (c *Client) ListChannelTypes() (map[string]*ChannelType, error) {
	var resp channelTypeResponse

	err := c.makeRequest(http.MethodGet, "channeltypes", nil, nil, &resp)

	return resp.ChannelTypes, err
}

func (c *Client) UpdateChannelType(name string, data map[string]interface{}) error {
	switch {
	case name == "":
		return errors.New("channel type name is empty")
	case len(data) == 0:
		return errors.New("options are empty")
	}

	p := path.Join("channeltypes", url.PathEscape(name))

	return c.makeRequest(http.MethodPut, p, nil, nil, nil)
}

func (c *Client) DeleteChannelType(name string) error {
	if name == "" {
		return errors.New("channel type name is empty")
	}

	p := path.Join("channeltypes", url.PathEscape(name))

	return c.makeRequest(http.MethodDelete, p, nil, nil, nil)
}
