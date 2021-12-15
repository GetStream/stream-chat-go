package stream_chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	if resp.Body == nil {
		return errors.New("http body is nil")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("failed to read HTTP response")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 399 {
		var apiErr Error
		err := json.Unmarshal(b, &apiErr)
		if err != nil {
			// IP rate limit errors sent by our Edge infrastructure are not JSON encoded.
			// If decode fails here, we need to handle this manually.
			apiErr.Message = string(b)
			apiErr.StatusCode = resp.StatusCode
			return apiErr
		}

		apiErr.RateLimit = NewRateLimitFromHeaders(resp.Header)
		return apiErr
	}

	if result != nil {
		return json.Unmarshal(b, result)
	}
	return nil
}

func (c *Client) requestURL(path string, values url.Values) (string, error) {
	u, err := url.Parse(c.BaseURL + "/" + path)
	if err != nil {
		return "", errors.New("url.Parse: " + err.Error())
	}

	if values == nil {
		values = make(url.Values)
	}

	values.Add("api_key", c.apiKey)

	u.RawQuery = values.Encode()

	return u.String(), nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, params url.Values, data interface{}) (*http.Request, error) {
	u, err := c.requestURL(path, params)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, method, u, nil)
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
