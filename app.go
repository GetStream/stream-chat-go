package stream

import (
	"net/http"
	"time"
)

type AppSettings struct {
	DisableAuth        *bool           `json:"disable_auth_checks,omitempty"`
	DisablePermissions *bool           `json:"disable_permissions_checks,omitempty"`
	APNConfig          *APNConfig      `json:"apn_config,omitempty"`
	FirebaseConfig     *FirebaseConfig `json:"firebase_config,omitempty"`
	WebhookURL         *string         `json:"webhook_url,omitempty"`
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
	NotificationTemplate string `json:"notification_template"`
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
	Name                 string                    `json:"name"`
	OrganizationName     string                    `json:"organization"`
	WebhookURL           string                    `json:"webhook_url"`
	SuspendedExplanation string                    `json:"suspended_explanation"`
	PushNotifications    PushNotificationFields    `json:"push_notifications"`
	ConfigNameMap        map[string]*ChannelConfig `json:"channel_configs"`
	Policies             map[string][]Policy       `json:"policies"`
	Suspended            bool                      `json:"suspended"`
	DisableAuth          bool                      `json:"disable_auth_checks"`
	DisablePermissions   bool                      `json:"disable_permissions_checks"`
}

type appResponse struct {
	App *AppConfig `json:"app"`
}

// GetAppConfig returns app settings
func (c *Client) GetAppConfig() (*AppConfig, error) {
	var resp appResponse

	err := c.makeRequest(http.MethodGet, "app", nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return resp.App, nil
}

// UpdateAppSettings makes request to update app settings
// Example of usage:
//  settings := NewAppSettings().SetDisableAuth(true)
//  err := client.UpdateAppSettings(settings)
func (c *Client) UpdateAppSettings(settings *AppSettings) error {
	return c.makeRequest(http.MethodPatch, "app", nil, settings, nil)
}
