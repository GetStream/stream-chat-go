package stream_chat

import (
	"fmt"
	"net/http"
)

const (
	channelPathFmt = "channels/%s/%s"
	messagePathFmt = "channels/%s/%s/message"
	eventPathFmt   = "channels/%s/%s/event"
)

type Channel struct {
	id         string
	_type      string
	client     *Client
	customData map[string]interface{}
}

func (ch Channel) formatPath(path string, params ...interface{}) string {
	if ch.id == "" {
		return ""
	}

	params = append([]interface{}{ch._type, ch.id}, params...)

	return fmt.Sprintf(path, params)
}

func addUserID(hash map[string]interface{}, userID string) map[string]interface{} {
	hash["user"] = userID
	return hash
}

// SendMessage sends a message to this channel
//
// message: the Message object
// userID: the ID of the user that created the message
func (ch *Channel) SendMessage(message Message, userID string) error {
	data := map[string]interface{}{
		"message": addUserID(message.toHash(), userID),
	}

	return ch.client.makeRequest(http.MethodPost, ch.formatPath(messagePathFmt), nil, data, nil)
}

// SendEvent sends an event on this channel
//
// event: event data, ie {type: 'message.read'}
// userID: the ID of the user sending the event
func (ch *Channel) SendEvent(event Event, userID string) error {
	data := map[string]interface{}{
		"event": addUserID(event.toHash(), userID),
	}

	path := ch.formatPath(eventPathFmt)
	return ch.client.makeRequest(http.MethodPost, path, nil, data, nil)
}

// SendReaction sends a reaction about a message
//
// messageID: the message id
// reaction: the reaction object, ie {type: 'love'}
// userID: the ID of the user that created the reaction
func (ch *Channel) SendReaction(messageID string, reaction Reaction, userID string) error {
	data := map[string]interface{}{
		"reaction": addUserID(reaction.toHash(), userID),
	}

	path := "messages/" + messageID + "/reaction"
	return ch.client.makeRequest(http.MethodPost, path, nil, data, nil)
}

// DeleteReaction removes a reaction by user and type
//
// messageID: the id of the message from which te remove the reaction
// reaction_type: the type of reaction that should be removed
// userID: the id of the user
func (ch *Channel) DeleteReaction(messageID string, reactionType string, userID string) error {

	path := "messages/" + messageID + "/reaction/" + reactionType

	return ch.client.makeRequest(http.MethodDelete, path, nil, nil, nil)
}

// Create creates the channel
//
// userID: the ID of the user creating this channel
func (ch *Channel) Create(userID string) (map[string]interface{}, error) {
	ch.customData["created_by"] = map[string]interface{}{"id": userID}

	options := map[string]interface{}{
		"watch":    false,
		"state":    false,
		"presence": false,
	}

	return ch.Query(options)
}

// Query queries the API for this channel, get messages, members or other channel fields
//
// options: the query options, check docs on https://getstream.io/chat/docs/
func (ch *Channel) Query(options map[string]interface{}) (result map[string]interface{}, err error) {
	payload := map[string]interface{}{
		"state": true,
		"data":  ch.customData,
	}

	for k, v := range options {
		payload[k] = v
	}

	path := "channels/" + ch._type
	if ch.id != "" {
		path += "/" + ch.id
	}
	path += "/query"

	err = ch.client.makeRequest(http.MethodPost, path, nil, payload, &result)

	// TODO: set ch.id from result

	return result, err
}

// Update edits the channel's custom properties
//
// options: the object to update the custom properties of this channel with
// message: optional update message
func (ch *Channel) Update(options map[string]interface{}, message string) error {
	payload := map[string]interface{}{
		"data":    options,
		"message": message,
	}
	return ch.client.makeRequest(http.MethodPost, ch.formatPath(channelPathFmt), nil, payload, nil)
}

// Delete removes the channel. Messages are permanently removed.
func (ch *Channel) Delete() error {
	return ch.client.makeRequest(http.MethodDelete, ch.formatPath(channelPathFmt), nil, nil, nil)
}

// Truncate removes all messages from the channel
func (ch *Channel) Truncate() error {
	path := ch.formatPath(channelPathFmt) + "/truncate"
	return ch.client.makeRequest(http.MethodPost, path, nil, nil, nil)
}

// Adds members to the channel
//
// users: user IDs to add as members
func (ch *Channel) AddMembers(users []string) error {
	data := map[string]interface{}{
		"add_members": users,
	}

	return ch.client.makeRequest(http.MethodPost, ch.formatPath(channelPathFmt), nil, data, nil)
}

//  RemoveMembers deletes members from the channel
//
//  users: user IDs to remove from the member list
func (ch *Channel) RemoveMembers(users []string) error {
	data := map[string]interface{}{
		"remove_members": users,
	}

	return ch.client.makeRequest(http.MethodPost, ch.formatPath(channelPathFmt), nil, data, nil)
}

// AddModerators adds moderators to the channel
//
// users: user IDs to add as moderators
func (ch *Channel) AddModerators(users []string) error {
	data := map[string]interface{}{
		"add_moderators": users,
	}

	return ch.client.makeRequest(http.MethodPost, ch.formatPath(channelPathFmt), nil, data, nil)
}

// DemoteModerators moderators from the channel
//
// users: user IDs to demote
func (ch *Channel) DemoteModerators(users []string) error {
	data := map[string]interface{}{
		"demote_moderators": users,
	}

	return ch.client.makeRequest(http.MethodPost, ch.formatPath(channelPathFmt), nil, data, nil)
}

//  MarkRead send the mark read event for this user, only works if the `read_events` setting is enabled
//
//  userID: the user ID for the event
//  options: additional data, ie {"messageID": last_messageID}
func (ch *Channel) MarkRead(userID string, options map[string]interface{}) error {
	path := ch.formatPath(channelPathFmt) + "/read"

	options = addUserID(options, userID)
	return ch.client.makeRequest(http.MethodPost, path, nil, options, nil)
}

// GetReplies returns list of the message replies for a parent message
//
// parenID: The message parent id, ie the top of the thread
// options: Pagination params, ie {limit:10, idlte: 10}
func (ch *Channel) GetReplies(parentID string, options map[string][]string) (resp map[string]interface{}, err error) {
	path := "messages/" + parentID + "/replies"

	err = ch.client.makeRequest(http.MethodGet, path, options, nil, &resp)

	return resp, err
}

// GetReactions returns list of the reactions, supports pagination
//
// messageID: The message id
// options: Pagination params, ie {"limit":10, "idlte": 10}
func (ch *Channel) GetReactions(messageID string, options map[string][]string) (resp []Reaction, err error) {
	path := "messages/" + messageID + "/reactions"

	err = ch.client.makeRequest(http.MethodGet, path, options, nil, &resp)

	return resp, err
}

// BanUser bans a user from this channel
//
// targetID: the ID of the user to ban
// options: additional ban options, ie {"timeout": 3600, "reason": "offensive language is not allowed here"}
func (ch *Channel) BanUser(targetID string, options map[string]interface{}) error {
	options["type"] = ch._type
	options["id"] = ch.id

	return ch.client.BanUser(targetID, options)
}

// UnBanUser removes the ban for a user on this channel
//
// targetID: the ID of the user to unban
func (ch *Channel) UnBanUser(targetID string, options map[string]string) error {
	options["type"] = ch._type
	options["id"] = ch.id

	return ch.client.UnBanUser(targetID, options)
}

// NewChannel returns new channel struct
func (c *Client) NewChannel(chanType string, chanID string, data map[string]interface{}) *Channel {
	return &Channel{
		_type:      chanType,
		client:     c,
		customData: data,
	}
}
