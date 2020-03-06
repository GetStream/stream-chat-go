package stream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
TODO: clean this up

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

*/

func ExampleNewOption() {
	client, _ := NewClient("XXXX", []byte("XXXX"))
	opt := NewOption("new_awesome_feature", true)

	client.BanUser("badUser", "awesomeMod", opt)
}

func TestOptionTimeout(t *testing.T) {
	type testCase struct {
		input    time.Duration
		expected int
	}

	for _, c := range []testCase{
		testCase{input: 5 * time.Second, expected: 5},
		testCase{input: 60 * time.Minute, expected: 3600},
		testCase{input: time.Second / 2, expected: 1}, // TODO: is this correct behaviour?
	} {
		opt := OptionTimeout(c.input)

		assert.Equal(t, c.expected, opt.Value())
	}
}
