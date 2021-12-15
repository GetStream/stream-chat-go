package stream_chat //nolint: golint

import (
	"context"
	"encoding/json"
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
	ID   string `json:"id"`
	Type string `json:"type"`
	CID  string `json:"cid"` // full id in format channel_type:channel_ID
	Team string `json:"team"`

	Config ChannelConfig `json:"config"`

	CreatedBy *User `json:"created_by"`
	Disabled  bool  `json:"disabled"`
	Frozen    bool  `json:"frozen"`

	MemberCount int              `json:"member_count"`
	Members     []*ChannelMember `json:"members"`

	Messages []*Message     `json:"messages"`
	Read     []*ChannelRead `json:"read"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastMessageAt time.Time `json:"last_message_at"`

	ExtraData map[string]interface{} `json:"-"`

	client *Client
}

func (ch Channel) cid() string {
	if ch.CID != "" {
		return ch.CID
	}
	return ch.Type + ":" + ch.ID
}

type PartialUpdate struct {
	Set   map[string]interface{} `json:"set"`
	Unset []string               `json:"unset"`
}

type channelForJSON Channel

// UnmarshalJSON implements json.Unmarshaler.
func (ch *Channel) UnmarshalJSON(data []byte) error {
	var ch2 channelForJSON
	if err := json.Unmarshal(data, &ch2); err != nil {
		return err
	}
	*ch = Channel(ch2)

	if err := json.Unmarshal(data, &ch.ExtraData); err != nil {
		return err
	}

	removeFromMap(ch.ExtraData, *ch)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (ch Channel) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(ch.ExtraData, channelForJSON(ch))
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

type ImportChannelMessagesResponse struct {
	Messages []Message `json:"messages"`
}

// query makes request to channel api and updates channel internal state.
func (ch *Channel) query(ctx context.Context, options, data map[string]interface{}) (err error) {
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

	err = ch.client.makeRequest(ctx, http.MethodPost, p, nil, payload, &resp)
	if err != nil {
		return err
	}

	resp.updateChannel(ch)
	return nil
}

// Update edits the channel's custom properties.
//
// options: the object to update the custom properties of this channel with
// message: optional update message
func (ch *Channel) Update(ctx context.Context, options map[string]interface{}, message *Message) error {
	payload := map[string]interface{}{
		"data": options,
	}

	if message != nil {
		payload["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, payload, nil)
}

//  PartialUpdate set and unset specific fields when it is necessary to retain additional custom data fields on the object. AKA a patch style update.
// options: the object to update the custom properties of the channel
func (ch *Channel) PartialUpdate(ctx context.Context, update PartialUpdate) error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))
	return ch.client.makeRequest(ctx, http.MethodPatch, p, nil, update, nil)
}

// Delete removes the channel. Messages are permanently removed.
func (ch *Channel) Delete(ctx context.Context) error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))
	return ch.client.makeRequest(ctx, http.MethodDelete, p, nil, nil, nil)
}

type truncateOptions struct {
	HardDelete  bool       `json:"hard_delete,omitempty"`
	SkipPush    bool       `json:"skip_push,omitempty"`
	TruncatedAt *time.Time `json:"truncated_at,omitempty"`
	Message     *Message   `json:"message,omitempty"`
}

type TruncateOption func(*truncateOptions)

func TruncateWithHardDelete(hardDelete bool) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.HardDelete = hardDelete
	}
}

func TruncateWithSkipPush(skipPush bool) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.SkipPush = skipPush
	}
}

func TruncateWithMessage(message *Message) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.Message = message
	}
}

func TruncateWithTruncatedAt(truncatedAt *time.Time) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.TruncatedAt = truncatedAt
	}
}

// Truncate removes all messages from the channel.
// You can pass in options such as hard_delete, skip_push
// or a custom message.
func (ch *Channel) Truncate(ctx context.Context, options ...TruncateOption) error {
	option := &truncateOptions{}

	for _, fn := range options {
		fn(option)
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "truncate")
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, option, nil)
}

// AddMembers adds members with given user IDs to the channel.
// You can set a message for channel object notifications.
// If you want to hide history of the channel for new members, you can pass "hide_history": true to options parameter.
func (ch *Channel) AddMembers(ctx context.Context, userIDs []string, message *Message, options map[string]interface{}) error {
	if len(userIDs) == 0 {
		return errors.New("user IDs are empty")
	}

	if options == nil {
		options = map[string]interface{}{}
	}

	options["add_members"] = userIDs

	if message != nil {
		options["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, options, nil)
}

// RemoveMembers deletes members with given IDs from the channel.
func (ch *Channel) RemoveMembers(ctx context.Context, userIDs []string, message *Message) error {
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

	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	if err != nil {
		return err
	}

	resp.updateChannel(ch)
	return nil
}

// ImportMessages is a batch endpoint for inserting multiple messages.
func (ch *Channel) ImportMessages(ctx context.Context, messages ...*Message) (*ImportChannelMessagesResponse, error) {
	for _, m := range messages {
		if m.User == nil || m.User.ID == "" {
			return nil, errors.New("message.user is a required field")
		}
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "import")

	var resp ImportChannelMessagesResponse
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, map[string]interface{}{
		"messages": messages,
	}, &resp)
	return &resp, err
}

type queryMembersResponse struct {
	Members []*ChannelMember `json:"members"`
}

// QueryMembers queries members of a channel.
func (ch *Channel) QueryMembers(ctx context.Context, q *QueryOption, sorters ...*SortOption) ([]*ChannelMember, error) {
	qp := map[string]interface{}{
		"id":                ch.ID,
		"type":              ch.Type,
		"filter_conditions": q.Filter,
		"limit":             q.Limit,
		"offset":            q.Offset,
		"sort":              sorters,
	}

	if ch.ID == "" && len(ch.Members) > 0 {
		members := make([]*ChannelMember, 0, len(ch.Members))
		for _, m := range ch.Members {
			if m.User != nil {
				members = append(members, &ChannelMember{UserID: m.User.ID})
			} else {
				members = append(members, &ChannelMember{UserID: m.UserID})
			}
		}
		qp["members"] = members
	}

	data, err := json.Marshal(&qp)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Set("payload", string(data))

	var resp queryMembersResponse
	err = ch.client.makeRequest(ctx, http.MethodGet, "members", values, nil, &resp)
	return resp.Members, err
}

// AddModerators adds moderators with given IDs to the channel.
func (ch *Channel) AddModerators(ctx context.Context, userIDs ...string) error {
	return ch.addModerators(ctx, userIDs, nil)
}

// AddModerators adds moderators with given IDs to the channel and produce system message.
func (ch *Channel) AddModeratorsWithMessage(ctx context.Context, userIDs []string, msg *Message) error {
	return ch.addModerators(ctx, userIDs, msg)
}

// AddModerators adds moderators with given IDs to the channel.
func (ch *Channel) addModerators(ctx context.Context, userIDs []string, msg *Message) error {
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
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

// InviteMembers invites users with given IDs to the channel.
func (ch *Channel) InviteMembers(ctx context.Context, userIDs ...string) error {
	return ch.inviteMembers(ctx, userIDs, nil)
}

// InviteMembers invites users with given IDs to the channel and produce system message.
func (ch *Channel) InviteMembersWithMessage(ctx context.Context, userIDs []string, msg *Message) error {
	return ch.inviteMembers(ctx, userIDs, msg)
}

// InviteMembers invites users with given IDs to the channel.
func (ch *Channel) inviteMembers(ctx context.Context, userIDs []string, msg *Message) error {
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
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

// DemoteModerators moderators with given IDs from the channel.
func (ch *Channel) DemoteModerators(ctx context.Context, userIDs ...string) error {
	return ch.demoteModerators(ctx, userIDs, nil)
}

// DemoteModerators moderators with given IDs from the channel and produce system message.
func (ch *Channel) DemoteModeratorsWithMessage(ctx context.Context, userIDs []string, msg *Message) error {
	return ch.demoteModerators(ctx, userIDs, msg)
}

// DemoteModerators moderators with given IDs from the channel.
func (ch *Channel) demoteModerators(ctx context.Context, userIDs []string, msg *Message) error {
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
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

// MarkRead send the mark read event for user with given ID,
// only works if the `read_events` setting is enabled.
// options: additional data, ie {"messageID": last_messageID}
func (ch *Channel) MarkRead(ctx context.Context, userID string, options map[string]interface{}) error {
	switch {
	case userID == "":
		return errors.New("user ID must be not empty")
	case options == nil:
		options = map[string]interface{}{}
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "read")

	options["user"] = map[string]interface{}{"id": userID}
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, options, nil)
}

// BanUser bans target user ID from this channel
// userID: user who bans target.
// options: additional ban options, ie {"timeout": 3600, "reason": "offensive language is not allowed here"}.
func (ch *Channel) BanUser(ctx context.Context, targetID, userID string, options map[string]interface{}) error {
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

	return ch.client.BanUser(ctx, targetID, userID, options)
}

// UnBanUser removes the ban for target user ID on this channel.
func (ch *Channel) UnBanUser(ctx context.Context, targetID string, options map[string]string) error {
	switch {
	case targetID == "":
		return errors.New("target ID must be not empty")
	case options == nil:
		options = map[string]string{}
	}

	options["type"] = ch.Type
	options["id"] = ch.ID

	return ch.client.UnBanUser(ctx, targetID, options)
}

// ShadowBan shadow bans userID from this channel
// bannedByID: user who shadow bans userID.
// options: additional shadow ban options, ie {"timeout": 3600, "reason": "offensive language is not allowed here"}.
func (ch *Channel) ShadowBan(ctx context.Context, userID, bannedByID string, options map[string]interface{}) error {
	if options == nil {
		options = map[string]interface{}{}
	}

	options["type"] = ch.Type
	options["id"] = ch.ID

	return ch.client.ShadowBan(ctx, userID, bannedByID, options)
}

// RemoveShadowBan removes the shadow ban for target user ID on this channel.
func (ch *Channel) RemoveShadowBan(ctx context.Context, userID string) error {
	options := map[string]string{
		"type": ch.Type,
		"id":   ch.ID,
	}

	return ch.client.RemoveShadowBan(ctx, userID, options)
}

// Query fills channel info with state (messages, members, reads).
func (ch *Channel) Query(ctx context.Context, data map[string]interface{}) error {
	options := map[string]interface{}{
		"state": true,
	}

	return ch.query(ctx, options, data)
}

// Show makes channel visible for userID.
func (ch *Channel) Show(ctx context.Context, userID string) error {
	data := map[string]interface{}{
		"user_id": userID,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "show")
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

// Hide makes channel hidden for userID.
func (ch *Channel) Hide(ctx context.Context, userID string) error {
	return ch.hide(ctx, userID, false)
}

// HideWithHistoryClear clear marks channel as hidden and remove all messages for user.
func (ch *Channel) HideWithHistoryClear(ctx context.Context, userID string) error {
	return ch.hide(ctx, userID, true)
}

func (ch *Channel) hide(ctx context.Context, userID string, clearHistory bool) error {
	data := map[string]interface{}{
		"user_id":       userID,
		"clear_history": clearHistory,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "hide")
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

// CreateChannel creates new channel of given type and id or returns already created one.
func (c *Client) CreateChannel(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*Channel, error) {
	_, membersPresent := data["members"]

	switch {
	case chanType == "":
		return nil, errors.New("channel type is empty")
	case chanID == "" && !membersPresent:
		return nil, errors.New("either channel ID or members must be provided")
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

	if data == nil {
		data = make(map[string]interface{}, 1)
	}

	data["created_by"] = map[string]string{"id": userID}

	if err := ch.query(ctx, options, data); err != nil {
		return nil, err
	}
	return ch, nil
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

// SendFile sends file to the channel. Returns file url or error.
func (ch *Channel) SendFile(ctx context.Context, request SendFileRequest) (string, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "file")

	return ch.client.sendFile(ctx, p, request)
}

// SendFile sends image to the channel. Returns file url or error.
func (ch *Channel) SendImage(ctx context.Context, request SendFileRequest) (string, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "image")

	return ch.client.sendFile(ctx, p, request)
}

// DeleteFile removes uploaded file.
func (ch *Channel) DeleteFile(ctx context.Context, location string) error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "file")

	params := url.Values{}
	params.Set("url", location)
	return ch.client.makeRequest(ctx, http.MethodDelete, p, params, nil, nil)
}

// DeleteImage removes uploaded image.
func (ch *Channel) DeleteImage(ctx context.Context, location string) error {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "image")

	params := url.Values{}
	params.Set("url", location)
	return ch.client.makeRequest(ctx, http.MethodDelete, p, params, nil, nil)
}

func (ch *Channel) AcceptInvite(ctx context.Context, userID string, message *Message) error {
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
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

func (ch *Channel) RejectInvite(ctx context.Context, userID string, message *Message) error {
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
	return ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, nil)
}

func (ch *Channel) Mute(ctx context.Context, userID string, expiration *time.Duration) (*ChannelMuteResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"user_id":     userID,
		"channel_cid": ch.cid(),
	}
	if expiration != nil {
		data["expiration"] = int(expiration.Milliseconds())
	}

	mute := &ChannelMuteResponse{}
	err := ch.client.makeRequest(ctx, http.MethodPost, "moderation/mute/channel", nil, data, mute)
	return mute, err
}

func (ch *Channel) Unmute(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"user_id":     userID,
		"channel_cid": ch.cid(),
	}
	return ch.client.makeRequest(ctx, http.MethodPost, "moderation/unmute/channel", nil, data, nil)
}

func (ch *Channel) refresh(ctx context.Context) error {
	options := map[string]interface{}{
		"watch":    false,
		"state":    true,
		"presence": false,
	}
	return ch.query(ctx, options, nil)
}
