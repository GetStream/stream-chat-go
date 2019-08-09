package stream_chat

import (
	"time"

	"github.com/francoispqt/gojay"
)

func (m *Message) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	if m.ExtraData == nil {
		m.ExtraData = map[string]interface{}{}
	}

	switch key {
	case "id":
		return dec.String(&m.ID)
	case "text":
		return dec.String(&m.Text)
	case "html":
		return dec.String(&m.HTML)
	case "type":
		return dec.String(&m.Type)

	case "user":
		return dec.ObjectNull(&m.User)
	case "mentioned_users":
		return dec.Array(&m.MentionedUsers)

	case "attachments":
		return dec.Array(&m.Attachments)
	case "latest_reactions":
		return dec.Array(&m.LatestReactions)
	case "own_reactions":
		return dec.Array(&m.OwnReactions)
	case "reaction_counts":
		m.ReactionCounts = reactionsCount{}
		return dec.Object(&m.ReactionCounts)

	case "reply_count":
		return dec.Int(&m.ReplyCount)

	case "created_at":
		return dec.Time(&m.CreatedAt, time.RFC3339)
	case "updated_at":
		return dec.Time(&m.UpdatedAt, time.RFC3339)

	default:
		var i interface{}
		if err := dec.Interface(&i); err != nil {
			return err
		}
		m.ExtraData[key] = i
	}

	return nil
}

func (m Message) NKeys() int {
	return 0
}
