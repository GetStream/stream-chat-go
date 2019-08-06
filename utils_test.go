package stream_chat

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	APIKey     = os.Getenv("STREAM_API_KEY")
	APISecret  = os.Getenv("STREAM_API_SECRET")
	StreamHost = os.Getenv("STREAM_HOST")

	testUsers = []User{
		{ID: "frodo-baggins", Name: "Frodo Baggins", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 50}},
		{ID: "sam-gamgee", Name: "Samwise Gamgee", ExtraData: map[string]interface{}{"race": "Hobbit", "age": 38}},
		{ID: "gandalf", Name: "Gandalf the Grey", ExtraData: map[string]interface{}{"race": "Istari"}},
		{ID: "legolas", Name: "Legolas", ExtraData: map[string]interface{}{"race": "Elf", "age": 500}},
	}
)

func randomUser() User {
	return testUsers[rand.Intn(len(testUsers)-1)]
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func mustNoError(t *testing.T, err error) {
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}

func mustError(t *testing.T, err error) {
	if !assert.Error(t, err) {
		t.FailNow()
	}
}
