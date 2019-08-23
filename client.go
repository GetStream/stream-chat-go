package stream_chat

import (
	"errors"
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
	HTTP    *http.Client

	apiKey    string
	apiSecret []byte
	token     string

	header http.Header
}

func (c *Client) requestURL(path string, values url.Values) (string, error) {
	_url, err := url.Parse(c.BaseURL + "/" + path)
	if err != nil {
		return "", errors.New("url.Parse: " + err.Error())
	}

	if len(values) == 0 {
		values = make(url.Values, 1)
	}

	values.Add("api_key", c.apiKey)

	_url.RawQuery = values.Encode()

	return _url.String(), nil
}

func (c *Client) Get(path string, urlParams url.Values, result easyjson.Unmarshaler) error {
	_url, err := c.requestURL(path, urlParams)
	if err != nil {
		return err
	}

	return MakeRequest(c.HTTP, c.header, http.MethodGet, _url, nil, result)
}

func (c *Client) Post(path string, urlParams url.Values, body interface{}, result easyjson.Unmarshaler) error {
	_url, err := c.requestURL(path, urlParams)
	if err != nil {
		return err
	}

	return MakeRequest(c.HTTP, c.header, http.MethodPost, _url, body, result)
}

func (c *Client) Delete(path string, urlParams url.Values, result easyjson.Unmarshaler) error {
	_url, err := c.requestURL(path, urlParams)
	if err != nil {
		return err
	}

	return MakeRequest(c.HTTP, c.header, http.MethodDelete, _url, nil, result)
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
	claims.Expires = jwt.NewNumericTime(expire)

	return claims.HMACSign(jwt.HS256, c.apiSecret)
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
		header: make(http.Header),
	}

	token, err := client.createToken(map[string]interface{}{"server": true}, time.Time{})
	if err != nil {
		return nil, err
	}

	client.header.Set("Content-Type", "application/json")
	client.header.Set("X-Stream-Client", "stream-go-client")
	client.header.Set("Authorization", string(token))
	client.header.Set("Stream-Auth-Type", "jwt")

	return client, nil
}
