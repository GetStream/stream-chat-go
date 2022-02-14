package stream_chat

import (
	"context"
	"net/http"
	"time"
)

type AppSettings struct {
	DisableAuth               *bool               `json:"disable_auth_checks,omitempty"`
	DisablePermissions        *bool               `json:"disable_permissions_checks,omitempty"`
	APNConfig                 *APNConfig          `json:"apn_config,omitempty"`
	FirebaseConfig            *FirebaseConfig     `json:"firebase_config,omitempty"`
	WebhookURL                string              `json:"webhook_url,omitempty"`
	MultiTenantEnabled        *bool               `json:"multi_tenant_enabled,omitempty"`
	AsyncURLEnrichEnabled     *bool               `json:"async_url_enrich_enabled,omitempty"`
	AutoTranslationEnabled    *bool               `json:"auto_translation_enabled,omitempty"`
	Grants                    map[string][]string `json:"grants,omitempty"`
	MigratePermissionsToV2    *bool               `json:"migrate_permissions_to_v2,omitempty"`
	PermissionVersion         string              `json:"permission_version,omitempty"`
	FileUploadConfig          *FileUploadConfig   `json:"file_upload_config,omitempty"`
	ImageUploadConfig         *FileUploadConfig   `json:"image_upload_config,omitempty"`
	ImageModerationLabels     []string            `json:"image_moderation_labels,omitempty"`
	ImageModerationEnabled    *bool               `json:"image_moderation_enabled,omitempty"`
	BeforeMessageSendHookURL  string              `json:"before_message_send_hook_url,omitempty"`
	CustomActionHandlerURL    string              `json:"custom_action_handler_url,omitempty"`
	UserSearchDisallowedRoles []string            `json:"user_search_disallowed_roles,omitempty"`
	EnforceUniqueUsernames    string              `json:"enforce_unique_usernames,omitempty"`
	SqsURL                    string              `json:"sqs_url,omitempty"`
	SqsKey                    string              `json:"sqs_key,omitempty"`
	SqsSecret                 string              `json:"sqs_secret,omitempty"`
	WebhookEvents             []string            `json:"webhook_events,omitempty"`
	ChannelHideMembersOnly    *bool               `json:"channel_hide_members_only,omitempty"`
}

func (a *AppSettings) SetDisableAuth(b bool) *AppSettings {
	a.DisableAuth = &b
	return a
}

func (a *AppSettings) SetDisablePermissions(b bool) *AppSettings {
	a.DisablePermissions = &b
	return a
}

func (a *AppSettings) SetAPNConfig(c APNConfig) *AppSettings {
	a.APNConfig = &c
	return a
}

func (a *AppSettings) SetFirebaseConfig(c FirebaseConfig) *AppSettings {
	a.FirebaseConfig = &c
	return a
}

func (a *AppSettings) SetWebhookURL(s string) *AppSettings {
	a.WebhookURL = s
	return a
}

func (a *AppSettings) SetMultiTenant(b bool) *AppSettings {
	a.MultiTenantEnabled = &b
	return a
}

func (a *AppSettings) SetGrants(g map[string][]string) *AppSettings {
	a.Grants = g
	return a
}

func NewAppSettings() *AppSettings {
	return &AppSettings{}
}

type FileUploadConfig struct {
	AllowedFileExtensions []string `json:"allowed_file_extensions,omitempty"`
	BlockedFileExtensions []string `json:"blocked_file_extensions,omitempty"`
	AllowedMimeTypes      []string `json:"allowed_mime_types,omitempty"`
	BlockedMimeTypes      []string `json:"blocked_mime_types,omitempty"`
}

type APNConfig struct {
	Enabled              bool   `json:"enabled"`
	Development          bool   `json:"development"`
	AuthType             string `json:"auth_type,omitempty"`
	AuthKey              []byte `json:"auth_key,omitempty"`
	NotificationTemplate string `json:"notification_template"`
	Host                 string `json:"host,omitempty"`
	BundleID             string `json:"bundle_id,omitempty"`
	TeamID               string `json:"team_id,omitempty"`
	KeyID                string `json:"key_id,omitempty"`
}

type FirebaseConfig struct {
	Enabled              bool   `json:"enabled"`
	ServerKey            string `json:"server_key"`
	NotificationTemplate string `json:"notification_template,omitempty"`
	DataTemplate         string `json:"data_template,omitempty"`
}

type PushNotificationFields struct {
	APNConfig      APNConfig      `json:"apn"`
	FirebaseConfig FirebaseConfig `json:"firebase"`
}

type Policy struct {
	Name      string    `json:"name"`
	Resources []string  `json:"resources"`
	Roles     []string  `json:"roles"`
	Action    int       `json:"action"` // allow: 1, deny: 0
	Owner     bool      `json:"owner"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AppConfig struct {
	Name                      string                    `json:"name"`
	OrganizationName          string                    `json:"organization"`
	WebhookURL                string                    `json:"webhook_url"`
	SuspendedExplanation      string                    `json:"suspended_explanation"`
	PushNotifications         PushNotificationFields    `json:"push_notifications"`
	ConfigNameMap             map[string]*ChannelConfig `json:"channel_configs"`
	Policies                  map[string][]Policy       `json:"policies"`
	Suspended                 bool                      `json:"suspended"`
	DisableAuth               bool                      `json:"disable_auth_checks"`
	DisablePermissions        bool                      `json:"disable_permissions_checks"`
	MultiTenantEnabled        bool                      `json:"multi_tenant_enabled"`
	RevokeTokensIssuedBefore  *time.Time                `json:"revoke_tokens_issued_before"`
	AsyncURLEnrichEnabled     bool                      `json:"async_url_enrich_enabled"`
	Grants                    map[string][]string       `json:"grants"`
	PermissionVersion         string                    `json:"permission_version"`
	ImageModerationLabels     []string                  `json:"image_moderation_labels"`
	ImageModerationEnabled    *bool                     `json:"image_moderation_enabled"`
	BeforeMessageSendHookURL  string                    `json:"before_message_send_hook_url"`
	CustomActionHandlerURL    string                    `json:"custom_action_handler_url"`
	UserSearchDisallowedRoles []string                  `json:"user_search_disallowed_roles"`
	EnforceUniqueUsernames    string                    `json:"enforce_unique_usernames"`
	SqsURL                    string                    `json:"sqs_url"`
	SqsKey                    string                    `json:"sqs_key"`
	SqsSecret                 string                    `json:"sqs_secret"`
	WebhookEvents             []string                  `json:"webhook_events"`
}

type AppResponse struct {
	App *AppConfig `json:"app"`
	Response
}

// GetAppConfig returns app settings.
func (c *Client) GetAppConfig(ctx context.Context) (*AppResponse, error) {
	var resp AppResponse

	err := c.makeRequest(ctx, http.MethodGet, "app", nil, nil, &resp)
	return &resp, err
}

// UpdateAppSettings makes request to update app settings
// Example of usage:
//  settings := NewAppSettings().SetDisableAuth(true)
//  err := client.UpdateAppSettings(settings)
func (c *Client) UpdateAppSettings(ctx context.Context, settings *AppSettings) (*Response, error) {
	var resp Response
	err := c.makeRequest(ctx, http.MethodPatch, "app", nil, settings, &resp)
	return &resp, err
}

type CheckSQSRequest struct {
	SqsURL    string `json:"sqs_url"`
	SqsKey    string `json:"sqs_key"`
	SqsSecret string `json:"sqs_secret"`
}

type CheckSQSResponse struct {
	Status string                 `json:"status"`
	Error  string                 `json:"error"`
	Data   map[string]interface{} `json:"data"`
	Response
}

// CheckSqs checks whether the AWS credentials are valid for the SQS queue access.
func (c *Client) CheckSqs(ctx context.Context, req *CheckSQSRequest) (*CheckSQSResponse, error) {
	var resp CheckSQSResponse
	err := c.makeRequest(ctx, http.MethodPost, "check_sqs", nil, req, &resp)
	return &resp, err
}

type CheckPushRequest struct {
	MessageID            string `json:"message_id,omitempty"`
	ApnTemplate          string `json:"apn_template,omitempty"`
	FirebaseTemplate     string `json:"firebase_template,omitempty"`
	FirebaseDataTemplate string `json:"firebase_data_template,omitempty"`
	SkipDevices          *bool  `json:"skip_devices,omitempty"`
	UserID               string `json:"user_id,omitempty"`
	User                 *User  `json:"user,omitempty"`
}

type DeviceError struct {
	Provider     string `json:"provider"`
	ErrorMessage string `json:"error_message"`
}

type CheckPushResponse struct {
	DeviceErrors             map[string]DeviceError `json:"device_errors"`
	GeneralErrors            []string               `json:"general_errors"`
	SkipDevices              *bool                  `json:"skip_devices"`
	RenderedApnTemplate      string                 `json:"rendered_apn_template"`
	RenderedFirebaseTemplate string                 `json:"rendered_firebase_template"`
	RenderedMessage          map[string]string      `json:"rendered_message"`
	Response
}

// CheckPush initiates a push test.
func (c *Client) CheckPush(ctx context.Context, req *CheckPushRequest) (*CheckPushResponse, error) {
	var resp CheckPushResponse
	err := c.makeRequest(ctx, http.MethodPost, "check_push", nil, req, &resp)
	return &resp, err
}

// RevokeTokens revokes all tokens for an application issued before given time.
func (c *Client) RevokeTokens(ctx context.Context, before *time.Time) (*Response, error) {
	setting := make(map[string]interface{})
	if before == nil {
		setting["revoke_tokens_issued_before"] = nil
	} else {
		setting["revoke_tokens_issued_before"] = before.Format(time.RFC3339)
	}

	var resp Response
	err := c.makeRequest(ctx, http.MethodPatch, "app", nil, setting, &resp)
	return &resp, err
}
