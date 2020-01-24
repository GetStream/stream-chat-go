# stream-chat-go

[![build](https://github.com/GetStream/stream-chat-go/workflows/build/badge.svg)](https://github.com/GetStream/stream-chat-go/actions)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/GetStream/stream-chat-go)

the official Golang API client for [Stream chat](https://getstream.io/chat/) a service for building chat applications.

You can sign up for a Stream account at https://getstream.io/chat/get_started/.

You can use this library to access chat API endpoints server-side, for the client-side integrations (web and mobile) have a look at the Javascript, iOS and Android SDK libraries (https://getstream.io/chat/).

### Installation

```bash
go get github.com/GetStream/stream-chat-go@v2
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

	stream "github.com/GetStream/stream-chat-go/v2"
)

var APIKey = os.Getenv("STREAM_API_KEY")
var APISecret = os.Getenv("STREAM_API_SECRET")
var userID = "" // your server user id

func main() {
	client, err := stream.NewClient(APIKey, []byte(APISecret))
	// use client methods

	// create channel with users
	users := []string{"id1", "id2", "id3"}
	channel, err := client.CreateChannel("messaging", "channel-id", userID, map[string]interface{}{
		"users": users,
	})

	// use channel methods
	msg, err := channel.SendMessage(&stream.Message{Text: "hello"}, userID)
}
```

### Contributing

Contributions to this project are very much welcome, please make sure that your code changes are tested and that follow
Go best-practices.
