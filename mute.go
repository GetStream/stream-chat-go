package stream_chat

import "time"

type Mute struct {
	User      User
	Target    User
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Mutes []Mute

func (m *Mutes) unmarshalSlice(slice []interface{}) {
	*m = make([]Mute, 0, len(slice))
	for _, v := range slice {
		switch val := v.(type) {
		case map[string]interface{}:
			var mut Mute
			mut.unmarshalMap(val)
			*m = append(*m, mut)

		default:
			//TODO: error
		}
	}
}

func (m *Mute) timeMap() map[string]*time.Time {
	return map[string]*time.Time{
		"created_at": &m.CreatedAt,
		"updated_at": &m.UpdatedAt,
	}
}

func (m *Mute) structMap() map[string]unmarshalMap {
	return map[string]unmarshalMap{
		"user":   &m.User,
		"target": &m.Target,
	}
}

func (m *Mute) unmarshalMap(data map[string]interface{}) {
	structMap := m.structMap()
	timeMap := m.timeMap()

	for k, v := range data {
		switch val := v.(type) {
		case string:
			if p, ok := timeMap[k]; ok {
				// todo: error handling
				t, _ := time.Parse(time.RFC3339, val)
				*p = t
			} else {
				// TODO: logging
			}

		case map[string]interface{}:
			if p, ok := structMap[k]; ok {
				p.unmarshalMap(val)
			} else {
				// TODO: logging
			}

		default:
			// TODO: logging
		}
	}
}
