package stream_chat

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type ChannelRead struct {
	User                   *User      `json:"user"`
	LastRead               time.Time  `json:"last_read"`
	UnreadMessages         int        `json:"unread_messages"`
	LastDeliveredAt        *time.Time `json:"last_delivered_at,omitempty"`
	LastDeliveredMessageID *string    `json:"last_delivered_message_id,omitempty"`
}

type ChannelMember struct {
	UserID      string `json:"user_id,omitempty"`
	User        *User  `json:"user,omitempty"`
	IsModerator bool   `json:"is_moderator,omitempty"`

	Invited            bool       `json:"invited,omitempty"`
	InviteAcceptedAt   *time.Time `json:"invite_accepted_at,omitempty"`
	InviteRejectedAt   *time.Time `json:"invite_rejected_at,omitempty"`
	Status             string     `json:"status,omitempty"`
	Role               string     `json:"role,omitempty"`
	ChannelRole        string     `json:"channel_role"`
	Banned             bool       `json:"banned"`
	BanExpires         *time.Time `json:"ban_expires,omitempty"`
	ShadowBanned       bool       `json:"shadow_banned"`
	ArchivedAt         *time.Time `json:"archived_at,omitempty"`
	PinnedAt           *time.Time `json:"pinned_at,omitempty"`
	NotificationsMuted bool       `json:"notifications_muted"`

	ExtraData map[string]interface{} `json:"-"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// newChannelMembersFromStrings creates a ChannelMembers from a slice of strings
func newChannelMembersFromStrings(members []string) []*ChannelMember {
	channelMembers := make([]*ChannelMember, len(members))
	for i, m := range members {
		channelMembers[i] = &ChannelMember{
			UserID: m,
		}
	}
	return channelMembers
}

type channelMemberForJSON ChannelMember

// UnmarshalJSON implements json.Unmarshaler.
func (m *ChannelMember) UnmarshalJSON(data []byte) error {
	var m2 channelMemberForJSON
	if err := json.Unmarshal(data, &m2); err != nil {
		return err
	}
	*m = ChannelMember(m2)

	if err := json.Unmarshal(data, &m.ExtraData); err != nil {
		return err
	}

	removeFromMap(m.ExtraData, *m)
	flattenExtraData(m.ExtraData)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (m ChannelMember) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(m.ExtraData, channelMemberForJSON(m))
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

	MessageCount    *int           `json:"message_count"`
	Messages        []*Message     `json:"messages"`
	PinnedMessages  []*Message     `json:"pinned_messages"`
	PendingMessages []*Message     `json:"pending_messages"`
	Read            []*ChannelRead `json:"read"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastMessageAt time.Time `json:"last_message_at"`

	TruncatedBy *User      `json:"truncated_by"`
	TruncatedAt *time.Time `json:"truncated_at"`

	ExtraData map[string]interface{} `json:"-"`

	WatcherCount int     `json:"watcher_count"`
	Watchers     []*User `json:"watchers"`

	PushPreferences *ChannelPushPreferences `json:"push_preferences"`
	Hidden          bool                    `json:"hidden"`

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
	flattenExtraData(ch.ExtraData)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (ch Channel) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(ch.ExtraData, channelForJSON(ch))
}

type QueryResponse struct {
	Channel        *Channel         `json:"channel,omitempty"`
	Messages       []*Message       `json:"messages,omitempty"`
	PinnedMessages []*Message       `json:"pinned_messages,omitempty"`
	Members        []*ChannelMember `json:"members,omitempty"`
	Read           []*ChannelRead   `json:"read,omitempty"`

	Response
}

type ChannelRequest struct {
	CreatedBy               *User                  `json:"created_by,omitempty"`
	Team                    string                 `json:"team,omitempty"`
	AutoTranslationEnabled  bool                   `json:"auto_translation_enabled,omitempty"`
	AutoTranslationLanguage string                 `json:"auto_translation_language,omitempty"`
	Frozen                  *bool                  `json:"frozen,omitempty"`
	Disabled                *bool                  `json:"disabled,omitempty"`
	ChannelMembers          []*ChannelMember       `json:"members,omitempty"`
	Members                 []string               `json:"-"`
	Invites                 []string               `json:"invites,omitempty"`
	ExtraData               map[string]interface{} `json:"-"`
	FilterTags              []string               `json:"filter_tags,omitempty"`
}

type channelRequestForJSON ChannelRequest

// UnmarshalJSON implements json.Unmarshaler.
func (c *ChannelRequest) UnmarshalJSON(data []byte) error {
	var ch2 channelRequestForJSON
	if err := json.Unmarshal(data, &ch2); err != nil {
		return err
	}
	*c = ChannelRequest(ch2)

	if err := json.Unmarshal(data, &c.ExtraData); err != nil {
		return err
	}

	removeFromMap(c.ExtraData, *c)
	flattenExtraData(c.ExtraData)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (c ChannelRequest) MarshalJSON() ([]byte, error) {
	return addToMapAndMarshal(c.ExtraData, channelRequestForJSON(c))
}

type PaginationParamsRequest struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	IDGTE  string `json:"id_gte,omitempty"`
	IDGT   string `json:"id_gt,omitempty"`
	IDLTE  string `json:"id_lte,omitempty"`
	IDLT   string `json:"id_lt,omitempty"`
}

type MessagePaginationParamsRequest struct {
	PaginationParamsRequest
	CreatedAtAfterEq  *time.Time `json:"created_at_after_or_equal,omitempty"`
	CreatedAtAfter    *time.Time `json:"created_at_after,omitempty"`
	CreatedAtBeforeEq *time.Time `json:"created_at_before_or_equal,omitempty"`
	CreatedAtBefore   *time.Time `json:"created_at_before,omitempty"`
	IDAround          string     `json:"id_around,omitempty"`
	CreatedAtAround   *time.Time `json:"created_at_around,omitempty"`
}

type QueryRequest struct {
	Data           *ChannelRequest                 `json:"data,omitempty"`
	Watch          bool                            `json:"watch,omitempty"`
	State          bool                            `json:"state,omitempty"`
	Presence       bool                            `json:"presence,omitempty"`
	Messages       *MessagePaginationParamsRequest `json:"messages,omitempty"`
	Members        *PaginationParamsRequest        `json:"members,omitempty"`
	Watchers       *PaginationParamsRequest        `json:"watchers,omitempty"`
	HideForCreator bool                            `json:"hide_for_creator,omitempty"`
}

func (q QueryResponse) updateChannel(ch *Channel) {
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
	if q.PinnedMessages != nil {
		ch.PinnedMessages = q.PinnedMessages
	}
}

// Query makes request to channel api and updates channel internal state.
func (ch *Channel) Query(ctx context.Context, q *QueryRequest) (*QueryResponse, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "query")

	var resp QueryResponse

	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, &q, &resp)
	if err != nil {
		return nil, err
	}

	resp.updateChannel(ch)
	return &resp, nil
}

// Update edits the channel's custom properties.
//
// properties: the object to update the custom properties of this channel with
// message: optional update message
func (ch *Channel) Update(ctx context.Context, properties map[string]interface{}, message *Message) (*Response, error) {
	payload := map[string]interface{}{
		"data": properties,
	}

	if message != nil {
		payload["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))
	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, payload, &resp)
	return &resp, err
}

// PartialUpdate set and unset specific fields when it is necessary to retain additional custom data fields on the object. AKA a patch style update.
func (ch *Channel) PartialUpdate(ctx context.Context, update PartialUpdate) (*Response, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPatch, p, nil, update, &resp)
	return &resp, err
}

// Delete removes the channel. Messages are permanently removed.
func (ch *Channel) Delete(ctx context.Context) (*Response, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodDelete, p, nil, nil, &resp)
	return &resp, err
}

type truncateOptions struct {
	HardDelete  bool       `json:"hard_delete,omitempty"`
	SkipPush    bool       `json:"skip_push,omitempty"`
	TruncatedAt *time.Time `json:"truncated_at,omitempty"`
	Message     *Message   `json:"message,omitempty"`
	UserID      string     `json:"user_id,omitempty"`
	User        *User      `json:"user,omitempty"`
}

type TruncateOption func(*truncateOptions)

func TruncateWithHardDelete() func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.HardDelete = true
	}
}

func TruncateWithSkipPush() func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.SkipPush = true
	}
}

func TruncateWithMessage(message *Message) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.Message = message
	}
}

func TruncateWithUserID(userID string) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.UserID = userID
	}
}

func TruncateWithUser(user *User) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.User = user
	}
}

func TruncateWithTruncatedAt(truncatedAt *time.Time) func(*truncateOptions) {
	return func(o *truncateOptions) {
		o.TruncatedAt = truncatedAt
	}
}

type TruncateResponse struct {
	Response
	Channel *Channel `json:"channel"`
	Message *Message `json:"message"`
}

// Truncate removes all messages from the channel.
// You can pass in options such as hard_delete, skip_push
// or a custom message.
func (ch *Channel) Truncate(ctx context.Context, options ...TruncateOption) (*Response, error) {
	option := &truncateOptions{}

	for _, fn := range options {
		fn(option)
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "truncate")

	var resp TruncateResponse
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, option, &resp)
	return &resp.Response, err
}

type GetMessagesResponse struct {
	Messages []*Message `json:"messages"`
	Response
}

// GetMessages returns messages for multiple message ids.
func (ch *Channel) GetMessages(ctx context.Context, messageIDs []string) (*GetMessagesResponse, error) {
	params := url.Values{}
	params.Set("ids", strings.Join(messageIDs, ","))
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "messages")

	var resp GetMessagesResponse
	err := ch.client.makeRequest(ctx, http.MethodGet, p, params, nil, &resp)
	return &resp, err
}

type addMembersOptions struct {
	MemberIDs      []string         `json:"-"`
	ChannelMembers []*ChannelMember `json:"add_members"`

	RolesAssignement  []*RoleAssignment `json:"assign_roles"`
	HideHistory       bool              `json:"hide_history"`
	HideHistoryBefore *time.Time        `json:"hide_history_before"`
	Message           *Message          `json:"message,omitempty"`
}

type AddMembersOptions func(*addMembersOptions)

func AddMembersWithMessage(message *Message) func(*addMembersOptions) {
	return func(opt *addMembersOptions) {
		opt.Message = message
	}
}

func AddMembersWithHideHistory() func(*addMembersOptions) {
	return func(opt *addMembersOptions) {
		opt.HideHistory = true
	}
}

func AddMembersWithHideHistoryBefore(before time.Time) func(*addMembersOptions) {
	return func(opt *addMembersOptions) {
		opt.HideHistoryBefore = &before
	}
}

func AddMembersWithRolesAssignment(assignements []*RoleAssignment) func(*addMembersOptions) {
	return func(opt *addMembersOptions) {
		opt.RolesAssignement = assignements
	}
}

// AddMembers adds members with given user IDs to the channel. If you want to add members with ChannelMember objects, use AddChannelMembers instead.
func (ch *Channel) AddMembers(ctx context.Context, userIDs []string, options ...AddMembersOptions) (*Response, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("user IDs are empty")
	}

	opts := &addMembersOptions{
		ChannelMembers: newChannelMembersFromStrings(userIDs),
	}

	for _, fn := range options {
		fn(opts)
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, opts, &resp)
	return &resp, err
}

// AddChannelMembers adds members with given []*ChannelMember to the channel.
func (ch *Channel) AddChannelMembers(ctx context.Context, members []*ChannelMember, options ...AddMembersOptions) (*Response, error) {
	if len(members) == 0 {
		return nil, errors.New("members are empty")
	}

	opts := &addMembersOptions{
		ChannelMembers: members,
	}

	for _, fn := range options {
		fn(opts)
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, opts, &resp)
	return &resp, err
}

// RemoveMembers deletes members with given IDs from the channel.
func (ch *Channel) RemoveMembers(ctx context.Context, userIDs []string, message *Message) (*Response, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"remove_members": userIDs,
	}

	if message != nil {
		data["message"] = message
	}
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp QueryResponse

	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	if err != nil {
		return nil, err
	}

	resp.updateChannel(ch)
	return &resp.Response, nil
}

type RoleAssignment struct {
	// UserID is the ID of the user to assign the role to.
	UserID string `json:"user_id"`

	// ChannelRole is the role to assign to the user.
	ChannelRole string `json:"channel_role"`
}

// AssignRoles assigns roles to members with given IDs.
func (ch *Channel) AssignRole(ctx context.Context, assignments []*RoleAssignment, msg *Message) (*Response, error) {
	if len(assignments) == 0 {
		return nil, errors.New("assignments are empty")
	}
	ids := make([]string, 0, len(assignments))
	for _, a := range assignments {
		ids = append(ids, a.UserID)
	}

	data := map[string]interface{}{"assign_roles": assignments, "add_members": ids}
	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

type QueryMembersResponse struct {
	Members []*ChannelMember `json:"members"`

	Response
}

// QueryMembers queries members of a channel.
func (ch *Channel) QueryMembers(ctx context.Context, q *QueryOption, sorters ...*SortOption) (*QueryMembersResponse, error) {
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

	var resp QueryMembersResponse
	err = ch.client.makeRequest(ctx, http.MethodGet, "members", values, nil, &resp)
	return &resp, err
}

// AddModerators adds moderators with given IDs to the channel.
func (ch *Channel) AddModerators(ctx context.Context, userIDs ...string) (*Response, error) {
	return ch.addModerators(ctx, userIDs, nil)
}

// AddModerators adds moderators with given IDs to the channel and produce system message.
func (ch *Channel) AddModeratorsWithMessage(ctx context.Context, userIDs []string, msg *Message) (*Response, error) {
	return ch.addModerators(ctx, userIDs, msg)
}

// AddModerators adds moderators with given IDs to the channel.
func (ch *Channel) addModerators(ctx context.Context, userIDs []string, msg *Message) (*Response, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"add_moderators": userIDs,
	}

	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// InviteMembers invites users with given IDs to the channel.
func (ch *Channel) InviteMembers(ctx context.Context, userIDs ...string) (*Response, error) {
	return ch.inviteMembers(ctx, userIDs, nil)
}

// InviteMembers invites users with given IDs to the channel and produce system message.
func (ch *Channel) InviteMembersWithMessage(ctx context.Context, userIDs []string, msg *Message) (*Response, error) {
	return ch.inviteMembers(ctx, userIDs, msg)
}

// InviteMembers invites users with given IDs to the channel.
func (ch *Channel) inviteMembers(ctx context.Context, userIDs []string, msg *Message) (*Response, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"invites": userIDs,
	}

	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// DemoteModerators moderators with given IDs from the channel.
func (ch *Channel) DemoteModerators(ctx context.Context, userIDs ...string) (*Response, error) {
	return ch.demoteModerators(ctx, userIDs, nil)
}

// DemoteModerators moderators with given IDs from the channel and produce system message.
func (ch *Channel) DemoteModeratorsWithMessage(ctx context.Context, userIDs []string, msg *Message) (*Response, error) {
	return ch.demoteModerators(ctx, userIDs, msg)
}

// DemoteModerators moderators with given IDs from the channel.
func (ch *Channel) demoteModerators(ctx context.Context, userIDs []string, msg *Message) (*Response, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("user IDs are empty")
	}

	data := map[string]interface{}{
		"demote_moderators": userIDs,
	}

	if msg != nil {
		data["message"] = msg
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

type markReadOption struct {
	MessageID string `json:"message_id"`
	ThreadID  string `json:"thread_id"`

	UserID string `json:"user_id"`
}
type MarkReadOption func(*markReadOption)

func MarkReadUntilMessage(id string) func(*markReadOption) {
	return func(opt *markReadOption) {
		opt.MessageID = id
	}
}

func MarkReadThread(id string) func(*markReadOption) {
	return func(opt *markReadOption) {
		opt.ThreadID = id
	}
}

// MarkRead sends the mark read event for user with given ID,
// only works if the `read_events` setting is enabled.
func (ch *Channel) MarkRead(ctx context.Context, userID string, options ...MarkReadOption) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "read")

	opts := &markReadOption{
		UserID: userID,
	}

	for _, fn := range options {
		fn(opts)
	}

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, opts, &resp)
	return &resp, err
}

type markUnreadOption struct {
	MessageID        string     `json:"message_id,omitempty"`
	ThreadID         string     `json:"thread_id,omitempty"`
	MessageTimestamp *time.Time `json:"message_timestamp,omitempty"`

	UserID string `json:"user_id"`
}

type MarkUnreadOption func(option *markUnreadOption)

// Specify ID of the message from where the channel is marked unread.
func MarkUnreadFromMessage(id string) func(*markUnreadOption) {
	return func(opt *markUnreadOption) {
		opt.MessageID = id
	}
}

func MarkUnreadThread(id string) func(*markUnreadOption) {
	return func(opt *markUnreadOption) {
		opt.ThreadID = id
	}
}

func MarkUnreadFromTimestamp(timestamp time.Time) func(*markUnreadOption) {
	return func(opt *markUnreadOption) {
		opt.MessageTimestamp = &timestamp
	}
}

// MarkUnread message or thread (not both) for specified user.
func (ch *Channel) MarkUnread(ctx context.Context, userID string, options ...MarkUnreadOption) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "unread")

	opts := &markUnreadOption{
		UserID: userID,
	}

	for _, fn := range options {
		fn(opts)
	}

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, opts, &resp)
	return &resp, err
}

// RefreshState makes request to channel api and updates channel internal state.
func (ch *Channel) RefreshState(ctx context.Context) (*QueryResponse, error) {
	q := &QueryRequest{State: true}

	resp, err := ch.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	resp.updateChannel(ch)

	return resp, nil
}

// Show makes channel visible for userID.
func (ch *Channel) Show(ctx context.Context, userID string) (*Response, error) {
	data := map[string]interface{}{
		"user_id": userID,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "show")

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// Hide makes channel hidden for userID.
func (ch *Channel) Hide(ctx context.Context, userID string) (*Response, error) {
	return ch.hide(ctx, userID, false)
}

// HideWithHistoryClear clear marks channel as hidden and remove all messages for user.
func (ch *Channel) HideWithHistoryClear(ctx context.Context, userID string) (*Response, error) {
	return ch.hide(ctx, userID, true)
}

func (ch *Channel) hide(ctx context.Context, userID string, clearHistory bool) (*Response, error) {
	data := map[string]interface{}{
		"user_id":       userID,
		"clear_history": clearHistory,
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "hide")

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

type CreateChannelResponse struct {
	Channel *Channel
	*Response
}

type CreateChannelOptions struct {
	HideForCreator bool
}

type CreateChannelOptionFunc func(*CreateChannelOptions)

func HideForCreator(hideForCreator bool) CreateChannelOptionFunc {
	return func(options *CreateChannelOptions) {
		options.HideForCreator = hideForCreator
	}
}

// CreateChannel creates new channel of given type and id or returns already created one.
func (c *Client) CreateChannel(ctx context.Context, chanType, chanID, userID string, data *ChannelRequest, opts ...CreateChannelOptionFunc) (*CreateChannelResponse, error) {
	switch {
	case chanType == "":
		return nil, errors.New("channel type is empty")
	case chanID == "" && (data == nil || (len(data.Members) == 0 && len(data.ChannelMembers) == 0)):
		return nil, errors.New("either channel ID or members must be provided")
	case userID == "":
		return nil, errors.New("user ID is empty")
	}

	options := CreateChannelOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	ch := &Channel{
		Type:      chanType,
		ID:        chanID,
		client:    c,
		CreatedBy: &User{ID: userID},
	}

	if data == nil {
		data = &ChannelRequest{CreatedBy: &User{ID: userID}}
	} else {
		data.CreatedBy = &User{ID: userID}
		if len(data.ChannelMembers) == 0 {
			data.ChannelMembers = newChannelMembersFromStrings(data.Members)
		}
	}

	q := &QueryRequest{
		Watch:          false,
		State:          true,
		Presence:       false,
		Data:           data,
		HideForCreator: options.HideForCreator,
	}

	resp, err := ch.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	return &CreateChannelResponse{Channel: ch, Response: &resp.Response}, nil
}

// CreateChannelWithMembers creates new channel of given type and id or returns already created one.
func (c *Client) CreateChannelWithMembers(ctx context.Context, chanType, chanID, userID string, memberIDs ...string) (*CreateChannelResponse, error) {
	return c.CreateChannel(ctx, chanType, chanID, userID, &ChannelRequest{Members: memberIDs})
}

type SendFileRequest struct {
	Reader io.Reader `json:"-"`
	// name of the file would be stored
	FileName string
	// User object; required
	User *User
}

// SendFile sends file to the channel. Returns file url or error.
func (ch *Channel) SendFile(ctx context.Context, request SendFileRequest) (*SendFileResponse, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "file")

	return ch.client.sendFile(ctx, p, request)
}

// SendFile sends image to the channel. Returns file url or error.
func (ch *Channel) SendImage(ctx context.Context, request SendFileRequest) (*SendFileResponse, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "image")

	return ch.client.sendFile(ctx, p, request)
}

// DeleteFile removes uploaded file.
func (ch *Channel) DeleteFile(ctx context.Context, location string) (*Response, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "file")

	params := url.Values{}
	params.Set("url", location)

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodDelete, p, params, nil, &resp)
	return &resp, err
}

// DeleteImage removes uploaded image.
func (ch *Channel) DeleteImage(ctx context.Context, location string) (*Response, error) {
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "image")

	params := url.Values{}
	params.Set("url", location)

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodDelete, p, params, nil, &resp)
	return &resp, err
}

// AcceptInvite accepts an invite to the channel.
func (ch *Channel) AcceptInvite(ctx context.Context, userID string, message *Message) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"accept_invite": true,
		"user_id":       userID,
	}

	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// RejectInvite rejects an invite to the channel.
func (ch *Channel) RejectInvite(ctx context.Context, userID string, message *Message) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"reject_invite": true,
		"user_id":       userID,
	}

	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

type ChannelMuteResponse struct {
	ChannelMute ChannelMute `json:"channel_mute"`
	Response
}

// Mute mutes the channel. The user will stop receiving messages from the channel.
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

// Unmute removes a mute from a channel so the user will receive messages again.
func (ch *Channel) Unmute(ctx context.Context, userID string) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	data := map[string]interface{}{
		"user_id":     userID,
		"channel_cid": ch.cid(),
	}

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, "moderation/unmute/channel", nil, data, &resp)
	return &resp, err
}

type ChannelMemberResponse struct {
	ChannelMember ChannelMember `json:"channel_member"`
	Response
}

// Pin pins the channel for the user.
func (ch *Channel) Pin(ctx context.Context, userID string) (*ChannelMemberResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "member", url.PathEscape(userID))

	data := map[string]interface{}{
		"set": map[string]interface{}{
			"pinned": true,
		},
	}

	resp := &ChannelMemberResponse{}
	err := ch.client.makeRequest(ctx, http.MethodPatch, p, nil, data, resp)
	return resp, err
}

// Unpin unpins the channel for the user.
func (ch *Channel) Unpin(ctx context.Context, userID string) (*ChannelMemberResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "member", url.PathEscape(userID))

	data := map[string]interface{}{
		"set": map[string]interface{}{
			"pinned": false,
		},
	}

	resp := &ChannelMemberResponse{}
	err := ch.client.makeRequest(ctx, http.MethodPatch, p, nil, data, resp)
	return resp, err
}

// Archive archives the channel for the user.
func (ch *Channel) Archive(ctx context.Context, userID string) (*ChannelMemberResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "member", url.PathEscape(userID))

	data := map[string]interface{}{
		"set": map[string]interface{}{
			"archived": true,
		},
	}

	resp := &ChannelMemberResponse{}
	err := ch.client.makeRequest(ctx, http.MethodPatch, p, nil, data, resp)
	return resp, err
}

// Unarchive unarchives the channel for the user.
func (ch *Channel) Unarchive(ctx context.Context, userID string) (*ChannelMemberResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "member", url.PathEscape(userID))

	data := map[string]interface{}{
		"set": map[string]interface{}{
			"archived": false,
		},
	}

	resp := &ChannelMemberResponse{}
	err := ch.client.makeRequest(ctx, http.MethodPatch, p, nil, data, resp)
	return resp, err
}

// PartialUpdateMember set and unset specific fields when it is necessary to retain additional custom data fields on the object. AKA a patch style update.
func (ch *Channel) PartialUpdateMember(ctx context.Context, userID string, update PartialUpdate) (*ChannelMemberResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "member", url.PathEscape(userID))

	resp := &ChannelMemberResponse{}
	err := ch.client.makeRequest(ctx, http.MethodPatch, p, nil, update, resp)
	return resp, err
}

// CreateDraft creates or updates a draft message in a channel.
func (ch *Channel) CreateDraft(ctx context.Context, userID string, message *messageRequestMessage) (*CreateDraftResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}
	if message == nil {
		return nil, errors.New("message is required")
	}

	// Set the userID in the message
	message.UserID = userID

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "draft")

	data := map[string]interface{}{
		"message": message,
	}

	var resp CreateDraftResponse
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// DeleteDraft deletes a draft message from a channel.
func (ch *Channel) DeleteDraft(ctx context.Context, userID string, parentID *string) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "draft")

	// Convert to url.Values
	values := url.Values{"user_id": []string{userID}}
	if parentID != nil {
		values.Set("parent_id", *parentID)
	}

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodDelete, p, values, nil, &resp)
	return &resp, err
}

// GetDraft retrieves a draft message from a channel.
func (ch *Channel) GetDraft(ctx context.Context, parentID *string, userID string) (*GetDraftResponse, error) {
	if userID == "" {
		return nil, errors.New("user ID must be not empty")
	}
	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID), "draft")

	// Convert to url.Values
	values := url.Values{"user_id": []string{userID}}
	if parentID != nil {
		values.Set("parent_id", *parentID)
	}

	var resp GetDraftResponse
	err := ch.client.makeRequest(ctx, http.MethodGet, p, values, nil, &resp)
	return &resp, err
}

// DraftMessage represents a draft message.
type DraftMessage struct {
	ID              string                 `json:"id"`
	Text            string                 `json:"text"`
	HTML            *string                `json:"html,omitempty"`
	MML             *string                `json:"mml,omitempty"`
	ParentID        *string                `json:"parent_id,omitempty"`
	ShowInChannel   *bool                  `json:"show_in_channel,omitempty"`
	Attachments     []Attachment           `json:"attachments,omitempty"`
	MentionedUsers  []User                 `json:"mentioned_users,omitempty"`
	Custom          map[string]interface{} `json:"custom,omitempty"`
	QuotedMessageID *string                `json:"quoted_message_id,omitempty"`
	Type            string                 `json:"type,omitempty"`
	Silent          *bool                  `json:"silent,omitempty"`
	PollID          string                 `json:"poll_id,omitempty"`
}

// Draft represents a draft message and its associated channel and optionally parent and quoted message.
type Draft struct {
	ChannelCID    string        `json:"channel_cid"`
	CreatedAt     time.Time     `json:"created_at"`
	Message       *DraftMessage `json:"message"`
	Channel       *Channel      `json:"channel,omitempty"`
	ParentID      string        `json:"parent_id,omitempty"`
	ParentMessage *Message      `json:"parent_message,omitempty"`
	QuotedMessage *Message      `json:"quoted_message,omitempty"`
}

// CreateDraftResponse is the response from CreateDraft.
type CreateDraftResponse struct {
	Draft Draft `json:"draft"`
	Response
}

// GetDraftResponse is the response from GetDraft.
type GetDraftResponse struct {
	Draft Draft `json:"draft"`
	Response
}

// QueryDraftsResponse is the response from QueryDrafts.
type QueryDraftsResponse struct {
	Drafts []Draft `json:"drafts"`
	// Next is to be used as the 'next' parameter to get the next page.
	Next *string `json:"next,omitempty"`
	// Prev is to be used as the 'prev' parameter to get the previous page.
	Prev *string `json:"prev,omitempty"`
	Response
}

// QueryDraftsOptions represents the options for the QueryDrafts request.
type QueryDraftsOptions struct {
	UserID string `json:"user_id"`

	// Filter is the filter to be used to query the drafts based on 'channel_cid', 'parent_id' or 'created_at'..
	Filter map[string]any `json:"filter,omitempty"`

	// Sort is the sort to be used to query the drafts. Can sort by 'created_at'.
	Sort []*SortOption `json:"sort,omitempty"`

	// Limit the number of drafts returned.
	Limit int `json:"limit,omitempty"`
	// Pagination parameter. Pass the 'next' value from a previous response to continue from that point.
	Next string `json:"next,omitempty"`
	// Pagination parameter. Pass the 'prev' value from a previous response to continue from that point.
	Prev string `json:"prev,omitempty"`
}

// QueryDrafts retrieves all drafts for the current user.
func (c *Client) QueryDrafts(ctx context.Context, options *QueryDraftsOptions) (*QueryDraftsResponse, error) {
	if options.UserID == "" {
		return nil, errors.New("user ID must be not empty")
	}

	p := "drafts/query"

	var resp QueryDraftsResponse
	err := c.makeRequest(ctx, http.MethodPost, p, nil, options, &resp)
	return &resp, err
}

// AddFilterTags adds filter tags to the channel.
func (ch *Channel) AddFilterTags(ctx context.Context, tags []string, message *Message) (*Response, error) {
	if len(tags) == 0 {
		return nil, errors.New("tags are empty")
	}

	data := map[string]interface{}{
		"add_filter_tags": tags,
	}
	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp Response
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	return &resp, err
}

// RemoveFilterTags removes filter tags from the channel and refreshes channel state.
func (ch *Channel) RemoveFilterTags(ctx context.Context, tags []string, message *Message) (*Response, error) {
	if len(tags) == 0 {
		return nil, errors.New("tags are empty")
	}

	data := map[string]interface{}{
		"remove_filter_tags": tags,
	}
	if message != nil {
		data["message"] = message
	}

	p := path.Join("channels", url.PathEscape(ch.Type), url.PathEscape(ch.ID))

	var resp QueryResponse
	err := ch.client.makeRequest(ctx, http.MethodPost, p, nil, data, &resp)
	if err != nil {
		return nil, err
	}
	resp.updateChannel(ch)
	return &resp.Response, nil
}

// ChannelBatchOperation represents the type of batch update operation.
type ChannelBatchOperation string

const (
	BatchUpdateOperationAddMembers       ChannelBatchOperation = "addMembers"
	BatchUpdateOperationRemoveMembers    ChannelBatchOperation = "removeMembers"
	BatchUpdateOperationInviteMembers    ChannelBatchOperation = "inviteMembers"
	BatchUpdateOperationAssignRoles      ChannelBatchOperation = "assignRoles"
	BatchUpdateOperationAddModerators    ChannelBatchOperation = "addModerators"
	BatchUpdateOperationDemoteModerators ChannelBatchOperation = "demoteModerators"
	BatchUpdateOperationHide             ChannelBatchOperation = "hide"
	BatchUpdateOperationShow             ChannelBatchOperation = "show"
	BatchUpdateOperationArchive          ChannelBatchOperation = "archive"
	BatchUpdateOperationUnarchive        ChannelBatchOperation = "unarchive"
	BatchUpdateOperationUpdateData       ChannelBatchOperation = "updateData"
	BatchUpdateOperationAddFilterTags    ChannelBatchOperation = "addFilterTags"
	BatchUpdateOperationRemoveFilterTags ChannelBatchOperation = "removeFilterTags"
)

// ChannelDataUpdate represents data that can be updated on channels in batch.
type ChannelDataUpdate struct {
	Frozen                  *bool                  `json:"frozen,omitempty"`
	Disabled                *bool                  `json:"disabled,omitempty"`
	Custom                  map[string]interface{} `json:"custom,omitempty"`
	Team                    string                 `json:"team,omitempty"`
	ConfigOverrides         map[string]interface{} `json:"config_overrides,omitempty"`
	AutoTranslationEnabled  *bool                  `json:"auto_translation_enabled,omitempty"`
	AutoTranslationLanguage string                 `json:"auto_translation_language,omitempty"`
}

// ChannelsBatchFilters represents filters for batch channel updates.
type ChannelsBatchFilters struct {
	CIDs       interface{} `json:"cids,omitempty"`
	Types      interface{} `json:"types,omitempty"`
	FilterTags interface{} `json:"filter_tags,omitempty"`
}

// ChannelsBatchOptions represents options for batch channel updates.
type ChannelsBatchOptions struct {
	Operation        ChannelBatchOperation `json:"operation"`
	Filter           ChannelsBatchFilters  `json:"filter"`
	Members          interface{}           `json:"members,omitempty"`
	Data             *ChannelDataUpdate    `json:"data,omitempty"`
	FilterTagsUpdate []string              `json:"filter_tags_update,omitempty"`
}
