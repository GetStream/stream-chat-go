package stream_chat //nolint: golint

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

type Task struct {
	TaskID    string     `json:"task_id"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	Result map[string]interface{} `json:"result,omitempty"`
}

// GetTask returns the status of a task that has been ran asynchronously.
func (c *Client) GetTask(ctx context.Context, id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("id should not be empty")
	}

	p := path.Join("tasks", url.PathEscape(id))

	var task Task
	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, &task)
	if err != nil {
		return nil, fmt.Errorf("cannot get task status: %v", err)
	}

	return &task, nil
}

type AsyncTaskResponse struct {
	TaskID string `json:"task_id"`
}

// DeleteChannels deletes channels asynchronously.
// Channels and messages will be hard deleted if hardDelete is true.
// It returns a task ID, the status of the task can be check with client.GetTask method.
func (c *Client) DeleteChannels(ctx context.Context, cids []string, hardDelete bool) (string, error) {
	if len(cids) == 0 {
		return "", fmt.Errorf("cids parameter should not be empty")
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
	if err != nil {
		return "", fmt.Errorf("cannot delete channels: %v", err)
	}

	return resp.TaskID, nil
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
// It returns a task ID, the status of the task can be check with client.GetTask method.
func (c *Client) DeleteUsers(ctx context.Context, userIDs []string, options DeleteUserOptions) (string, error) {
	if len(userIDs) == 0 {
		return "", fmt.Errorf("userIDs parameter should not be empty")
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
	if err != nil {
		return "", fmt.Errorf("cannot delete users: %v", err)
	}

	return resp.TaskID, nil
}

type ExportableChannel struct {
	Type          string     `json:"type"`
	ID            string     `json:"id"`
	MessagesSince *time.Time `json:"messages_since,omitempty"`
	MessagesUntil *time.Time `json:"messages_until,omitempty"`
}

// ExportChannels requests an asynchronous export of the provided channels and returns
// the ID of task.
func (c *Client) ExportChannels(ctx context.Context, channels []*ExportableChannel, clearDeletedMessageText, includeTruncatedMessages *bool) (string, error) {
	if len(channels) == 0 {
		return "", errors.New("number of channels must be at least one")
	}

	err := verifyExportableChannels(channels)
	if err != nil {
		return "", err
	}

	req := struct {
		Channels                 []*ExportableChannel `json:"channels"`
		ClearDeletedMessageText  *bool                `json:"clear_deleted_message_text,omitempty"`
		IncludeTruncatedMessages *bool                `json:"include_truncated_messages,omitempty"`
	}{
		Channels:                 channels,
		ClearDeletedMessageText:  clearDeletedMessageText,
		IncludeTruncatedMessages: includeTruncatedMessages,
	}

	var resp AsyncTaskResponse
	if err := c.makeRequest(ctx, http.MethodPost, "export_channels", nil, req, &resp); err != nil {
		return "", err
	}

	return resp.TaskID, nil
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
func (c *Client) GetExportChannelsTask(ctx context.Context, taskID string) (*Task, error) {
	task := &Task{}

	if taskID == "" {
		return task, errors.New("task ID must be not empty")
	}

	p := path.Join("export_channels", url.PathEscape(taskID))

	err := c.makeRequest(ctx, http.MethodGet, p, nil, nil, task)
	return task, err
}
