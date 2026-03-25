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

func TestExportUserResponse_UnmarshalJSON(t *testing.T) {
	// This is the actual response format returned by the API.
	// Previously, ExportUserResponse embedded *User directly, which caused
	// User.UnmarshalJSON to consume the entire response body, losing the
	// user, messages, and reactions data.
	apiResponse := `{
		"user": {
			"id": "103415720",
			"name": "Batman",
			"language": "",
			"role": "user",
			"teams": [],
			"created_at": "2025-05-06T19:41:07.894092Z",
			"updated_at": "2025-05-06T20:15:52.812595Z",
			"banned": false,
			"online": false,
			"last_active": "2026-03-10T14:09:24.664584Z",
			"blocked_user_ids": [],
			"shadow_banned": false,
			"invisible": false
		},
		"messages": [
			{
				"id": "msg1",
				"cid": "messaging:general",
				"text": "Hello world",
				"user": {"id": "103415720"},
				"user_id": "103415720"
			}
		],
		"reactions": [
			{
				"message_id": "msg1",
				"user_id": "103415720",
				"type": "like"
			}
		],
		"duration": "117.63ms"
	}`

	var resp ExportUserResponse
	err := json.Unmarshal([]byte(apiResponse), &resp)
	require.NoError(t, err)

	require.NotNil(t, resp.User)
	require.Equal(t, "103415720", resp.User.ID)
	require.Equal(t, "Batman", resp.User.Name)
	require.Equal(t, "user", resp.User.Role)

	require.Len(t, resp.Messages, 1)
	require.Equal(t, "msg1", resp.Messages[0].ID)
	require.Equal(t, "Hello world", resp.Messages[0].Text)

	require.Len(t, resp.Reactions, 1)
	require.Equal(t, "msg1", resp.Reactions[0].MessageID)
	require.Equal(t, "like", resp.Reactions[0].Type)
}
