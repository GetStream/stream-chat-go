Freezing a channel prevents users from sending new messages and adding or deleting reactions.

Sending a message to a frozen channel returns an error message. Attempting to add or delete reactions returns a `403 Not Allowed` error.

User roles with the `UseFrozenChannel` permission can still use frozen channels normally. By default, no user role has this permission.

## Freeze a Channel

```go
channel.Update(ctx, map[string]interface{}{"frozen": true}, nil)
```

## Unfreeze a Channel

```go
channel.Update(ctx, map[string]interface{}{"frozen": false}, nil)
```

## Granting the Frozen Channel Permission

Permissions are typically managed in the [Stream Dashboard](https://dashboard.getstream.io/) under your app's **Roles & Permissions** settings. This is the recommended approach for most use cases.

To grant permissions programmatically, update the channel type using a server-side API call. See [user permissions](/chat/docs/go-golang/chat_permission_policies/) for more details.

```go
resp, err := c.GetChannelType(ctx, "messaging")
adminGrants = append(resp.Grants["admin"], "use-frozen-channel")

c.UpdateChannelType(ctx, "messaging", map[string]interface{}{"grants": map[string][]string{"admin": adminGrants}})
```
