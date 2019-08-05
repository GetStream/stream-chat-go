package stream_chat

import (
	"net/http"
	"testing"
	"time"

	"github.com/pascaldekloe/jwt"

	"github.com/stretchr/testify/assert"
)

func initClient(t *testing.T) (c *Client, ch *Channel) {
	c, err := NewClient(APIKey, []byte(APISecret), WithBaseURL("http://localhost:3030"))
	mustNoError(t, err)

	err = c.UpdateUsers(testUsers...)
	mustNoError(t, err)

	members := make([]string, 0, len(testUsers))
	for i := range testUsers {
		members = append(members, testUsers[i].ID)
	}

	ch = c.NewChannel("team", "fellowship-of-the-ring", map[string]interface{}{
		"members": members,
	})

	_, err = ch.Create("gendalf")
	mustNoError(t, err)

	return c, ch
}

func TestNewClient(t *testing.T) {
	c, _ := initClient(t)

	assert.Equal(t, c.apiKey, APIKey)
	assert.Equal(t, c.apiSecret, []byte(APISecret))
	assert.NotEmpty(t, c.authToken)
	assert.Equal(t, c.timeout, defaultTimeout)
	//	assert.Equal(t, c.baseURL, defaultBaseURL)
	assert.Equal(t, c.http, http.DefaultClient)
	assert.Equal(t, c.http.Timeout, defaultTimeout)
}

func Test_client_CreateToken(t *testing.T) {
	c, _ := initClient(t)

	var expire = time.Now().Add(time.Hour)
	tt := []struct {
		name   string
		expire *time.Time
	}{
		{"token without expire", nil},
		{"token with expire", &expire},
	}

	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			token, err := c.CreateToken(testUsers[0].ID, test.expire)
			mustNoError(t, err)

			claims, err := jwt.HMACCheck(token, c.apiSecret)
			mustNoError(t, err)

			var expiresIn *jwt.NumericTime
			if test.expire != nil {
				expiresIn = jwt.NewNumericTime(*test.expire)
			}

			assert.Equal(t, expiresIn, claims.Expires)
			assert.Equal(t, testUsers[0].ID, claims.Set["user_id"])
		})
	}
}

func TestWithBaseURL(t *testing.T) {
	c, _ := initClient(t)

	u := "http://test:3030"
	WithBaseURL(u)(c)
	assert.Equal(t, u, c.baseURL)
}

func TestWithTimeout(t *testing.T) {
	c, _ := initClient(t)

	timeout := time.Hour

	WithTimeout(timeout)(c)

	assert.Equal(t, timeout, c.timeout)
	assert.Equal(t, timeout, c.http.Timeout)
}

func TestWithHTTPTransport(t *testing.T) {
	c, _ := initClient(t)

	tr := &http.Transport{}

	WithHTTPTransport(tr)(c)

	assert.Equal(t, tr, c.http.Transport)
}
