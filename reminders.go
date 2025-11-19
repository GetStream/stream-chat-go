package stream_chat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Reminder struct {
	ChannelID string     `json:"channel_id"`
	MessageID string     `json:"message_id"`
	Message   *Message   `json:"message,omitempty"`
	UserID    string     `json:"user_id"`
	User      *User      `json:"user,omitempty"`
	RemindAt  *time.Time `json:"remind_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type reminderForJSON Reminder

func (s *Reminder) UnmarshalJSON(data []byte) error {
	var s2 reminderForJSON
	if err := json.Unmarshal(data, &s2); err != nil {
		return err
	}
	*s = Reminder(s2)

	return nil
}

type ReminderRequest struct {
	UserID    string     `json:"user_id"`
	MessageID string     `json:"message_id"`
	RemindAt  *time.Time `json:"remind_at,omitempty"`
}

type ReminderResponse struct {
	Reminder *Reminder `json:"reminder"`
	Response
}

// CreateReminder creates a reminder for a message.
// messageID: The ID of the message to create a reminder for
// userID: The ID of the user creating the reminder
// remindAt: When to remind the user (optional)
func (c *Client) CreateReminder(ctx context.Context, messageID, userID string, remindAt *time.Time) (*ReminderResponse, error) {
	if messageID == "" {
		return nil, errors.New("message ID is empty")
	}
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	data := map[string]interface{}{
		"user_id": userID,
	}

	if remindAt != nil {
		data["remind_at"] = remindAt.Format(time.RFC3339)
	}

	p := path.Join("messages", url.PathEscape(messageID), "reminders")

	var resp ReminderResponse
	err := c.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// UpdateReminder updates a reminder for a message.
// messageID: The ID of the message with the reminder
// userID: The ID of the user who owns the reminder
// remindAt: When to remind the user (optional)
func (c *Client) UpdateReminder(ctx context.Context, messageID, userID string, remindAt *time.Time) (*ReminderResponse, error) {
	if messageID == "" {
		return nil, errors.New("message ID is empty")
	}
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	data := map[string]interface{}{
		"user_id": userID,
	}

	if remindAt != nil {
		data["remind_at"] = remindAt.Format(time.RFC3339)
	}

	p := path.Join("messages", url.PathEscape(messageID), "reminders")

	var resp ReminderResponse
	err := c.makeRequest(ctx, http.MethodPatch, p, nil, data, &resp)
	return &resp, err
}

// DeleteReminder deletes a reminder for a message.
// messageID: The ID of the message with the reminder
// userID: The ID of the user who owns the reminder
func (c *Client) DeleteReminder(ctx context.Context, messageID, userID string) (*Response, error) {
	if messageID == "" {
		return nil, errors.New("message ID is empty")
	}
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	p := path.Join("messages", url.PathEscape(messageID), "reminders")

	params := url.Values{}
	params.Set("user_id", userID)

	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, p, params, nil, &resp)
	return &resp, err
}

type QueryRemindersResponse struct {
	Reminders []*Reminder `json:"reminders"`
	Response
}

// QueryReminders queries reminders based on filter conditions.
// userID: The ID of the user whose reminders to query
// filterConditions: Conditions to filter reminders
// sort: Sort parameters (default: [{field: "remind_at", direction: 1}])
// options: Additional query options like limit, offset
func (c *Client) QueryReminders(ctx context.Context, userID string, filterConditions map[string]interface{}, sort []*SortOption, options map[string]interface{}) (*QueryRemindersResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	if sort == nil {
		sort = []*SortOption{
			{
				Field:     "remind_at",
				Direction: 1,
			},
		}
	}

	data := map[string]interface{}{
		"user_id": userID,
		"filter":  filterConditions,
		"sort":    sort,
	}

	// Add additional options
	for k, v := range options {
		data[k] = v
	}

	var resp QueryRemindersResponse
	err := c.makeRequest(ctx, http.MethodPost, "reminders/query", nil, data, &resp)
	return &resp, err
}
