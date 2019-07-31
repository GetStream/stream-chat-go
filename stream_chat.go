package stream_chat

type StreamChatAPI interface {
	DeviceAPI
	UserAPI
	ChannelAPI
	UpdateAppSettings()
}
