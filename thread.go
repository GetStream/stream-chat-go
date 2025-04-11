package stream_chat

import (
	"context"
	"net/http"
	"time"
)

type QueryThreadsRequest struct {
	User   *User  `json:"user,omitempty"`
	UserID string `json:"user_id,omitempty"`

	Filter map[string]any        `json:"filter,omitempty"`
	Sort   *SortParamRequestList `json:"sort,omitempty"`
	Watch  *bool                 `json:"watch,omitempty"`
	PagerRequest
}

type QueryThreadsResponse struct {
	Threads []ThreadResponse `json:"threads"`
	Response
	PagerResponse
}

type ThreadResponse struct {
	ChannelCID             string             `json:"channel_cid"`
	Channel                *Channel           `json:"channel,omitempty"`
	ParentMessageID        string             `json:"parent_message_id"`
	ParentMessage          *MessageResponse   `json:"parent_message,omitempty"`
	CreatedByUserID        string             `json:"created_by_user_id"`
	CreatedBy              *UsersResponse     `json:"created_by,omitempty"`
	ReplyCount             int                `json:"reply_count,omitempty"`
	ParticipantCount       int                `json:"participant_count,omitempty"`
	ActiveParticipantCount int                `json:"active_participant_count,omitempty"`
	Participants           ThreadParticipants `json:"thread_participants,omitempty"`
	LastMessageAt          *time.Time         `json:"last_message_at,omitempty"`
	CreatedAt              *time.Time         `json:"created_at"`
	UpdatedAt              *time.Time         `json:"updated_at"`
	DeletedAt              *time.Time         `json:"deleted_at,omitempty"`
	Title                  string             `json:"title"`
	Custom                 map[string]any     `json:"custom"`

	LatestReplies []MessageResponse `json:"latest_replies,omitempty"`
	Read          ChannelRead       `json:"read,omitempty"`
	Draft         Draft
}

type Thread struct {
	AppPK int `json:"app_pk"`

	ChannelCID string   `json:"channel_cid"`
	Channel    *Channel `json:"channel,omitempty"`

	ParentMessageID string   `json:"parent_message_id"`
	ParentMessage   *Message `json:"parent_message,omitempty"`

	CreatedByUserID string `json:"created_by_user_id"`
	CreatedBy       *User  `json:"created_by,omitempty"`

	ReplyCount             int                `json:"reply_count,omitempty"`
	ParticipantCount       int                `json:"participant_count,omitempty"`
	ActiveParticipantCount int                `json:"active_participant_count,omitempty"`
	Participants           ThreadParticipants `json:"thread_participants,omitempty"`

	LastMessageAt time.Time  `json:"last_message_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`

	Title  string         `json:"title"`
	Custom map[string]any `json:"custom"`
}

type ThreadParticipant struct {
	AppPK int `json:"app_pk"`

	ChannelCID string `json:"channel_cid"`

	LastThreadMessageAt *time.Time `json:"last_thread_message_at"`
	ThreadID            string     `json:"thread_id,omitempty"`

	UserID string `json:"user_id,omitempty"`
	User   *User  `json:"user,omitempty"`

	CreatedAt    time.Time  `json:"created_at"`
	LeftThreadAt *time.Time `json:"left_thread_at,omitempty"`
	LastReadAt   time.Time  `json:"last_read_at"`

	Custom map[string]interface{} `json:"custom"`
}

type ThreadParticipants []*ThreadParticipant

func (c *Client) QueryThreads(ctx context.Context, query *QueryThreadsRequest) (*QueryThreadsResponse, error) {
	var resp QueryThreadsResponse
	err := c.makeRequest(ctx, http.MethodPost, "threads", nil, query, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
