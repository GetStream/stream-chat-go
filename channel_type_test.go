package stream_chat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareChannelType(t *testing.T, c *Client) *ChannelType {
	ct := NewChannelType(randomString(10))

	err := c.CreateChannelType(&ct)
	mustNoError(t, err)

	return &ct
}

func TestClient_GetChannelType(t *testing.T) {
	c := initClient(t)

	ct := prepareChannelType(t, c)
	defer c.DeleteChannelType(ct.Name)

	got, err := c.GetChannelType(ct.Name)
	mustNoError(t, err)

	assert.Equal(t, ct.Name, got.Name)
	assert.Equal(t, len(ct.Commands), len(got.Commands))
	assert.Equal(t, ct.Permissions, got.Permissions)
}

func TestClient_ListChannelTypes(t *testing.T) {
	c := initClient(t)

	ct := prepareChannelType(t, c)
	defer c.DeleteChannelType(ct.Name)

	got, err := c.ListChannelTypes()
	mustNoError(t, err)

	assert.Contains(t, got, ct.Name)
}
