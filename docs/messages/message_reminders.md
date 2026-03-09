Message reminders let users schedule notifications for specific messages, making it easier to follow up later. When a reminder includes a timestamp, it's like saying "remind me later about this message," and the user who set it will receive a notification at the designated time. If no timestamp is provided, the reminder functions more like a bookmark, allowing the user to save the message for later reference.

Reminders require Push V3 to be enabled - see details [here](/chat/docs/go-golang/push_template/)

## Enabling Reminders

The Message Reminders feature must be activated at the channel level before it can be used. You have two configuration options: activate it for a single channel using configuration overrides, or enable it globally for all channels of a particular type.

```go
// Enabling it for a channel
channel.PartialUpdate(ctx, PartialUpdate{
    Set: map[string]interface{}{
        "config_overrides": map[string]interface{}{
            "user_message_reminders": true,
        },
    },
})

// Enabling it for a channel type
client.UpdateChannelType(ctx, "messaging", map[string]interface{}{
    "user_message_reminders": true,
})
```

Message reminders allow users to:

- schedule a notification after given amount of time has elapsed
- bookmark a message without specifying a deadline

## Limits

- A user cannot have more than 250 reminders scheduled
- A user can only have one reminder created per message

## Creating a Message Reminder

You can create a reminder for any message. When creating a reminder, you can specify a reminder time or save it for later without a specific time.

```go
import "time"

// Create a reminder with a specific due date
remindAt := time.Now().Add(time.Hour)
reminder, err := client.CreateReminder(ctx, "message-id", "user-id", &remindAt)

// Create a "Save for later" reminder without a specific time
reminder, err := client.CreateReminder(ctx, "message-id", "user-id", nil)
```

## Updating a Message Reminder

You can update an existing reminder for a message to change the reminder time.

```go
import "time"

// Update a reminder with a new due date
remindAt := time.Now().Add(2 * time.Hour)
updatedReminder, err := client.UpdateReminder(ctx, "message-id", "user-id", &remindAt)

// Convert a timed reminder to "Save for later"
updatedReminder, err := client.UpdateReminder(ctx, "message-id", "user-id", nil)
```

## Deleting a Message Reminder

You can delete a reminder for a message when it's no longer needed.

```go
// Delete the reminder for the message
err := client.DeleteReminder(ctx, "message-id", "user-id")
```

## Querying Message Reminders

The SDK allows you to fetch all reminders of the current user. You can filter, sort, and paginate through all the user's reminders.

```go
// Query reminders for a user
reminders, err := client.QueryReminders(ctx, "user-id", nil)

// Query reminders with filters
filter := map[string]interface{}{
    "channel_cid": "messaging:general",
}
reminders, err := client.QueryReminders(ctx, "user-id", filter)
```

### Filtering Reminders

You can filter the reminders based on different criteria:

- `message_id` - Filter by the message that the reminder is created on.
- `remind_at` - Filter by the reminder time.
- `created_at` - Filter by the creation date.
- `channel_cid` - Filter by the channel ID.

The most common use case would be to filter by the reminder time. Like filtering overdue reminders, upcoming reminders, or reminders with no due date (saved for later).

```go
import "time"

// Filter overdue reminders
overdueFilter := map[string]interface{}{
    "remind_at": map[string]interface{}{"$lt": time.Now()},
}
overdueReminders, err := client.QueryReminders(ctx, "user-id", overdueFilter)

// Filter upcoming reminders
upcomingFilter := map[string]interface{}{
    "remind_at": map[string]interface{}{"$gt": time.Now()},
}
upcomingReminders, err := client.QueryReminders(ctx, "user-id", upcomingFilter)

// Filter reminders with no due date (saved for later)
savedFilter := map[string]interface{}{
    "remind_at": nil,
}
savedReminders, err := client.QueryReminders(ctx, "user-id", savedFilter)
```

### Pagination

If you have many reminders, you can paginate the results.

```go
// Load reminders with pagination
options := map[string]interface{}{
    "limit":  10,
    "offset": 0,
}
reminders, err := client.QueryReminders(ctx, "user-id", nil, options)

// Load next page
nextPageOptions := map[string]interface{}{
    "limit":  10,
    "offset": 10,
}
nextReminders, err := client.QueryReminders(ctx, "user-id", nil, nextPageOptions)
```

## Events

The following WebSocket events are available for message reminders:

- `reminder.created` - Triggered when a reminder is created
- `reminder.updated` - Triggered when a reminder is updated
- `reminder.deleted` - Triggered when a reminder is deleted
- `notification.reminder_due` - Triggered when a reminder's due time is reached

When a reminder's due time is reached, the server also sends a push notification to the user. Ensure push notifications are configured in your app.


## Webhooks

The same events are available as webhooks to notify your backend systems:

- `reminder.created`
- `reminder.updated`
- `reminder.deleted`
- `notification.reminder_due`

These webhook events contain the same payload structure as their WebSocket counterparts. For more information on configuring webhooks, see the [Webhooks documentation](/chat/docs/go-golang/webhook_events/).
