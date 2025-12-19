package stream_chat

import (
	"encoding/json"
	"reflect"
	"strings"
)

func copyMap(m map[string]interface{}) map[string]interface{} {
	m2 := make(map[string]interface{}, len(m))
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

func removeFromMap(m map[string]interface{}, obj interface{}) {
	t := reflect.TypeOf(obj)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if tag := f.Tag.Get("json"); tag != "" {
			tag = strings.Split(tag, ",")[0]
			delete(m, tag)
		} else {
			delete(m, f.Name)
		}
	}
}

// flattenExtraData flattens the nested "extra_data" key if it exists.
// The API may return custom fields nested under an "extra_data" key,
// which should be flattened to the root level of the ExtraData map.
func flattenExtraData(m map[string]interface{}) {
	if extraData, ok := m["extra_data"]; ok {
		if extraDataMap, ok := extraData.(map[string]interface{}); ok {
			// Copy all fields from nested extra_data to root level
			for k, v := range extraDataMap {
				m[k] = v
			}
			// Remove the nested extra_data key
			delete(m, "extra_data")
		}
	}
}

func addToMapAndMarshal(m map[string]interface{}, obj interface{}) ([]byte, error) {
	m2 := copyMap(m)

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &m2); err != nil {
		return nil, err
	}
	return json.Marshal(m2)
}
