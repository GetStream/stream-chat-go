package stream_chat

type channelType string

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

func (*client) AddChannelType(data interface{}) error {
	panic("implement me")
}

func (*client) GetChannelType(chanType channelType) {
	panic("implement me")
}

func (*client) ListChannelTypes() {
	panic("implement me")
}

func (*client) NewChannel(chanType channelType, chanId string, data map[string]interface{}) {
	panic("implement me")
}
