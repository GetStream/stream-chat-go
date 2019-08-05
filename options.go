package stream_chat

import (
	"net/http"
	"time"
)

// WithTimeout sets http requests timeout to the client
func WithTimeout(t time.Duration) func(*client) {
	return func(c *client) {
		c.timeout = t
	}
}

// WithBaseURL sets base url to the client
func WithBaseURL(url string) func(*client) {
	return func(c *client) {
		c.baseURL = url
	}
}

// WithHTTPTransport sets custom transport for http client.
// Useful to set proxy, timeouts, tests etc.
func WithHTTPTransport(tr *http.Transport) func(*client) {
	return func(c *client) {
		c.http.Transport = tr
	}
}
