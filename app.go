package stream_chat

import (
	"context"
	"net/http"
	"time"
)

type AppSettings struct {
	Name                     string                    `json:"name"`
	OrganizationName         string                    `json:"organization"`
	Suspended                bool                      `json:"suspended"`
	SuspendedExplanation     string                    `json:"suspended_explanation"`
	ConfigNameMap            map[string]*ChannelConfig `json:"channel_configs"`
	RevokeTokensIssuedBefore *time.Time                `json:"revoke_tokens_issued_before"`

	DisableAuth        *bool `json:"disable_auth_checks,omitempty"`
	DisablePermissions *bool `json:"disable_permissions_checks,omitempty"`

	PushNotifications        PushNotificationFields `json:"push_notifications"`
	PushConfig               *PushConfigRequest     `json:"push_config,omitempty"`
	APNConfig                *APNConfig             `json:"apn_config,omitempty"`
	FirebaseConfig           *FirebaseConfigRequest `json:"firebase_config,omitempty"`
	XiaomiConfig             *XiaomiConfigRequest   `json:"xiaomi_config,omitempty"`
	HuaweiConfig             *HuaweiConfigRequest   `json:"huawei_config,omitempty"`
	WebhookURL               *string                `json:"webhook_url,omitempty"`
	WebhookEvents            []string               `json:"webhook_events,omitempty"`
	SqsURL                   *string                `json:"sqs_url,omitempty"`
	SqsKey                   *string                `json:"sqs_key,omitempty"`
	SqsSecret                *string                `json:"sqs_secret,omitempty"`
	BeforeMessageSendHookURL *string                `json:"before_message_send_hook_url,omitempty"`
	CustomActionHandlerURL   *string                `json:"custom_action_handler_url,omitempty"`

	FileUploadConfig       *FileUploadConfig `json:"file_upload_config,omitempty"`
	ImageUploadConfig      *FileUploadConfig `json:"image_upload_config,omitempty"`
	ImageModerationLabels  []string          `json:"image_moderation_labels,omitempty"`
	ImageModerationEnabled *bool             `json:"image_moderation_enabled,omitempty"`

	PermissionVersion      *string             `json:"permission_version,omitempty"`
	MigratePermissionsToV2 *bool               `json:"migrate_permissions_to_v2,omitempty"`
	Policies               map[string][]Policy `json:"policies"`
	Grants                 map[string][]string `json:"grants,omitempty"`

	MultiTenantEnabled        *bool    `json:"multi_tenant_enabled,omitempty"`
	AsyncURLEnrichEnabled     *bool    `json:"async_url_enrich_enabled,omitempty"`
	AutoTranslationEnabled    *bool    `json:"auto_translation_enabled,omitempty"`
	RemindersInterval         int      `json:"reminders_interval,omitempty"`
	UserSearchDisallowedRoles []string `json:"user_search_disallowed_roles,omitempty"`
	EnforceUniqueUsernames    *string  `json:"enforce_unique_usernames,omitempty"`
	ChannelHideMembersOnly    *bool    `json:"channel_hide_members_only,omitempty"`
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

func (a *AppSettings) SetFirebaseConfig(c FirebaseConfigRequest) *AppSettings {
	a.FirebaseConfig = &c
	return a
}

func (a *AppSettings) SetWebhookURL(s string) *AppSettings {
	a.WebhookURL = &s
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
	AuthKey              string `json:"auth_key,omitempty"`
	NotificationTemplate string `json:"notification_template"`
	Host                 string `json:"host,omitempty"`
	BundleID             string `json:"bundle_id,omitempty"`
	TeamID               string `json:"team_id,omitempty"`
	KeyID                string `json:"key_id,omitempty"`
}

type PushNotificationFields struct {
	Version        string         `json:"version"`
	OfflineOnly    bool           `json:"offline_only"`
	APNConfig      APNConfig      `json:"apn"`
	FirebaseConfig FirebaseConfig `json:"firebase"`
	HuaweiConfig   HuaweiConfig   `json:"huawei"`
	XiaomiConfig   XiaomiConfig   `json:"xiaomi"`
}

type FirebaseConfigRequest struct {
	ServerKey            string  `json:"server_key"`
	NotificationTemplate string  `json:"notification_template,omitempty"`
	DataTemplate         string  `json:"data_template,omitempty"`
	APNTemplate          *string `json:"apn_template,omitempty"`
	CredentialsJSON      string  `json:"credentials_json,omitempty"`
}

type FirebaseConfig struct {
	Enabled              bool   `json:"enabled"`
	NotificationTemplate string `json:"notification_template"`
	DataTemplate         string `json:"data_template"`
}

type XiaomiConfigRequest struct {
	PackageName string `json:"package_name"`
	Secret      string `json:"secret"`
}

type XiaomiConfig struct {
	Enabled bool `json:"enabled"`
}

type HuaweiConfigRequest struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

type HuaweiConfig struct {
	Enabled bool `json:"enabled"`
}

type PushConfigRequest struct {
	Version     string `json:"version,omitempty"`
	OfflineOnly bool   `json:"offline_only,omitempty"`
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

type AppResponse struct {
	App *AppSettings `json:"app"`
	Response
}

// GetAppSettings returns app settings.
func (c *Client) GetAppSettings(ctx context.Context) (*AppResponse, error) {
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
	PushProviderName     string `json:"push_provider_name,omitempty"`
	PushProviderType     string `json:"push_provider_type,omitempty"`
	UserID               string `json:"user_id,omitempty"`
	User                 *User  `json:"user,omitempty"`
}

type DeviceError struct {
	Provider     string `json:"provider"`
	ProviderName string `json:"provider_name"`
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

type PushProvider struct {
	Type           PushProviderType `json:"type"`
	Name           string           `json:"name"`
	Description    string           `json:"description,omitempty"`
	DisabledAt     *time.Time       `json:"disabled_at,omitempty"`
	DisabledReason string           `json:"disabled_reason,omitempty"`

	APNAuthKey string `json:"apn_auth_key,omitempty"`
	APNKeyID   string `json:"apn_key_id,omitempty"`
	APNTeamID  string `json:"apn_team_id,omitempty"`
	APNTopic   string `json:"apn_topic,omitempty"`

	FirebaseCredentials          string  `json:"firebase_credentials,omitempty"`
	FirebaseNotificationTemplate *string `json:"firebase_notification_template,omitempty"`
	FirebaseAPNTemplate          *string `json:"firebase_apn_template,omitempty"`

	HuaweiAppID     string `json:"huawei_app_id,omitempty"`
	HuaweiAppSecret string `json:"huawei_app_secret,omitempty"`

	XiaomiPackageName string `json:"xiaomi_package_name,omitempty"`
	XiaomiAppSecret   string `json:"xiaomi_app_secret,omitempty"`
}

// UpsertPushProvider inserts or updates a push provider.
func (c *Client) UpsertPushProvider(ctx context.Context, provider *PushProvider) (*Response, error) {
	body := map[string]PushProvider{"push_provider": *provider}
	var resp Response
	err := c.makeRequest(ctx, http.MethodPost, "push_providers", nil, body, &resp)
	return &resp, err
}

// DeletePushProvider deletes a push provider.
func (c *Client) DeletePushProvider(ctx context.Context, providerType, name string) (*Response, error) {
	var resp Response
	err := c.makeRequest(ctx, http.MethodDelete, "push_providers/"+providerType+"/"+name, nil, nil, &resp)
	return &resp, err
}

type PushProviderListResponse struct {
	Response
	PushProviders []PushProvider `json:"push_providers"`
}

// ListPushProviders returns the list of push providers.
func (c *Client) ListPushProviders(ctx context.Context) (*PushProviderListResponse, error) {
	var providers PushProviderListResponse
	err := c.makeRequest(ctx, http.MethodGet, "push_providers", nil, nil, &providers)
	return &providers, err
}
