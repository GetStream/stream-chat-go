package stream_chat // nolint: golint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareChannelType(t *testing.T, c *Client) *ChannelType {
	ct := NewChannelType(randomString(10))

	ct, err := c.CreateChannelType(ct)
	mustNoError(t, err, "create channel type")

	return ct
}

func TestClient_GetChannelType(t *testing.T) {
	c := initClient(t)

	ct := prepareChannelType(t, c)
	defer func() {
		mustNoError(t, c.DeleteChannelType(ct.Name), "delete channel type")
	}()

	got, err := c.GetChannelType(ct.Name)
	mustNoError(t, err, "get channel type")

	assert.Equal(t, ct.Name, got.Name)
	assert.Equal(t, len(ct.Commands), len(got.Commands))
	assert.Equal(t, ct.Permissions, got.Permissions)
}

func TestClient_ListChannelTypes(t *testing.T) {
	c := initClient(t)

	ct := prepareChannelType(t, c)
	defer func() {
		mustNoError(t, c.DeleteChannelType(ct.Name), "delete channel type")
	}()

	got, err := c.ListChannelTypes()
	mustNoError(t, err, "list channel types")

	assert.Contains(t, got, ct.Name)
}
