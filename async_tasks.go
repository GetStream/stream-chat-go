package stream_chat //nolint: golint

import (
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
