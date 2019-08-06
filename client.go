package stream_chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/pascaldekloe/jwt"
)

const (
	defaultBaseURL = "https://chat-us-east-1.stream-io-api.com"
	defaultTimeout = 6 * time.Second
)

type Client struct {
	baseURL   string
	apiKey    string
	apiSecret []byte
	authToken string
	timeout   time.Duration
	http      *http.Client
}

func (c *Client) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Stream-Client", "stream-go-client")
	r.Header.Set("Authorization", c.authToken)
	r.Header.Set("Stream-Auth-Type", "jwt")
}

func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 399 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("chat-client: HTTP %s %s status %s: %s", resp.Request.Method, resp.Request.URL, resp.Status, string(msg))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) requestURL(path string, params map[string][]string) (string, error) {
	_url, err := url.Parse(c.baseURL + "/" + path)
	if err != nil {
		return "", errors.New("url.Parse: " + err.Error())
	}

	values := url.Values{}
	// set request params to url
	for key, vv := range params {
		for _, v := range vv {
			values.Add(key, v)
		}
	}

	values.Add("api_key", c.apiKey)

	_url.RawQuery = values.Encode()

	return _url.String(), nil
}

func (c *Client) makeRequest(method string, path string, params map[string][]string, data interface{}, result interface{}) error {
	_url, err := c.requestURL(path, params)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r, err := http.NewRequest(method, _url, bytes.NewReader(body))
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

// CreateToken creates new token for user with optional expire time
func (c *Client) CreateToken(userID string, expire ...time.Time) ([]byte, error) {
	params := map[string]interface{}{
		"user_id": userID,
	}

	return c.createToken(params, expire...)
}

func (c *Client) createToken(params map[string]interface{}, expire ...time.Time) ([]byte, error) {
	var claims = jwt.Claims{
		Set: params,
	}

	if len(expire) > 0 {
		claims.Expires = jwt.NewNumericTime(expire[0])
	}

	return claims.HMACSign(jwt.HS256, c.apiSecret)
}

// NewClient creates new stream chat api client
func NewClient(apiKey string, apiSecret []byte, options ...func(*Client)) (*Client, error) {
	client := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		timeout:   defaultTimeout,
		baseURL:   defaultBaseURL,
		http:      http.DefaultClient,
	}

	token, err := client.createToken(map[string]interface{}{"server": true})
	if err != nil {
		return nil, err
	}

	client.authToken = string(token)
	for _, opt := range options {
		opt(client)
	}

	client.http.Timeout = client.timeout

	return client, nil
}
