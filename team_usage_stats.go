package stream_chat

import (
	"context"
	"net/http"
)

// DailyValue represents a metric value for a specific date.
type DailyValue struct {
	// Date in YYYY-MM-DD format.
	Date string `json:"date"`
	// Value is the metric value for this date.
	Value int64 `json:"value"`
}

// MetricStats represents statistics for a single metric with optional daily breakdown.
type MetricStats struct {
	// Daily contains per-day values (only present in daily mode).
	Daily []DailyValue `json:"daily,omitempty"`
	// Total is the aggregated total value.
	Total int64 `json:"total"`
}

// TeamUsageStats represents team-level usage statistics for multi-tenant apps.
type TeamUsageStats struct {
	// Team identifier (empty string for users not assigned to any team).
	Team string `json:"team"`

	// Daily activity metrics (total = SUM of daily values)

	// UsersDaily is the daily active users.
	UsersDaily MetricStats `json:"users_daily"`
	// MessagesDaily is the daily messages sent.
	MessagesDaily MetricStats `json:"messages_daily"`
	// TranslationsDaily is the daily translations.
	TranslationsDaily MetricStats `json:"translations_daily"`
	// ImageModerationsDaily is the daily image moderations.
	ImageModerationsDaily MetricStats `json:"image_moderations_daily"`

	// Peak metrics (total = MAX of daily values)

	// ConcurrentUsers is the peak concurrent users.
	ConcurrentUsers MetricStats `json:"concurrent_users"`
	// ConcurrentConnections is the peak concurrent connections.
	ConcurrentConnections MetricStats `json:"concurrent_connections"`

	// Rolling/cumulative metrics (total = LATEST daily value)

	// UsersTotal is the total users.
	UsersTotal MetricStats `json:"users_total"`
	// UsersLast24Hours is users active in last 24 hours.
	UsersLast24Hours MetricStats `json:"users_last_24_hours"`
	// UsersLast30Days is MAU - users active in last 30 days.
	UsersLast30Days MetricStats `json:"users_last_30_days"`
	// UsersMonthToDate is users active this month.
	UsersMonthToDate MetricStats `json:"users_month_to_date"`
	// UsersEngagedLast30Days is engaged MAU.
	UsersEngagedLast30Days MetricStats `json:"users_engaged_last_30_days"`
	// UsersEngagedMonthToDate is engaged users this month.
	UsersEngagedMonthToDate MetricStats `json:"users_engaged_month_to_date"`

	// MessagesTotal is total messages.
	MessagesTotal MetricStats `json:"messages_total"`
	// MessagesLast24Hours is messages in last 24 hours.
	MessagesLast24Hours MetricStats `json:"messages_last_24_hours"`
	// MessagesLast30Days is messages in last 30 days.
	MessagesLast30Days MetricStats `json:"messages_last_30_days"`
	// MessagesMonthToDate is messages this month.
	MessagesMonthToDate MetricStats `json:"messages_month_to_date"`
}

// QueryTeamUsageStatsRequest contains the parameters for querying team usage stats.
type QueryTeamUsageStatsRequest struct {
	// Month in YYYY-MM format (e.g., '2026-01'). Mutually exclusive with StartDate/EndDate.
	// Returns aggregated monthly values.
	Month string `json:"month,omitempty"`
	// StartDate in YYYY-MM-DD format. Used with EndDate for custom date range.
	// Returns daily breakdown.
	StartDate string `json:"start_date,omitempty"`
	// EndDate in YYYY-MM-DD format. Used with StartDate for custom date range.
	// Returns daily breakdown.
	EndDate string `json:"end_date,omitempty"`
	// Limit is the maximum number of teams to return per page (default: 30, max: 30).
	Limit *int `json:"limit,omitempty"`
	// Next is the cursor for pagination to fetch next page of teams.
	Next string `json:"next,omitempty"`
}

// QueryTeamUsageStatsResponse is the response from QueryTeamUsageStats.
type QueryTeamUsageStatsResponse struct {
	// Teams contains the array of team usage statistics.
	Teams []TeamUsageStats `json:"teams"`
	// Next is the cursor for pagination to fetch next page.
	Next string `json:"next,omitempty"`
	Response
}

// QueryTeamUsageStats queries team-level usage statistics from the warehouse database.
//
// Returns usage metrics grouped by team with cursor-based pagination.
//
// Date Range Options (mutually exclusive):
//   - Use 'Month' parameter (YYYY-MM format) for monthly aggregated values
//   - Use 'StartDate'/'EndDate' parameters (YYYY-MM-DD format) for daily breakdown
//   - If neither provided, defaults to current month (monthly mode)
//
// This endpoint is server-side only.
func (c *Client) QueryTeamUsageStats(ctx context.Context, req *QueryTeamUsageStatsRequest) (*QueryTeamUsageStatsResponse, error) {
	if req == nil {
		req = &QueryTeamUsageStatsRequest{}
	}

	var resp QueryTeamUsageStatsResponse
	err := c.makeRequest(ctx, http.MethodPost, "stats/team_usage", nil, req, &resp)
	return &resp, err
}
