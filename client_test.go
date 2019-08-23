package stream_chat

import (
	"testing"
	"time"

	"github.com/pascaldekloe/jwt"

	"github.com/stretchr/testify/assert"
)

func initClient(t *testing.T) *Client {
	c, err := NewClient(APIKey, []byte(APISecret))
	mustNoError(t, err, "new client")

	// set hostname to client from env if present
	if StreamHost != "" {
		c.BaseURL = StreamHost
	}

	return c
}

func initChannel(t *testing.T, c *Client) *Channel {
	_, err := c.UpdateUsers(testUsers...)
	mustNoError(t, err, "update users")

	members := make([]string, 0, len(testUsers))
	for i := range testUsers {
		members = append(members, testUsers[i].ID)
	}

	ch, err := CreateChannel(c, ChannelOptions{ID: "fellowship-of-the-ring", Type: "team", Data: map[string]interface{}{
		"members": members,
	}}, serverUser.ID)

	mustNoError(t, err, "create channel")

	return ch
}

func TestNewClient(t *testing.T) {
	c := initClient(t)

	assert.Equal(t, c.apiKey, APIKey)
	assert.Equal(t, c.apiSecret, []byte(APISecret))
	assert.NotEmpty(t, c.header)
	assert.Equal(t, defaultTimeout, c.HTTP.Timeout)
	//	assert.Equal(t, defaultBaseURL, c.BaseURL, )
}

func Test_client_CreateToken(t *testing.T) {
	c := initClient(t)

	var expire = time.Now().Add(time.Hour)
	tt := []struct {
		name   string
		expire time.Time
	}{
		{"token without expire", time.Time{}},
		{"token with expire", expire},
	}

	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			token, err := c.CreateToken(testUsers[0].ID, test.expire)
			mustNoError(t, err, "create token")

			claims, err := jwt.HMACCheck(token, c.apiSecret)
			mustNoError(t, err, "jwt check")

			var expiresIn *jwt.NumericTime
			if !test.expire.IsZero() {
				expiresIn = jwt.NewNumericTime(test.expire)
			}

			assert.Equal(t, expiresIn, claims.Expires)
			assert.Equal(t, testUsers[0].ID, claims.Set["user_id"])
		})
	}
}
