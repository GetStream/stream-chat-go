Disabling a channel is a visibility and access toggle. The channel and all its data remain intact, but client-side read and write operations return a `403 Not Allowed` error. Server-side access is preserved for admin operations like moderation and data export.

Disabled channels still appear in query results by default. This means users see the channel in their list but receive errors when attempting to open it. To hide disabled channels from users, filter them out in your queries:


Re-enabling a channel restores full client-side access with all historical messages intact.

## Disable a Channel

```go
// disable a channel with full update
channel.Update(ctx, map[string]interface{}{"disabled": true}, nil)

// disable a channel with partial update
channel.PartialUpdate(ctx, PartialUpdate{
	Set: map[string]interface{}{
		"disabled": true,
	},
})

// enable a channel with full update
channel.Update(ctx, map[string]interface{}{"disabled": false}, nil)

// enable a channel with partial update
channel.PartialUpdate(ctx, PartialUpdate{
	Set: map[string]interface{}{
		"disabled": false,
	},
})
```

> [!NOTE]
> To prevent new messages while still allowing users to read existing messages, use [freeze the channel](/chat/docs/go-golang/freezing_channels/) instead.
