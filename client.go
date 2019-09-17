// Package stream_chat provides chat via stream api
//nolint: golint
package stream_chat

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/getstream/easyjson"
	"github.com/pascaldekloe/jwt"
)

const (
	defaultBaseURL = "https://chat-us-east-1.stream-io-api.com"
	defaultTimeout = 6 * time.Second
)

type Client struct {
	BaseURL string
	HTTP    *http.Client `json:"-"`

	apiKey    string
	apiSecret []byte
	authToken string
}

func (c *Client) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("X-Stream-Client", "stream-go-client")
	r.Header.Set("Authorization", c.authToken)
	r.Header.Set("Stream-Auth-Type", "jwt")
}

func (c *Client) parseResponse(resp *http.Response, result easyjson.Unmarshaler) error {
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode >= 399 {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("chat-client: HTTP %s %s status %s: %s",
			resp.Request.Method, resp.Request.URL, resp.Status, string(msg))
	}

	if result != nil {
		return easyjson.UnmarshalFromReader(resp.Body, result)
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

func (c *Client) makeRequest(method, path string,
	params url.Values, data interface{}, result easyjson.Unmarshaler) error {
	_url, err := c.requestURL(path, params)
	if err != nil {
		return err
	}

	var body []byte
	if m, ok := data.(easyjson.Marshaler); ok {
		body, err = easyjson.Marshal(m)
	} else {
		body, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	r, err := http.NewRequest(method, _url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	c.setHeaders(r)

	resp, err := c.HTTP.Do(r)
	if err != nil {
		return err
	}

	return c.parseResponse(resp, result)
}

// CreateToken creates new token for user with optional expire time
func (c *Client) CreateToken(userID string, expire time.Time) ([]byte, error) {
	if userID == "" {
		return nil, errors.New("user ID is empty")
	}

	params := map[string]interface{}{
		"user_id": userID,
	}

	return c.createToken(params, expire)
}

func (c *Client) createToken(params map[string]interface{}, expire time.Time) ([]byte, error) {
	var claims = jwt.Claims{
		Set: params,
	}

	claims.Expires = jwt.NewNumericTime(expire.Round(time.Second))

	return claims.HMACSign(jwt.HS256, c.apiSecret)
}

// VerifyWebhook validates if hmac signature is correct for message body
func (c *Client) VerifyWebhook(body []byte, signature []byte) (valid bool) {
	mac := hmac.New(crypto.SHA256.New, c.apiSecret)
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(signature, expectedMAC)
}

// NewClient creates new stream chat api client
func NewClient(apiKey string, apiSecret []byte) (*Client, error) {
	switch {
	case apiKey == "":
		return nil, errors.New("API key is empty")
	case len(apiSecret) == 0:
		return nil, errors.New("API secret is empty")
	}

	client := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		BaseURL:   defaultBaseURL,
		HTTP: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	token, err := client.createToken(map[string]interface{}{"server": true}, time.Time{})
	if err != nil {
		return nil, err
	}

	client.authToken = string(token)

	return client, nil
}
