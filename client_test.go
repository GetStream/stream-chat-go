package stream_chat

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func initClient(t *testing.T) *Client {
	t.Helper()

	c, err := NewClientFromEnvVars()
	require.NoError(t, err, "new client")

	return c
}

func initChannel(t *testing.T, c *Client, membersID ...string) *Channel {
	return initChannelWithType(t, c, "team", membersID...)
}

func initChannelWithType(t *testing.T, c *Client, channelType string, membersID ...string) *Channel {
	t.Helper()

	owner := randomUser(t, c)
	ctx := context.Background()

	resp, err := c.CreateChannelWithMembers(ctx, channelType, randomString(12), owner.ID, membersID...)
	require.NoError(t, err, "create channel")

	t.Cleanup(func() {
		_, _ = c.DeleteChannels(ctx, []string{resp.Channel.CID}, true)
	})

	return resp.Channel
}

func TestClient_SwapHttpClient(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	tr := http.DefaultTransport.(*http.Transport).Clone() //nolint:forcetypeassert
	proxyURL, _ := url.Parse("http://getstream.io")
	tr.Proxy = http.ProxyURL(proxyURL)
	cl := &http.Client{Transport: tr}
	c.SetClient(cl)
	_, err := c.GetAppSettings(ctx)
	require.Error(t, err)

	cl = &http.Client{}
	c.SetClient(cl)
	_, err = c.GetAppSettings(ctx)
	require.NoError(t, err)
}

func TestClient_CreateToken(t *testing.T) {
	type args struct {
		userID string
		expire time.Time
		iat    time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"simple without expiration and iat",
			args{"tommaso", time.Time{}, time.Time{}},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidG9tbWFzbyJ9.v-x-jt3ZnBXXbQ0GoWloIZtVnat2IE74U1a4Yuxd63M",
			false,
		},
		{
			"simple with expiration and iat",
			args{"tommaso", time.Unix(1566941272, 123121), time.Unix(1566941272, 123121)},
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjY5NDEyNzIsImlhdCI6MTU2Njk0MTI3MiwidXNlcl9pZCI6InRvbW1hc28ifQ.3HY2O_7o5ZjZ-6KCXLzyPpHZOlNEDy6_m3iNb5DKAMY",
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient("key", "secret")
			require.NoError(t, err)

			got, err := c.CreateToken(tt.args.userID, tt.args.expire, tt.args.iat)
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
