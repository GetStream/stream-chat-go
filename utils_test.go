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

	serverUser = &User{ID: randomString(10), Name: "Gandalf the Grey", ExtraData: map[string]interface{}{"race": "Istari"}}

	testUsers = []*User{
		{ID: randomString(10), Name: "Frodo Baggins", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 50}},
		{ID: randomString(10), Name: "Samwise Gamgee", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 38}},
		{ID: randomString(10), Name: "Legolas", ExtraData: map[string]interface{}{"race": "Elf", "age": 500}},
		serverUser,
	}
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
