Channel members can pin a channel for themselves. This is a per-user setting that does not affect other members.

Pinned channels function identically to regular channels via the API, but your UI can display them separately. When a channel is pinned, the timestamp is recorded and returned as `pinned_at` in the response.

When querying channels, filter by `pinned: true` to retrieve only pinned channels, or `pinned: false` to exclude them. You can also sort by `pinned_at` to display pinned channels first.

## Pin a Channel

```go
ctx := context.Background()

// Get a channel
client.channel("messaging", "general")

// Pin the channel for user amy.
userID := "amy"
channelMemberResp, err := channel.Pin(ctx, userID)

// Query for channels that are pinned.
resp, err := client.QueryChannels(ctx, &QueryOption{
		UserID: userID,
		Filter: map[string]interface{}{
			"pinned": true,
		},
	})

// Query for channels for specific members and show pinned first.
resp, err = client.QueryChannels(ctx, &QueryOption{
		UserID: userID,
		Filter: map[string]interface{}{
			"members": map[string]interface{}{
        "$in": []string{"amy", "ben"},
      },
		},
		Sort: []*SortOption{
			{Field: "pinned_at", Direction: -1},
		},
})

channelMemberResp, err := channel.Unpin(ctx, userID)
```

## Global Pinning

Channels are pinned for a specific member. If the channel should instead be pinned for all users, this can be stored as custom data in the channel itself. The value cannot collide with existing fields, so use a value such as `globally_pinned: true`.
