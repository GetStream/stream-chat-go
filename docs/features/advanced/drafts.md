Draft messages allow users to save messages as drafts for later use. This feature is useful when users want to compose a message but aren't ready to send it yet.

## Creating a draft message

It is possible to create a draft message for a channel or a thread. Only one draft per channel/thread can exist at a time, so a newly created draft overrides the existing one.

```go
// Create/update a draft message in a channel
resp, err := channel.CreateDraft(ctx, user.ID, &messageRequestMessage{
    Text: "This is a draft message",
})

// Create/update a draft message in a thread (parent message)
resp, err := channel.CreateDraft(ctx, user.ID, &messageRequestMessage{
    Text:     "This is a draft message",
    ParentID: parentID,
})

```

## Deleting a draft message

You can delete a draft message for a channel or a thread as well.

```go
// Channel draft
resp, err := channel.DeleteDraft(ctx, user.ID, nil)

// Thread draft
resp, err := channel.DeleteDraft(ctx, user.ID, parentMessageID)
```

## Loading a draft message

It is also possible to load a draft message for a channel or a thread. Although, when querying channels, each channel will contain the draft message payload, in case there is one. The same for threads (parent messages). So, for the most part this function will not be needed.

```go
// Channel draft
resp, err := channel.GetDraft(ctx, nil, user.ID)

// Thread draft
resp, err := channel.GetDraft(ctx, parentMessageID, user.ID)
```

## Querying draft messages

The Stream Chat SDK provides a way to fetch all the draft messages for the current user. This can be useful to for the current user to manage all the drafts they have in one place.

```go
// Query all user drafts
resp, err := c.QueryDrafts(ctx, &QueryDraftsOptions{UserID: user.ID, Limit: 10})

// Query drafts for certain channels and sort
resp, err := client.QueryDrafts(ctx, &QueryDraftsRequest{
    UserID: user.ID,
    Filter: map[string]interface{}{
        "channel_cid": map[string][]string{
            "$in": {"messaging:channel-1", "messaging:channel-2"},
        },
    },
    Sort: []*SortOption{
        {Field: "created_at", Direction: 1},
    },
})
```

Filtering is possible on the following fields:

| Name        | Type                       | Description                    | Supported operations      | Example                                                |
| ----------- | -------------------------- | ------------------------------ | ------------------------- | ------------------------------------------------------ |
| channel_cid | string                     | the ID of the message          | $in, $eq                  | { channel_cid: { $in: [ 'channel-1', 'channel-2' ] } } |
| parent_id   | string                     | the ID of the parent message   | $in, $eq, $exists         | { parent_id: 'parent-message-id' }                     |
| created_at  | string (RFC3339 timestamp) | the time the draft was created | $eq, $gt, $lt, $gte, $lte | { created_at: { $gt: '2024-04-24T15:50:00.00Z' }       |

Sorting is possible on the `created_at` field. By default, draft messages are returned with the newest first.

### Pagination

In case the user has a lot of draft messages, you can paginate the results.

```go
// Query drafts with a limit
resp, err := client.QueryDrafts(ctx, &QueryDraftsOptions{UserID: user.ID, Limit: 5})

// Query the next page
resp, err = client.QueryDrafts(ctx, &QueryDraftsOptions{
    UserID: user.ID,
    Limit:  5,
    Next:   *resp.Next,
})
```

## Events

The following WebSocket events are available for draft messages:

- `draft.updated`, triggered when a draft message is updated.
- `draft.deleted`, triggered when a draft message is deleted.

You can subscribe to these events using the Stream Chat SDK.
