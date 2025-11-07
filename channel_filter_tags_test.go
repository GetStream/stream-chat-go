package stream_chat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannel_AddFilterTags(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	chanID := randomString(12)

	resp, err := c.CreateChannel(ctx, "messaging", chanID, randomUser(t, c).ID, nil)
	require.NoError(t, err)
	ch := resp.Channel

	tags := []string{"sports", "news"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ensure correct endpoint
		assert.Equal(t, "/channels/messaging/"+chanID, r.URL.Path)
		// ensure POST
		assert.Equal(t, http.MethodPost, r.Method)
		// read body
		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, tags, toStringSlice(body["add_filter_tags"]))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"duration":"0.01"}`))
	}))
	defer srv.Close()

	// redirect client's baseURL to mock server
	ch.client.BaseURL = srv.URL

	_, err = ch.AddFilterTags(ctx, tags, nil)
	require.NoError(t, err)
}

func TestChannel_RemoveFilterTags(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	chanID := randomString(12)

	resp, err := c.CreateChannel(ctx, "messaging", chanID, randomUser(t, c).ID, nil)
	require.NoError(t, err)
	ch := resp.Channel

	tags := []string{"sports"}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/channels/messaging/"+chanID, r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, tags, toStringSlice(body["remove_filter_tags"]))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"duration":"0.01", "members":[], "messages":[]}`))
	}))
	defer srv.Close()

	ch.client.BaseURL = srv.URL

	_, err = ch.RemoveFilterTags(ctx, tags, nil)
	require.NoError(t, err)
}

func toStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	if s, ok := v.([]string); ok {
		return s
	}
	if arr, ok := v.([]interface{}); ok {
		out := make([]string, len(arr))
		for i, a := range arr {
			out[i] = a.(string)
		}
		return out
	}
	return nil
}
