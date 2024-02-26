package stream_chat

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type UnreadCountsChannel struct {
	ChannelID   string    `json:"channel_id"`
	UnreadCount int       `json:"unread_count"`
	LastRead    time.Time `json:"last_read"`
}

type UnreadCountsChannelType struct {
	ChannelType  string `json:"channel_type"`
	ChannelCount int    `json:"channel_count"`
	UnreadCount  int    `json:"unread_count"`
}

type UnreadCountsThread struct {
	UnreadCount       int       `json:"unread_count"`
	LastRead          time.Time `json:"last_read"`
	LastReadMessageID string    `json:"last_read_message_id"`
	ParentMessageID   string    `json:"parent_message_id"`
}

type UnreadCountsResponse struct {
	TotalUnreadCount        int                       `json:"total_unread_count"`
	TotalUnreadThreadsCount int                       `json:"total_unread_threads_count"`
	Channels                []UnreadCountsChannel     `json:"channels"`
	ChannelType             []UnreadCountsChannelType `json:"channel_type"`
	Threads                 []UnreadCountsThread      `json:"threads"`
	Response
}

func (c *Client) UnreadCounts(ctx context.Context, userID string) (*UnreadCountsResponse, error) {
	var resp UnreadCountsResponse
	err := c.makeRequest(ctx, http.MethodGet, "unread", url.Values{"user_id": []string{userID}}, nil, &resp)
	return &resp, err
}

type UnreadCountsBatchResponse struct {
	CountsByUser map[string]*UnreadCountsResponse `json:"counts_by_user"`
	Response
}

func (c *Client) UnreadCountsBatch(ctx context.Context, userIDs []string) (*UnreadCountsBatchResponse, error) {
	var resp UnreadCountsBatchResponse
	err := c.makeRequest(ctx, http.MethodPost, "unread_batch", nil, map[string][]string{"user_ids": userIDs}, &resp)
	return &resp, err
}
