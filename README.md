# stream-chat-go

[![build](https://github.com/GetStream/stream-chat-go/workflows/build/badge.svg)](https://github.com/GetStream/stream-chat-go/actions)
[![godoc](https://pkg.go.dev/badge/GetStream/stream-chat-go)](https://pkg.go.dev/github.com/GetStream/stream-chat-go/v3?tab=doc)

the official Golang API client for [Stream chat](https://getstream.io/chat/) a service for building chat applications.

You can sign up for a Stream account at https://getstream.io/chat/get_started/.

You can use this library to access chat API endpoints server-side, for the client-side integrations (web and mobile) have a look at the Javascript, iOS and Android SDK libraries (https://getstream.io/chat/).

### Installation

```bash
go get github.com/GetStream/stream-chat-go/v3
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

	stream "github.com/GetStream/stream-chat-go/v3"
)

var APIKey = os.Getenv("STREAM_CHAT_API_KEY")
var APISecret = os.Getenv("STREAM_CHAT_API_SECRET")
var userID = "" // your server user id
var baseURL = "https://chat-proxy-singapore.stream-io-api.com"

func main() {
	client, err := stream.NewClient(APIKey, APISecret)
	// handle error

	// set a base url if other than US, singapore in this example
	client.BaseURL = baseURL

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

Contributions to this project are very much welcome, please make sure that your code changes are tested and that follow
Go best-practices.
