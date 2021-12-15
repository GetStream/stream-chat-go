package stream_chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Error struct {
	Code            int               `json:"code"`
	Message         string            `json:"message"`
	ExceptionFields map[string]string `json:"exception_fields,omitempty"`
	StatusCode      int               `json:"StatusCode"`
	Duration        string            `json:"duration"`
	MoreInfo        string            `json:"more_info"`

	RateLimit *RateLimit `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

const (
	HeaderRateLimit     = "X-Ratelimit-Limit"
	HeaderRateRemaining = "X-Ratelimit-Remaining"
	HeaderRateReset     = "X-Ratelimit-Reset"
)

type RateLimit struct {
	Reset     time.Time
	Limit     int
	Remaining int
}

func NewRateLimitFromHeaders(headers http.Header) *RateLimit {
	var rl RateLimit

	limit, err := strconv.Atoi(headers.Get(HeaderRateLimit))
	if err == nil {
		rl.Limit = limit
	}
	remaining, err := strconv.Atoi(headers.Get(HeaderRateRemaining))
	if err == nil {
		rl.Remaining = remaining
	}
	reset, err := strconv.ParseInt(headers.Get(HeaderRateReset), 10, 64)
	if err == nil && reset > 0 {
		rl.Reset = time.Unix(reset, 0)
	}

	return &rl
}

func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 399 {
		var apiErr Error
		err := json.NewDecoder(resp.Body).Decode(&apiErr)
		if err != nil {
			apiErr.Message = fmt.Sprintf("cannot decode error: %v", err)
			return apiErr
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			apiErr.RateLimit = NewRateLimitFromHeaders(resp.Header)
		}
		return apiErr
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) requestURL(path string, values url.Values) (string, error) {
	_url, err := url.Parse(c.BaseURL + "/" + path)
	if err != nil {
		return "", errors.New("url.Parse: " + err.Error())
	}

	if values == nil {
		values = make(url.Values)
	}

	values.Add("api_key", c.apiKey)

	_url.RawQuery = values.Encode()

	return _url.String(), nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, params url.Values, data interface{}) (*http.Request, error) {
	_url, err := c.requestURL(path, params)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, method, _url, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(r)
	switch t := data.(type) {
	case nil:
		r.Body = nil

	case io.ReadCloser:
		r.Body = t

	case io.Reader:
		r.Body = ioutil.NopCloser(t)

	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(b))
	}

	return r, nil
}

func (c *Client) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Stream-Client", versionHeader())
	r.Header.Set("Authorization", c.authToken)
	r.Header.Set("Stream-Auth-Type", "jwt")
}

func (c *Client) makeRequest(ctx context.Context, method, path string, params url.Values, data, result interface{}) error {
	r, err := c.newRequest(ctx, method, path, params, data)
	if err != nil {
		return err
	}

	resp, err := c.HTTP.Do(r)
	if err != nil {
		select {
		case <-ctx.Done():
			// If we got an error, and the context has been canceled,
			// return context's error which is more useful.
			return ctx.Err()
		default:
		}
		return err
	}

	return c.parseResponse(resp, result)
}
