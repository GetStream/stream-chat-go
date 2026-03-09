Location sharing allows users to send a static position or share their real-time location with other participants in a channel. Stream Chat supports both static and live location sharing.

There are two types of location sharing:

- **Static Location**: A one-time location share that does not update over time.
- **Live Location**: A real-time location sharing that updates over time.

> [!NOTE]
> The SDK handles location message creation and updates, but location tracking must be implemented by the application using device location services.


## Enabling location sharing

The location sharing feature must be activated at the channel level before it can be used. You have two configuration options: activate it for a single channel using configuration overrides, or enable it globally for all channels of a particular type via [channel type settings](/chat/docs/go-golang/channel_features/).

```go
// Enabling it for a channel type
client.UpdateChannelType("messaging", map[string]interface{}{
	"shared_locations": true,
})
```

## Sending static location

Static location sharing allows you to send a message containing a static location.

```go
channel := chatClient.CreateChannelWithMembers(ctx, "messaging", channelID, userID)

// Create a SharedLocation Object
location := &SharedLocation{
	Longitude:         &longitude,
	Latitude:          &latitude,
	CreatedByDeviceID: "test-device",
}

// Send a message with the SharedLocation object
message := channel.SendMessage(ctx, &Message{
    ShareSharedLocation: location,
})
```

## Starting live location sharing

Live location sharing enables real-time location updates for a specified duration. The SDK manages the location message lifecycle, but your application is responsible for providing location updates.

```go
channel := chatClient.CreateChannelWithMembers(ctx, "messaging", channelID, userID)

// Create a SharedLocation Object with end_at
end_at := time.Now().Add(1 * time.Hour)
location := &SharedLocation{
	Longitude:         &longitude,
	Latitude:          &latitude,
    EndAt:             &end_at,
	CreatedByDeviceID: "test-device",
}

// Send a message with the SharedLocation object
message := channel.SendMessage(ctx, &Message{
    SharedLocation: location,
})
```

## Stopping live location sharing

You can stop live location sharing for a specific message using the message controller:

```go
// Set end_at to now
end_at := time.Now()
location := &SharedLocation{
    MessageID:         message_id,
	Longitude:         &longitude,
	Latitude:          &latitude,
    EndAt:             &end_at,
	CreatedByDeviceID: "test-device",
}

// Update the live location
c.UpdateUserActiveLocation(ctx, user.ID, newLocation)
```

## Updating live location

Your application must implement location tracking and provide updates to the SDK. The SDK handles updating all the current user's active live location messages and provides a throttling mechanism to prevent excessive API calls.

```go
// Get current user active live locations
userActiveLiveLocations, err := c.GetUserActiveLocations(ctx, user.ID)

// New location to set
end_at := time.Now()
newLocation := &SharedLocation{
	Longitude:         &longitude,
	Latitude:          &latitude,
    EndAt:             &end_at,
	CreatedByDeviceID: "test-device",
}

// Update active live location of the current user
for _, location := range userActiveLiveLocations.ActiveLiveLocations {
    newLocation.MessageID = location.MessageID
    c.UpdateUserActiveLocation(ctx, user.ID, newLocation)
}
```

Whenever the location is updated, the message will automatically be updated with the new location.

The SDK will also notify your application when it should start or stop location tracking as well as when the active live location messages change.


## Events

Whenever a location is created or updated, the following WebSocket events will be sent:

- `message.new`: When a new location message is created.
- `message.updated`: When a location message is updated.

> [!NOTE]
> In Dart, these events are resolved to more specific location events:
>
> - `location.shared`: When a new location message is created.
> - `location.updated`: When a location message is updated.


You can easily check if a message is a location message by checking the `message.sharedLocation` property. For example, you can use this events to render the locations in a map view.
