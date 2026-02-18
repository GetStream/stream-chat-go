# Official Go SDK for [Stream Chat](https://getstream.io/chat/)

[![build](https://github.com/GetStream/stream-chat-go/workflows/build/badge.svg)](https://github.com/GetStream/stream-chat-go/actions)
[![godoc](https://pkg.go.dev/badge/GetStream/stream-chat-go)](https://pkg.go.dev/github.com/GetStream/stream-chat-go/v8?tab=doc)

<p align="center">
    <img src="./assets/logo.svg" width="50%" height="50%">
</p>
<p align="center">
    Official Go API client for Stream Chat, a service for building chat applications.
    <br />
    <a href="https://getstream.io/chat/docs/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/GetStream/stream-chat-go/issues">Report Bug</a>
    ·
    <a href="https://github.com/GetStream/stream-chat-go/issues">Request Feature</a>
</p>

## 📝 About Stream

You can sign up for a Stream account at our [Get Started](https://getstream.io/chat/get_started/) page.

You can use this library to access chat API endpoints server-side.

For the client-side integrations (web and mobile) have a look at the JavaScript, iOS and Android SDK libraries ([docs](https://getstream.io/chat/)).

## ⚙️ Installation

```shell
go get github.com/GetStream/stream-chat-go/v8
```

## ✨ Getting started

```go
package main

import (
	"os"

	stream "github.com/GetStream/stream-chat-go/v8"
)

var APIKey = os.Getenv("STREAM_KEY")
var APISecret = os.Getenv("STREAM_SECRET")
var userID = "" // your server user id

func main() {
	// Initialize client
	client, err := stream.NewClient(APIKey, APISecret)
	
	// Or with a specific timeout
	client, err := stream.NewClient(APIKey, APISecret, WithTimeout(3 * time.Second))

	// Or using only environmental variables: (required) STREAM_KEY, (required) STREAM_SECRET,
	// (optional) STREAM_CHAT_TIMEOUT
	client, err := stream.NewClientFromEnvVars()

	// handle error

	// Define a context
	ctx := context.Background()

	// use client methods

	// create channel with users
	users := []string{"id1", "id2", "id3"}
	userID := "id1"
	channel, err := client.CreateChannelWithMembers(ctx, "messaging", "channel-id", userID, users...)

	// use channel methods
	msg, err := channel.SendMessage(ctx, &stream.Message{Text: "hello"}, userID)
}
```

## ✍️ Contributing

We welcome code changes that improve this library or fix a problem, please make sure to follow all best practices and add tests if applicable before submitting a Pull Request on Github. We are very happy to merge your code in the official repository. Make sure to sign our [Contributor License Agreement (CLA)](https://docs.google.com/forms/d/e/1FAIpQLScFKsKkAJI7mhCr7K9rEIOpqIDThrWxuvxnwUq2XkHyG154vQ/viewform) first. See our [license file](./LICENSE) for more details.

Head over to [CONTRIBUTING.md](./CONTRIBUTING.md) for some development tips.

## 🧑‍💻 We are hiring!

We've recently closed a [$38 million Series B funding round](https://techcrunch.com/2021/03/04/stream-raises-38m-as-its-chat-and-activity-feed-apis-power-communications-for-1b-users/) and we keep actively growing.
Our APIs are used by more than a billion end-users, and you'll have a chance to make a huge impact on the product within a team of the strongest engineers all over the world.

Check out our current openings and apply via [Stream's website](https://getstream.io/team/#jobs).
