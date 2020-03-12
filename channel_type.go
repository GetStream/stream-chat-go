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

// ChannelTypeLabel marks a string as a type of channel. The builtin settings are all
// available as consts.
type ChannelTypeLabel = string

const (
	// ChannelTypeLabelLivestream is sensible defaults in case you want to build chat like YouTube or Twitch.
	ChannelTypeLabelLivestream = "livestream"
	// ChannelTypeLabelMessaging is configured for apps such as WhatsApp or Facebook Messenger.
	ChannelTypeLabelMessaging = "messaging"
	// ChannelTypeLabelTeam is for If you want to build your own version of Slack or something similar.
	ChannelTypeLabelTeam = "team"
	// ChannelTypeLabelGaming is configured for in-game chat.
	ChannelTypeLabelGaming = "gaming"
	// ChannelTypeLabelCommerce is good defaults for building something like your own version of Intercom or Drift.
	ChannelTypeLabelCommerce = "commerce"
)

// ActionPermission is a type alias to assist in making sure permissions are
// correctly set.
type ActionPermission = string

const (
	// PermissionAllow sets a permission structure to allow access to the
	// resources to the roles in question.
	PermissionAllow = ActionPermission("Allow")

	// PermissionDeny sets a permission structure to deny access to the
	// resources to the roles in question.
	PermissionDeny = ActionPermission("Deny")
)

type modType string
type modBehaviour string

type Permission struct {
	Name   string           `json:"name"`   // required
	Action ActionPermission `json:"action"` // one of: Deny Allow

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

// AddPermissions helps to add permissions when building a channel type.
func (ct *ChannelType) AddPermissions(perms ...*Permission) *ChannelType {
	ct.Permissions = append(ct.Permissions, perms...)

	return ct
}

// AddCommands is a helper function to assist when building a channel type.
func (ct *ChannelType) AddCommands(cmds ...*Command) *ChannelType {
	ct.Commands = append(ct.Commands, cmds...)

	return ct
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
