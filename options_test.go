package stream

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleNewOption() {
	client, _ := NewClient("XXXX", "XXXX")
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
