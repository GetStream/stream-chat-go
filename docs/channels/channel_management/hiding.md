Hiding a channel removes it from query channel requests for that user until a new message is added. Only channel members can hide a channel.

Hidden channels may still have unread messages. Consider [marking the channel as read](/chat/docs/go-golang/unread/) before hiding it.

You can optionally clear the message history for that user when hiding. When a new message is received, it will be the only message visible to that user.

## Hide a Channel

```go
// hides the channel until a new message is added there
channel.hide(ctx, "john");

// shows a previously hidden channel
channel.show(ctx, "john");
```

> [!NOTE]
> You can still retrieve the list of hidden channels using the `{ "hidden" : true }` query parameter.
