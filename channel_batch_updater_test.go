package stream_chat

import (
	"context"
	"strings"
	"testing"
	"time"

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

	time.Sleep(2 * time.Second)

	for i := 0; i < 120; i++ {
		task, err := c.GetTask(ctx, resp.TaskID)
		if err != nil {
			if i < 10 {
				time.Sleep(time.Second)
				continue
			}
			require.NoError(t, err, "failed to get task status")
		}
		require.Equal(t, resp.TaskID, task.TaskID)
		
		if task.Status == TaskStatusWaiting || task.Status == TaskStatusPending || task.Status == TaskStatusRunning {
			time.Sleep(time.Second)
			continue
		}

		if task.Status == TaskStatusCompleted {
			for j := 0; j < 120; j++ {
				time.Sleep(time.Second)

				err = ch1.refresh(ctx)
				if err != nil {
					continue
				}
				err = ch2.refresh(ctx)
				if err != nil {
					continue
				}

				ch1MemberIDs := make([]string, len(ch1.Members))
				for i, m := range ch1.Members {
					ch1MemberIDs[i] = m.UserID
				}
				allFound := true
				for _, userID := range usersToAdd {
					if !contains(ch1MemberIDs, userID) {
						allFound = false
						break
					}
				}
				if allFound {
					return
				}
			}
			t.Fatal("changes not visible after 2 minutes")
		}
		if task.Status == TaskStatusFailed {
			if len(task.Result) == 0 {
				time.Sleep(2 * time.Second)
				continue
			}
			if desc, ok := task.Result["description"].(string); ok {
				if strings.Contains(strings.ToLower(desc), "rate limit") {
					time.Sleep(2 * time.Second)
					continue
				}
			}
			t.Fatalf("task failed with result: %v", task.Result)
		}

		time.Sleep(time.Second)
	}
	t.Fatal("task did not complete in 2 minutes")
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

	time.Sleep(2 * time.Second)

	for i := 0; i < 120; i++ {
		task, err := c.GetTask(ctx, resp.TaskID)
		if err != nil {
			if i < 10 {
				time.Sleep(time.Second)
				continue
			}
			require.NoError(t, err, "failed to get task status")
		}
		require.Equal(t, resp.TaskID, task.TaskID)
		
		if task.Status == TaskStatusWaiting || task.Status == TaskStatusPending || task.Status == TaskStatusRunning {
			time.Sleep(time.Second)
			continue
		}

		if task.Status == TaskStatusCompleted {
			var ch1MemberIDs []string
			for j := 0; j < 120; j++ {
				time.Sleep(time.Second)

				err = ch1.refresh(ctx)
				if err != nil {
					continue
				}

				ch1MemberIDs = make([]string, len(ch1.Members))
				for i, m := range ch1.Members {
					ch1MemberIDs[i] = m.UserID
				}
				if !contains(ch1MemberIDs, memberToRemove) {
					return
				}
			}
			t.Fatalf("changes not visible after 2 minutes. Channel 1 still has members: %v", ch1MemberIDs)
		}
		if task.Status == TaskStatusFailed {
			if len(task.Result) == 0 {
				time.Sleep(2 * time.Second)
				continue
			}
			if desc, ok := task.Result["description"].(string); ok {
				if strings.Contains(strings.ToLower(desc), "rate limit") {
					time.Sleep(2 * time.Second)
					continue
				}
			}
			t.Fatalf("task failed with result: %v", task.Result)
		}

		time.Sleep(time.Second)
	}
	t.Fatal("task did not complete in 2 minutes")
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

	time.Sleep(2 * time.Second)

	for i := 0; i < 120; i++ {
		task, err := c.GetTask(ctx, resp.TaskID)
		if err != nil {
			if i < 10 {
				time.Sleep(time.Second)
				continue
			}
			require.NoError(t, err, "failed to get task status")
		}
		require.Equal(t, resp.TaskID, task.TaskID)
		
		if task.Status == TaskStatusWaiting || task.Status == TaskStatusPending || task.Status == TaskStatusRunning {
			time.Sleep(time.Second)
			continue
		}

		if task.Status == TaskStatusCompleted {
			for j := 0; j < 120; j++ {
				time.Sleep(time.Second)

				err = ch1.refresh(ctx)
				if err != nil {
					continue
				}

				for _, m := range ch1.Members {
					if m.UserID == membersID[0] {
						if m.ArchivedAt != nil {
							return
						}
						break
					}
				}
			}
			t.Fatal("changes not visible after 2 minutes")
		}
		if task.Status == TaskStatusFailed {
			if len(task.Result) == 0 {
				time.Sleep(2 * time.Second)
				continue
			}
			if desc, ok := task.Result["description"].(string); ok {
				if strings.Contains(strings.ToLower(desc), "rate limit") {
					time.Sleep(2 * time.Second)
					continue
				}
			}
			t.Fatalf("task failed with result: %v", task.Result)
		}

		time.Sleep(time.Second)
	}
	t.Fatal("task did not complete in 2 minutes")
}
