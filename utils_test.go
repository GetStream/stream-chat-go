package stream_chat

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

	ctx := context.Background()

	resp, err := c.ListChannelTypes(ctx)
	if err != nil {
		return err
	}

	for _, ct := range resp.ChannelTypes {
		if contains(defaultChannelTypes, ct.Name) {
			continue
		}
		filter := map[string]interface{}{"type": ct.Name}
		resp, _ := c.QueryChannels(ctx, &QueryOption{Filter: filter})

		hasChannel := false
		for _, ch := range resp.Channels {
			if _, err := ch.Delete(ctx); err != nil {
				hasChannel = true
				break
			}
		}

		if !hasChannel {
			_, _ = c.DeleteChannelType(ctx, ct.Name)
		}
	}
	return nil
}

func randomUser(t *testing.T, c *Client) *User {
	t.Helper()

	ctx := context.Background()
	resp, err := c.UpsertUser(ctx, &User{ID: randomString(10)})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = c.DeleteUsers(ctx, []string{resp.User.ID}, DeleteUserOptions{
			User:          HardDelete,
			Messages:      HardDelete,
			Conversations: HardDelete,
		})
	})

	return resp.User
}

func randomUserWithRole(t *testing.T, c *Client, role string) *User {
	t.Helper()

	ctx := context.Background()
	resp, err := c.UpsertUser(ctx, &User{
		ID:   randomString(10),
		Role: role,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = c.DeleteUsers(ctx, []string{resp.User.ID}, DeleteUserOptions{
			User:          HardDelete,
			Messages:      HardDelete,
			Conversations: HardDelete,
		})
	})

	return resp.User
}

func randomUsers(t *testing.T, c *Client, n int) []*User {
	t.Helper()

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
	t.Helper()

	users := randomUsers(t, c, n)
	ids := make([]string, n)
	for i, u := range users {
		ids[i] = u.ID
	}
	return ids
}

func randomUsersChannelMember(t *testing.T, c *Client, n int) []ChannelMember {
	t.Helper()

	users := randomUsers(t, c, n)
	members := make([]ChannelMember, n)
	for i, u := range users {
		members[i] = ChannelMember{UserID: u.ID}
	}
	return members
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
