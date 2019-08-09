package stream_chat

import (
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/francoispqt/gojay"
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
	ChannelConfig

	Commands    Commands     `json:"commands"`
	Permissions []Permission `json:"permissions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewChannelType returns initialized ChannelType with default values
func NewChannelType(name string) ChannelType {
	ct := ChannelType{
		ChannelConfig: ChannelConfig{
			Name:             name,
			Automod:          AutoModDisabled,
			ModBehavior:      ModBehaviourFlag,
			MaxMessageLength: defaultMessageLength,
			MessageRetention: MessageRetentionForever,
		},
	}

	return ct
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

type channelTypes map[string]ChannelType

func (c *channelTypes) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	var ct ChannelType
	if err := dec.Object(&ct); err != nil {
		return err
	}
	(*c)[key] = ct
	return nil
}

func (c *channelTypes) NKeys() int {
	return 0
}

type channelTypeResponse struct {
	ChannelTypes channelTypes `json:"channel_types"`
}

func (c *channelTypeResponse) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	if key == "channel_types" {
		c.ChannelTypes = channelTypes{}
		return dec.Object(&c.ChannelTypes)
	}
	return nil
}

func (c *channelTypeResponse) NKeys() int {
	return 1
}

// ListChannelTypes returns all channel types
func (c *Client) ListChannelTypes() (map[string]ChannelType, error) {
	var resp channelTypeResponse

	err := c.makeRequest(http.MethodGet, "channeltypes", nil, nil, &resp)

	return resp.ChannelTypes, err
}

func (c *Client) DeleteChannelType(ct string) error {
	p := path.Join("channeltypes", url.PathEscape(ct))

	return c.makeRequest(http.MethodDelete, p, nil, nil, nil)
}
