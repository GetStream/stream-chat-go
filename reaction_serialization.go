package stream_chat

func (r *Reaction) stringsMap() map[string]*string {
	return map[string]*string{
		"message_id": &r.MessageID,
		"user_id":    &r.UserID,
		"type":       &r.Type,
	}
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

func (r *Reaction) unmarshalMap(data map[string]interface{}) {
	stringsMap := r.stringsMap()

	for k, v := range data {
		switch val := v.(type) {
		case string:
			if p, ok := stringsMap[k]; ok {
				*p = val
			} else {
				r.ExtraData[k] = val
			}

		default:
			r.ExtraData[k] = val
		}
	}
}
