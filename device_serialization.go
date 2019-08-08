package stream_chat

func (d *Device) stringsMap() map[string]*string {
	return map[string]*string{
		"user_id":       &d.UserID,
		"id":            &d.ID,
		"push_provider": &d.PushProvider,
	}
}

func (d *Device) unmarshalMap(data map[string]interface{}) {
	stringsMap := d.stringsMap()

	for k, v := range data {
		switch val := v.(type) {
		case string:
			if p, ok := stringsMap[k]; ok {
				*p = val
			} else {
				// TODO: logging
			}

		default:
			//todo: logging
		}
	}
}

func (d Device) marshalMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":       d.UserID,
		"id":            d.ID,
		"push_provider": d.PushProvider,
	}
}
