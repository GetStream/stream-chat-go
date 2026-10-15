package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_UpdateChannelsBatch(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	t.Run("Return error if options is nil", func(t *testing.T) {
		_, err := c.UpdateChannelsBatch(ctx, nil)
		require.Error(t, err)
	})

	t.Run("Batch update channels with valid options", func(t *testing.T) {
		ch1 := initChannel(t, c)
		ch2 := initChannel(t, c)
		user := randomUser(t, c)

		resp, err := c.UpdateChannelsBatch(ctx, &ChannelsBatchOptions{
			Operation: BatchUpdateOperationAddMembers,
			Filter: ChannelsBatchFilters{
				CIDs: map[string]interface{}{
					"$in": []string{ch1.CID, ch2.CID},
				},
			},
			Members: []ChannelBatchMemberRequest{
				{UserID: user.ID},
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.TaskID)
	})
}

func TestChannelBatchUpdater_AddMembers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	ch1 := initChannel(t, c)
	ch2 := initChannel(t, c)
	usersToAdd := randomUsersID(t, c, 2)

	updater := c.ChannelBatchUpdater()

	members := make([]ChannelBatchMemberRequest, len(usersToAdd))
	for i, userID := range usersToAdd {
		members[i] = ChannelBatchMemberRequest{UserID: userID}
	}

	resp, err := updater.AddMembers(ctx, ChannelsBatchFilters{
		CIDs: map[string]interface{}{
			"$in": []string{ch1.CID, ch2.CID},
		},
	}, members)
	require.NoError(t, err)
	require.NotEmpty(t, resp.TaskID)
}

func TestChannelBatchUpdater_RemoveMembers(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	membersID := randomUsersID(t, c, 2)
	ch1 := initChannel(t, c, membersID...)
	ch2 := initChannel(t, c, membersID...)

	err := ch1.refresh(ctx)
	require.NoError(t, err, "failed to refresh channel 1")
	require.Len(t, ch1.Members, 2, "channel 1 should have 2 members before removal")

	err = ch2.refresh(ctx)
	require.NoError(t, err, "failed to refresh channel 2")
	require.Len(t, ch2.Members, 2, "channel 2 should have 2 members before removal")

	ch1MemberIDs := make([]string, len(ch1.Members))
	for i, m := range ch1.Members {
		ch1MemberIDs[i] = m.UserID
	}
	ch2MemberIDs := make([]string, len(ch2.Members))
	for i, m := range ch2.Members {
		ch2MemberIDs[i] = m.UserID
	}
	require.ElementsMatch(t, membersID, ch1MemberIDs, "channel 1 should have the expected members")
	require.ElementsMatch(t, membersID, ch2MemberIDs, "channel 2 should have the expected members")

	updater := c.ChannelBatchUpdater()

	memberToRemove := membersID[0]
	resp, err := updater.RemoveMembers(ctx, ChannelsBatchFilters{
		CIDs: map[string]interface{}{
			"$in": []string{ch1.CID, ch2.CID},
		},
	}, []ChannelBatchMemberRequest{{UserID: memberToRemove}})
	require.NoError(t, err)
	require.NotEmpty(t, resp.TaskID)
}

func TestChannelBatchUpdater_Archive(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	membersID := randomUsersID(t, c, 2)
	ch1 := initChannel(t, c, membersID...)
	ch2 := initChannel(t, c, membersID...)

	updater := c.ChannelBatchUpdater()

	resp, err := updater.Archive(ctx, ChannelsBatchFilters{
		CIDs: map[string]interface{}{
			"$in": []string{ch1.CID, ch2.CID},
		},
	}, []ChannelBatchMemberRequest{{UserID: membersID[0]}})
	require.NoError(t, err)
	require.NotEmpty(t, resp.TaskID)
}
