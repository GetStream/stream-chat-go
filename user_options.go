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

/*
TODO: clean this up

const (
	hardDeleteKey          = "hard_delete"
	markMessagesDeletedKey = "mark_messages_deleted"
)

// UserOptions is a helper for setting various options for the user actions.
//easyjson:nojson
type UserOptions struct {
	// HardDelete is used by DeleteUser to indicate that the user should be
	// completely deleted with no records retained.
	HardDelete *bool `json:"hard_delete,omitempty"`

	// MarkMessagesDeleted is used by DeleteUser to indicate that all of the
	// Users messages should be deleted as well.
	MarkMessagesDeleted *bool `json:"mark_messages_deleted"`

	// Extra provides a trapdoor of sorts. If the value/setting you want isn't
	// builtin you can add it to Extra and it will be sent to the server.
	Extra map[string]interface{} `json:"-"`
}

func (u UserOptions) output() output {
	out := output{}
	out.addBool(hardDeleteKey, u.HardDelete)
	out.addBool(markMessagesDeletedKey, u.MarkMessagesDeleted)

	for k, v := range u.Extra {
		out.add(k, v)
	}
	return out
}

func (u UserOptions) len() int {
	return len(u.output().cast())
}

// URLValues returns the UserOptions as url.Values.
func (u UserOptions) URLValues() url.Values {
	return u.output().urlValues()
}

// MarshalJSON takes a UserOption and preps it to be sent to the server.
func (u UserOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.output())
}

// output is a set of helper functions to take Option structs and serialise them
// for sending to the API.
type output map[string]interface{}

func (o output) cast() map[string]interface{} {
	return ((map[string]interface{})(o))
}

// addBool takes a pointer to a bool and if the value is set adds it to the
// output as a string.
func (o output) addBool(key string, b *bool) {
	if b == nil {
		return
	}

	o.add(key, *b)
}

func (o output) add(key string, value interface{}) {
	o.cast()[key] = value
}

// urlValues converts the output into the appropriate url.Values.
func (o output) urlValues() url.Values {
	values := url.Values{}
	for k, v := range o.cast() {
		switch x := v.(type) {
		case string:
			values.Add(k, x)
		default:
			values.Add(k, fmt.Sprintf("%v", v))
		}
	}

	return values
}

*/
