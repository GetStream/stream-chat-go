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
	GetChannelType(chanType channelType)
	// List all channel types
	ListChannelTypes()
	// Creates channel object
	NewChannel(chanType channelType, chanId string, data map[string]interface{})
	//TODO: Search for a channel
	// QueryChannels()
}
