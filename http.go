package stream_chat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Error struct {
	Code            int               `json:"code"`
	Message         string            `json:"message"`
	ExceptionFields map[string]string `json:"exception_fields,omitempty"`
	StatusCode      int               `json:"StatusCode"`
	Duration        string            `json:"duration"`
	MoreInfo        string            `json:"more_info"`

	RateLimit *RateLimitInfo `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

// Response is the base response returned to client. It contains rate limit information.
// All specific response returned to the client should embed this type.
type Response struct {
	RateLimitInfo *RateLimitInfo `json:"ratelimit"`
}

func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	if resp.Body == nil {
		return errors.New("http body is nil")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response: %w", err)
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

		// Include rate limit information.
		apiErr.RateLimit = NewRateLimitFromHeaders(resp.Header)
		return apiErr
	}

	_, ok := result.(*Response)
	if !ok {
		// Unmarshal the body only when it is expected.
		err = json.Unmarshal(b, result)
		if err != nil {
			return fmt.Errorf("cannot unmarshal body: %w", err)
		}
	}

	return c.addRateLimitInfo(resp.Header, result)
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

	r, err := http.NewRequestWithContext(ctx, method, u, http.NoBody)
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
		r.Body = io.NopCloser(t)

	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(bytes.NewReader(b))
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

func (c *Client) addRateLimitInfo(headers http.Header, result interface{}) error {
	rl := map[string]interface{}{
		"ratelimit": NewRateLimitFromHeaders(headers),
	}

	b, err := json.Marshal(rl)
	if err != nil {
		return fmt.Errorf("cannot marshal rate limit info: %w", err)
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return fmt.Errorf("cannot unmarshal rate limit info: %w", err)
	}
	return nil
}
