package stream_chat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (c *client) parseResponse(resp *http.Response) (result interface{}, err error) {
	defer resp.Body.Close()
	if resp.StatusCode >= 399 {
		err = fmt.Errorf("response code: %d; %s", resp.StatusCode, bufio.NewScanner(resp.Body).Text())
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	return
}

func (c *client) makeRequest(method string, path string, params map[string][]string, data io.Reader) (interface{}, error) {

	_url, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, errors.New("url.Parse:" + err.Error())
	}

	// set request params to url
	for key, vv := range params {
		for _, v := range vv {
			_url.Query().Add(key, v)
		}
	}

	_url.Query().Set("api_key", c.apiKey)

	r, err := http.NewRequest(method, _url.String(), data)
	if err != nil {
		return nil, err
	}

	c.setHeaders(r)

	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}

	return c.parseResponse(resp)
}

func (c *client) Get(path string, params map[string][]string) (interface{}, error) {
	return c.makeRequest(http.MethodGet, path, params, nil)
}

func (c *client) Post(path string, params map[string][]string, data io.Reader) (interface{}, error) {
	return c.makeRequest(http.MethodPost, path, params, data)
}

func (c *client) Put(path string, params map[string][]string, data io.Reader) (interface{}, error) {
	return c.makeRequest(http.MethodPut, path, params, data)
}

func (c *client) Patch(path string, params map[string][]string, data io.Reader) (interface{}, error) {
	return c.makeRequest(http.MethodPatch, path, params, data)
}

func (c *client) Delete(path string, params map[string][]string) (interface{}, error) {
	return c.makeRequest(http.MethodDelete, path, params, nil)
}

// NewStreamChat creates new stream chat api client
func NewStreamChat(apiKey string, apiSecret []byte, options ...func(*client)) (interface{}, error) {
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
