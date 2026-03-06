Silent and system messages provide ways to add informational content to channels without disrupting users.

## Silent Messages

Silent messages do not increment unread counts or mark channels as unread. Use them for transactional or automated content such as:

- "Your ride is waiting"
- "James updated the trip information"
- "You and Jane are now matched"

Set `silent: true` when sending a message.

```go
msg := &Message{Text: "You completed your trip"}
channel.SendMessage(ctx, msg, user.ID, func(msg *messageRequest) {
  msg.Message.Silent = true
})
```

> [!NOTE]
> Existing messages cannot be converted to silent messages.


> [!NOTE]
> Silent replies still increment the parent message's `reply_count`.


> [!NOTE]
> Silent messages still trigger push notifications by default. Set `skip_push: true` to disable push notifications. See [Messages Overview](/chat/docs/go-golang/send_message/) for details.


## System Messages

System messages have a distinct visual presentation, typically styled differently from user messages. Set `type: "system"` to create a system message.

Stream's UI SDKs include default styling for system messages. However, customizing their appearance is common to match your application's design. See your platform's UI customization documentation for details on styling system messages.

Client-side users require the `Create System Message` permission. Server-side system messages do not require this permission.

```go
msg := &Message{
  Text: "You completed your trip",
  Type: "system",
}
resp, err := ch.SendMessage(ctx, msg, systemUserID)
```
