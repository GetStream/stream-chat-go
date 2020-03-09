package stream

import (
	"fmt"
	"net/url"
	"time"
)

var (
	// OptionHardDelete tells the API to do a hard delete instead of a
	// normal soft delete.
	OptionHardDelete = NewOption("hard_delete", true)

	// OptionMarkMessagesDeleted tells the API to mark all messages belonging to
	// the user as deleted in addition to deleting the user.
	OptionMarkMessagesDeleted = NewOption("mark_messages_deleted", true)
)

const (
	optionKeyType         = "type"
	optionKeyID           = "id"
	optionKeyUserID       = "user_id"
	optionKeyTargetUserID = "target_user_id"
	optionKeyTimeout      = "timeout"
	optionKeyLocation     = "url"
)

func compileOptions(opts ...Option) url.Values {
	val := url.Values{}
	for _, opt := range opts {
		switch v := opt.Value().(type) {
		case string:
			val.Add(opt.Key(), v)
		default:
			val.Add(opt.Key(), fmt.Sprintf("%v", v))
		}
	}

	return val
}

// Option represents a optional value that can be sent with an API call.
type Option interface {
	Key() string
	Value() interface{}
}

// NewOption is provided as a trapdoor for easily adding options that aren't
// explicitly supported by the library.
func NewOption(key string, value interface{}) Option {
	return option{key: key, value: value}
}

// makeTargetID is a helper function for making a target ID option.
func makeTargetID(targetID string) Option {
	return NewOption(optionKeyTargetUserID, targetID)
}

// makeTargetID is a helper function for making a userID option.
func makeUserID(userID string) Option {
	return NewOption(optionKeyUserID, userID)
}

// OptionTimeout returns a timeout option for banning or muting a user. If the
// time is less than one second then TODO: something here!!!! fatal or round up?
func OptionTimeout(duration time.Duration) Option {
	return NewOption(optionKeyTimeout, int(duration.Seconds()))
}

// optionLocation returns an option for a location.
func optionURL(location string) Option {
	return NewOption(optionKeyLocation, location)
}

type option struct {
	key   string
	value interface{}
}

func (o option) Key() string {
	return o.key
}

func (o option) Value() interface{} {
	return o.value
}
