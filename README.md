# stream-chat-go

[![build](https://github.com/GetStream/stream-chat-go/workflows/build/badge.svg)](https://github.com/GetStream/stream-chat-go/actions)
[![godoc](https://pkg.go.dev/badge/GetStream/stream-chat-go)](https://pkg.go.dev/github.com/GetStream/stream-chat-go/v4?tab=doc)

the official Golang API client for [Stream chat](https://getstream.io/chat/) a service for building chat applications.

You can sign up for a Stream account at https://getstream.io/chat/get_started/.

You can use this library to access chat API endpoints server-side, for the client-side integrations (web and mobile) have a look at the Javascript, iOS and Android SDK libraries (https://getstream.io/chat/).

### Installation

```bash
go get github.com/GetStream/stream-chat-go/v4
```

### Documentation

[Official API docs](https://getstream.io/chat/docs/)

### Supported features

- [x] Chat channels
- [x] Messages
- [x] Chat channel types
- [x] User management
- [x] Moderation API
- [x] Push configuration
- [x] User devices
- [x] User search
- [x] Channel search
- [x] Message search

### Quickstart

```go
package main

import (
	"os"

	stream "github.com/GetStream/stream-chat-go/v4"
)

var APIKey = os.Getenv("STREAM_KEY")
var APISecret = os.Getenv("STREAM_SECRET")
var userID = "" // your server user id

func main() {
	client, err := stream.NewClient(APIKey, APISecret)
	// handle error

	// use client methods

	// create channel with users
	users := []string{"id1", "id2", "id3"}
	channel, err := client.CreateChannel("messaging", "channel-id", userID, map[string]interface{}{
		"members": users,
	})

	// use channel methods
	msg, err := channel.SendMessage(&stream.Message{Text: "hello"}, userID)
}
```

### Contributing

Contributions to this project are very much welcome, please make sure that your code changes are tested and that they follow
Go best-practices. You can find some tips in [CONTRIBUTING.md](./CONTRIBUTING.md).

## We are hiring!

We've recently closed a [$38 million Series B funding round](https://techcrunch.com/2021/03/04/stream-raises-38m-as-its-chat-and-activity-feed-apis-power-communications-for-1b-users/) and we keep actively growing.
Our APIs are used by more than a billion end-users, and you'll have a chance to make a huge impact on the product within a team of the strongest engineers all over the world.

Check out our current openings and apply via [Stream's website](https://getstream.io/team/#jobs).
