package stream_chat //nolint:golint

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleNewOption() {
	client, _ := NewClient("XXXX", "XXXX")
	opt := NewOption("new_awesome_feature", true)

	_ = client.BanUser("badUser", "awesomeMod", opt)
}

func TestOptionTimeout(t *testing.T) {
	type testCase struct {
		input    time.Duration
		expected int
	}

	for _, c := range []testCase{
		{input: 5 * time.Second, expected: 5},
		{input: 60 * time.Minute, expected: 3600},
		{input: time.Second / 2, expected: 1},
	} {
		opt := OptionTimeout(c.input)

		assert.Equal(t, c.expected, opt.Value())
	}
}
