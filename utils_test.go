package stream_chat

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	APIKey     = os.Getenv("STREAM_API_KEY")
	APISecret  = os.Getenv("STREAM_API_SECRET")
	StreamHost = os.Getenv("STREAM_HOST")

	serverUser = &User{ID: randomString(10), Name: "Gandalf the Grey", ExtraData: map[string]interface{}{"race": "Istari"}}

	testUsers = []*User{
		{ID: randomString(10), Name: "Frodo Baggins", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 50}},
		{ID: randomString(10), Name: "Samwise Gamgee", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 38}},
		{ID: randomString(10), Name: "Legolas", ExtraData: map[string]interface{}{"race": "Elf", "age": 500}},
		serverUser,
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomUser() *User {
	return testUsers[rand.Intn(len(testUsers)-1)]
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func mustNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if !assert.NoError(t, err, msgAndArgs...) {
		t.FailNow()
	}
}

func mustError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if !assert.Error(t, err, msgAndArgs) {
		t.FailNow()
	}
}
