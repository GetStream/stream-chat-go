package stream_chat //nolint: golint

// ChannelConfig is the configuration for a channel.
type ChannelConfig struct {
	Name string `json:"name"`

	// features
	// show typing indicators or not (probably auto disable if more than X users in a channel)
	TypingEvents bool `json:"typing_events"`
	// store who has read the message, or at least when they last viewed the chat
	ReadEvents bool `json:"read_events"`
	// connect events can get very noisy for larger chat groups
	ConnectEvents bool `json:"connect_events"`
	// make messages searchable
	Search    bool `json:"search"`
	Reactions bool `json:"reactions"`
	Replies   bool `json:"replies"`
	Mutes     bool `json:"mutes"`
	// enable/disable push notifications
	PushNotifications bool `json:"push_notifications"`
	Uploads           bool `json:"uploads"`
	URLEnrichment     bool `json:"url_enrichment"`
	CustomEvents      bool `json:"custom_events"`

	// number of days to keep messages, must be MessageRetentionForever or numeric string
	MessageRetention string `json:"message_retention"`
	MaxMessageLength int    `json:"max_message_length"`

	Automod     modType      `json:"automod"` // disabled, simple or AI
	ModBehavior modBehaviour `json:"automod_behavior"`

	BlockList         string       `json:"blocklist"`
	BlockListBehavior modBehaviour `json:"blocklist_behavior"`
	AutomodThresholds *Thresholds  `json:"automod_thresholds"`
}

type LabelThresholds struct {
	Flag  float32 `json:"flag"`
	Block float32 `json:"block"`
}

type Thresholds struct {
	Explicit *LabelThresholds `json:"explicit"`
	Spam     *LabelThresholds `json:"spam"`
	Toxic    *LabelThresholds `json:"toxic"`
}

// DefaultChannelConfig is the default channel configuration.
var DefaultChannelConfig = ChannelConfig{
	Automod:           AutoModDisabled,
	ModBehavior:       ModBehaviourFlag,
	MaxMessageLength:  defaultMessageLength,
	MessageRetention:  MessageRetentionForever,
	PushNotifications: true,
}
