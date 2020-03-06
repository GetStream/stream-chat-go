package stream

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	hardDeleteKey          = "hard_delete"
	markMessagesDeletedKey = "mark_messages_deleted"
)

// UserOptions is a helper for setting various options for the user actions.
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
	out.addBool("hard_delete", u.HardDelete)
	out.addBool(markMessagesDeletedKey, u.MarkMessagesDeleted)

	for k, v := range u.Extra {
		out.add(k, v)
	}
	return out
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
