package stream_chat

type messageType = string

const (
	MessageTypeRegular   messageType = "regular"
	MessageTypeError     messageType = "error"
	MessageTypeReply     messageType = "reply"
	MessageTypeSystem    messageType = "system"
	MessageTypeEphemeral messageType = "ephemeral"
)
