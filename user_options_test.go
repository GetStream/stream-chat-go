package stream

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserOptions_MarshalJSON_basic(t *testing.T) {
	input := &UserOptions{
		HardDelete: Bool(true),
	}
	expected := `{"hard_delete": true}`

	bytes, err := input.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, expected, string(bytes))
}

func TestUserOptions_MarshalJSON_extra(t *testing.T) {
	input := &UserOptions{
		HardDelete: Bool(true),
		Extra: map[string]interface{}{
			"foo": "bar",
		},
	}
	expected := `{"hard_delete": true, "foo": "bar"}`

	bytes, err := input.MarshalJSON()
	require.NoError(t, err)

	require.JSONEq(t, expected, string(bytes))
}

func TestUserOptions_URLValues(t *testing.T) {
	input := &UserOptions{
		MarkMessagesDeleted: Bool(true),
	}
	expected := url.Values{}
	expected.Add("mark_messages_deleted", "true")

	values := input.URLValues()
	require.Equal(t, expected, values)
}
