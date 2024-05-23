package stream_chat

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEventSupportsAllFields that we can decode all of the keys in the
// examples. We do this via the DisallowUnknownFields flag.
func TestEventSupportsAllFields(t *testing.T) {
	// Tests are taken from https://getstream.io/chat/docs/webhook_events/ and
	// compressed with `jq -c .`
	events := map[EventType]string{
		"message.new":      `{"cid":"messaging:fun","type":"message.new","message":{"id":"fff0d7c0-60bd-4835-833b-3843007817bf","text":"8b780762-4830-4e2a-aa43-18aabaf1732d","html":"<p>8b780762-4830-4e2a-aa43-18aabaf1732d</p>\n","type":"regular","user":{"id":"97b49906-0b98-463b-aa47-0aa945677eb2","role":"user","created_at":"2019-04-24T08:48:38.440123Z","updated_at":"2019-04-24T08:48:38.440708Z","online":false},"attachments":[],"latest_reactions":[],"own_reactions":[],"reaction_counts":null,"reply_count":0,"created_at":"2019-04-24T08:48:39.918761Z","updated_at":"2019-04-24T08:48:39.918761Z","mentioned_users":[]},"user":{"id":"97b49906-0b98-463b-aa47-0aa945677eb2","role":"user","created_at":"2019-04-24T08:48:38.440123Z","updated_at":"2019-04-24T08:48:38.440708Z","online":false,"channel_unread_count":1,"channel_last_read_at":"2019-04-24T08:48:39.900585Z","total_unread_count":1,"unread_channels":1,"unread_count":1},"created_at":"2019-04-24T08:48:38.949986Z","members":[{"user_id":"97b49906-0b98-463b-aa47-0aa945677eb2","user":{"id":"97b49906-0b98-463b-aa47-0aa945677eb2","role":"user","created_at":"2019-04-24T08:48:38.440123Z","updated_at":"2019-04-24T08:48:38.440708Z","online":false,"channel_unread_count":1,"channel_last_read_at":"2019-04-24T08:48:39.900585Z","total_unread_count":1,"unread_channels":1,"unread_count":1},"created_at":"2019-04-24T08:48:39.652296Z","updated_at":"2019-04-24T08:48:39.652296Z"}]}`,
		"message.read":     `{"cid":"messaging:fun","type":"message.read","user":{"id":"a6e21b36-798b-408a-9cd1-0cf6c372fc7f","role":"user","created_at":"2019-04-24T08:49:58.170034Z","updated_at":"2019-04-24T08:49:59.345304Z","last_active":"2019-04-24T08:49:59.344201Z","online":true,"total_unread_count":0,"unread_channels":0,"unread_count":0,"channel_unread_count":0,"channel_last_read_at":"2019-04-24T08:49:59.365498Z"},"created_at":"2019-04-24T08:49:59.365489Z"}`,
		"message.updated":  `{"cid":"messaging:fun","type":"message.updated","message":{"id":"93163f53-4174-4be8-90cd-e59bef78da00","text":"new stuff","html":"<p>new stuff</p>\n","type":"regular","user":{"id":"75af03a7-fe83-4a2a-a447-9ed4fac2ea36","role":"user","created_at":"2019-04-24T08:51:26.846395Z","updated_at":"2019-04-24T08:51:27.973941Z","last_active":"2019-04-24T08:51:27.972713Z","online":false},"attachments":[],"latest_reactions":[],"own_reactions":[],"reaction_counts":null,"reply_count":0,"created_at":"2019-04-24T08:51:28.005691Z","updated_at":"2019-04-24T08:51:28.138422Z","mentioned_users":[]},"user":{"id":"75af03a7-fe83-4a2a-a447-9ed4fac2ea36","role":"user","created_at":"2019-04-24T08:51:26.846395Z","updated_at":"2019-04-24T08:51:27.973941Z","last_active":"2019-04-24T08:51:27.972713Z","online":true,"channel_unread_count":1,"channel_last_read_at":"2019-04-24T08:51:27.994245Z","total_unread_count":2,"unread_channels":2,"unread_count":2},"created_at":"2019-04-24T10:51:28.142291+02:00"}`,
		"message.deleted":  `{"cid":"messaging:fun","type":"message.deleted","message":{"id":"268d121f-82e0-4de1-8c8b-ef1201efd7a3","text":"new stuff","html":"<p>new stuff</p>\n","type":"regular","user":{"id":"76cd8430-2f91-4059-90e5-02dffb910297","role":"user","created_at":"2019-04-24T09:44:21.390868Z","updated_at":"2019-04-24T09:44:22.537305Z","last_active":"2019-04-24T09:44:22.535872Z","online":false},"attachments":[],"latest_reactions":[],"own_reactions":[],"reaction_counts":{},"reply_count":0,"created_at":"2019-04-24T09:44:22.57073Z","updated_at":"2019-04-24T09:44:22.717078Z","deleted_at":"2019-04-24T09:44:22.730524Z","mentioned_users":[]},"created_at":"2019-04-24T09:44:22.733305Z"}`,
		"reaction.new":     `{"cid":"messaging:fun","type":"reaction.new","message":{"id":"4b3c7b6c-a39d-4069-9450-2a3716cf4ca6","text":"new stuff","html":"<p>new stuff</p>\n","type":"regular","user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T09:49:48.300566Z","online":false},"attachments":[],"latest_reactions":[{"message_id":"4b3c7b6c-a39d-4069-9450-2a3716cf4ca6","user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T09:49:48.300566Z","online":true},"type":"lol","created_at":"2019-04-24T09:49:48.481994Z"}],"own_reactions":[],"reaction_counts":{"lol":1},"reply_count":0,"created_at":"2019-04-24T09:49:48.334808Z","updated_at":"2019-04-24T09:49:48.483028Z","mentioned_users":[]},"reaction":{"message_id":"4b3c7b6c-a39d-4069-9450-2a3716cf4ca6","user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T09:49:48.300566Z","online":true},"type":"lol","created_at":"2019-04-24T09:49:48.481994Z"},"user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T09:49:48.300566Z","online":true,"unread_channels":2,"unread_count":2,"channel_unread_count":1,"channel_last_read_at":"2019-04-24T09:49:48.321138Z","total_unread_count":2},"created_at":"2019-04-24T09:49:48.488497Z"}`,
		"reaction.deleted": `{"cid":"messaging:fun","type":"reaction.deleted","message":{"id":"4b3c7b6c-a39d-4069-9450-2a3716cf4ca6","text":"new stuff","html":"<p>new stuff</p>\n","type":"regular","user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T09:49:48.300566Z","online":false},"attachments":[],"latest_reactions":[],"own_reactions":[],"reaction_counts":{},"reply_count":0,"created_at":"2019-04-24T09:49:48.334808Z","updated_at":"2019-04-24T09:49:48.511631Z","mentioned_users":[]},"reaction":{"message_id":"4b3c7b6c-a39d-4069-9450-2a3716cf4ca6","user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T11:49:48.497656+02:00","online":true},"type":"lol","created_at":"2019-04-24T09:49:48.481994Z"},"user":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T11:49:48.497656+02:00","online":true,"total_unread_count":2,"unread_channels":2,"unread_count":2,"channel_unread_count":1,"channel_last_read_at":"2019-04-24T09:49:48.321138Z"},"created_at":"2019-04-24T09:49:48.511082Z"}`,
		"member.added":     `{"cid":"messaging:fun","type":"member.added","member":{"user_id":"d4d7b21a-78d4-4148-9830-eb2d3b99c1ec","user":{"id":"d4d7b21a-78d4-4148-9830-eb2d3b99c1ec","role":"user","created_at":"2019-04-24T09:49:47.149933Z","updated_at":"2019-04-24T09:49:47.151159Z","online":false},"created_at":"2019-04-24T09:49:48.534412Z","updated_at":"2019-04-24T09:49:48.534412Z"},"user":{"id":"d4d7b21a-78d4-4148-9830-eb2d3b99c1ec","role":"user","created_at":"2019-04-24T09:49:47.149933Z","updated_at":"2019-04-24T09:49:47.151159Z","online":false,"channel_last_read_at":"2019-04-24T09:49:48.537084Z","total_unread_count":0,"unread_channels":0,"unread_count":0,"channel_unread_count":0},"created_at":"2019-04-24T09:49:48.537082Z"}`,
		"member.updated":   `{"cid":"messaging:fun","type":"member.updated","member":{"user_id":"d4d7b21a-78d4-4148-9830-eb2d3b99c1ec","user":{"id":"d4d7b21a-78d4-4148-9830-eb2d3b99c1ec","role":"user","created_at":"2019-04-24T09:49:47.149933Z","updated_at":"2019-04-24T09:49:47.151159Z","online":false},"is_moderator":true,"created_at":"2019-04-24T09:49:48.534412Z","updated_at":"2019-04-24T09:49:48.547034Z"},"user":{"id":"d4d7b21a-78d4-4148-9830-eb2d3b99c1ec","role":"user","created_at":"2019-04-24T09:49:47.149933Z","updated_at":"2019-04-24T09:49:47.151159Z","online":false,"total_unread_count":0,"unread_channels":0,"unread_count":0,"channel_unread_count":0,"channel_last_read_at":"2019-04-24T09:49:48.549211Z"},"created_at":"2019-04-24T09:49:48.54921Z"}`,
		"member.removed":   `{"cid":"messaging:fun","type":"member.removed","user":{"id":"6585dbbb-3d46-4943-9b14-a645aca11df4","role":"user","created_at":"2019-03-22T14:22:04.581208Z","online":false},"created_at":"2019-03-22T14:22:07.040496Z"}`,
		"channel.updated":  `{"cid":"messaging:fun","type":"channel.updated","channel":{"cid":"messaging:fun","id":"fun","type":"messaging","last_message_at":"2019-04-24T09:49:48.576202Z","created_by":{"id":"57fabaed-446a-40b4-a6ec-e0ac8cad57e3","role":"user","created_at":"2019-04-24T09:49:47.158005Z","updated_at":"2019-04-24T09:49:48.301933Z","last_active":"2019-04-24T09:49:48.497656Z","online":true},"created_at":"2019-04-24T09:49:48.180908Z","updated_at":"2019-04-24T09:49:48.180908Z","frozen":false,"config":{"created_at":"2016-08-18T16:42:30.586808Z","updated_at":"2016-08-18T16:42:30.586808Z","name":"messaging","typing_events":true,"read_events":true,"connect_events":true,"search":true,"reactions":true,"replies":true,"mutes":true,"message_retention":"infinite","max_message_length":5000,"automod":"disabled","commands":["giphy","flag","ban","unban","mute","unmute"]},"awesome":"yes"},"created_at":"2019-04-24T09:49:48.594316Z"}`,
		"channel.deleted":  `{"cid":"messaging:fun","type":"channel.deleted","channel":{"cid":"messaging:fun","id":"fun","type":"messaging","created_at":"2019-04-24T09:49:48.180908Z","updated_at":"2019-04-24T09:49:48.180908Z","deleted_at":"2019-04-24T09:49:48.626704Z","frozen":false,"config":{"created_at":"2016-08-18T18:42:30.586808+02:00","updated_at":"2016-08-18T18:42:30.586808+02:00","name":"messaging","typing_events":true,"read_events":true,"connect_events":true,"search":true,"reactions":true,"replies":true,"mutes":true,"message_retention":"infinite","max_message_length":5000,"automod":"disabled","commands":["giphy","flag","ban","unban","mute","unmute"]}},"created_at":"2019-04-24T09:49:48.630913Z"}`,
		"user.updated":     `{"type":"user.updated","user":{"id":"thierry-7b690297-98fa-42dd-b999-a75dd4c7c993","role":"user","online":false,"awesome":true},"created_at":"2019-04-24T12:54:58.956621Z","members":[]}`,
	}

	for name, blob := range events {
		dec := json.NewDecoder(bytes.NewBufferString(blob))
		dec.DisallowUnknownFields()

		result := Event{}
		if err := dec.Decode(&result); err != nil {
			t.Errorf("Error unmarshaling %q: %v", name, err)
		}

		assert.Equal(t, name, result.Type)
	}
}

func TestSendUserCustomEvent(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	tests := []struct {
		name         string
		event        *UserCustomEvent
		targetUserID string
		expectedErr  string
	}{
		{
			name: "ok",
			event: &UserCustomEvent{
				Type: "custom_event",
			},
			targetUserID: "user1",
		},
		{
			name:         "error: event is nil",
			event:        nil,
			targetUserID: "user1",
			expectedErr:  "event is nil",
		},
		{
			name:         "error: empty targetUserID",
			event:        &UserCustomEvent{},
			targetUserID: "",
			expectedErr:  "targetUserID should not be empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expectedErr == "" {
				_, err := c.UpsertUser(ctx, &User{ID: test.targetUserID})
				require.NoError(t, err)
			}

			_, err := c.SendUserCustomEvent(ctx, test.targetUserID, test.event)

			if test.expectedErr == "" {
				require.NoError(t, err)
				return
			}
			require.EqualError(t, err, test.expectedErr)
		})
	}
}

func TestMarshalUnmarshalUserCustomEvent(t *testing.T) {
	ev1 := UserCustomEvent{
		Type: "custom_event",
		ExtraData: map[string]interface{}{
			"name":   "John Doe",
			"age":    99.0,
			"hungry": true,
			"fruits": []interface{}{},
		},
	}

	b, err := json.Marshal(ev1)
	require.NoError(t, err)

	ev2 := UserCustomEvent{}
	err = json.Unmarshal(b, &ev2)
	require.NoError(t, err)

	require.Equal(t, ev1, ev2)
}
