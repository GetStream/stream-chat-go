package stream_chat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type TaskStatus string

const (
	TaskStatusWaiting   TaskStatus = "waiting"
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type TaskResponse struct {
	TaskID    string     `json:"task_id"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	Result map[string]interface{} `json:"result,omitempty"`
	Response
}

// GetTask returns the status of a task that has been ran asynchronously.
func (c *Client) GetTask(ctx context.Context, id string) (*TaskResponse, error) {
	if id == "" {
		return nil, errors.New("id should not be empty")
	}

	p := path.Join("tasks", url.PathEscape(id))

	var task TaskResponse
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &task)
	return &task, err
}

type AsyncTaskResponse struct {
	TaskID string `json:"task_id"`
	Response
}

// DeleteChannels deletes channels asynchronously.
// Channels and messages will be hard deleted if hardDelete is true.
// It returns an AsyncTaskResponse object which contains the task ID, the status of the task can be check with client.GetTask method.
func (c *Client) DeleteChannels(ctx context.Context, cids []string, hardDelete bool) (*AsyncTaskResponse, error) {
	if len(cids) == 0 {
		return nil, errors.New("cids parameter should not be empty")
	}

	data := struct {
		CIDs       []string `json:"cids"`
		HardDelete bool     `json:"hard_delete"`
	}{
		CIDs:       cids,
		HardDelete: hardDelete,
	}

	var resp AsyncTaskResponse
	err := c.makeRequest(ctx, http.MethodPost, "channels/delete", nil, data, &resp)
	return &resp, err
}

type DeleteType string

const (
	HardDelete DeleteType = "hard"
	SoftDelete DeleteType = "soft"
)

type DeleteUserOptions struct {
	User              DeleteType `json:"user"`
	Messages          DeleteType `json:"messages,omitempty"`
	Conversations     DeleteType `json:"conversations,omitempty"`
	NewChannelOwnerID string     `json:"new_channel_owner_id,omitempty"`
}

// DeleteUsers deletes users asynchronously.
// User will be deleted either "hard" or "soft"
// Conversations (1to1 channels) will be deleted if either "hard" or "soft"
// Messages will be deleted if either "hard" or "soft"
// NewChannelOwnerID any channels owned by the hard-deleted user will be transferred to this user ID
// It returns an AsyncTaskResponse object which contains the task ID, the status of the task can be check with client.GetTask method.
func (c *Client) DeleteUsers(ctx context.Context, userIDs []string, options DeleteUserOptions) (*AsyncTaskResponse, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("userIDs parameter should not be empty")
	}

	data := struct {
		DeleteUserOptions
		UserIDs []string `json:"user_ids"`
	}{
		DeleteUserOptions: options,
		UserIDs:           userIDs,
	}

	var resp AsyncTaskResponse
	err := c.makeRequest(ctx, http.MethodPost, "users/delete", nil, data, &resp)
	return &resp, err
}

type ExportableChannel struct {
	Type          string     `json:"type"`
	ID            string     `json:"id"`
	MessagesSince *time.Time `json:"messages_since,omitempty"`
	MessagesUntil *time.Time `json:"messages_until,omitempty"`
}

type ExportChannelOptions struct {
	ClearDeletedMessageText  *bool  `json:"clear_deleted_message_text,omitempty"`
	IncludeTruncatedMessages *bool  `json:"include_truncated_messages,omitempty"`
	ExportUsers              *bool  `json:"export_users,omitempty"`
	Version                  string `json:"version,omitempty"`
}

// ExportChannels requests an asynchronous export of the provided channels.
// It returns an AsyncTaskResponse object which contains the task ID, the status of the task can be check with client.GetTask method.
func (c *Client) ExportChannels(ctx context.Context, channels []*ExportableChannel, options *ExportChannelOptions) (*AsyncTaskResponse, error) {
	if len(channels) == 0 {
		return nil, errors.New("number of channels must be at least one")
	}

	err := verifyExportableChannels(channels)
	if err != nil {
		return nil, err
	}

	req := struct {
		Channels                 []*ExportableChannel `json:"channels"`
		ClearDeletedMessageText  *bool                `json:"clear_deleted_message_text,omitempty"`
		IncludeTruncatedMessages *bool                `json:"include_truncated_messages,omitempty"`
		ExportUsers              *bool                `json:"export_users,omitempty"`
		Version                  string               `json:"version,omitempty"`
	}{
		Channels: channels,
	}

	if options != nil {
		req.ClearDeletedMessageText = options.ClearDeletedMessageText
		req.IncludeTruncatedMessages = options.IncludeTruncatedMessages
		req.ExportUsers = options.ExportUsers
		req.Version = options.Version
	}

	var resp AsyncTaskResponse
	err = c.makeRequest(ctx, http.MethodPost, "export_channels", nil, req, &resp)
	return &resp, err
}

func verifyExportableChannels(channels []*ExportableChannel) error {
	for i, ch := range channels {
		if ch.Type == "" || ch.ID == "" {
			return fmt.Errorf("channel type and id must not be empty for index: %d", i)
		}
	}
	return nil
}

// GetExportChannelsTask returns current state of the export task.
func (c *Client) GetExportChannelsTask(ctx context.Context, taskID string) (*TaskResponse, error) {
	if taskID == "" {
		return nil, errors.New("task ID must be not empty")
	}

	p := path.Join("export_channels", url.PathEscape(taskID))

	var task TaskResponse
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &task)
	return &task, err
}
