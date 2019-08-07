package stream_chat

import (
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	AutomodDisabled autoMod = "disabled"
	AutomodSimple   autoMod = "simple"
	AutomodAI       autoMod = "AI"
)

type autoMod string

type Permission struct {
	// required
	Name string `json:"name"`
	// one of: Deny Allow
	Action string `json:"action"`
	// required
	Resources []string `json:"resources"`
	Roles     []string `json:"roles"`
	Owner     bool     `json:"owner"`
	// required
	Priority int `json:"priority"`
}

type ChannelType struct {
	// required fields
	Name string `json:"name"`
	// one of Automod* constant. Required
	Automod autoMod `json:"auto_mod"`

	// one of: flag block
	AutomodBehaviour string `json:"automod_behaviour,omitempty"`

	TypingEvents  bool `json:"typing_events"`
	ReadEvents    bool `json:"read_events"`
	ConnectEvents bool `json:"connect_events"`
	Search        bool `json:"search"`
	Reactions     bool `json:"reactions"`
	Replies       bool `json:"replies"`
	Mutes         bool `json:"mutes"`

	// one of: infinite numeric
	MessageRetention string `json:"message_retention,omitempty"`
	MaxMessageLength int    `json:"max_message_length,omitempty"`

	Commands    Commands     `json:"commands"`
	Permissions []Permission `json:"permissions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateChannelType adds new channel type
func (c *Client) CreateChannelType(chType *ChannelType) (err error) {
	var resp ChannelType
	err = c.makeRequest(http.MethodPost, "channeltypes", nil, chType, &resp)
	if err != nil {
		return err
	}

	*chType = resp

	return err
}

// GetChannelType returns information about channel type
func (c *Client) GetChannelType(chanType string) (ct ChannelType, err error) {
	p := path.Join("channeltypes", url.PathEscape(chanType))

	err = c.makeRequest(http.MethodGet, p, nil, nil, &ct)

	return ct, err
}

// ListChannelTypes returns all channel types
func (c *Client) ListChannelTypes() (map[string]ChannelType, error) {
	var resp struct {
		ChannelTypes map[string]ChannelType `json:"channel_types"`
	}

	err := c.makeRequest(http.MethodGet, "channeltypes", nil, nil, &resp)

	return resp.ChannelTypes, err
}

func (c *Client) DeleteChannelType(ct string) error {
	p := path.Join("channeltypes", url.PathEscape(ct))

	return c.makeRequest(http.MethodDelete, p, nil, nil, nil)
}
