package stream_chat

import (
	"time"

	"github.com/francoispqt/gojay"
)

func (u *User) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	if u.ExtraData == nil {
		u.ExtraData = map[string]interface{}{}
	}

	switch key {
	// strings
	case "id":
		return dec.String(&u.ID)
	case "name":
		return dec.String(&u.Name)
	case "image":
		return dec.String(&u.Image)
	case "role":
		return dec.String(&u.Role)

		//time
	case "last_active":
		return dec.Time(&u.LastActive, time.RFC3339)
	case "created_at":
		return dec.Time(&u.CreatedAt, time.RFC3339)
	case "updated_at":
		return dec.Time(&u.UpdatedAt, time.RFC3339)

		// bool
	case "online":
		return dec.Bool(&u.Online)
	case "invisible":
		return dec.Bool(&u.Invisible)
	case "mutes":
		return dec.DecodeArray(&u.Mutes)

	// extra fields
	default:
		var i interface{}
		err := dec.Interface(&i)
		u.ExtraData[key] = i
		return err
	}
}

func (u *User) NKeys() int {
	return 0
}

func (u *User) marshalMap() map[string]interface{} {
	var res = map[string]interface{}{}

	// fill extra data first to avoid main fields overwrite
	for k, v := range u.ExtraData {
		res[k] = v
	}

	res["id"] = u.ID
	res["name"] = u.Name
	res["image"] = u.Image
	res["role"] = u.Role
	res["online"] = u.Online
	res["invisible"] = u.Invisible
	res["mutes"] = u.Mutes

	return res
}
