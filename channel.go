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
	User           *User     `json:"user"`
	LastRead       time.Time `json:"last_read"`
	UnreadMessages int       `json:"unread_messages"`
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
	Members                 []string               `json:"members,omitempty"`
	Invites                 []string               `json:"invites,omitempty"`
	ExtraData               map[string]interface{} `json:"-"`
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
	MemberIDs []string `json:"add_members"`

	RolesAssignement []*RoleAssignment `json:"assign_roles"`
	HideHistory      bool              `json:"hide_history"`
	Message          *Message          `json:"message,omitempty"`
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

func AddMembersWithRolesAssignment(assignements []*RoleAssignment) func(*addMembersOptions) {
	return func(opt *addMembersOptions) {
		opt.RolesAssignement = assignements
	}
}

// AddMembers adds members with given user IDs to the channel.
func (ch *Channel) AddMembers(ctx context.Context, userIDs []string, options ...AddMembersOptions) (*Response, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("user IDs are empty")
	}

	opts := &addMembersOptions{
		MemberIDs: userIDs,
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
	MessageID string `json:"message_id"`
	ThreadID  string `json:"thread_id"`

	UserID string `json:"user_id"`
}

type MarkUnreadOption func(option *markUnreadOption)

// Specify ID of the message from where the channel is marked unread
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
	case chanID == "" && (data == nil || len(data.Members) == 0):
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
