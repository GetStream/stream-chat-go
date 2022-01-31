package stream_chat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initClient(t *testing.T) *Client {
	c, err := NewClient(APIKey, APISecret)
	require.NoError(t, err, "new client")

	// set hostname to client from env if present
	if StreamHost != "" {
		c.BaseURL = StreamHost
	}

	return c
}

func initChannel(t *testing.T, c *Client, membersID ...string) *Channel {
	owner := randomUser(t, c)
	resp, err := c.CreateChannel(context.Background(), "team", randomString(12), owner.ID, map[string]interface{}{
		"members": membersID,
	})

	require.NoError(t, err, "create channel")
	return resp.Channel
}

func TestNewClient(t *testing.T) {
	c := initClient(t)

	assert.Equal(t, c.apiKey, APIKey)
	assert.Equal(t, string(c.apiSecret), APISecret)
	assert.NotEmpty(t, c.authToken)
	assert.Equal(t, defaultTimeout, c.HTTP.Timeout)
}

//nolint: lll
func TestClient_CreateToken(t *testing.T) {
	type args struct {
		userID string
		expire time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"simple without expiration",
			args{"tommaso", time.Time{}},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidG9tbWFzbyJ9.v-x-jt3ZnBXXbQ0GoWloIZtVnat2IE74U1a4Yuxd63M",
			false,
		},
		{
			"simple with expiration",
			args{"tommaso", time.Unix(1566941272, 123121)},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5NDEyNzIsInVzZXJfaWQiOiJ0b21tYXNvIn0.jF4ZbAIEuzS2jRH0uiu3HW9n0NHwT96QkzGlywcG9HU",
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient("key", "secret")
			require.NoError(t, err)

			got, err := c.CreateToken(tt.args.userID, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
