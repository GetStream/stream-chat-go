package stream_chat //nolint: golint

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//nolint: gochecknoglobals
var (
	APIKey     = os.Getenv("STREAM_CHAT_API_KEY")
	APISecret  = os.Getenv("STREAM_CHAT_API_SECRET")
	StreamHost = os.Getenv("STREAM_CHAT_API_HOST")

	serverUser *User
	testUsers  []*User
)

//nolint: gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())

	if err := clearOldChannelTypes(); err != nil {
		panic(err) // app has bad data from previous runs
	}

	serverUser = &User{ID: randomString(10), Name: "Gandalf the Grey", ExtraData: map[string]interface{}{"race": "Istari"}}

	testUsers = []*User{
		{ID: randomString(10), Name: "Frodo Baggins", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 50}},
		{ID: randomString(10), Name: "Samwise Gamgee", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 38}},
		{ID: randomString(10), Name: "Legolas", ExtraData: map[string]interface{}{"race": "Elf", "age": 500}},
		serverUser,
	}
}

func clearOldChannelTypes() error {
	c, err := NewClient(APIKey, APISecret)
	if err != nil {
		return err
	}
	c.BaseURL = defaultBaseURL

	got, err := c.ListChannelTypes()
	if err != nil {
		return err
	}

	for _, ct := range got {
		if contains(defaultChannelTypes, ct.Name) {
			continue
		}
		filter := map[string]interface{}{"type": ct.Name}
		chs, _ := c.QueryChannels(&QueryOption{Filter: filter})

		hasChannel := false
		for _, ch := range chs {
			if err := ch.Delete(); err != nil {
				hasChannel = true
				break
			}
		}

		if !hasChannel {
			_ = c.DeleteChannelType(ct.Name)
		}
	}
	return nil
}

func randomUser(t *testing.T, c *Client) *User {
	u, err := c.UpsertUser(&User{ID: randomString(10)})
	require.NoError(t, err)
	return u
}

func randomUsers(t *testing.T, c *Client, n int) []*User {
	users := make([]*User, 0, n)
	for i := 0; i < n; i++ {
		users = append(users, &User{ID: randomString(10)})
	}

	userss, err := c.UpsertUsers(users...)
	require.NoError(t, err)
	users = users[:0]
	for _, user := range userss {
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
