package stream_chat

import (
	"net/http"
	"time"
)

// WithTimeout sets http requests timeout to the client
func WithTimeout(t time.Duration) func(*Client) {
	return func(c *Client) {
		c.timeout = t
		c.http.Timeout = t
	}
}

// WithBaseURL sets base url to the client
func WithBaseURL(url string) func(*Client) {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPTransport sets custom transport for http client.
// Useful to set proxy, timeouts, tests etc.
func WithHTTPTransport(tr *http.Transport) func(*Client) {
	return func(c *Client) {
		c.http.Transport = tr
	}
}
