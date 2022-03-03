package stream_chat

import (
	"bytes"
	"context"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// DefaultBaseURL is the default base URL for the stream chat api.
	// It works like CDN style and connects you to the closest production server.
	// By default, there is no real reason to change it. Use it only if you know what you are doing.
	DefaultBaseURL = "https://chat.stream-io-api.com"
	defaultTimeout = 6 * time.Second
)

type Client struct {
	BaseURL string
	HTTP    *http.Client `json:"-"`

	apiKey    string
	apiSecret []byte
	authToken string
}

// NewClientFromEnvVars creates a new Client where the API key
// is retrieved from STREAM_KEY and the secret from STREAM_SECRET
// environmental variables.
func NewClientFromEnvVars() (*Client, error) {
	return NewClient(os.Getenv("STREAM_KEY"), os.Getenv("STREAM_SECRET"))
}

// NewClient creates new stream chat api client.
func NewClient(apiKey, apiSecret string) (*Client, error) {
	switch {
	case apiKey == "":
		return nil, errors.New("API key is empty")
	case apiSecret == "":
		return nil, errors.New("API secret is empty")
	}

	baseURL := DefaultBaseURL
	if baseURLEnv := os.Getenv("STREAM_CHAT_URL"); strings.HasPrefix(baseURLEnv, "http") {
		baseURL = baseURLEnv
	}

	timeout := defaultTimeout
	if timeoutEnv := os.Getenv("STREAM_CHAT_TIMEOUT"); timeoutEnv != "" {
		i, err := strconv.Atoi(timeoutEnv)
		if err != nil {
			return nil, err
		}
		timeout = time.Duration(i) * time.Second
	}

	tr := http.DefaultTransport.(*http.Transport).Clone() //nolint:forcetypeassert
	tr.MaxIdleConnsPerHost = 5
	tr.IdleConnTimeout = 59 * time.Second // load balancer's idle timeout is 60 sec
	tr.ExpectContinueTimeout = 2 * time.Second

	client := &Client{
		apiKey:    apiKey,
		apiSecret: []byte(apiSecret),
		BaseURL:   baseURL,
		HTTP: &http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
	}

	token, err := client.createToken(jwt.MapClaims{"server": true})
	if err != nil {
		return nil, err
	}

	client.authToken = token

	return client, nil
}

// SetClient sets a new underlying HTTP client.
func (c *Client) SetClient(client *http.Client) {
	c.HTTP = client
}

// Channel returns a Channel object for future API calls.
func (c *Client) Channel(channelType, channelID string) *Channel {
	return &Channel{
		client: c,

		ID:   channelID,
		Type: channelType,
	}
}

// Permissions returns a client for handling app permissions.
func (c *Client) Permissions() *PermissionClient {
	return &PermissionClient{client: c}
}

// CreateToken creates a new token for user with optional expire time.
// Zero time is assumed to be no expire.
func (c *Client) CreateToken(userID string, expire time.Time, issuedAt ...time.Time) (string, error) {
	if userID == "" {
		return "", errors.New("user ID is empty")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
	}
	if !expire.IsZero() {
		claims["exp"] = expire.Unix()
	}
	if len(issuedAt) > 0 {
		claims["iat"] = issuedAt[0].Unix()
	}

	return c.createToken(claims)
}

func (c *Client) createToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(c.apiSecret)
}

// VerifyWebhook validates if hmac signature is correct for message body.
func (c *Client) VerifyWebhook(body, signature []byte) (valid bool) {
	mac := hmac.New(crypto.SHA256.New, c.apiSecret)
	_, _ = mac.Write(body)

	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return bytes.Equal(signature, []byte(expectedMAC))
}

// this makes possible to set content type.
type multipartForm struct {
	*multipart.Writer
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name, file name and content type.
func (form *multipartForm) CreateFormFile(fieldName, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name=%q; filename=%q`, fieldName, filename))

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h.Set("Content-Type", contentType)

	return form.Writer.CreatePart(h)
}

func (form *multipartForm) setData(fieldName string, data interface{}) error {
	field, err := form.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	return json.NewEncoder(field).Encode(data)
}

func (form *multipartForm) setFile(fieldName string, r io.Reader, fileName, contentType string) error {
	file, err := form.CreateFormFile(fieldName, fileName, contentType)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, r)

	return err
}

type SendFileResponse struct {
	File string `json:"file"`
	Response
}

func (c *Client) sendFile(ctx context.Context, link string, opts SendFileRequest) (*SendFileResponse, error) {
	if opts.User == nil {
		return nil, errors.New("user is nil")
	}

	tmpfile, err := ioutil.TempFile("", opts.FileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
	}()

	form := multipartForm{multipart.NewWriter(tmpfile)}

	if err := form.setData("user", opts.User); err != nil {
		return nil, err
	}

	err = form.setFile("file", opts.Reader, opts.FileName, opts.ContentType)
	if err != nil {
		return nil, err
	}

	err = form.Close()
	if err != nil {
		return nil, err
	}

	if _, err = tmpfile.Seek(0, 0); err != nil {
		return nil, err
	}

	r, err := c.newRequest(ctx, http.MethodPost, link, nil, tmpfile)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", form.FormDataContentType())

	res, err := c.HTTP.Do(r)
	if err != nil {
		return nil, err
	}

	var resp SendFileResponse
	err = c.parseResponse(res, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, err
}
