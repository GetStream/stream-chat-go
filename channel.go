package stream_chat //nolint:golint

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

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
	ID   string           `json:"id"`
	Type ChannelTypeLabel `json:"type"`
	CID  string           `json:"cid"` // full id in format channel_type:channel_ID

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

type queryResponse struct {
	Channel  *Channel         `json:"channel,omitempty"`
	Messages []*Message       `json:"messages,omitempty"`
	Members  []*ChannelMember `json:"members,omitempty"`
	Read     []*ChannelRead   `json:"read,omitempty"`
}

func (q queryResponse) updateChannel(ch *Channel) {
	if q.Channel != nil {
		// save client pointer but update channel information
		client := ch.client
		*ch = *q.Channel
		ch.client = client
	}

	if q.Members != nil {
		ch.Members = q.Members
	}
	if q.Messages != nil {
		ch.Messages = q.Messages
	}
	if q.Read != nil {
		ch.Read = q.Read
	}
}

// query makes request to channel api and updates channel internal state
func (ch *Channel) query(options, data map[string]interface{}) (err error) {
	payload := map[string]interface{}{
		"state": true,
	}

	for k, v := range options {
		payload[k] = v
	}

	if data == nil {
		data = map[string]interface{}{}
	}

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
func (ch *Channel) Update(data map[string]interface{}, message *Message) error {
	payload := map[string]interface{}{
		"data": data,
	}

	if message != nil {
		payload["message"] = message
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

// AddMembers adds members with given user IDs to the channel
func (ch *Channel) AddMembers(userIDs []string, message *Message) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"add_members": userIDs,
	}

	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// RemoveMembers deletes members with given IDs from the channel
func (ch *Channel) RemoveMembers(userIDs []string, message *Message) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"remove_members": userIDs,
	}

	if message != nil {
		data["message"] = message
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
func (ch *Channel) AddModerators(userIDs ...string) error {
	return ch.addModerators(userIDs, nil)
}

// AddModeratorsWithMessage adds moderators with given IDs to the channel and produce system message
func (ch *Channel) AddModeratorsWithMessage(userIDs []string, msg *Message) error {
	return ch.addModerators(userIDs, msg)
}

// AddModerators adds moderators with given IDs to the channel
func (ch *Channel) addModerators(userIDs []string, msg *Message) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"add_moderators": userIDs,
	}

	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// InviteMembers invites users with given IDs to the channel
func (ch *Channel) InviteMembers(userIDs ...string) error {
	return ch.inviteMembers(userIDs, nil)
}

// InviteMembersWithMessage invites users with given IDs to the channel and produce system message
func (ch *Channel) InviteMembersWithMessage(userIDs []string, msg *Message) error {
	return ch.inviteMembers(userIDs, msg)
}

// InviteMembers invites users with given IDs to the channel
func (ch *Channel) inviteMembers(userIDs []string, msg *Message) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"invites": userIDs,
	}

	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// DemoteModerators moderators with given IDs from the channel
func (ch *Channel) DemoteModerators(userIDs ...string) error {
	return ch.demoteModerators(userIDs, nil)
}

// DemoteModeratorsWithMessage moderators with given IDs from the channel and produce system message
func (ch *Channel) DemoteModeratorsWithMessage(userIDs []string, msg *Message) error {
	return ch.demoteModerators(userIDs, msg)
}

// DemoteModerators moderators with given IDs from the channel
func (ch *Channel) demoteModerators(userIDs []string, msg *Message) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"demote_moderators": userIDs,
	}

	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// MarkRead sends the mark read event for user with given ID, only works if the `read_events` setting is enabled
// options: additional data, ie {"messageID": last_messageID}
func (ch *Channel) MarkRead(userID string, options ...Option) error {
	if userID == "" {
		return ErrorMissingUserID
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "read")
	options = append(options, NewOption(optionKeyID, userID))

	return ch.client.makeRequestWithOptions(http.MethodPost, p, nil, options, nil)
}

// BanUser bans target user ID from this channel
// userID: user who bans target
// options: additional ban options, e.g. BanUser("badUser", "admin",
// OptionTimeout(time.Hour), NewOption("reason", "offensive language is
// not allowed here"))
func (ch *Channel) BanUser(targetID, userID string, options ...Option) error {
	return ch.client.banUser(&banUserInput{
		TargetID: targetID,
		UserID:   userID,

		ChannelType: &ch.Type,
		ChannelID:   &ch.ID,
	}, options...)
}

// UnBanUser removes the ban for target user ID on this channel
func (ch *Channel) UnBanUser(targetID string, options ...Option) error {
	if targetID == "" {
		return errors.New("target ID must be not empty")
	}

	options = append(options, NewOption(optionKeyType, ch.Type))
	options = append(options, NewOption(optionKeyID, ch.ID))

	return ch.client.UnBanUser(targetID, options...)
}

// Query fills channel info without state (messages, members, reads)
func (ch *Channel) Query(data map[string]interface{}) error {
	options := map[string]interface{}{
		"watch":    false,
		"state":    false,
		"presence": false,
	}

	return ch.query(options, data)
}

// Show makes channel visible for userID
func (ch *Channel) Show(userID string) error {
	data := map[string]interface{}{
		"user_id": userID,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "show")

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// Hide makes channel hidden for userID
func (ch *Channel) Hide(userID string) error {
	return ch.hide(userID, false)
}

// HideWithHistoryClear clear marks channel as hidden and remove all messages for user
func (ch *Channel) HideWithHistoryClear(userID string) error {
	return ch.hide(userID, true)
}

func (ch *Channel) hide(userID string, clearHistory bool) error {
	data := map[string]interface{}{
		"user_id":       userID,
		"clear_history": clearHistory,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "hide")

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

// CreateChannel creates new channel of given type and id or returns already created one
func (c *Client) CreateChannel(chanType ChannelTypeLabel, chanID, userID string, data map[string]interface{}) (*Channel, error) {
	_, membersPresent := data["members"]

	switch {
	case chanType == "":
		return nil, ErrorMissingChannelType
	case chanID == "" && !membersPresent:
		return nil, errors.New("either channel ID or members must be provided")
	case userID == "":
		return nil, ErrorMissingUserID
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

	if data == nil {
		data = make(map[string]interface{}, 1)
	}

	data["created_by"] = map[string]string{"id": userID}

	err := ch.query(options, data)

	return ch, err
}

type SendFileRequest struct {
	Reader io.Reader `json:"-"`
	// name of the file would be stored
	FileName string
	// User object; required
	User *User
	// file content type, required for SendImage
	ContentType string
}

// SendFile sends file to the channel. Returns file url or error
func (ch *Channel) SendFile(request SendFileRequest) (string, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "file")

	return ch.client.sendFile(p, request)
}

// SendFile sends image to the channel. Returns file url or error
func (ch *Channel) SendImage(request SendFileRequest) (string, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "image")

	return ch.client.sendFile(p, request)
}

// DeleteFile removes uploaded file
func (ch *Channel) DeleteFile(location string) error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "file")

	options := []Option{optionURL(location)}

	return ch.client.makeRequestWithOptions(http.MethodDelete, p, options, nil, nil)
}

// DeleteImage removes uploaded image
func (ch *Channel) DeleteImage(location string) error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "image")

	options := []Option{optionURL(location)}

	return ch.client.makeRequestWithOptions(http.MethodDelete, p, options, nil, nil)
}

func (ch *Channel) AcceptInvite(userID string, message *Message) error {
	if userID == "" {
		return errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"accept_invite": true,
		"user_id":       userID,
	}

	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

func (ch *Channel) RejectInvite(userID string, message *Message) error {
	if userID == "" {
		return errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"reject_invite": true,
		"user_id":       userID,
	}

	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	return ch.client.makeRequest(http.MethodPost, p, nil, data, nil)
}

func (ch *Channel) refresh() error {
	options := map[string]interface{}{
		"watch":    false,
		"state":    true,
		"presence": false,
	}

	return ch.query(options, nil)
}
