There are two ways to update a channel with the Stream API: partial updates and full updates. A partial update preserves existing custom key–value data, while a full update replaces the entire channel object and removes any fields not included in the request.

## Partial Update

A partial update lets you set or unset specific fields without affecting the rest of the channel’s custom data — essentially a patch-style update.

```go
resp, err := c.CreateChannel(ctx, "messaging", "general", "thierry", map[string]interface{}{
	"source": "user",
	"source_detail": map[string]interface{}{"user_id": "123"},
	"channel_detail": map[string]interface{}{"topic": "Plants and Animals", "rating": "pg"},
})

channel := resp.Channel

// let's change the source of this channel
channel.PartialUpdate(ctx, PartialUpdate{
	Set: map[string]interface{}{
		"source": "system",
	},
})

// since it's system generated we no longer need source_detail
channel.PartialUpdate(ctx, PartialUpdate{
	Unset: []string{"source_detail"},
})

// and finally update one of the nested fields in the channel_detail
channel.PartialUpdate(ctx, PartialUpdate{
	Set: map[string]interface{}{
		"channel_detail.topic": "Nature",
	},
})

// and maybe we decide we no longer need a rating
channel.PartialUpdate(ctx, PartialUpdate{
	Unset: []string{"channel_detail.rating"},
})
```

## Full Update

The `update` function updates all of the channel data. **Any data that is present on the channel and not included in a full update will be deleted.**

```go
channel.Update(ctx, map[string]interface{}{"color": "green"},)
```

### Request Params

| Name         | Type   | Description                                                                                                                                                                          | Optional |
| ------------ | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| channel data | object | Object with the new channel information. One special field is "frozen". Setting this field to true will freeze the channel. Read more about freezing channels in "Freezing Channels" |          |
| text         | object | Message object allowing you to show a system message in the Channel that something changed.                                                                                          | Yes      |

> [!NOTE]
> Updating a channel using these methods cannot be used to add or remove members. For this, you must use specific methods for adding/removing members, more information can be found [here](/chat/docs/go-golang/channel_members/).
