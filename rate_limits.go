package stream_chat // nolint: golint

import (
	"net/http"
	"net/url"
	"time"
)

// RateLimitInfo represents the quota and usage for a single endpoint.
type RateLimitInfo struct {
	// Limit is the maximum number of API calls for a single time window (1 minute).
	Limit int64 `json:"limit"`
	// Remaining is the number of API calls remaining in the current time window (1 minute).
	Remaining int64 `json:"remaining"`
	// Reset is the Unix timestamp of the expiration of the current rate limit time window.
	Reset int64 `json:"reset"`
}

// RateLimitsMap holds the rate limit information, where the keys are the names of the endpoints and the values are
// the related RateLimitInfo containing the quota, usage, and reset data.
type RateLimitsMap map[string]RateLimitInfo

// ResetTime is a simple helper to get the time.Time representation of the Reset field of the given limit window.
func (i RateLimitInfo) ResetTime() time.Time {
	return time.Unix(i.Reset, 0)
}

// GetRateLimitsResponse is the response of the Client.GetRateLimits call. It includes, if present, the rate
// limits for the supported platforms, namely server-side, Android, iOS, and web.
type GetRateLimitsResponse struct {
	ServerSide RateLimitsMap `json:"server_side,omitempty"`
	Android    RateLimitsMap `json:"android,omitempty"`
	IOS        RateLimitsMap `json:"ios,omitempty"`
	Web        RateLimitsMap `json:"web,omitempty"`
}

type getRateLimitsParams struct {
	serverSide bool
	android    bool
	iOS        bool
	web        bool
	endpoints  []string
}

// GetRateLimitsOption configures the Client.GetRateLimits call.
type GetRateLimitsOption func(p *getRateLimitsParams)

// WithServerSide restricts the returned limits to server-side clients only.
func WithServerSide() GetRateLimitsOption {
	return func(p *getRateLimitsParams) {
		p.serverSide = true
	}
}

// WithAndroid restricts the returned limits to Android clients only.
func WithAndroid() GetRateLimitsOption {
	return func(p *getRateLimitsParams) {
		p.android = true
	}
}

// WithIOS restricts the returned limits to iOS clients only.
func WithIOS() GetRateLimitsOption {
	return func(p *getRateLimitsParams) {
		p.iOS = true
	}
}

// WithWeb restricts the returned limits to web clients only.
func WithWeb() GetRateLimitsOption {
	return func(p *getRateLimitsParams) {
		p.web = true
	}
}

// WithEndpoints restricts the returned limits info to the specified endpoints.
func WithEndpoints(endpoints ...string) GetRateLimitsOption {
	return func(p *getRateLimitsParams) {
		p.endpoints = append(p.endpoints, endpoints...)
	}
}

// GetRateLimits returns the current rate limit quotas and usage. If no options are passed, all the limits
// for all platforms are returned.
func (c *Client) GetRateLimits(options ...GetRateLimitsOption) (GetRateLimitsResponse, error) {
	var resp GetRateLimitsResponse

	params := url.Values{}

	rlParams := getRateLimitsParams{}
	for _, opt := range options {
		opt(&rlParams)
	}
	if rlParams.serverSide {
		params.Set("server_side", "true")
	}
	if rlParams.android {
		params.Set("android", "true")
	}
	if rlParams.iOS {
		params.Set("ios", "true")
	}
	if rlParams.web {
		params.Set("web", "true")
	}
	if rlParams.endpoints != nil {
		for _, e := range rlParams.endpoints {
			params.Add("endpoints", e)
		}
	}

	err := c.makeRequest(http.MethodGet, "rate_limits", params, nil, &resp)
	if err != nil {
		return GetRateLimitsResponse{}, err
	}
	return resp, nil
}
