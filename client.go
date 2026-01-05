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
	"net/url"
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

type ClientOption func(c *Client)

func WithTimeout(t time.Duration) func(c *Client) {
	return func(c *Client) {
		c.HTTP.Timeout = t
	}
}

// NewClientFromEnvVars creates a new Client where the API key
// is retrieved from STREAM_KEY and the secret from STREAM_SECRET
// environmental variables.
func NewClientFromEnvVars() (*Client, error) {
	return NewClient(os.Getenv("STREAM_KEY"), os.Getenv("STREAM_SECRET"))
}

// NewClient creates new stream chat api client.
func NewClient(apiKey, apiSecret string, options ...ClientOption) (*Client, error) {
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

	for _, fn := range options {
		fn(client)
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
	if len(issuedAt) > 0 && !issuedAt[0].IsZero() {
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
func (form *multipartForm) CreateFormFile(fieldName, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name=%q; filename=%q`, fieldName, filename))

	return form.Writer.CreatePart(h)
}

func (form *multipartForm) setData(fieldName string, data interface{}) error {
	field, err := form.CreateFormField(fieldName)
	if err != nil {
		return err
	}
	return json.NewEncoder(field).Encode(data)
}

func (form *multipartForm) setFile(fieldName string, r io.Reader, fileName string) error {
	file, err := form.CreateFormFile(fieldName, fileName)
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

	err = form.setFile("file", opts.Reader, opts.FileName)
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

type DeliveredMessageConfirmation struct {
	ChannelCID string `json:"cid"`
	MessageID  string `json:"id"`
}

// MarkDeliveredOptions represents the options for marking messages as delivered.
type MarkDeliveredOptions struct {
	LatestDeliveredMessages []DeliveredMessageConfirmation `json:"latest_delivered_messages"`
	User                    *User                          `json:"user,omitempty"`
	UserID                  string                         `json:"user_id,omitempty"`
}

// MarkDelivered sends the mark delivered event for this user, only works if the `delivery_receipts` setting is enabled.
// Note: Unlike the JavaScript SDK, this method doesn't automatically check delivery receipts settings
// as the Go SDK doesn't maintain user state. You should check this manually if needed.
func (c *Client) MarkDelivered(ctx context.Context, options *MarkDeliveredOptions) (*Response, error) {
	if options == nil {
		return nil, errors.New("options must not be nil")
	}

	if len(options.LatestDeliveredMessages) == 0 {
		return nil, errors.New("latest_delivered_messages must not be empty")
	}

	params := url.Values{}
	if options.User == nil && options.UserID == "" {
		return nil, errors.New("either user or user_id must be provided")
	}
	if options.User == nil {
		params.Set("user_id", options.UserID)
	} else {
		params.Set("user_id", options.User.ID)
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "channels/delivered", params, options, &resp)
	return &resp, err
}

// MarkDeliveredSimple is a convenience method to mark a message as delivered for a specific user.
func (c *Client) MarkDeliveredSimple(ctx context.Context, userID, messageID, channelCID string) (*Response, error) {
	if userID == "" {
		return nil, errors.New("user ID must not be empty")
	}
	if messageID == "" {
		return nil, errors.New("message ID must not be empty")
	}
	if channelCID == "" {
		return nil, errors.New("channel CID must not be empty")
	}

	options := &MarkDeliveredOptions{
		LatestDeliveredMessages: []DeliveredMessageConfirmation{
			{
				ChannelCID: channelCID,
				MessageID:  messageID,
			},
		},
		UserID: userID,
	}

	return c.MarkDelivered(ctx, options)
}

// UpdateChannelsBatch updates channels in batch based on the provided options.
func (c *Client) UpdateChannelsBatch(ctx context.Context, options *ChannelsBatchOptions) (*AsyncTaskResponse, error) {
	if options == nil {
		return nil, errors.New("options must not be nil")
	}

	var resp AsyncTaskResponse
	err := c.makeRequest(ctx, http.MethodPut, "channels/batch", nil, options, &resp)
	return &resp, err
}

// ChannelBatchUpdater returns a ChannelBatchUpdater instance for batch channel operations.
func (c *Client) ChannelBatchUpdater() *ChannelBatchUpdater {
	return &ChannelBatchUpdater{client: c}
}
