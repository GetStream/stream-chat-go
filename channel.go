package stream_chat

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"
)

type ChannelMember struct {
	UserID      string `json:"user_id,omitempty"`
	User        *User  `json:"user,omitempty"`
	IsModerator bool   `json:"is_moderator,omitempty"`

	Invited          bool       `json:"invited,omitempty"`
	InviteAcceptedAt *time.Time `json:"invite_accepted_at,omitempty"`
	InviteRejectedAt *time.Time `json:"invite_rejected_at,omitempty"`
	Role             string     `json:"role,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Channel struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	// full id in format channel_type:channel_ID
	CID string `json:"cid"`

	CreatedBy User `json:"created_by"`
	Frozen    bool `json:"frozen"`

	MemberCount int             `json:"member_count"`
	Members     []ChannelMember `json:"members"`

	Messages []Message `json:"messages"`
	Read     []User    `json:"read"`

	Config ChannelConfig `json:"config"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastMessageAt time.Time `json:"last_message_at"`

	client *Client
}

func addUserID(hash map[string]interface{}, userID string) map[string]interface{} {
	hash["user"] = map[string]interface{}{"id": userID}
	return hash
}

type messageResponse struct {
	Message  Message `json:"message"`
	duration string
}

type messageRequest struct {
	Message Message `json:"message"`
}

// SendMessage sends a message to this channel.
// *Message will be updated from response body
func (ch *Channel) SendMessage(message *Message, userID string) error {
	var resp messageResponse

	message.ExtraData["user"] = map[string]interface{}{"id": userID}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "message")

	req := messageRequest{Message: *message}

	err := ch.client.makeRequest(http.MethodPost, p, nil, req, &resp)
	if err != nil {
		return err
	}

	*message = resp.Message

	return nil
}

type eventRequest struct {
	Event Event `json:"event"`
}

// SendEvent sends an event on this channel
//
// event: event data, ie {type: 'message.read'}
// userID: the ID of the user sending the event
func (ch *Channel) SendEvent(event Event, userID string) error {
	if event.User == nil {
		event.User = &User{ID: userID}
	}

	req := eventRequest{Event: event}
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "event")

	return ch.client.makeRequest(http.MethodPost, p, nil, req, nil)
}

type reactionResponse struct {
	Message  Message  `json:"message"`
	Reaction Reaction `json:"reaction"`
}

type reactionRequest struct {
	Reaction *Reaction `json:"reaction"`
}

// SendReaction sends a reaction about a message
//
// message: pointer to the message struct
// reaction: the reaction object, ie {type: 'love'}
// userID: the ID of the user that created the reaction
func (ch *Channel) SendReaction(msg *Message, reaction *Reaction, userID string) error {
	var resp reactionResponse

	reaction.ExtraData["user"] = map[string]interface{}{"id": userID}

	p := path.Join("messages", url.PathEscape(msg.ID), "reaction")

	req := reactionRequest{Reaction: reaction}
	err := ch.client.makeRequest(http.MethodPost, p, nil, req, &resp)

	*msg = resp.Message
	*reaction = resp.Reaction

	return err
}

// DeleteReaction removes a reaction by user and type
//
// message:  pointer to the message from which we remove the reaction. Message will be updated from response body
// reaction_type: the type of reaction that should be removed
// userID: the id of the user
func (ch *Channel) DeleteReaction(message *Message, reactionType string, userID string) error {
	if message.ID == "" {
		return errors.New("message ID must be not empty")
	}
	if reactionType == "" {
		return errors.New("reaction type must be not empty")
	}

	p := path.Join("messages", url.PathEscape(message.ID), "reaction", url.PathEscape(reactionType))

	params := map[string][]string{
		"user_id": {userID},
	}

	var resp reactionResponse

	err := ch.client.makeRequest(http.MethodDelete, p, params, nil, &resp)
	if err != nil {
		return err
	}

	*message = resp.Message

	return nil
}

type queryResponse struct {
	Channel  *Channel  `json:"channel,omitempty"`
	Messages *messages `json:"messages,omitempty"`
	Members  *members  `json:"members,omitempty"`
	Read     *users    `json:"read,omitempty"`
}

type messages []Message

type users []User

type members []ChannelMember

func (q queryResponse) updateChannel(ch *Channel) {
	if q.Channel != nil {
		// save client pointer but update channel information
		client := ch.client
		*ch = *q.Channel
		ch.client = client
	}

	if q.Members != nil {
		ch.Members = *q.Members
	}
	if q.Messages != nil {
		ch.Messages = *q.Messages
	}
	if q.Read != nil {
		ch.Read = *q.Read
	}
}

// query makes request to channel api and updates channel internal state
func (ch *Channel) query(options map[string]interface{}, data map[string]interface{}) (err error) {
	payload := map[string]interface{}{
		"state": true,
	}

	for k, v := range options {
		payload[k] = v
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	data["created_by"] = map[string]interface{}{"id": ch.CreatedBy.ID}

	payload["data"] = data

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "query")

	var resp queryResponse

	err = ch.client.makeRequest(http.MethodPost, p, nil, payload, &resp)
	if err != nil {
		return err
	}

	resp.updateChannel(ch)

	return nil
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

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, payload, nil)
}

// Delete removes the channel. Messages are permanently removed.
func (ch *Channel) Delete() error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodDelete, p, nil, nil, nil)
}

// Truncate removes all messages from the channel
func (ch *Channel) Truncate() error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "truncate")

	return ch.client.makeRequest(http.MethodPost, p, nil, nil, nil)
}

// Adds members to the channel
//
// users: user IDs to add as members
func (ch *Channel) AddMembers(users []string) error {
	data := map[string]interface{}{
		"add_members": users,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

//  RemoveMembers deletes members with given IDs from the channel
func (ch *Channel) RemoveMembers(userIDs []string) error {
	data := map[string]interface{}{
		"remove_members": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp queryResponse

	err := ch.client.makeRequest(http.MethodPost, p, nil, data, &resp)
	if err != nil {
		return err
	}

	resp.updateChannel(ch)

	return nil
}

// AddModerators adds moderators with given IDs to the channel
func (ch *Channel) AddModerators(userIDs []string) error {
	data := map[string]interface{}{
		"add_moderators": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// DemoteModerators moderators with given IDs from the channel
func (ch *Channel) DemoteModerators(userIDs []string) error {
	data := map[string]interface{}{
		"demote_moderators": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

//  MarkRead send the mark read event for this user, only works if the `read_events` setting is enabled
//
//  userID: the user ID for the event
//  options: additional data, ie {"messageID": last_messageID}
func (ch *Channel) MarkRead(userID string, options map[string]interface{}) error {
	if userID == "" {
		return errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "read")

	options = addUserID(options, userID)

	return ch.client.makeRequest(http.MethodPost, p, nil, options, nil)
}

// GetReplies returns list of the message replies for a parent message
//
// parenID: The message parent id, ie the top of the thread
// options: Pagination params, ie {limit:10, idlte: 10}
func (ch *Channel) GetReplies(parentID string, options map[string][]string) (replies []Message, err error) {
	if parentID == "" {
		return nil, errors.New("parent ID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(parentID), "replies")

	err = ch.client.makeRequest(http.MethodGet, p, options, nil, nil)

	return replies, err
}

type reactionsResponse struct {
	Reactions []Reaction `json:"reactions"`
}

// GetReactions returns list of the reactions, supports pagination
//
// messageID: The message id
// options: Pagination params, ie {"limit":10, "idlte": 10}
func (ch *Channel) GetReactions(messageID string, options map[string][]string) ([]Reaction, error) {
	if messageID == "" {
		return nil, errors.New("messageID must be not empty")
	}

	p := path.Join("messages", url.PathEscape(messageID), "reactions")

	var resp reactionsResponse

	err := ch.client.makeRequest(http.MethodGet, p, options, nil, &resp)

	return resp.Reactions, err
}

// BanUser bans target user ID from this channel
// userID: user who bans target
// options: additional ban options, ie {"timeout": 3600, "reason": "offensive language is not allowed here"}
func (ch *Channel) BanUser(targetID string, userID string, options map[string]interface{}) error {
	if targetID == "" {
		return errors.New("targetID must be not empty")
	}
	if options == nil {
		options = map[string]interface{}{}
	}

	options["type"] = ch.Type
	options["id"] = ch.ID

	return ch.client.BanUser(targetID, userID, options)
}

// UnBanUser removes the ban for target user ID on this channel
func (ch *Channel) UnBanUser(targetID string, options map[string]string) error {
	if targetID == "" {
		return errors.New("target ID must be not empty")
	}
	if options == nil {
		options = map[string]string{}
	}

	options["type"] = ch.Type
	options["id"] = ch.ID

	return ch.client.UnBanUser(targetID, options)
}

// CreateChannel creates new channel of given type and id or returns already created one
func (c *Client) CreateChannel(chanType string, chanID string, userID string, data map[string]interface{}) (*Channel, error) {
	ch := &Channel{
		Type:      chanType,
		ID:        chanID,
		client:    c,
		CreatedBy: User{ID: userID},
	}

	options := map[string]interface{}{
		"watch":    false,
		"state":    true,
		"presence": false,
	}

	err := ch.query(options, data)

	return ch, err
}

// todo: cleanup this
func (ch *Channel) refresh() error {
	options := map[string]interface{}{
		"watch":    false,
		"state":    true,
		"presence": false,
	}

	err := ch.query(options, nil)

	return err
}
