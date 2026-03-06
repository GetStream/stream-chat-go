Channel members are users who have been added to a channel and can participate in conversations. This page covers how to manage channel membership, including adding and removing members, controlling message history visibility, and managing member roles.

## Adding and Removing Members

### Adding Members

Using the `addMembers()` method adds the given users as members to a channel.

```go
err := channel.AddMembers([]string{"thierry"}, nil, nil)
```

> [!NOTE]
> **Note:** You can only add/remove up to 100 members at once.


Members can also be added when creating a channel:


### Removing Members

Using the `removeMembers()` method removes the given users from the channel.

```go
resp, err := ch.RemoveMembers(ctx, []string{"my_user_id"})
```

### Leaving a Channel

Users can leave a channel without moderator-level permissions. Ensure channel members have the `Leave Own Channel` permission enabled.


> [!NOTE]
> You can familiarize yourself with all permissions in the [Permissions section](/chat/docs/go-golang/chat_permission_policies/).


## Hide History

When members join a channel, you can specify whether they have access to the channel's message history. By default, new members can see the history. Set `hide_history` to `true` to hide it for new members.

```go
err := channel.AddMembers([]string{"thierry"}, nil, stream_chat.AddMembersWithHideHistory())
```

### Hide History Before a Specific Date

Alternatively, `hide_history_before` can be used to hide any history before a given timestamp while giving members access to later messages. The value must be a timestamp in the past in RFC 3339 format. If both parameters are defined, `hide_history_before` takes precedence over `hide_history`.

```go
cutoff := time.Now().Add(-7 * 24 * time.Hour)
err := channel.AddMembers([]string{"thierry"}, nil, stream_chat.AddMembersWithHideHistoryBefore(cutoff))
```

## System Message Parameter

You can optionally include a message object when adding or removing members that client-side SDKs will use to display a system message. This works for both adding and removing members.

```go
err := channel.AddMembers([]string{"tommaso"}, &Message{Text: "Tommaso joined the channel.", User: &User{ID: "tommaso"}}, nil)
```

## Adding and Removing Moderators

Using the `addModerators()` method adds the given users as moderators (or updates their role to moderator if already members), while `demoteModerators()` removes the moderator status.

### Add Moderators

```go
newModerators := []string{"bob", "sue"}
err = channel.AddModerators("thierry", "josh")
err = channel.AddModerators(newModerators...)
```

### Remove Moderators

```go
err = channel.DemoteModerators(newModerators...)
```

> [!NOTE]
> These operations can only be performed server-side, and a maximum of 100 moderators can be added or removed at once.


## Member Custom Data

Custom data can be added at the channel member level. This is useful for storing member-specific information that is separate from user-level data. Ensure custom data does not exceed 5KB.

### Adding Custom Data


### Updating Member Data

Channel members can be partially updated. Only custom data and channel roles are eligible for modification. You can set or unset fields, either separately or in the same call.

```go
// Set some fields
member, err := ch.PartialUpdateMember(ctx, members[0], PartialUpdate{
	Set: map[string]interface{}{
		"color": "red",
	},
})

// Unset some fields
member, err := ch.PartialUpdateMember(ctx, members[0], PartialUpdate{
	Unset: []string{"age"},
})

// Set and unset in the same call
member, err := ch.PartialUpdateMember(ctx, members[0], PartialUpdate{
	Set: map[string]interface{}{
		"color": "red",
	},
	Unset: []string{"age"},
})
```
