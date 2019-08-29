package stream_chat

import (
	"errors"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"
)

// global lock storage for channels
var channelsLock sync.Map

type ChannelRead struct {
	User     *User     `json:"user"`
	LastRead time.Time `json:"last_read"`
}

type ChannelMember struct {
	UserID      string `json:"user_id,omitempty"`
	User        *User  `json:"user,omitempty"`
	IsModerator bool   `json:"is_moderator,omitempty"`

	Invited          bool       `json:"invited,omitempty"`
	InviteAcceptedAt *time.Time `json:"invite_accepted_at,omitempty"`
	InviteRejectedAt *time.Time `json:"invite_rejected_at,omitempty"`
	Role             string     `json:"role,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Channel struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	CID  string `json:"cid"` // full id in format channel_type:channel_ID

	Config ChannelConfig `json:"config"`

	CreatedBy *User `json:"created_by"`
	Frozen    bool  `json:"frozen"`

	MemberCount int              `json:"member_count"`
	Members     []*ChannelMember `json:"members"`

	Messages []*Message     `json:"messages"`
	Read     []*ChannelRead `json:"read"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastMessageAt time.Time `json:"last_message_at"`

	client *Client
}

type channelQueryResponse struct {
	Channel  *Channel         `json:"channel,omitempty"`
	Messages []*Message       `json:"messages,omitempty"`
	Members  []*ChannelMember `json:"members,omitempty"`
	Read     []*ChannelRead   `json:"read,omitempty"`
}

func (ch *Channel) update(q channelQueryResponse) {
	// atomically add or create new lock in channel lock store
	lock, _ := channelsLock.LoadOrStore(ch.CID, &sync.Mutex{})

	lock.(*sync.Mutex).Lock()
	defer lock.(*sync.Mutex).Unlock()

	// update channel
	if q.Channel != nil {
		//save client pointer from being overwritten
		client := ch.client
		*ch = *q.Channel
		ch.client = client
	}

	// update members
	if q.Members != nil {
		ch.Members = q.Members
	}

	// update messages
	if q.Messages != nil {
		ch.Messages = q.Messages
	}

	// update read
	if q.Read != nil {
		ch.Read = q.Read
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

	if ch.CreatedBy != nil {
		data["created_by"] = map[string]interface{}{"id": ch.CreatedBy.ID}
	}

	payload["data"] = data

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "query")

	var resp channelQueryResponse

	err = ch.client.makeRequest(http.MethodPost, p, nil, payload, &resp)
	if err != nil {
		return err
	}

	ch.update(resp)

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

	var resp channelQueryResponse
	if err := ch.client.makeRequest(http.MethodPost, p, nil, payload, &resp); err != nil {
		return err
	}

	ch.update(resp)

	return nil
}

// Delete removes the channel. Messages are permanently removed.
func (ch *Channel) Delete() error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodDelete, p, nil, nil, nil)
}

// Truncate removes all messages from the channel
func (ch *Channel) Truncate() error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "truncate")

	var resp channelQueryResponse

	if err := ch.client.makeRequest(http.MethodPost, p, nil, nil, &resp); err != nil {
		return err
	}

	ch.update(resp)

	return nil
}

// AddMembers adds members with given user IDs to the channel
func (ch *Channel) AddMembers(userIDs ...string) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"add_members": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp channelQueryResponse

	err := ch.client.makeRequest(http.MethodPost, p, nil, data, &resp)
	if err != nil {
		return err
	}

	ch.update(resp)

	return nil
}

//  RemoveMembers deletes members with given IDs from the channel
func (ch *Channel) RemoveMembers(userIDs ...string) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"remove_members": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp channelQueryResponse

	err := ch.client.makeRequest(http.MethodPost, p, nil, data, &resp)
	if err != nil {
		return err
	}

	ch.update(resp)

	return nil
}

// AddModerators adds moderators with given IDs to the channel
func (ch *Channel) AddModerators(userIDs ...string) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"add_moderators": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp channelQueryResponse

	if err := ch.client.makeRequest(http.MethodPost, p, nil, data, &resp); err != nil {
		return err
	}

	ch.update(resp)

	return nil
}

// DemoteModerators moderators with given IDs from the channel
func (ch *Channel) DemoteModerators(userIDs ...string) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"demote_moderators": userIDs,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp channelQueryResponse

	if err := ch.client.makeRequest(http.MethodPost, p, nil, data, &resp); err != nil {
		return err
	}

	ch.update(resp)

	return nil
}

//  MarkRead send the mark read event for user with given ID, only works if the `read_events` setting is enabled
//  options: additional data, ie {"messageID": last_messageID}
func (ch *Channel) MarkRead(userID string, options map[string]interface{}) error {
	switch {
	case userID == "":
		return errors.New("user ID must be not empty")
	case options == nil:
		options = map[string]interface{}{}
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "read")

	options["user"] = map[string]interface{}{"id": userID}

	return ch.client.makeRequest(http.MethodPost, p, nil, options, nil)
}

// BanUser bans target user ID from this channel
// userID: user who bans target
// options: additional ban options, ie {"timeout": 3600, "reason": "offensive language is not allowed here"}
func (ch *Channel) BanUser(targetID string, userID string, options map[string]interface{}) error {
	switch {
	case targetID == "":
		return errors.New("target ID is empty")
	case userID == "":
		return errors.New("user ID is empty")
	case options == nil:
		options = map[string]interface{}{}
	}

	options["type"] = ch.Type
	options["id"] = ch.ID

	return ch.client.BanUser(targetID, userID, options)
}

// UnBanUser removes the ban for target user ID on this channel
func (ch *Channel) UnBanUser(targetID string, options map[string]string) error {
	switch {
	case targetID == "":
		return errors.New("target ID must be not empty")
	case options == nil:
		options = map[string]string{}
	}

	options["type"] = ch.Type
	options["id"] = ch.ID

	return ch.client.UnBanUser(targetID, options)
}

// CreateChannel creates new channel of given type and id or returns already created one
func (c *Client) CreateChannel(chanType string, chanID string, userID string, data map[string]interface{}) (*Channel, error) {
	switch {
	case chanType == "":
		return nil, errors.New("channel type is empty")
	case chanID == "":
		return nil, errors.New("channel ID is empty")
	case userID == "":
		return nil, errors.New("user ID is empty")
	}

	ch := &Channel{
		Type:      chanType,
		ID:        chanID,
		client:    c,
		CreatedBy: &User{ID: userID},
	}

	options := map[string]interface{}{
		"watch":    false,
		"state":    true,
		"presence": false,
	}

	err := ch.query(options, data)

	return ch, err
}

// Reload updates channel data from server
func (ch *Channel) Reload() error {
	options := map[string]interface{}{
		"watch":    false,
		"state":    true,
		"presence": false,
	}

	err := ch.query(options, nil)

	return err
}
