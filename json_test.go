package stream_chat

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomExtraData(in interface{}) {
	v := reflect.ValueOf(in).Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	f := v.FieldByName("ExtraData")
	f.Set(reflect.ValueOf(map[string]interface{}{
		"mystring":    randomString(10),
		"mybool":      rand.Float64() < 0.5,
		"data":        "custom",
		"custom_data": "really_custom",
		"extra": map[string]interface{}{
			randomString(10): randomString(10),
		},
		"stream":   randomString(10),
		"my_score": float64(rand.Intn(100)),
	}))
}

func testInvariantJSON(t *testing.T, in, in2 interface{}) {
	t.Helper()

	// put random
	randomExtraData(in)

	// marshal given
	data, err := json.Marshal(in)
	require.NoError(t, err)

	// unmarshal again
	require.NoError(t, json.Unmarshal(data, in2))

	// ensure they are same
	require.Equal(t, in, in2)
}

func TestJSON(t *testing.T) {
	var c, c2 Channel
	testInvariantJSON(t, &c, &c2)

	var u, u2 User
	testInvariantJSON(t, &u, &u2)

	var e, e2 Event
	testInvariantJSON(t, &e, &e2)

	var m, m2 Message
	testInvariantJSON(t, &m, &m2)

	var mr, mr2 messageRequestMessage
	testInvariantJSON(t, &mr, &mr2)

	var a, a2 Attachment
	testInvariantJSON(t, &a, &a2)

	var r, r2 Reaction
	testInvariantJSON(t, &r, &r2)
}

// TestFlattenExtraData tests the flattenExtraData function directly
func TestFlattenExtraData(t *testing.T) {
	t.Run("Flatten nested extra_data", func(t *testing.T) {
		m := map[string]interface{}{
			"field1": "value1",
			"extra_data": map[string]interface{}{
				"custom_field":  "custom_value",
				"another_field": 123,
			},
		}

		flattenExtraData(m)

		// Fields should be flattened
		require.Equal(t, "custom_value", m["custom_field"])
		require.Equal(t, 123, m["another_field"])
		require.Equal(t, "value1", m["field1"])
		// The nested "extra_data" key should not exist
		require.NotContains(t, m, "extra_data")
	})

	t.Run("No extra_data key", func(t *testing.T) {
		m := map[string]interface{}{
			"field1": "value1",
			"field2": 123,
		}

		flattenExtraData(m)

		// Map should be unchanged
		require.Equal(t, "value1", m["field1"])
		require.Equal(t, 123, m["field2"])
		require.Len(t, m, 2)
	})

	t.Run("extra_data is not a map", func(t *testing.T) {
		m := map[string]interface{}{
			"field1":     "value1",
			"extra_data": "not_a_map",
		}

		flattenExtraData(m)

		// extra_data should remain unchanged if it's not a map
		require.Equal(t, "not_a_map", m["extra_data"])
		require.Equal(t, "value1", m["field1"])
	})

	t.Run("Empty extra_data map", func(t *testing.T) {
		m := map[string]interface{}{
			"field1":     "value1",
			"extra_data": map[string]interface{}{},
		}

		flattenExtraData(m)

		// Empty extra_data should be removed
		require.NotContains(t, m, "extra_data")
		require.Equal(t, "value1", m["field1"])
	})
}
