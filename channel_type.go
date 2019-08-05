package stream_chat

import "net/http"

type ChannelType struct {
}

// CreateChannelType adds new channel type
func (c *client) CreateChannelType(data map[string]interface{}) (err error) {
	if data["commands"] == "" {
		data["commands"] = "all"
	}
	err = c.makeRequest(http.MethodPost, "channeltypes", nil, data, nil)

	return
}

// GetChannelType returns information about channel type
func (c *client) GetChannelType(chanType string) (resp map[string]interface{}, err error) {
	err = c.makeRequest(http.MethodGet, "channeltypes/"+chanType, nil, nil, &resp)

	return
}

// ListChannelTypes returns all channel types
func (c *client) ListChannelTypes() (resp map[string]interface{}, err error) {
	err = c.makeRequest(http.MethodGet, "channeltypes", nil, nil, &resp)

	return
}
