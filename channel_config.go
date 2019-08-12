package stream_chat

type ChannelConfig struct {
	Name string `json:"name"`

	// features
	TypingEvents  bool `json:"typing_events"`  // show typing indicators or not (probably auto disable if more than X users in a channel)
	ReadEvents    bool `json:"read_events"`    // store who has read the message, or at least when they last viewed the chat
	ConnectEvents bool `json:"connect_events"` // connect events can get very noisy for larger chat groups
	Search        bool `json:"search"`         // make messages searchable
	Reactions     bool `json:"reactions"`
	Replies       bool `json:"replies"`
	Mutes         bool `json:"mutes"`

	MessageRetention string `json:"message_retention"` // number of days to keep messages, must be MessageRetentionForever or numeric string
	MaxMessageLength int    `json:"max_message_length"`

	Automod     modType      `json:"automod"` // disabled, simple or AI
	ModBehavior modBehaviour `json:"automod_behavior"`
}
