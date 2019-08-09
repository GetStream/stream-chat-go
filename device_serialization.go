package stream_chat

import "github.com/francoispqt/gojay"

func (d *Device) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "user_id":
		return dec.String(&d.UserID)
	case "id":
		return dec.String(&d.ID)
	case "push_provider":
		return dec.String(&d.PushProvider)
	}

	return nil
}

func (d Device) NKeys() int {
	return 0
}

func (d Device) marshalMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":       d.UserID,
		"id":            d.ID,
		"push_provider": d.PushProvider,
	}
}
