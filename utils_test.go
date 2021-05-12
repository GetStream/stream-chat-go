package stream_chat //nolint: golint

import (
	"math/rand"
	"os"
	"time"
)

//nolint: gochecknoglobals
var (
	APIKey     = os.Getenv("STREAM_CHAT_API_KEY")
	APISecret  = os.Getenv("STREAM_CHAT_API_SECRET")
	StreamHost = os.Getenv("STREAM_CHAT_API_HOST")

	serverUser *User
	testUsers  []*User
)

//nolint: gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())

	if err := clearOldChannelTypes(); err != nil {
		panic(err) // app has bad data from previous runs
	}

	serverUser = &User{ID: randomString(10), Name: "Gandalf the Grey", ExtraData: map[string]interface{}{"race": "Istari"}}

	testUsers = []*User{
		{ID: randomString(10), Name: "Frodo Baggins", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 50}},
		{ID: randomString(10), Name: "Samwise Gamgee", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 38}},
		{ID: randomString(10), Name: "Legolas", ExtraData: map[string]interface{}{"race": "Elf", "age": 500}},
		serverUser,
	}
}

func clearOldChannelTypes() error {
	c, err := NewClient(APIKey, APISecret)
	if err != nil {
		return err
	}
	c.BaseURL = defaultBaseURL

	got, err := c.ListChannelTypes()
	if err != nil {
		return err
	}

	for _, ct := range got {
		if contains(defaultChannelTypes, ct.Name) {
			continue
		}
		filter := map[string]interface{}{"type": ct.Name}
		chs, _ := c.QueryChannels(&QueryOption{Filter: filter})

		hasChannel := false
		for _, ch := range chs {
			if err := ch.Delete(); err != nil {
				hasChannel = true
				break
			}
		}

		if !hasChannel {
			_ = c.DeleteChannelType(ct.Name)
		}
	}
	return nil
}

func randomUser() *User {
	return testUsers[rand.Intn(len(testUsers)-1)]
}

func randomString(n int) string {
	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}

func contains(ls []string, s string) bool {
	for _, item := range ls {
		if item == s {
			return true
		}
	}
	return false
}
