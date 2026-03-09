Invites allow you to add users to a channel with a pending state. The invited user receives a notification and can accept or reject the invite.

Unread counts are not incremented for channels with a pending invite.

## Invite Users

```go
channel.InviteMembers(ctx, "thierry")
```

## Accept an Invite

Call `acceptInvite` to accept a pending invite. You can optionally include a `message` parameter to post a system message when the user joins (e.g., "Nick joined this channel!").

```go
channel.AcceptInvite(ctx, "thierry", nil)
```

## Reject an Invite

Call `rejectInvite` to decline a pending invite. Client-side calls use the currently connected user. Server-side calls require a `user_id` parameter.

```go
channel.RejectInvite(ctx, "thierry", nil)
```

## Query Invites by Status

Use `queryChannels` with the `invite` filter to retrieve channels based on invite status. Valid values are `pending`, `accepted`, and `rejected`.

### Query Accepted Invites

```go
client.QueryChannels(ctx, &QueryOption{
	Filter: map[string]interface{}{
		"invite": "accepted",
	},
})
```

### Query Rejected Invites

```go
client.QueryChannels(ctx, &QueryOption{
	Filter: map[string]interface{}{
		"invite": "rejected",
	},
})
```

### Query Pending Invites

```go
client.QueryChannels(ctx, &QueryOption{
	Filter: map[string]interface{}{
		"invite": "pending",
	},
})
```
