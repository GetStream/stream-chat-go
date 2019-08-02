package stream_chat

import "time"

type Mute struct {
	User      *User     `json:"user"`
	Target    *User     `json:"target"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
