package stream_chat //nolint: golint

import (
	"context"
	"net/http"
	"time"
)

type AppSettings struct {
	DisableAuth           *bool           `json:"disable_auth_checks,omitempty"`
	DisablePermissions    *bool           `json:"disable_permissions_checks,omitempty"`
	APNConfig             *APNConfig      `json:"apn_config,omitempty"`
	FirebaseConfig        *FirebaseConfig `json:"firebase_config,omitempty"`
	WebhookURL            *string         `json:"webhook_url,omitempty"`
	MultiTenantEnabled    *bool           `json:"multi_tenant_enabled,omitempty"`
	AsyncURLEnrichEnabled *bool           `json:"async_url_enrich_enabled,omitempty"`
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
	a.WebhookURL = &s
	return a
}

func (a *AppSettings) SetMultiTenant(b bool) *AppSettings {
	a.MultiTenantEnabled = &b
	return a
}

func NewAppSettings() *AppSettings {
	return &AppSettings{}
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
	Name                     string                    `json:"name"`
	OrganizationName         string                    `json:"organization"`
	WebhookURL               string                    `json:"webhook_url"`
	SuspendedExplanation     string                    `json:"suspended_explanation"`
	PushNotifications        PushNotificationFields    `json:"push_notifications"`
	ConfigNameMap            map[string]*ChannelConfig `json:"channel_configs"`
	Policies                 map[string][]Policy       `json:"policies"`
	Suspended                bool                      `json:"suspended"`
	DisableAuth              bool                      `json:"disable_auth_checks"`
	DisablePermissions       bool                      `json:"disable_permissions_checks"`
	MultiTenantEnabled       bool                      `json:"multi_tenant_enabled"`
	RevokeTokensIssuedBefore *time.Time                `json:"revoke_tokens_issued_before"`
	AsyncURLEnrichEnabled    bool                      `json:"async_url_enrich_enabled"`
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
