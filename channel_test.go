package stream_chat

import (
	"context"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (ch *Channel) refresh(ctx context.Context) error {
	_, err := ch.RefreshState(ctx)
	return err
}

func TestClient_TestQuery(t *testing.T) {
	c := initClient(t)
	membersID := randomUsersID(t, c, 3)
	ch := initChannel(t, c, membersID...)
	ctx := context.Background()
	msg, err := ch.SendMessage(ctx, &Message{Text: "test message", Pinned: true}, ch.CreatedBy.ID)
	require.NoError(t, err)

	q := &QueryRequest{
		State:    true,
		Messages: &MessagePaginationParamsRequest{PaginationParamsRequest: PaginationParamsRequest{Limit: 1, IDLT: msg.Message.ID}},
		Members:  &PaginationParamsRequest{Limit: 3, Offset: 0},
		Watchers: &PaginationParamsRequest{Limit: 3, Offset: 0},
	}
	resp, err := ch.Query(ctx, q)
	require.NoError(t, err)
	require.Len(t, resp.Members, 3)

	for _, read := range resp.Read {
		if ch.CreatedBy.ID == read.User.ID {
			require.Equal(t, 0, read.UnreadMessages)
			continue
		}

		require.Equal(t, 1, read.UnreadMessages)
	}
}

func TestClient_CreateChannel(t *testing.T) {
	c := initClient(t)
	userID := randomUser(t, c).ID
	ctx := context.Background()

	t.Run("get existing channel", func(t *testing.T) {
		membersID := randomUsersID(t, c, 3)
		ch := initChannel(t, c, membersID...)
		resp, err := c.CreateChannel(ctx, ch.Type, ch.ID, userID, nil)
		require.NoError(t, err, "create channel", ch)

		channel := resp.Channel
		assert.Equal(t, c, channel.client, "client link")
		assert.Equal(t, ch.Type, channel.Type, "channel type")
		assert.Equal(t, ch.ID, channel.ID, "channel id")
		assert.Equal(t, ch.MemberCount, channel.MemberCount, "member count")
		assert.Len(t, channel.Members, ch.MemberCount, "members length")
	})

	tests := []struct {
		name        string
		channelType string
		id          string
		userID      string
		data        *ChannelRequest
		options     []CreateChannelOptionFunc
		wantErr     bool
	}{
		{"create channel with ID", "messaging", randomString(12), userID, nil, nil, false},
		{"create channel without ID and members", "messaging", "", userID, nil, nil, true},
		{
			"create channel without ID but with members", "messaging", "", userID,
			&ChannelRequest{Members: randomUsersID(t, c, 2)},
			nil, false,
		},
		{
			"create channel with HideForCreator", "messaging", "", userID,
			&ChannelRequest{Members: []string{userID, randomUsersID(t, c, 1)[0]}},
			[]CreateChannelOptionFunc{HideForCreator(true)},
			false,
		},
		{"create channel with ChannelMembers", "messaging", "", userID,
			&ChannelRequest{
				ChannelMembers: NewChannelMembersFromStrings([]string{userID, randomUsersID(t, c, 1)[0]}),
			}, nil, false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.CreateChannel(ctx, tt.channelType, tt.id, tt.userID, tt.data, tt.options...)
			if tt.wantErr {
				require.Error(t, err, "create channel", tt)
				return
			}
			require.NoError(t, err, "create channel", tt)

			channel := resp.Channel
			assert.Equal(t, tt.channelType, channel.Type, "channel type")
			assert.NotEmpty(t, channel.ID)
			if tt.id != "" {
				assert.Equal(t, tt.id, channel.ID, "channel id")
			}

			assert.Equal(t, tt.userID, channel.CreatedBy.ID, "channel created by")
		})
	}
}

func TestChannel_GetManyMessages(t *testing.T) {
	ctx := context.Background()
	c := initClient(t)
	userA := randomUser(t, c)
	userB := randomUser(t, c)
	ch := initChannel(t, c, userA.ID, userB.ID)

	msg := &Message{Text: "test message"}
	messageResp, err := ch.SendMessage(ctx, msg, userB.ID)
	require.NoError(t, err)

	getMsgResp, err := ch.GetMessages(ctx, []string{messageResp.Message.ID})
	require.NoError(t, err)
	require.Len(t, getMsgResp.Messages, 1)
	require.Equal(t, messageResp.Message.ID, getMsgResp.Messages[0].ID)
}

func TestChannel_AddMembers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	chanID := randomString(12)
	resp, err := c.CreateChannel(ctx, "messaging", chanID, randomUser(t, c).ID, nil)
	require.NoError(t, err, "create channel")
	ch := resp.Channel

	assert.Empty(t, ch.Members, "members are empty")

	user := randomUser(t, c)

	msg := &Message{Text: "some members", User: &User{ID: user.ID}}
	_, err = ch.AddMembers(ctx,
		[]string{user.ID},
		AddMembersWithMessage(msg),
		AddMembersWithHideHistory(),
	)
	require.NoError(t, err, "add members")

	// refresh channel state
	require.NoError(t, ch.refresh(ctx), "refresh channel")
	assert.Equal(t, user.ID, ch.Members[0].User.ID, "members contain user id")
}

func TestChannel_AddChannelMembers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	chanID := randomString(12)
	AddChannelMemberUser := randomUser(t, c)
	resp, err := c.CreateChannel(ctx, "messaging", chanID, AddChannelMemberUser.ID, nil)
	require.NoError(t, err, "create channel")
	ch := resp.Channel

	assert.Empty(t, ch.Members, "members are empty")

	channelModeratorID := randomUser(t, c).ID
	channelAdminID := randomUser(t, c).ID

	tests := []struct {
		name          string
		members       []*ChannelMember
		options       []AddMembersOptions
		expectedCount int
		expectedRoles map[string]string // userID -> expected role
		expectedUsers []string          // expected user IDs
	}{
		{
			name: "Add members with ID only",
			members: []*ChannelMember{
				{UserID: randomUser(t, c).ID},
				{UserID: randomUser(t, c).ID},
			},
			options: []AddMembersOptions{
				AddMembersWithMessage(&Message{Text: "adding members with ID only", User: AddChannelMemberUser}),
				AddMembersWithHideHistory(),
			},
			expectedCount: 2,
			expectedRoles: map[string]string{},
			expectedUsers: []string{},
		},
		{
			name: "Add members with ID and role",
			members: []*ChannelMember{
				{UserID: randomUser(t, c).ID, ChannelRole: "channel_moderator"},
				{UserID: randomUser(t, c).ID, ChannelRole: "channel_member"},
			},
			options: []AddMembersOptions{
				AddMembersWithMessage(&Message{Text: "adding members with roles", User: AddChannelMemberUser}),
			},
			expectedCount: 2,
			expectedRoles: map[string]string{
				channelModeratorID: "channel_moderator",
				channelAdminID:     "channel_member",
			},
			expectedUsers: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store user IDs for verification
			userIDs := make([]string, len(tt.members))
			for i, member := range tt.members {
				userIDs[i] = member.UserID
			}

			// Add members
			_, err := ch.AddChannelMembers(ctx, tt.members, tt.options...)
			require.NoError(t, err, "add channel members")

			// Refresh channel state
			require.NoError(t, ch.refresh(ctx), "refresh channel")

			// Verify member count
			assert.Len(t, ch.Members, tt.expectedCount, "member count should match")

			// Verify each member
			for i, userID := range userIDs {
				found := false
				for _, member := range ch.Members {
					if member.User.ID == userID {
						found = true

						// Check role if expected
						if tt.members[i].ChannelRole != "" {
							assert.Equal(t, tt.members[i].ChannelRole, member.ChannelRole,
								"user %s should have role %s", userID, tt.members[i].ChannelRole)
						}
						break
					}
				}
				assert.True(t, found, "user %s should be found in members", userID)
			}
		})
	}
}

func TestChannel_AssignRoles(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	owner := randomUser(t, c)
	other := randomUser(t, c)
	chanID := randomString(12)

	resp, err := c.CreateChannel(ctx, "messaging", chanID, owner.ID, nil)
	require.NoError(t, err, "create channel")
	ch := resp.Channel

	a := []*RoleAssignment{{ChannelRole: "channel_moderator", UserID: other.ID}}
	_, err = ch.AssignRole(ctx, a, nil)
	require.NoError(t, err)
}

func TestChannel_QueryMembers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	chanID := randomString(12)

	resp1, err := c.CreateChannel(ctx, "messaging", chanID, randomUser(t, c).ID, nil)
	require.NoError(t, err, "create channel")
	ch := resp1.Channel

	assert.Empty(t, ch.Members, "members are empty")

	prefix := randomString(12)
	names := []string{"paul", "george", "john", "jessica", "john2"}

	for _, name := range names {
		id := prefix + name
		_, err := c.UpsertUser(ctx, &User{ID: id, Name: id})
		require.NoError(t, err)
		_, err = ch.AddMembers(ctx, []string{id})
		require.NoError(t, err)
	}

	resp2, err := ch.QueryMembers(ctx, &QueryOption{
		Filter: map[string]interface{}{
			"name": map[string]interface{}{"$autocomplete": prefix + "j"},
		},
		Offset: 1,
		Limit:  10,
	}, &SortOption{Field: "created_at", Direction: 1})

	members := resp2.Members
	require.NoError(t, err)
	require.Len(t, members, 2)
	require.Equal(t, prefix+"jessica", members[0].User.ID)
	require.Equal(t, prefix+"john2", members[1].User.ID)
}

// See https://getstream.io/chat/docs/channel_members/ for more details.
func ExampleChannel_AddModerators() {
	channel := &Channel{}
	newModerators := []string{"bob", "sue"}
	ctx := context.Background()

	_, _ = channel.AddModerators(ctx, "thierry", "josh")
	_, _ = channel.AddModerators(ctx, newModerators...)
	_, _ = channel.DemoteModerators(ctx, newModerators...)
}

func TestChannel_InviteMembers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	chanID := randomString(12)

	resp, err := c.CreateChannel(ctx, "messaging", chanID, randomUser(t, c).ID, nil)
	require.NoError(t, err, "create channel")
	ch := resp.Channel

	assert.Empty(t, ch.Members, "members are empty")

	user := randomUser(t, c)

	_, err = ch.InviteMembers(ctx, user.ID)
	require.NoError(t, err, "invite members")

	// refresh channel state
	require.NoError(t, ch.refresh(ctx), "refresh channel")

	assert.Equal(t, user.ID, ch.Members[0].User.ID, "members contain user id")
	assert.True(t, ch.Members[0].Invited, "member is invited")
	assert.Nil(t, ch.Members[0].InviteAcceptedAt, "invite is not accepted")
	assert.Nil(t, ch.Members[0].InviteRejectedAt, "invite is not rejected")
}

func TestChannel_Moderation(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// init random channel
	chanID := randomString(12)
	resp, err := c.CreateChannel(ctx, "messaging", chanID, randomUser(t, c).ID, nil)
	require.NoError(t, err, "create channel")
	ch := resp.Channel

	assert.Empty(t, ch.Members, "members are empty")

	user := randomUser(t, c)

	_, err = ch.AddModeratorsWithMessage(ctx,
		[]string{user.ID},
		&Message{Text: "accepted", User: &User{ID: user.ID}},
	)

	require.NoError(t, err, "add moderators")

	// refresh channel state
	require.NoError(t, ch.refresh(ctx), "refresh channel")

	assert.Equal(t, user.ID, ch.Members[0].User.ID, "user exists")
	assert.Equal(t, "moderator", ch.Members[0].Role, "user role is moderator")

	_, err = ch.DemoteModerators(ctx, user.ID)
	require.NoError(t, err, "demote moderators")

	// refresh channel state
	require.NoError(t, ch.refresh(ctx), "refresh channel")

	assert.Equal(t, user.ID, ch.Members[0].User.ID, "user exists")
	assert.Equal(t, "member", ch.Members[0].Role, "user role is member")
}

func TestChannel_Delete(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	_, err := ch.Delete(ctx)
	require.NoError(t, err, "delete channel")
}

func TestChannel_GetReplies(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	msg := &Message{Text: "test message"}

	resp, err := ch.SendMessage(ctx, msg, randomUser(t, c).ID, MessageSkipPush)
	require.NoError(t, err, "send message")

	msg = resp.Message

	reply := &Message{Text: "test reply", ParentID: msg.ID}
	resp, err = ch.SendMessage(ctx, reply, randomUser(t, c).ID)
	require.NoError(t, err, "send reply")
	require.Equal(t, MessageTypeReply, resp.Message.Type, "message type is reply")

	repliesResp, err := ch.GetReplies(ctx, msg.ID, nil)
	require.NoError(t, err, "get replies")
	assert.Len(t, repliesResp.Messages, 1)
}

func TestChannel_RemoveMembers(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	user := randomUser(t, c)
	_, err := ch.RemoveMembers(ctx,
		[]string{user.ID},
		&Message{Text: "some members", User: &User{ID: user.ID}},
	)

	require.NoError(t, err, "remove members")

	for _, member := range ch.Members {
		assert.NotEqual(t, member.User.ID, user.ID, "member is not present")
	}
}

func TestChannel_SendEvent(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	u := randomUser(t, c)
	ctx := context.Background()
	t.Cleanup(func() {
		_, _ = ch.Delete(ctx)
		_, _ = c.DeleteUser(ctx, u.ID)
	})

	_, err := ch.SendEvent(ctx, &Event{
		Type: "typing.start",
	}, u.ID)
	require.NoError(t, err)
}

func TestChannel_SendMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	user1 := randomUser(t, c)
	user2 := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		User: user1,
	}

	resp, err := ch.SendMessage(ctx, msg, user2.ID)
	require.NoError(t, err, "send message")

	// check that message was updated
	msg = resp.Message
	assert.NotEmpty(t, msg.ID, "message has ID")
	assert.NotEmpty(t, msg.HTML, "message has HTML body")

	msg2 := &Message{
		Text:   "text message 2",
		User:   user1,
		Silent: true,
	}
	resp, err = ch.SendMessage(ctx, msg2, user2.ID)
	require.NoError(t, err, "send message 2")

	// check that message was updated
	msg2 = resp.Message
	assert.NotEmpty(t, msg2.ID, "message has ID")
	assert.NotEmpty(t, msg2.HTML, "message has HTML body")
	assert.True(t, msg2.Silent, "message silent flag is set")
}

func TestChannel_SendSystemMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	user := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		Type: MessageTypeSystem,
	}

	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err, "send message")

	// check that message was updated
	msg = resp.Message
	assert.NotEmpty(t, msg.ID, "message has ID")
	assert.Equal(t, MessageTypeSystem, msg.Type, "message type is system")
}

func TestChannel_SendRestrictedVisibilityMessage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()
	adminUser := randomUserWithRole(t, c, "admin")
	user := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		RestrictedVisibility: []string{
			user.ID,
		},
	}

	resp, err := ch.SendMessage(ctx, msg, adminUser.ID)
	require.NoError(t, err, "send message")
	assert.Equal(t, msg.RestrictedVisibility, resp.Message.RestrictedVisibility)
}

func TestChannel_Truncate(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	user := randomUser(t, c)
	msg := &Message{
		Text: "test message",
		User: user,
	}

	// Make sure we have one message in the channel
	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err, "send message")
	require.NoError(t, ch.refresh(ctx), "refresh channel")
	assert.Equal(t, ch.Messages[0].ID, resp.Message.ID, "message exists")

	// Now truncate it
	_, err = ch.Truncate(ctx)
	require.NoError(t, err, "truncate channel")
	require.NoError(t, ch.refresh(ctx), "refresh channel")
	assert.Empty(t, ch.Messages, "channel is empty")
}

func TestChannel_TruncateWithOptions(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	user := randomUser(t, c)
	truncaterUser := randomUser(t, c)
	ctx := context.Background()
	msg := &Message{
		Text: "test message",
		User: user,
	}

	// Make sure we have one message in the channel
	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err, "send message")
	require.NoError(t, ch.refresh(ctx), "refresh channel")
	assert.Equal(t, ch.Messages[0].ID, resp.Message.ID, "message exists")

	// Now truncate it
	_, err = ch.Truncate(ctx,
		TruncateWithSkipPush(),
		TruncateWithMessage(&Message{Text: "truncated channel", User: &User{ID: user.ID}}),
		TruncateWithUser(&User{ID: truncaterUser.ID}),
	)
	require.NoError(t, err, "truncate channel")
	require.NoError(t, ch.refresh(ctx), "refresh channel")
	require.Len(t, ch.Messages, 1, "channel has one message")
	require.Equal(t, "truncated channel", ch.Messages[0].Text)
	require.NotNil(t, ch.TruncatedBy)
	require.Equal(t, truncaterUser.ID, ch.TruncatedBy.ID)
	require.NotNil(t, ch.TruncatedAt)
}

func TestChannel_Update(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	_, err := ch.Update(ctx, map[string]interface{}{"color": "blue"},
		&Message{Text: "color is blue", User: &User{ID: randomUser(t, c).ID}})
	require.NoError(t, err)
}

func TestChannel_PartialUpdate(t *testing.T) {
	c := initClient(t)
	users := randomUsers(t, c, 5)
	ctx := context.Background()

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}

	req := &ChannelRequest{Members: members, ExtraData: map[string]interface{}{"color": "blue", "age": 30}}
	resp, err := c.CreateChannel(ctx, "team", randomString(12), randomUser(t, c).ID, req)
	require.NoError(t, err)

	ch := resp.Channel
	_, err = ch.PartialUpdate(ctx, PartialUpdate{
		Set: map[string]interface{}{
			"color": "red",
			"config_override": map[string]interface{}{
				"typing_events": false,
			},
		},
		Unset: []string{"age"},
	})
	require.NoError(t, err)
	err = ch.refresh(ctx)
	require.NoError(t, err)
	require.Equal(t, "red", ch.ExtraData["color"])
	require.Nil(t, ch.ExtraData["age"])
}

func TestChannel_MemberPartialUpdate(t *testing.T) {
	c := initClient(t)
	users := randomUsers(t, c, 5)
	ctx := context.Background()

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}

	req := &ChannelRequest{Members: members}
	resp, err := c.CreateChannel(ctx, "team", randomString(12), randomUser(t, c).ID, req)
	require.NoError(t, err)

	ch := resp.Channel
	member, err := ch.PartialUpdateMember(ctx, members[0], PartialUpdate{
		Set: map[string]interface{}{
			"color": "red",
		},
		Unset: []string{"age"},
	})
	require.NoError(t, err)
	require.Equal(t, "red", member.ChannelMember.ExtraData["color"])

	member, err = ch.PartialUpdateMember(ctx, members[0], PartialUpdate{
		Set: map[string]interface{}{
			"age": "18",
		},
		Unset: []string{"color"},
	})
	require.NoError(t, err)
	require.Equal(t, "18", member.ChannelMember.ExtraData["age"])
	require.Nil(t, member.ChannelMember.ExtraData["color"])
}

func TestChannel_SendFile(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	var url string

	t.Run("Send file", func(t *testing.T) {
		file, err := os.Open(path.Join("testdata", "helloworld.txt"))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := ch.SendFile(ctx, SendFileRequest{
			Reader:   file,
			FileName: "HelloWorld.txt",
			User:     randomUser(t, c),
		})
		url = resp.File
		if err != nil {
			t.Fatalf("send file failed: %s", err)
		}
		if url == "" {
			t.Fatal("upload file returned empty url")
		}
	})

	t.Run("Delete file", func(t *testing.T) {
		_, err := ch.DeleteFile(ctx, url)
		if err != nil {
			t.Fatalf("delete file failed: %s", err.Error())
		}
	})
}

func TestChannel_SendImage(t *testing.T) {
	c := initClient(t)
	ch := initChannel(t, c)
	ctx := context.Background()

	var url string

	t.Run("Send image", func(t *testing.T) {
		file, err := os.Open(path.Join("testdata", "helloworld.jpg"))
		if err != nil {
			t.Fatal(err)
		}

		resp, err := ch.SendImage(ctx, SendFileRequest{
			Reader:   file,
			FileName: "HelloWorld.jpg",
			User:     randomUser(t, c),
		})
		if err != nil {
			t.Fatalf("Send image failed: %s", err.Error())
		}

		url = resp.File
		if url == "" {
			t.Fatal("upload image returned empty url")
		}
	})

	t.Run("Delete image", func(t *testing.T) {
		_, err := ch.DeleteImage(ctx, url)
		if err != nil {
			t.Fatalf("delete image failed: %s", err.Error())
		}
	})
}

func TestChannel_AcceptInvite(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	users := randomUsers(t, c, 5)

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}

	req := &ChannelRequest{Members: members, Invites: []string{members[0]}}
	resp, err := c.CreateChannel(ctx, "team", randomString(12), randomUser(t, c).ID, req)

	require.NoError(t, err, "create channel")
	_, err = resp.Channel.AcceptInvite(ctx, members[0], &Message{Text: "accepted", User: &User{ID: members[0]}})
	require.NoError(t, err, "accept invite")
}

func TestChannel_RejectInvite(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	users := randomUsers(t, c, 5)

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}

	req := &ChannelRequest{Members: members, Invites: []string{members[0]}}
	resp, err := c.CreateChannel(ctx, "team", randomString(12), randomUser(t, c).ID, req)

	require.NoError(t, err, "create channel")
	_, err = resp.Channel.RejectInvite(ctx, members[0], &Message{Text: "rejected", User: &User{ID: members[0]}})
	require.NoError(t, err, "reject invite")
}

func TestChannel_Mute_Unmute(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	users := randomUsers(t, c, 5)

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}
	ch := initChannel(t, c, members...)

	// mute the channel
	mute, err := ch.Mute(ctx, members[0], nil)
	require.NoError(t, err, "mute channel")

	require.Equal(t, ch.CID, mute.ChannelMute.Channel.CID)
	require.Equal(t, members[0], mute.ChannelMute.User.ID)
	// query for muted the channel
	queryChannResp, err := c.QueryChannels(ctx, &QueryOption{
		UserID: members[0],
		Filter: map[string]interface{}{
			"muted": true,
			"cid":   ch.CID,
		},
	})

	channels := queryChannResp.Channels
	require.NoError(t, err, "query muted channel")
	require.Len(t, channels, 1)
	require.Equal(t, channels[0].CID, ch.CID)

	// unmute the channel
	_, err = ch.Unmute(ctx, members[0])
	require.NoError(t, err, "mute channel")

	// query for unmuted the channel should return 1 results
	queryChannResp, err = c.QueryChannels(ctx, &QueryOption{
		UserID: members[0],
		Filter: map[string]interface{}{
			"muted": false,
			"cid":   ch.CID,
		},
	})

	require.NoError(t, err, "query muted channel")
	require.Len(t, queryChannResp.Channels, 1)
}

func TestChannel_Pin(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	users := randomUsers(t, c, 5)

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}
	ch := initChannel(t, c, members...)

	//pin the channel
	now := time.Now()
	member, err := ch.Pin(ctx, users[0].ID)
	require.NoError(t, err, "pin channel")
	require.NotNil(t, member.ChannelMember.PinnedAt)
	require.GreaterOrEqual(t, member.ChannelMember.PinnedAt.Unix(), now.Unix())

	// query for pinned the channel
	queryChannResp, err := c.QueryChannels(ctx, &QueryOption{
		UserID: users[0].ID,
		Filter: map[string]interface{}{
			"pinned": true,
			"cid":    ch.CID,
		},
	})

	channels := queryChannResp.Channels
	require.NoError(t, err, "query pinned channel")
	require.Len(t, channels, 1)
	require.Equal(t, channels[0].CID, ch.CID)

	member, err = ch.Unpin(ctx, users[0].ID)
	require.NoError(t, err, "unpin channel")
	require.Nil(t, member.ChannelMember.PinnedAt)

	// query for pinned the channel
	queryChannResp, err = c.QueryChannels(ctx, &QueryOption{
		UserID: users[0].ID,
		Filter: map[string]interface{}{
			"pinned": false,
			"cid":    ch.CID,
		},
	})

	channels = queryChannResp.Channels
	require.NoError(t, err, "query pinned channel")
	require.Len(t, channels, 1)
	require.Equal(t, channels[0].CID, ch.CID)
}

func TestChannel_Archive(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	users := randomUsers(t, c, 5)

	members := make([]string, 0, len(users))
	for i := range users {
		members = append(members, users[i].ID)
	}
	ch := initChannel(t, c, members...)

	//archive the channel
	now := time.Now()
	member, err := ch.Archive(ctx, users[0].ID)
	require.NoError(t, err, "archive channel")
	require.NotNil(t, member.ChannelMember.ArchivedAt)
	require.GreaterOrEqual(t, member.ChannelMember.ArchivedAt.Unix(), now.Unix())

	// query for pinned the channel
	queryChannResp, err := c.QueryChannels(ctx, &QueryOption{
		UserID: users[0].ID,
		Filter: map[string]interface{}{
			"archived": true,
			"cid":      ch.CID,
		},
	})

	channels := queryChannResp.Channels
	require.NoError(t, err, "query archived channel")
	require.Len(t, channels, 1)
	require.Equal(t, channels[0].CID, ch.CID)

	member, err = ch.Unarchive(ctx, users[0].ID)
	require.NoError(t, err, "unarchive channel")
	require.Nil(t, member.ChannelMember.ArchivedAt)

	// query for the archived channel
	queryChannResp, err = c.QueryChannels(ctx, &QueryOption{
		UserID: users[0].ID,
		Filter: map[string]interface{}{
			"archived": false,
			"cid":      ch.CID,
		},
	})

	channels = queryChannResp.Channels
	require.NoError(t, err, "query archived channel")
	require.Len(t, channels, 1)
	require.Equal(t, channels[0].CID, ch.CID)
}

func ExampleChannel_Update() {
	client := &Client{}
	ctx := context.Background()

	data := map[string]interface{}{
		"image":      "https://path/to/image",
		"created_by": "elon",
		"roles":      map[string]string{"elon": "admin", "gwynne": "moderator"},
	}

	spacexChannel := client.Channel("team", "spacex")
	if _, err := spacexChannel.Update(ctx, data, nil); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func (c *Client) ExampleClient_CreateChannel() {
	client, _ := NewClient("XXXX", "XXXX")
	ctx := context.Background()

	resp, _ := client.CreateChannel(ctx, "team", "stream", "tommaso", nil)
	_, _ = resp.Channel.SendMessage(ctx, &Message{
		User: &User{ID: "tomosso"},
		Text: "hi there!",
	}, "tomosso")
}

func ExampleChannel_Query() {
	ctx := context.Background()
	channel := &Channel{}
	msg, _ := channel.SendMessage(ctx, &Message{Text: "test message", Pinned: true}, channel.CreatedBy.ID)

	q := &QueryRequest{
		Messages: &MessagePaginationParamsRequest{PaginationParamsRequest: PaginationParamsRequest{Limit: 1, IDLT: msg.Message.ID}},
		Members:  &PaginationParamsRequest{Limit: 1, Offset: 0},
		Watchers: &PaginationParamsRequest{Limit: 1, Offset: 0},
	}
	_, _ = channel.Query(ctx, q)
}
