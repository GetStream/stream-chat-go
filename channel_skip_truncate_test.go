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

func newMockClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c, err := NewClient("key", "secret")
	require.NoError(t, err)
	c.BaseURL = srv.URL

	return c
}

func TestChannel_DeleteSkipTruncate(t *testing.T) {
	ctx := context.Background()

	var query string
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		query = r.URL.Query().Get("skip_truncate")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"duration":"0.01"}`))
	})
	ch := &Channel{Type: "messaging", ID: "chan", client: c}

	_, err := ch.Delete(ctx)
	require.NoError(t, err)
	assert.Empty(t, query)

	_, err = ch.Delete(ctx, DeleteWithSkipTruncate())
	require.NoError(t, err)
	assert.Equal(t, "true", query)
}

func TestClient_DeleteChannelsSkipTruncate(t *testing.T) {
	ctx := context.Background()

	var body map[string]interface{}
	c := newMockClient(t, func(w http.ResponseWriter, r *http.Request) {
		body = nil
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"duration":"0.01","task_id":"task"}`))
	})

	_, err := c.DeleteChannels(ctx, []string{"messaging:chan"}, false)
	require.NoError(t, err)
	assert.NotContains(t, body, "skip_truncate")

	_, err = c.DeleteChannels(ctx, []string{"messaging:chan"}, false, DeleteWithSkipTruncate())
	require.NoError(t, err)
	assert.Equal(t, true, body["skip_truncate"])
}
