Muting a channel prevents it from triggering push notifications, unhiding, or incrementing the unread count for that user.

By default, mutes remain active indefinitely until removed. You can optionally set an expiration time. The list of muted channels and their expiration times is returned when the user connects.

## Mute a Channel

```go
channel.Mute(ctx, "john", nil)

// With expiration
exp := 6 * time.Hour
channel.Mute(ctx, "john", &exp)
```

> [!NOTE]
> Messages added to muted channels do not increase the unread messages count.


### Query Muted Channels

Muted channels can be filtered or excluded by using the `muted` in your query channels filter.

```go
client.QueryChannels(ctx, &QueryOption{
	Filter: map[string]interface{}{
		"muted": true,
	},
})
```

### Remove a Channel Mute

Use the unmute method to restore normal notifications and unread behavior for a channel.

```go
channel.Unmute(ctx, "john")
```
