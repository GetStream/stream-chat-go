The channel query endpoint allows you to paginate messages, watchers, and members for a channel. Messages use ID-based pagination for consistency, while members and watchers use offset-based pagination.

## Message Pagination

Message pagination uses ID-based parameters rather than simple offset/limit. This approach improves performance and prevents issues when the message list changes while paginating.

For example, if you fetched the first 100 messages and want to load the next 100, pass the ID of the oldest message (when paginating in descending order) or the newest message (when paginating in ascending order).

### Pagination Parameters

| Parameter   | Description                                        |
| ----------- | -------------------------------------------------- |
| `id_lt`     | Retrieve messages older than (less than) the ID    |
| `id_gt`     | Retrieve messages newer than (greater than) the ID |
| `id_lte`    | Retrieve messages older than or equal to the ID    |
| `id_gte`    | Retrieve messages newer than or equal to the ID    |
| `id_around` | Retrieve messages around a specific message ID     |

```go
// Get the ID of the oldest message on the current page
lastMessageId := messages[0].Id

// Fetch older messages
options := map[string]interface{}{
  "messages": map[string]interface{}{"limit": 50, "id_lt": lastMessageId},
}

ch.query(ctx, options, nil)
```

## Member and Watcher Pagination

Members and watchers use `limit` and `offset` parameters for pagination.

| Parameter | Description                 | Maximum |
| --------- | --------------------------- | ------- |
| `limit`   | Number of records to return | 300     |
| `offset`  | Number of records to skip   | 10000   |

```go
options := map[string]interface{}{
  "members":  map[string]interface{}{"limit": 20, "offset": 0},
  "watchers": map[string]interface{}{"limit": 20, "offset": 0},
}

ch.query(ctx, options, nil)
```

> [!NOTE]
> To retrieve filtered and sorted members in a channel use the [Query Members](/chat/docs/go-golang/query_members/) API
