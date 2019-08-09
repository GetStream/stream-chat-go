package stream_chat

import (
	"time"

	"github.com/francoispqt/gojay"
)

type Mute struct {
	User      User
	Target    User
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Mute) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "user":
		return dec.Object(&m.User)
	case "target":
		return dec.Object(&m.Target)
	case "created_at":
		return dec.Time(&m.CreatedAt, time.RFC3339)
	case "updated_at":
		return dec.Time(&m.UpdatedAt, time.RFC3339)
	default:
		// just skip unknown fields
	}

	return nil
}

func (m *Mute) NKeys() int {
	return 0
}

type Mutes []Mute

func (m *Mutes) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var mut Mute
	err := dec.Object(&mut)
	if err != nil {
		return err
	}
	*m = append(*m, mut)
	return nil
}
