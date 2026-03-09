Stream Chat supports message reactions such as likes, hearts, and custom reaction types. Users can react to messages, and reactions can include custom data and scores for cumulative reactions.

## Sending a Reaction

Add a reaction to a message using the `sendReaction` method. Each user can have one reaction of each type per message.

```go
reaction := &stream_chat.Reaction{
  Type:      "love",
  ExtraData: map[string]interface{}{"custom_field": 123},
}
msg, err := channel.SendReaction(reaction, msgID, userID)
```

### Reaction Parameters

| Name           | Type    | Description                                                              | Default | Optional |
| -------------- | ------- | ------------------------------------------------------------------------ | ------- | -------- |
| message_id     | string  | ID of the message to react to                                            |         |          |
| type           | string  | Reaction type. Each user can have one reaction of each type per message. |         |          |
| score          | integer | Score for cumulative reactions                                           | 1       | ✓        |
| user_id        | string  | User ID (required for server-side calls)                                 |         | ✓        |
| enforce_unique | boolean | If true, replaces all existing reactions from this user with the new one | false   | ✓        |
| skip_push      | boolean | If true, do not send a push notification                                 | false   | ✓        |
| emoji_code     | string  | Unicode emoji for push notification display                              |         | ✓        |
| custom data    | object  | Custom fields for the reaction                                           |         | ✓        |

> [!WARNING]
> Custom data for reactions is limited to 1KB.


## Removing a Reaction

Remove a reaction by specifying the message ID and reaction type.

```go
_, err := channel.DeleteReaction(msgID, "love", userID)
```

## Retrieving Reactions

Reactions are included in the message object. Messages returned by the API include the 10 most recent reactions.

### Reaction Fields in Messages

| Field            | Type   | Description                                                                                     |
| ---------------- | ------ | ----------------------------------------------------------------------------------------------- |
| reaction_counts  | object | Count of reactions per type. Example: `{"love": 3, "fire": 2}`                                  |
| reaction_scores  | object | Sum of scores per type. Equals counts for standard reactions; differs for cumulative reactions. |
| reaction_groups  | object | Detailed statistics per type including count, sum_scores, first_reaction_at, last_reaction_at   |
| latest_reactions | array  | The 10 most recent reactions with type, user_id, and created_at                                 |
| own_reactions    | array  | The current user's reactions on this message                                                    |

<details>
<summary>Example Reaction Data</summary>

```json
{
  "reaction_counts": {
    "love": 3,
    "fire": 2,
    "thumbsup": 1
  },
  "reaction_scores": {
    "love": 3,
    "fire": 2,
    "thumbsup": 1
  },
  "reaction_groups": {
    "love": {
      "count": 3,
      "sum_scores": 3,
      "first_reaction_at": "2024-12-11T14:32:00.000Z",
      "last_reaction_at": "2024-12-11T15:18:00.000Z"
    },
    "fire": {
      "count": 2,
      "sum_scores": 2,
      "first_reaction_at": "2024-12-11T14:35:00.000Z",
      "last_reaction_at": "2024-12-11T14:52:00.000Z"
    },
    "thumbsup": {
      "count": 1,
      "sum_scores": 1,
      "first_reaction_at": "2024-12-11T16:05:00.000Z",
      "last_reaction_at": "2024-12-11T16:05:00.000Z"
    }
  },
  "latest_reactions": [
    {
      "type": "thumbsup",
      "user_id": "sarah-miller",
      "created_at": "2024-12-11T16:05:00.000Z"
    },
    {
      "type": "love",
      "user_id": "mike-johnson",
      "created_at": "2024-12-11T15:18:00.000Z"
    },
    {
      "type": "fire",
      "user_id": "emma-wilson",
      "created_at": "2024-12-11T14:52:00.000Z"
    }
  ],
  "own_reactions": []
}
```

</details>

> [!NOTE]
> Use `reaction_groups` instead of `reaction_counts` for if you're building a custom implementation. The `reaction_groups` field provides additional metadata including timestamps and is the recommended approach.


To retrieve more than 10 reactions, use pagination.

### Paginating Reactions

Retrieve reactions with pagination using `limit` and `offset` parameters.

| Parameter | Maximum Value |
| --------- | ------------- |
| limit     | 300           |
| offset    | 1000          |

```go
options := map[string][]string{
  "limit": {"10"},
}
_, err := channel.GetReactions(msgID, options)
```

## Querying Reactions

Filter reactions by type or user on a specific message. This endpoint requires the user to have read permission on the channel when called client-side.


## Cumulative Reactions

Cumulative reactions allow users to react multiple times to the same message, with the total score tracked. This is useful for features like Medium's "clap" functionality.

Set a `score` value when sending the reaction. The API returns:

- `sum_scores`: Total score across all users
- Individual user scores

```go
_, err = channel.SendReaction(&stream_chat.Reaction{
  Type:      "clap",
  ExtraData: map[string]interface{}{"score": 25},
}, msgID, userID)
```
