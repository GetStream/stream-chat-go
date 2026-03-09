User presence allows you to show when a user was last active and if they are online right now. This feature can be enabled or disabled per channel type in the [channel type settings](/chat/docs/go-golang/channel_features/).

## Listening to Presence Changes

To receive presence updates, you need to watch a channel or query channels with `presence: true`. This allows you to show a user as offline when they leave and update their status in real time.


A users online status change can be handled via event delegation by subscribing to the `user.presence.changed` event the same you do for any other event.

## Presence Data Format

Whenever you read a user the presence data will look like this:


> [!NOTE]
> The online field indicates if the user is online. The status field stores text indicating the current user status.


> [!NOTE]
> The last_active field is updated when a user connects and then refreshed every 15 minutes.


## Invisible

To mark your user as invisible, you can update your user to set the invisible property to _true_. Your user will remain invisible even if you disconnect and reconnect. You must explicitly set invisible to _false_ in order to become visible again.


You can also set your user to invisible when connecting by setting the invisible property to _true_. You can also set a custom status message at the same time:


> [!NOTE]
> When invisible is set to _true,_ the current user will appear as offline to other users.
