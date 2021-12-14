package stream_chat // nolint: golint

import (
	"context"
	"encoding/json"
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
	ch, err := c.CreateChannel(context.Background(), "team", randomString(12), owner.ID, map[string]interface{}{
		"members": membersID,
	})

	require.NoError(t, err, "create channel")
	return ch
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

func TestSendUserCustomEvent(t *testing.T) {
	c := initClient(t)

	tests := []struct {
		name         string
		event        *UserCustomEvent
		targetUserID string
		expectedErr  string
	}{
		{
			name: "ok",
			event: &UserCustomEvent{
				Type: "custom_event",
			},
			targetUserID: "user1",
		},
		{
			name:         "error: event is nil",
			event:        nil,
			targetUserID: "user1",
			expectedErr:  "event is nil",
		},
		{
			name:         "error: empty targetUserID",
			event:        &UserCustomEvent{},
			targetUserID: "",
			expectedErr:  "targetUserID should not be empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedErr == "" {
				_, err := c.UpsertUser(context.Background(), &User{ID: test.targetUserID})
				require.NoError(t, err)
			}

			err := c.SendUserCustomEvent(context.Background(), test.targetUserID, test.event)

			if test.expectedErr == "" {
				require.NoError(t, err)
				return
			}
			require.EqualError(t, err, test.expectedErr)
		})
	}
}

func TestMarshalUnmarshalUserCustomEvent(t *testing.T) {
	ev1 := UserCustomEvent{
		Type: "custom_event",
		ExtraData: map[string]interface{}{
			"name":   "John Doe",
			"age":    99.0,
			"hungry": true,
			"fruits": []interface{}{},
		},
	}

	b, err := json.Marshal(ev1)
	require.NoError(t, err)

	ev2 := UserCustomEvent{}
	err = json.Unmarshal(b, &ev2)
	require.NoError(t, err)

	require.Equal(t, ev1, ev2)
}
