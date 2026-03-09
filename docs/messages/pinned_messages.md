Pinned messages highlight important content in a channel. Use them for announcements, key information, or temporarily promoted content. Each channel can have multiple pinned messages, with optional expiration times.

## Pinning and Unpinning Messages

Pin an existing message using `pinMessage`, or create a pinned message by setting `pinned: true` when sending.

```go
msg := &Message{Text: "Important announcement", Pinned: true}
messageResp, err := channel.SendMessage(ctx, msg, user.ID)

// Pin message for 120 seconds
exp := time.Now().Add(time.Second * 120)
client.PinMessage(ctx, messageResp.Message.ID, user.ID, &exp)

// Unpin message
client.UnPinMessage(ctx, msg.ID, user.ID)
```

### Pin Parameters

| Name        | Type    | Description                                                            | Default | Optional |
| ----------- | ------- | ---------------------------------------------------------------------- | ------- | -------- |
| pinned      | boolean | Whether the message is pinned                                          | false   | ✓        |
| pinned_at   | string  | Timestamp when the message was pinned                                  | -       | ✓        |
| pin_expires | string  | Timestamp when the pin expires. Null means the message does not expire | null    | ✓        |
| pinned_by   | object  | The user who pinned the message                                        | -       | ✓        |

> [!NOTE]
> Pinning a message requires the `PinMessage` permission. See [Permission Resources](/chat/docs/go-golang/permissions_reference/) and [Default Permissions](/chat/docs/go-golang/chat_permission_policies/) for details.


## Retrieving Pinned Messages

Query a channel to retrieve the 10 most recent pinned messages from `pinned_messages`.

```go
_, err = channel.Query(ctx, map[string]interface{}{"watch": false, "state": true})
msgs := channel.PinnedMessages
```

## Paginating Pinned Messages

Use the dedicated pinned messages endpoint to retrieve all pinned messages with pagination.
