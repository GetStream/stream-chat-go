package stream_chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pascaldekloe/jwt"
)

const (
	defaultBaseURL = "https://chat-us-east-1.stream-io-api.com"
	defaultTimeout = 6 * time.Second
)

type client struct {
	baseURL   string
	apiKey    string
	apiSecret []byte
	authToken string
	timeout   time.Duration
	http      *http.Client
}

func (c *client) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Stream-Client", "stream-go-client")
	r.Header.Set("Authorization", c.authToken)
	r.Header.Set("stream-auth-type", "jwt")
}

func (c *client) parseResponse(resp *http.Response, result interface{}) error {
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 399 {
		msg := bufio.NewScanner(resp.Body).Text()
		return fmt.Errorf("response code: %d; %s", resp.StatusCode, msg)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *client) requestURL(path string, params map[string][]string) (string, error) {
	_url, err := url.Parse(c.baseURL + path)
	if err != nil {
		return "", errors.New("url.Parse:" + err.Error())
	}

	// set request params to url
	for key, vv := range params {
		for _, v := range vv {
			_url.Query().Add(key, v)
		}
	}

	_url.Query().Set("api_key", c.apiKey)
	return _url.String(), nil
}

func (c *client) makeRequest(method string, path string, params map[string][]string, data interface{}, result interface{}) error {
	path, err := c.requestURL(path, params)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	c.setHeaders(r)

	resp, err := c.http.Do(r)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, result)
}

// NewClient creates new stream chat api client
func NewClient(apiKey string, apiSecret []byte, options ...func(*client)) (interface{}, error) {
	var claims jwt.Claims
	claims.Set["server"] = true
	token, err := claims.HMACSign(jwt.ES256, apiSecret)
	if err != nil {
		return nil, err
	}

	client := &client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		authToken: string(token),
		timeout:   defaultTimeout,
		baseURL:   defaultBaseURL,
		http:      http.DefaultClient,
	}

	for _, opt := range options {
		opt(client)
	}

	client.http.Timeout = client.timeout

	return client, nil
}
