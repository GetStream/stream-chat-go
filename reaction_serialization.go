package stream_chat

import "github.com/francoispqt/gojay"

func (r *Reaction) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	if r.ExtraData == nil {
		r.ExtraData = map[string]interface{}{}
	}

	switch key {
	// strings
	case "message_id":
		return dec.String(&r.MessageID)
	case "user_id":
		return dec.String(&r.UserID)
	case "type":
		return dec.String(&r.Type)
	default:
		var i interface{}
		if err := dec.Interface(&i); err != nil {
			return err
		}
		r.ExtraData[key] = i
	}

	return nil
}

func (r Reaction) NKeys() int {
	return 0
}

func (r *Reaction) marshalMap() map[string]interface{} {
	resp := map[string]interface{}{}

	for k, v := range r.ExtraData {
		resp[k] = v
	}

	resp["type"] = r.Type
	resp["message_id"] = r.MessageID
	//optional
	if r.UserID != "" {
		resp["user_id"] = r.UserID
	}
	return resp
}
