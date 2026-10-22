Scheduled messages let a user write a message now and have it sent to a channel later. Until it is sent, a scheduled message is visible only to the user who scheduled it and to server-side callers. Server-side, you can schedule, update, cancel, send and list the scheduled messages of any user.

Scheduled messages must be enabled for your app. Requests return an error with the `scheduled_messages_disabled` reason otherwise.

## Limits

- The send time must be at least 1 minute and at most the app's maximum lead time (120 days by default) from now
- A user can have at most 100 scheduled messages, and at most 10 in one channel (app defaults)
- Slash commands cannot be scheduled

## Scheduling a Message

Pass `WithScheduledAt` to `SendMessage`. The message is stored instead of sent, and the response carries the scheduled message.

```go
resp, err := channel.SendMessage(ctx, &Message{Text: "See you tomorrow"}, "user-id",
	WithScheduledAt(time.Now().Add(24*time.Hour)))

scheduled := resp.ScheduledMessage
fmt.Println(scheduled.ID, scheduled.Status, scheduled.ScheduledAt)
```

`resp.Message` holds the content as it will be sent; it is not in the channel yet, and its ID is the scheduled message ID. The ID of the sent message is only known when it is sent (see the `scheduled_message.sent` event below).

## Updating a Scheduled Message

Change the message, the send time, or both. A new message replaces the scheduled one as a whole.

```go
newTime := time.Now().Add(48 * time.Hour)
resp, err := client.UpdateScheduledMessage(ctx, scheduled.ID, &Message{Text: "See you the day after"}, &newTime)

// Only move the send time
resp, err = client.UpdateScheduledMessage(ctx, scheduled.ID, nil, &newTime)
```

## Sending a Scheduled Message Now

```go
resp, err := client.SendScheduledMessageNow(ctx, scheduled.ID)
```

## Deleting a Scheduled Message

Cancels a scheduled message so it is never sent, or removes one that failed to send.

```go
resp, err := client.DeleteScheduledMessage(ctx, scheduled.ID)
```

Update, send now and delete return an error while the message is being sent. A message that failed to send can only be deleted.

## Querying Scheduled Messages

Lists the pending (`scheduled`, `sending`) and `failed` messages of one user across channels: pending ones first, then failed ones, each by send time.

```go
resp, err := client.QueryScheduledMessages(ctx, &QueryScheduledMessagesOptions{
	UserID: "user-id",
	Filter: &ScheduledMessagesFilter{
		CID:    "messaging:general",
		Status: []ScheduledMessageStatus{ScheduledMessageStatusFailed},
	},
	Limit: 20,
})

// Next page
if resp.Next != "" {
	resp, err = client.QueryScheduledMessages(ctx, &QueryScheduledMessagesOptions{
		UserID: "user-id",
		Next:   resp.Next,
	})
}
```

## Events

| Event | Webhook | Description |
| --- | --- | --- |
| `scheduled_message.created` | No | A message was scheduled |
| `scheduled_message.updated` | No | A scheduled message was updated or sent now |
| `scheduled_message.deleted` | No | A scheduled message was canceled |
| `scheduled_message.sent` | Yes | A scheduled message was sent. `message_id` holds the ID of the sent message |
| `scheduled_message.failed` | Yes | A scheduled message failed to send. `scheduled_message.failure_reason` says why |

The events carry the scheduled message in `Event.ScheduledMessage`.
