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
