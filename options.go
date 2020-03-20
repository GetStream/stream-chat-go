package stream_chat //nolint:golint

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
	optionKeyTargetUserID = "target_user_id"
	optionKeyTimeout      = "timeout"
	optionKeyLocation     = "url"
	optionKeyReason       = "reason"
)

func compileOptions(opts ...Option) url.Values {
	val := url.Values{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		val.Add(opt.Key(), optionAsString(opt))
	}

	return val
}

func optionAsString(o Option) string {
	switch v := o.Value().(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
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

// OptionTimeout returns a timeout option for banning or muting a user. If the
// time is less than one second then
func OptionTimeout(duration time.Duration) Option {
	if duration < time.Second {
		duration = time.Second
	}

	return NewOption(optionKeyTimeout, int(duration.Seconds()))
}

// OptionBanReason allows you to specify an optional reason for banning a user
// from chat.
func OptionBanReason(reason string) Option {
	return NewOption(optionKeyReason, reason)
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

// PaginateOffset returns an Option for paginating.
func PaginateOffset(offset int) Option {
	return NewOption("offset", offset)
}

// PaginateLimit returns an Option for paginating.
func PaginateLimit(limit int) Option {
	return NewOption("limit", limit)
}

// PaginateGreaterThanOrEqual returns an Option for paginating.
func PaginateGreaterThanOrEqual(id string) Option {
	return NewOption("id_gte", id)
}

// PaginateLessThan returns an Option for paginating.
func PaginateLessThan(id string) Option {
	return NewOption("id_lt", id)
}

// PaginateLessThanOrEqual returns an Option for paginating.
func PaginateLessThanOrEqual(id string) Option {
	return NewOption("id_lte", id)
}

// PaginateGreaterThan returns an Option for paginating.
func PaginateGreaterThan(id string) Option {
	return NewOption("id_gt", id)
}
