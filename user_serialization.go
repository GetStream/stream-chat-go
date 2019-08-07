package stream_chat

import (
	"time"
)

func (u *User) stringFieldsMap() map[string]*string {
	return map[string]*string{
		"id":    &u.ID,
		"name":  &u.Name,
		"image": &u.Image,
		"role":  &u.Role,
	}
}

func (u *User) timeFieldsMap() map[string]*time.Time {
	return map[string]*time.Time{
		"last_active": &u.LastActive,
		"created_at":  &u.CreatedAt,
		"updated_at":  &u.UpdatedAt,
	}
}

func (u *User) sliceFieldsMap() map[string]unmarshalSlice {
	return map[string]unmarshalSlice{
		"mutes": &u.Mutes,
	}
}

func (u *User) boolFieldsMap() map[string]*bool {
	return map[string]*bool{
		"online":    &u.Online,
		"invisible": &u.Invisible,
	}
}

func (u *User) unmarshalMap(data map[string]interface{}) {
	strMap := u.stringFieldsMap()
	timeMap := u.timeFieldsMap()
	slicesMap := u.sliceFieldsMap()
	boolMap := u.boolFieldsMap()

	for key, v := range data {
		switch val := v.(type) {
		case string:
			// try to parse time first from string
			if p, ok := timeMap[key]; ok {
				// todo handle error
				t, _ := time.Parse(time.RFC3339, val)
				*p = t

			} else if p, ok := strMap[key]; ok {
				*p = val
			} else {
				u.ExtraData[key] = val
			}

		case bool:
			if p, ok := boolMap[key]; ok {
				*p = val
			} else {
				u.ExtraData[key] = val
			}

		case []interface{}:
			if p, found := slicesMap[key]; found {
				p.unmarshalSlice(val)
			} else {
				u.ExtraData[key] = val
			}

		default:
			u.ExtraData[key] = val
		}
	}
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
