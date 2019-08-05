package stream_chat

type StreamChatAPI interface {
	DeviceAPI
	UserAPI
	ChannelAPI
	UpdateAppSettings()
}

type ChannelAPI interface {
	// Creates channel type
	AddChannelType(data interface{}) error
	// Get channel type
	GetChannelType(chanType string) (map[string]interface{}, error)
	// List all channel types
	ListChannelTypes() (map[string]interface{}, error)
	// Creates channel object
	NewChannel(chanType string, chanId string, data map[string]interface{})
	//TODO: Search for a channel
	// QueryChannels()
}
