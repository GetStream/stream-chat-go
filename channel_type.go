package stream_chat

import (
	"net/http"
	"time"
)

type ChannelType struct {
	Name             string
	TypingEvents     bool
	ReadEvents       bool
	ConnectEvents    bool
	Search           bool
	Reactions        bool
	Replies          bool
	Mutes            bool
	MessageRetention string
	MaxMessageLength int
	Automod          string
	AutomodBehaviour string
	Commands         []interface{}
	Permissions      []interface{}
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateChannelType adds new channel type
func (c *Client) CreateChannelType(chType ChannelType) (err error) {
	if len(chType.Commands) == 0 {
		chType.Commands = []interface{}{"all"}
	}

	err = c.makeRequest(http.MethodPost, "channeltypes", nil, data, nil)

	return
}

// GetChannelType returns information about channel type
func (c *Client) GetChannelType(chanType string) (ChannelType, error) {
	var resp struct {
		ChannelTypes map[string]ChannelType `json:"channel_types"`
	}

	err := c.makeRequest(http.MethodGet, "channeltypes/"+chanType, nil, nil, &resp)

	return resp.ChannelTypes[chanType], err
}

// ListChannelTypes returns all channel types
func (c *Client) ListChannelTypes() (map[string]ChannelType, error) {
	var resp struct {
		ChannelTypes map[string]ChannelType `json:"channel_types"`
	}

	err := c.makeRequest(http.MethodGet, "channeltypes", nil, nil, &resp)

	return resp.ChannelTypes, err
}
