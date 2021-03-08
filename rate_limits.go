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

// GetRateLimitsOptions configures the Client.GetRateLimits call.
type GetRateLimitsOptions struct {
	// ServerSide, if true, restricts the returned limits to server-side clients only.
	ServerSide bool
	// Android, if true, restricts the returned limits to Android clients only.
	Android bool
	// IOS, if true, restricts the returned limits to iOS clients only.
	IOS bool
	// Web, if true, restricts the returned limits to web clients only.
	Web bool
	// Endpoints restricts the returned limits info to the specified endpoints.
	Endpoints []string
}

// GetRateLimits returns the current rate limit quotas and usage. If no params are toggled, all the limits
// for all platforms are returned.
func (c *Client) GetRateLimits(options GetRateLimitsOptions) (GetRateLimitsResponse, error) {
	var resp GetRateLimitsResponse

	params := url.Values{}
	if options.ServerSide {
		params.Set("server_side", "true")
	}
	if options.Android {
		params.Set("android", "true")
	}
	if options.IOS {
		params.Set("ios", "true")
	}
	if options.Web {
		params.Set("web", "true")
	}
	if options.Endpoints != nil {
		for _, e := range options.Endpoints {
			params.Add("endpoints", e)
		}
	}

	err := c.makeRequest(http.MethodGet, "rate_limits", params, nil, &resp)
	if err != nil {
		return GetRateLimitsResponse{}, err
	}
	return resp, nil
}
