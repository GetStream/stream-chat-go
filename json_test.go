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
		"extra_data": map[string]interface{}{
			"mystring": randomString(10),
			"mybool":   rand.Float64() < 0.5,
		},
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
