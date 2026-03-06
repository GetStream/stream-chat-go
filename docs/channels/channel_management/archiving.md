Channel members can archive a channel for themselves. This is a per-user setting that does not affect other members.

Archived channels function identically to regular channels via the API, but your UI can display them separately. When a channel is archived, the timestamp is recorded and returned as `archived_at` in the response.

When querying channels, filter by `archived: true` to retrieve only archived channels, or `archived: false` to exclude them.

## Archive a Channel

```go
ctx := context.Background()

// Get a channel
channel := client.Channel("messaging", "general")

// Archive the channel for user amy
userID := "amy"
channelMemberResp, err := channel.Archive(ctx, userID)

// Query for channels that are archived
resp, err := client.QueryChannels(ctx, &QueryOption{
	UserID: userID,
	Filter: map[string]interface{}{
		"archived": true,
	},
})

// Unarchive the channel
channelMemberResp, err = channel.Unarchive(ctx, userID)
```

## Global Archiving

Channels are archived for a specific member. If the channel should instead be archived for all users, this can be stored as custom data in the channel itself. The value cannot collide with existing fields, so use a value such as `globally_archived: true`.
