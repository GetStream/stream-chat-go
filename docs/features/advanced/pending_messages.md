Pending Messages features lets you introduce asynchronous moderation on messages being sent on channel. To use this feature please get in touch with support so that we can enable it for your organisation.

## Sending Pending Messages

Messages can be made pending by default by setting the channel config property `mark_messages_pending` to true.

```go
_, err := client.UpdateChannelType(ctx, "messaging", map[string]interface{}{"mark_messages_pending": true})
```

You can also set the `pending` property on a message to mark it as pending on server side (this will override the channel configuration). **Please note that this is only server-side feature** .

```go
msg := &Message{Text: "test pending message"}
	metadata := map[string]string{"my": "metadata"}
	messageResp, err := Channel.SendMessage(ctx, msg, user.ID, MessagePending, MessagePendingMessageMetadata(metadata))
```

Pending messages will only be visible to the user that sent them. They will not be query-able by other users.

## Callbacks

When a pending message is either sent or deleted, the message and its associated pending message metadata are forwarded to your configured callback endpoint via HTTP(s). You may set up to two pending message hooks per application. Only the first commit to a pending message will succeed; any subsequent commit attempts will return an error, as the message is no longer pending. If multiple hooks specify a `timeout_ms`, the system will use the longest timeout value.

You can configure this callback using the dashboard or server-side SDKs.

### Using the Dashboard

1. Go to the [Stream Dashboard](https://getstream.io/dashboard/)
2. Select your app
3. Navigate to your app's settings until "Webhook & Event Configuration" section
4. Click on "Add Integration"
5. Add and configure pending message hook

![](@chat/_default/_assets/images/pending_message_dashboard.png)

### Using Server-Side SDKs

```go
// Note: Any previously existing hooks not included in event_hooks array will be deleted.
// Get current settings first to preserve your existing configuration.

// STEP 1: Get current app settings to preserve existing hooks
settings, err := client.GetAppSettings(ctx)
if err != nil {
    log.Fatal(err)
}
existingHooks := settings.App.EventHooks
fmt.Printf("Current event hooks: %+v\n", existingHooks)

// STEP 2: Add pending message hook while preserving existing hooks
newPendingMessageHook := EventHook{
    HookType:   PendingMessage,
    Enabled:    true,
    WebhookURL: "https://example.com/pending-messages",
    TimeoutMs:  10000, // how long messages should stay pending before being deleted
    Callback:   &Callback{Mode: CallbackModeREST},
}

// STEP 3: Update with complete array including existing hooks
allHooks := append(existingHooks, newPendingMessageHook)
_, err = client.UpdateAppSettings(ctx, NewAppSettings().SetEventHooks(allHooks))
if err != nil {
    log.Fatal(err)
}
```

See the [Webhooks](/chat/docs/go-golang/webhooks_overview/) documentation for complete details.

### Callback Request

For example, if your callback server url is <https://example.com>, we would send callbacks:

- When pending message is sent

`POST https://example.com/PassOnPendingMessage`

- When a pending message is deleted

`POST https://https://example.com/DeletedPendingMessage`

In both callbacks, the body of the POST request will be of the form:

```json
{
  "message": {
    // the message object
  },
  "metadata": {
    // keys and values that you passed as pending_message_metadata
  },
  "request_info": {
    // request info of the request that sent the pending message. Example:
    /*
    "type": "client",
    "ip": "127.0.0.1",
    "user_agent": "Mozilla/5.0...",
    "sdk": "stream-chat-js",
    "ext": "additional-data"
    */
  }
}
```

## Deleting pending messages

Pending messages can be deleted using the normal delete message endpoint. Users are only able to delete their own pending messages. The messages must be hard deleted. Soft deleting a pending message will return an error.

## Updating pending messages

Pending messages cannot be updated.

## Querying pending messages

A user can retrieve their own pending messages using the following endpoints:

```go
// To retrieve single message
messageResp, err = client.GetMessage(ctx, "message-id")

// To retrieve multiple messages
getMsgResp, err := channel.GetMessages(ctx, []string{"message-1", "message-2"})
```

## Query channels

Each channel that is returned from query channels will also have an array of `pending_messages` . These are pending messages that were sent to this channel, and belong to the user who made the query channels call. This array will contain a maximum of 100 messages and these will be the 100 most recently sent messages.

```go
filter := map[string]interface{}{"type": "messaging"}
resp, _ := c.QueryChannels(ctx, &QueryOption{Filter: filter})

fmt.Println("Pending Messages: ", resp.Channels[0].PendingMessages)
```

## Committing pending messages

Calling the commit message endpoint will promote a pending message into a normal message. This message will then be visible to other users and any events/push notifications associated with the message will be sent.

The commit message endpoint is server-side only.

```go
channel.CommitMessage(ctx, "message-id")
```

If a message has been in the pending state longer than the `timeout_ms` defined for your app, then the pending message will be deleted. The default timeout for a pending message is 3 days.
