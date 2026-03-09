If you use Stream's UI components, typing indicators are automatically handled.
The typing indicators can be turned on or off in the channel type settings.
This example below shows how to integrate typing indicators into your own message input UI.

## Sending Typing Events

When a user starts typing call the keystroke method. Optionally you can specify a thread id to have a thread specific typing indicator.
A few seconds after a user stops typing use stopTyping.


### Receiving typing indicator events

Listening to typing indicators uses the event system, an example is shown below


> [!NOTE]
> Because clients might fail at sending `typing.stop` event all Chat clients periodically prune the list of typing users.


### Typing Privacy Settings

Please take into account that `typing.start` and `typing.stop` events delivery can be controlled by user privacy settings:

```json
// user object with privacy settings where typing indicators are disabled
{
  // other user fields
  "privacy_settings": {
    "typing_indicators": {
      "enabled": false
    }
  }
}
```

If `privacy_settings.typing_indicators.enabled` is set to `false` , then `typing.start` and `typing.stop` events will be ignored for this user by Stream's server and these events will not be sent to other users. In other words other users will not know that the current user was typing.
