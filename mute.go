package stream_chat

import (
	"time"
)

type Mute struct {
	User      User
	Target    User
	CreatedAt time.Time
	UpdatedAt time.Time
}
