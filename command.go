package stream_chat // nolint: golint
import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"
)

// Command represents a custom command.
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Args        string `json:"args"`
	Set         string `json:"set"`
}

// CommandResponse represents an API response containing one Command.
type CommandResponse struct {
	Command *Command `json:"command"`
	Response
}

// CreateCommand registers a new custom command.
func (c *Client) CreateCommand(ctx context.Context, cmd *Command) (*CommandResponse, error) {
	if cmd == nil {
		return nil, errors.New("command is nil")
	}

	var resp CommandResponse

	err := c.makeRequest(ctx, http.MethodPost, "commands", nil, cmd, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Command == nil {
		return nil, errors.New("unexpected error: command response is nil")
	}

	return &resp, nil
}

type GetCommandResponse struct {
	*Command
	Response
}

// GetCommand retrieves a custom command referenced by cmdName.
func (c *Client) GetCommand(ctx context.Context, cmdName string) (*GetCommandResponse, error) {
	if cmdName == "" {
		return nil, errors.New("command name is empty")
	}

	p := path.Join("commands", url.PathEscape(cmdName))

	var resp GetCommandResponse
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &resp)
	return &resp, err
}

// DeleteCommand deletes a custom command referenced by cmdName.
func (c *Client) DeleteCommand(ctx context.Context, cmdName string) (*Response, error) {
	if cmdName == "" {
		return nil, errors.New("command name is empty")
	}

	p := path.Join("commands", url.PathEscape(cmdName))

	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, p, nil, nil, &resp)
	return &resp, err
}

// CommandsResponse represents an API response containing a list of Command.
type CommandsResponse struct {
	Commands []*Command
}

// ListCommands returns a list of custom commands.
func (c *Client) ListCommands(ctx context.Context) (*CommandsResponse, error) {
	var resp CommandsResponse
	err := c.makeRequest(ctx, http.MethodGet, "commands", nil, nil, &resp)
	return &resp, err
}

// UpdateCommand updates a custom command referenced by cmdName.
func (c *Client) UpdateCommand(ctx context.Context, cmdName string, update *Command) (*CommandResponse, error) {
	switch {
	case cmdName == "":
		return nil, errors.New("command name is empty")
	case update == nil:
		return nil, errors.New("update should not be nil")
	}

	p := path.Join("commands", url.PathEscape(cmdName))

	var resp CommandResponse
	err := c.makeRequest(ctx, http.MethodPut, p, nil, update, &resp)
	return &resp, err
}
