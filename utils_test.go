package stream_chat

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//nolint: gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())

	if err := clearOldChannelTypes(); err != nil {
		panic(err) // app has bad data from previous runs
	}
}

func clearOldChannelTypes() error {
	c, err := NewClientFromEnvVars()
	if err != nil {
		return err
	}
	c.BaseURL = defaultBaseURL

	resp, err := c.ListChannelTypes(context.Background())
	if err != nil {
		return err
	}

	for _, ct := range resp.ChannelTypes {
		if contains(defaultChannelTypes, ct.Name) {
			continue
		}
		filter := map[string]interface{}{"type": ct.Name}
		resp, _ := c.QueryChannels(context.Background(), &QueryOption{Filter: filter})

		hasChannel := false
		for _, ch := range resp.Channels {
			if _, err := ch.Delete(context.Background()); err != nil {
				hasChannel = true
				break
			}
		}

		if !hasChannel {
			_, _ = c.DeleteChannelType(context.Background(), ct.Name)
		}
	}
	return nil
}

func randomUser(t *testing.T, c *Client) *User {
	resp, err := c.UpsertUser(context.Background(), &User{ID: randomString(10)})
	require.NoError(t, err)
	return resp.User
}

func randomUsers(t *testing.T, c *Client, n int) []*User {
	users := make([]*User, 0, n)
	for i := 0; i < n; i++ {
		users = append(users, &User{ID: randomString(10)})
	}

	resp, err := c.UpsertUsers(context.Background(), users...)
	require.NoError(t, err)
	users = users[:0]
	for _, user := range resp.Users {
		users = append(users, user)
	}
	return users
}

func randomUsersID(t *testing.T, c *Client, n int) []string {
	users := randomUsers(t, c, n)
	ids := make([]string, n)
	for i, u := range users {
		ids[i] = u.ID
	}
	return ids
}

func randomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}

func contains(ls []string, s string) bool {
	for _, item := range ls {
		if item == s {
			return true
		}
	}
	return false
}
