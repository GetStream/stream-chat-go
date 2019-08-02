package stream_chat

import "time"

type Reaction struct {
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id,omitempty"`
	User      *User  `json:"user"`
	Type      string `json:"type"`

	CreatedAt time.Time `json:"created_at"`

	// any other fields the user wants to attach a reaction
	ExtraData map[string]interface{}
}
