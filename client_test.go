package stream_chat

import (
	"reflect"
	"testing"
	"time"

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

	ch, err := c.CreateChannel("team", "fellowship-of-the-ring", serverUser.ID, map[string]interface{}{
		"members": members,
	})

	mustNoError(t, err, "create channel")
	return ch
}

func TestNewClient(t *testing.T) {
	c := initClient(t)

	assert.Equal(t, c.apiKey, APIKey)
	assert.Equal(t, c.apiSecret, []byte(APISecret))
	assert.NotEmpty(t, c.authToken)
	assert.Equal(t, defaultTimeout, c.HTTP.Timeout)
	//	assert.Equal(t, defaultBaseURL, c.BaseURL, )
}

func TestClient_CreateToken(t *testing.T) {
	type args struct {
		userId string
		expire time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"simple without expiration", args{"tommaso", time.Time{}}, []byte("eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoidG9tbWFzbyJ9.oQLtgTc9_SIr3Rvrq-eW_WrLmdO1gAAYA335qTatxrU"), false},
		{"simple with expiration", args{"tommaso", time.Unix(1566941272, 123121)}, []byte("eyJhbGciOiJIUzI1NiJ9.eyJleHAiOjE1NjY5NDEyNzIsInVzZXJfaWQiOiJ0b21tYXNvIn0.bkMDhCJhzKKnSZO27QcP8n3o7u9C1TpoMt0MD-JCNnY"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewClient("key", []byte("secret"))
			got, err := c.CreateToken(tt.args.userId, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createToken() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
