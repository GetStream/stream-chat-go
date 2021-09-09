package stream_chat //nolint: golint

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type TaskStatus struct {
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Result map[string]interface{} `json:"result,omitempty"`
}

// GetTask returns the status of a task that has been ran asynchronously.
func (c *Client) GetTask(id string) (*TaskStatus, error) {
	if id == "" {
		return nil, fmt.Errorf("id should not be empty")
	}

	p := path.Join("tasks", url.PathEscape(id))

	var status TaskStatus
	err := c.makeRequest(http.MethodGet, p, nil, nil, &status)
	if err != nil {
		return nil, fmt.Errorf("cannot get task status: %v", err)
	}

	return &status, nil
}

type AsyncTaskResponse struct {
	TaskID string `json:"task_id"`
}

// DeleteChannels deletes channels asynchronously.
// Channels and messages will be hard deleted if hardDelete is true.
// It returns a task ID, the status of the task can be check with client.GetTask method.
func (c *Client) DeleteChannels(cids []string, hardDelete bool) (string, error) {
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
	err := c.makeRequest(http.MethodPost, "channels/delete", nil, data, &resp)
	if err != nil {
		return "", fmt.Errorf("cannot delete channels: %v", err)
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
func (c *Client) ExportChannels(channels []*ExportableChannel, clearDeletedMessageText, includeTruncatedMessages *bool) (string, error) {
	if len(channels) < 1 || len(channels) > 25 {
		return "", errors.New("number of channels must be between 1 and 25")
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
	if err := c.makeRequest(http.MethodPost, "export_channels", nil, req, &resp); err != nil {
		return "", err
	}

	return resp.TaskID, nil
}

func verifyExportableChannels(channels []*ExportableChannel) error {
	var err error
	for _, ch := range channels {
		switch {
		case ch.Type == "":
			err = errors.New("channel type must be not empty")
			break
		case ch.ID == "":
			err = errors.New("channel ID must be not empty")
			break
		}
	}

	return err
}

// GetExportChannelsTask returns current state of the export task.
func (c *Client) GetExportChannelsTask(taskID string) (TaskStatus, error) {
	if taskID == "" {
		return TaskStatus{}, errors.New("task ID must be not empty")
	}

	p := path.Join("export_channels", url.PathEscape(taskID))

	var resp TaskStatus
	if err := c.makeRequest(http.MethodGet, p, nil, nil, &resp); err != nil {
		return TaskStatus{}, err
	}

	return resp, nil
}
