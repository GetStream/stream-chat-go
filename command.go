package stream_chat // nolint: golint
import (
	"errors"
	"net/http"
	"net/url"
	"path"
)

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Args        string `json:"args"`
	Set         string `json:"set"`
}

type commandResponse struct {
	Command *Command
}

type commandsResponse struct {
	Commands []*Command
}

func (c *Client) CreateCommand(cmd *Command) (*Command, error) {
	if cmd == nil {
		return nil, errors.New("command is nil")
	}

	var resp commandResponse

	err := c.makeRequest(http.MethodPost, "commands", nil, cmd, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Command == nil {
		return nil, errors.New("unexpected error: command response is nil")
	}

	return resp.Command, nil
}

func (c *Client) GetCommand(cmdName string) (*Command, error) {
	if cmdName == "" {
		return nil, errors.New("command name is empty")
	}

	p := path.Join("commands", url.PathEscape(cmdName))

	cmd := Command{}

	err := c.makeRequest(http.MethodGet, p, nil, nil, &cmd)

	return &cmd, err
}

func (c *Client) DeleteCommand(cmdName string) error {
	if cmdName == "" {
		return errors.New("command name is empty")
	}

	p := path.Join("commands", url.PathEscape(cmdName))

	return c.makeRequest(http.MethodDelete, p, nil, nil, nil)
}

func (c *Client) ListCommands() ([]*Command, error) {
	var resp commandsResponse

	err := c.makeRequest(http.MethodGet, "commands", nil, nil, &resp)

	return resp.Commands, err
}

func (c *Client) UpdateCommand(cmdName string, options map[string]interface{}) (*Command, error) {
	switch {
	case cmdName == "":
		return nil, errors.New("command name is empty")
	case len(options) == 0:
		return nil, errors.New("options are empty")
	}

	p := path.Join("commands", url.PathEscape(cmdName))

	var resp commandResponse

	err := c.makeRequest(http.MethodPut, p, nil, options, &resp)
	return resp.Command, err

}

//type ChannelType struct {
//	ChannelConfig
//
//	Commands    []*Command    `json:"commands"`
//	Permissions []*Permission `json:"permissions"`
//
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}
//
//func (ct *ChannelType) toRequest() channelTypeRequest {
//	req := channelTypeRequest{ChannelType: ct}
//
//	if len(req.Commands) == 0 {
//		req.Commands = []string{"all"}
//	}
//
//	return req
//}
//
//// NewChannelType returns initialized ChannelType with default values.
//func NewChannelType(name string) *ChannelType {
//	ct := &ChannelType{ChannelConfig: DefaultChannelConfig}
//	ct.Name = name
//
//	return ct
//}
//
//type channelTypeRequest struct {
//	*ChannelType
//
//	Commands []string `json:"commands"`
//
//	CreatedAt time.Time `json:"-"`
//	UpdatedAt time.Time `json:"-"`
//}
//
//type channelTypeResponse struct {
//	ChannelTypes map[string]*ChannelType `json:"channel_types"`
//}
//
//// CreateChannelType adds new channel type.
//func (c *Client) CreateChannelType(chType *ChannelType) (*ChannelType, error) {
//	if chType == nil {
//		return nil, errors.New("channel type is nil")
//	}
//
//	var resp channelTypeRequest
//
//	err := c.makeRequest(http.MethodPost, "channeltypes", nil, chType.toRequest(), &resp)
//	if err != nil {
//		return nil, err
//	}
//	if resp.ChannelType == nil {
//		return nil, errors.New("unexpected error: channel type response is nil")
//	}
//
//	for _, cmd := range resp.Commands {
//		resp.ChannelType.Commands = append(resp.ChannelType.Commands, &Command{Name: cmd})
//	}
//
//	return resp.ChannelType, nil
//}
//
//// GetChannelType returns information about channel type.
//func (c *Client) GetChannelType(chanType string) (*ChannelType, error) {
//	if chanType == "" {
//		return nil, errors.New("channel type is empty")
//	}
//
//	p := path.Join("channeltypes", url.PathEscape(chanType))
//
//	ct := ChannelType{}
//
//	err := c.makeRequest(http.MethodGet, p, nil, nil, &ct)
//
//	return &ct, err
//}
//
//// ListChannelTypes returns all channel types.
//func (c *Client) ListChannelTypes() (map[string]*ChannelType, error) {
//	var resp channelTypeResponse
//
//	err := c.makeRequest(http.MethodGet, "channeltypes", nil, nil, &resp)
//
//	return resp.ChannelTypes, err
//}
//
//func (c *Client) UpdateChannelType(name string, options map[string]interface{}) error {
//	switch {
//	case name == "":
//		return errors.New("channel type name is empty")
//	case len(options) == 0:
//		return errors.New("options are empty")
//	}
//
//	p := path.Join("channeltypes", url.PathEscape(name))
//
//	return c.makeRequest(http.MethodPut, p, nil, nil, nil)
//}
//
//func (c *Client) DeleteChannelType(name string) error {
//	if name == "" {
//		return errors.New("channel type name is empty")
//	}
//
//	p := path.Join("channeltypes", url.PathEscape(name))
//
//	return c.makeRequest(http.MethodDelete, p, nil, nil, nil)
//}
