Push preferences allow users to control how they receive push notifications. You can set preferences at both the user level (global preferences) and channel level (per-channel preferences), giving users fine-grained control over their notification experience.

## How Push Preferences Work

> [!WARNING]
> Users must be channel members to receive push notifications regardless of their preferences.


### Chat push preferences operate on two levels

- **User-level preferences**: These act as global preferences that apply to all channels for a user.
- **Channel-level preferences**: These override the global preferences for specific channels.

### Chat push preferences support three levels of notifications

- **all**: Receive all push notifications **(default)**.
- **mentions**: Receive push notifications only when mentioned.
- **none**: Do not receive push notifications.

Additionally, you can temporarily disable push notifications until a specific time using the `disabled_until` parameter.

### The system evaluates preferences in the following priority order

1. **Channel-level preferences** are checked first (if they exist for the specific channel).
2. If no channel-level preference exists, **user-level (global) preferences** are used.
3. If no preferences are set at all, the default behavior is "all".
4. **Temporary disabling**: If `disabled_until` is set and the current time is before that timestamp, notifications are disabled regardless of other preferences.

## Setting Push Preferences

### User-Level Preferences

Set global push preferences that apply to all channels for a user:


### Channel-Level Preferences

Set preferences for specific channels, which override user-level preferences:


## Client-Side vs Server-Side Usage

### Client-Side Usage

When using client-side authentication, users can only update their own push preferences:


### Server-Side Usage

Server-side requests can update preferences for any user:


## Practical Examples

### 1: Creating a "Do Not Disturb" Mode


### 2: Channel-Specific Notification Settings

You can set different preferences for each individual channel, allowing users to customize their notification experience on a per-channel basis.


### 3: Temporarily Disabling Push Notifications

You can temporarily disable push notifications until a specific time using the `disabled_until` parameter. This is useful for implementing "Do Not Disturb" periods or scheduled quiet hours.


## Call Push Preferences

You can set preferences for call-related push notifications using the `call_level` field.

### Call push preferences support two levels of notifications

- **all**: Receive all call push notifications **(default)**.
- **none**: Do not receive call push notifications.
