This guide provides best practices for building chat experiences in marketplace applications, where communication between buyers and sellers is critical to successful transactions.

## Channel Types

The built-in **Messaging** channel type is designed for marketplace use cases and provides a good default configuration. It requires membership for channel access and supports all the features typically needed for buyer-seller communication.

For marketplace apps, we recommend:

- Using the `messaging` channel type as your starting point
- Creating channels with exactly 2 members (buyer and seller) for most transactions

## Recommended Features

### Unread Message Reminders

[Unread Reminders](/chat/docs/go-golang/unread-reminders/) notify users about messages they haven't read, helping ensure timely responses. Use them to trigger emails, push notifications, or SMS when a user has unread messages.

**Why this matters for marketplaces:**

- Buyers waiting for seller responses stay engaged
- Sellers don't miss potential sales opportunities
- Improves overall transaction completion rates

Enable reminders for your channel type:


> [!NOTE]
> Reminders work on channels with 10 or fewer members, making them ideal for typical buyer-seller conversations.


### User Average Response Time

The [User Average Response Time](/chat/docs/go-golang/user_average_response_time/) feature tracks how quickly users typically respond to messages. This is particularly valuable in marketplaces where prompt communication often determines transaction success.

**Benefits for marketplaces:**

- Buyers can see seller responsiveness before initiating contact
- Marketplaces can highlight responsive sellers with badges or sorting options
- Helps set expectations for communication timelines

Enable this feature in your app settings:


Once enabled, the `avg_response_time` field will be included in user responses, allowing you to display responsiveness indicators in seller profiles.

### Pending Messages

For marketplaces requiring message moderation, enable [pending messages](/chat/docs/go-golang/pending_messages/) to review content before delivery. This is useful for:

- Preventing spam or fraudulent messages
- Filtering contact information to keep transactions on-platform
- Ensuring compliance with marketplace policies


## Channel Configuration Tips

### Recommended Settings

For most marketplace apps, configure your channel type with these settings:


### Custom Data

Use channel and message custom data to enhance the marketplace experience:


## Trust and Safety

### Moderation

Implement moderation to maintain a safe marketplace environment:

- Use [block lists](/moderation/docs/engines/blocklists-and-regex-filters/) to filter prohibited content
- Enable [automod](/chat/docs/go-golang/moderation/) for automatic content filtering
- Set up [webhooks](/chat/docs/go-golang/webhooks_overview/) to monitor suspicious activity

## Performance Considerations

- Use [query channels best practices](/chat/docs/go-golang/query_channels/#best-practices) when loading conversation lists
- Filter by `cid` when querying specific channels for best performance
- Always include `members: { $in: [userID] }` in filters for consistent results

## Analytics and Insights

Track key metrics to improve your marketplace communication:

- **Response times**: Use `avg_response_time` to identify and reward responsive users
- **Message volume**: Monitor conversation activity to understand engagement
- **Unread rates**: Track how often reminders are triggered to optimize notification timing
