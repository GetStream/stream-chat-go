package stream_chat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type QueryMessageHistoryRequest struct {
	Filter map[string]any `json:"filter"`
	Sort   []*SortOption  `json:"sort,omitempty"`

	Limit int    `json:"limit,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

type MessageHistoryEntry struct {
	MessageID          string    `json:"message_id"`
	MessageUpdatedByID string    `json:"message_updated_by_id"`
	MessageUpdatedAt   time.Time `json:"message_updated_at"`

	Text        string                 `json:"text"`
	Attachments []*Attachment          `json:"attachments"`
	ExtraData   map[string]interface{} `json:"-"`
}

var (
	_ json.Unmarshaler = (*MessageHistoryEntry)(nil)
	_ json.Marshaler   = (*MessageHistoryEntry)(nil)
)

type messageHistoryJson MessageHistoryEntry

func (m *MessageHistoryEntry) UnmarshalJSON(data []byte) error {
	var m2 messageHistoryJson
	if err := json.Unmarshal(data, &m2); err != nil {
		return err
	}
	*m = MessageHistoryEntry(m2)

	if err := json.Unmarshal(data, &m.ExtraData); err != nil {
		return err
	}
	removeFromMap(m.ExtraData, *m)
	return nil
}

func (m MessageHistoryEntry) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(m.ExtraData, messageHistoryJson(m))
}

type QueryMessageHistoryResponse struct {
	MessageHistory []*MessageHistoryEntry `json:"message_history"`

	Next *string `json:"next,omitempty"`
	Prev *string `json:"prev,omitempty"`
	Response
}

func (c *Client) QueryMessageHistory(ctx context.Context, request QueryMessageHistoryRequest) (*QueryMessageHistoryResponse, error) {
	if len(request.Filter) == 0 {
		return nil, errors.New("you need specify one filter at least")
	}
	var resp QueryMessageHistoryResponse
	err := c.makeRequest(ctx, http.MethodPost, "messages/history", nil, request, &resp)
	return &resp, err
}
