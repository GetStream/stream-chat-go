package stream_chat

import (
	"github.com/francoispqt/gojay"
)

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

func (c *ChannelConfig) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "name":
		return dec.String(&c.Name)

	case "typing_events":
		return dec.Bool(&c.TypingEvents)
	case "read_events":
		return dec.Bool(&c.ReadEvents)
	case "connect_events":
		return dec.Bool(&c.ConnectEvents)
	case "search":
		return dec.Bool(&c.Search)
	case "reactions":
		return dec.Bool(&c.Reactions)
	case "replies":
		return dec.Bool(&c.Replies)
	case "mutes":
		return dec.Bool(&c.Mutes)
	case "message_retention":
		return dec.String(&c.MessageRetention)
	case "max_message_length":
		return dec.Int(&c.MaxMessageLength)
	case "automod":
		var mod string
		if err := dec.String(&mod); err != nil {
			return err
		}
		c.Automod = modType(mod)

	case "automod_behaviour":
		var mod string
		if err := dec.String(&mod); err != nil {
			return err
		}
		c.ModBehavior = modBehaviour(mod)
	}

	return nil
}

func (c *ChannelConfig) NKeys() int {
	return 0
}
