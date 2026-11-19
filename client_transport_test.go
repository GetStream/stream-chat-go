package stream_chat

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type wrappedTransport struct {
	next http.RoundTripper
}

func (w wrappedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return w.next.RoundTrip(r)
}

func TestNewClient_ReplacedDefaultTransport(t *testing.T) {
	orig := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = orig })
	http.DefaultTransport = wrappedTransport{next: orig}

	var (
		c   *Client
		err error
	)
	require.NotPanics(t, func() {
		c, err = NewClient("key", "secret")
	})
	require.NoError(t, err)

	tr, ok := c.HTTP.Transport.(*http.Transport)
	require.True(t, ok, "transport should be *http.Transport")
	require.Equal(t, 5, tr.MaxIdleConnsPerHost)
}
