package stream_chat

type Reaction struct {
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id,omitempty"`
	Type      string `json:"type"`

	// any other fields the user wants to attach a reaction
	ExtraData map[string]interface{}
}

func (r Reaction) toHash() map[string]interface{} {
	hash := r.ExtraData
	if hash == nil {
		hash = map[string]interface{}{}
	}

	hash["message_id"] = r.MessageID

	if r.UserID != "" {
		hash["user_id"] = r.UserID
	}

	hash["type"] = r.Type

	return hash
}
