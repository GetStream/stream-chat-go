package stream_chat

type Reaction struct {
	MessageID string
	UserID    string
	Type      string

	// any other fields the user wants to attach a reaction
	ExtraData map[string]interface{}
}
