package stream_chat

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// initUsageStatsClient creates a client using usage stats credentials.
// Falls back to standard credentials if usage stats specific ones aren't set.
func initUsageStatsClient(t *testing.T) *Client {
	t.Helper()

	apiKey := os.Getenv("STREAM_USAGE_STATS_KEY")
	apiSecret := os.Getenv("STREAM_USAGE_STATS_SECRET")

	// Fall back to standard credentials
	if apiKey == "" {
		apiKey = os.Getenv("STREAM_KEY")
	}
	if apiSecret == "" {
		apiSecret = os.Getenv("STREAM_SECRET")
	}

	if apiKey == "" || apiSecret == "" {
		t.Skip("Missing credentials. Set STREAM_USAGE_STATS_KEY/STREAM_USAGE_STATS_SECRET or STREAM_KEY/STREAM_SECRET")
	}

	c, err := NewClient(apiKey, apiSecret)
	require.NoError(t, err)
	return c
}

// findTeamByName finds a team by name in the response.
func findTeamByName(teams []TeamUsageStats, teamName string) *TeamUsageStats {
	for i := range teams {
		if teams[i].Team == teamName {
			return &teams[i]
		}
	}
	return nil
}

// assertAllMetricsExact verifies all 16 metrics have expected exact values.
func assertAllMetricsExact(t *testing.T, team *TeamUsageStats, teamName string) {
	t.Helper()

	// Daily activity metrics
	require.Equal(t, int64(0), team.UsersDaily.Total, "%s users_daily", teamName)
	require.Equal(t, int64(100), team.MessagesDaily.Total, "%s messages_daily", teamName)
	require.Equal(t, int64(0), team.TranslationsDaily.Total, "%s translations_daily", teamName)
	require.Equal(t, int64(0), team.ImageModerationsDaily.Total, "%s image_moderations_daily", teamName)

	// Peak metrics
	require.Equal(t, int64(0), team.ConcurrentUsers.Total, "%s concurrent_users", teamName)
	require.Equal(t, int64(0), team.ConcurrentConnections.Total, "%s concurrent_connections", teamName)

	// User rolling/cumulative metrics
	require.Equal(t, int64(5), team.UsersTotal.Total, "%s users_total", teamName)
	require.Equal(t, int64(5), team.UsersLast24Hours.Total, "%s users_last_24_hours", teamName)
	require.Equal(t, int64(5), team.UsersLast30Days.Total, "%s users_last_30_days", teamName)
	require.Equal(t, int64(5), team.UsersMonthToDate.Total, "%s users_month_to_date", teamName)
	require.Equal(t, int64(0), team.UsersEngagedLast30Days.Total, "%s users_engaged_last_30_days", teamName)
	require.Equal(t, int64(0), team.UsersEngagedMonthToDate.Total, "%s users_engaged_month_to_date", teamName)

	// Message rolling/cumulative metrics
	require.Equal(t, int64(100), team.MessagesTotal.Total, "%s messages_total", teamName)
	require.Equal(t, int64(100), team.MessagesLast24Hours.Total, "%s messages_last_24_hours", teamName)
	require.Equal(t, int64(100), team.MessagesLast30Days.Total, "%s messages_last_30_days", teamName)
	require.Equal(t, int64(100), team.MessagesMonthToDate.Total, "%s messages_month_to_date", teamName)
}

// findTeamAcrossPages searches for a team across multiple pages using pagination.
func findTeamAcrossPages(t *testing.T, c *Client, teamName string) *TeamUsageStats {
	t.Helper()

	ctx := context.Background()
	var nextCursor string
	maxPages := 10 // Safety limit
	limit := 5

	for page := 0; page < maxPages; page++ {
		req := &QueryTeamUsageStatsRequest{
			Limit: &limit,
			Next:  nextCursor,
		}

		resp, err := c.QueryTeamUsageStats(ctx, req)
		require.NoError(t, err)

		team := findTeamByName(resp.Teams, teamName)
		if team != nil {
			return team
		}

		nextCursor = resp.Next
		if nextCursor == "" {
			break // No more pages
		}
	}
	return nil
}

// ============================================================================
// Basic Queries
// ============================================================================

func TestQueryTeamUsageStats_NoParametersReturnsTeams(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Teams)
	require.Greater(t, len(resp.Teams), 0, "Should return at least one team")
}

func TestQueryTeamUsageStats_EmptyRequestReturnsTeams(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{})

	require.NoError(t, err)
	require.NotNil(t, resp.Teams)
}

// ============================================================================
// Month Parameter
// ============================================================================

func TestQueryTeamUsageStats_ValidMonthWorks(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "2026-02",
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Teams)
	require.Greater(t, len(resp.Teams), 0)
}

func TestQueryTeamUsageStats_PastMonthReturnsEmpty(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "2025-01",
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Teams)
	require.Equal(t, 0, len(resp.Teams))
}

func TestQueryTeamUsageStats_InvalidMonthThrows(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "invalid",
	})

	require.Error(t, err)
}

func TestQueryTeamUsageStats_WrongLengthMonthThrows(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "2026",
	})

	require.Error(t, err)
}

// ============================================================================
// Date Range Parameters
// ============================================================================

func TestQueryTeamUsageStats_ValidDateRangeWorks(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "2026-02-01",
		EndDate:   "2026-02-17",
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Teams)
	require.Greater(t, len(resp.Teams), 0)
}

func TestQueryTeamUsageStats_SingleDayRangeWorks(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "2026-02-17",
		EndDate:   "2026-02-17",
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Teams)
}

func TestQueryTeamUsageStats_InvalidStartDateThrows(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "bad",
	})

	require.Error(t, err)
}

func TestQueryTeamUsageStats_EndBeforeStartThrows(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "2026-02-20",
		EndDate:   "2026-02-10",
	})

	require.Error(t, err)
}

// ============================================================================
// Pagination
// ============================================================================

func TestQueryTeamUsageStats_LimitReturnsCorrectCount(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 3

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
	})

	require.NoError(t, err)
	require.Equal(t, 3, len(resp.Teams))
}

func TestQueryTeamUsageStats_LimitReturnsNextCursor(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 3

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Next)
}

func TestQueryTeamUsageStats_PaginationReturnsDifferentTeams(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 3

	page1, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
	})
	require.NoError(t, err)

	page2, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
		Next:  page1.Next,
	})
	require.NoError(t, err)

	// Verify no overlap between pages
	for _, t1 := range page1.Teams {
		for _, t2 := range page2.Teams {
			require.NotEqual(t, t1.Team, t2.Team, "Pages should not have overlapping teams")
		}
	}
}

func TestQueryTeamUsageStats_MaxLimitWorks(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 30

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Teams)
}

func TestQueryTeamUsageStats_OverMaxLimitThrows(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 31

	_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
	})

	require.Error(t, err)
}

func TestQueryTeamUsageStats_LimitWithMonthWorks(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 2

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit: &limit,
		Month: "2026-02",
	})

	require.NoError(t, err)
	require.Equal(t, 2, len(resp.Teams))
}

func TestQueryTeamUsageStats_LimitWithDateRangeWorks(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()
	limit := 2

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Limit:     &limit,
		StartDate: "2026-02-01",
		EndDate:   "2026-02-17",
	})

	require.NoError(t, err)
	require.Equal(t, 2, len(resp.Teams))
}

// ============================================================================
// Response Structure Validation
// ============================================================================

func TestQueryTeamUsageStats_TeamsHaveTeamField(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)
	require.Greater(t, len(resp.Teams), 0)
	// team field exists (may be empty string for default team)
	_ = resp.Teams[0].Team
}

func TestQueryTeamUsageStats_AllMetricsPresent(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)
	require.Greater(t, len(resp.Teams), 0)

	team := resp.Teams[0]

	// Daily activity metrics - verify they are initialized (Go doesn't have null, struct will be zero value)
	// We check that Total is accessible without panic
	_ = team.UsersDaily.Total
	_ = team.MessagesDaily.Total
	_ = team.TranslationsDaily.Total
	_ = team.ImageModerationsDaily.Total

	// Peak metrics
	_ = team.ConcurrentUsers.Total
	_ = team.ConcurrentConnections.Total

	// Rolling/cumulative metrics
	_ = team.UsersTotal.Total
	_ = team.UsersLast24Hours.Total
	_ = team.UsersLast30Days.Total
	_ = team.UsersMonthToDate.Total
	_ = team.UsersEngagedLast30Days.Total
	_ = team.UsersEngagedMonthToDate.Total
	_ = team.MessagesTotal.Total
	_ = team.MessagesLast24Hours.Total
	_ = team.MessagesLast30Days.Total
	_ = team.MessagesMonthToDate.Total
}

func TestQueryTeamUsageStats_MetricTotalsNonNegative(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)

	for _, team := range resp.Teams {
		require.GreaterOrEqual(t, team.MessagesTotal.Total, int64(0), "messages_total should be >= 0")
		require.GreaterOrEqual(t, team.UsersDaily.Total, int64(0), "users_daily should be >= 0")
		require.GreaterOrEqual(t, team.ConcurrentUsers.Total, int64(0), "concurrent_users should be >= 0")
	}
}

// ============================================================================
// Data Correctness - Date Range Query
// ============================================================================

func TestQueryTeamUsageStats_DateRange_SdkTestTeam1_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "2026-02-17",
		EndDate:   "2026-02-18",
	})

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-1")
	require.NotNil(t, team, "sdk-test-team-1 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-1")
}

func TestQueryTeamUsageStats_DateRange_SdkTestTeam2_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "2026-02-17",
		EndDate:   "2026-02-18",
	})

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-2")
	require.NotNil(t, team, "sdk-test-team-2 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-2")
}

func TestQueryTeamUsageStats_DateRange_SdkTestTeam3_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		StartDate: "2026-02-17",
		EndDate:   "2026-02-18",
	})

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-3")
	require.NotNil(t, team, "sdk-test-team-3 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-3")
}

// ============================================================================
// Data Correctness - Month Query
// ============================================================================

func TestQueryTeamUsageStats_Month_SdkTestTeam1_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "2026-02",
	})

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-1")
	require.NotNil(t, team, "sdk-test-team-1 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-1")
}

func TestQueryTeamUsageStats_Month_SdkTestTeam2_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "2026-02",
	})

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-2")
	require.NotNil(t, team, "sdk-test-team-2 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-2")
}

func TestQueryTeamUsageStats_Month_SdkTestTeam3_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
		Month: "2026-02",
	})

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-3")
	require.NotNil(t, team, "sdk-test-team-3 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-3")
}

// ============================================================================
// Data Correctness - No Parameters Query
// ============================================================================

func TestQueryTeamUsageStats_NoParams_SdkTestTeam1_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-1")
	require.NotNil(t, team, "sdk-test-team-1 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-1")
}

func TestQueryTeamUsageStats_NoParams_SdkTestTeam2_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-2")
	require.NotNil(t, team, "sdk-test-team-2 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-2")
}

func TestQueryTeamUsageStats_NoParams_SdkTestTeam3_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)
	ctx := context.Background()

	resp, err := c.QueryTeamUsageStats(ctx, nil)

	require.NoError(t, err)
	team := findTeamByName(resp.Teams, "sdk-test-team-3")
	require.NotNil(t, team, "sdk-test-team-3 should exist")
	assertAllMetricsExact(t, team, "sdk-test-team-3")
}

// ============================================================================
// Data Correctness - Pagination Query
// ============================================================================

func TestQueryTeamUsageStats_Pagination_SdkTestTeam1_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)

	team := findTeamAcrossPages(t, c, "sdk-test-team-1")
	require.NotNil(t, team, "sdk-test-team-1 should exist across paginated results")
	assertAllMetricsExact(t, team, "sdk-test-team-1")
}

func TestQueryTeamUsageStats_Pagination_SdkTestTeam2_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)

	team := findTeamAcrossPages(t, c, "sdk-test-team-2")
	require.NotNil(t, team, "sdk-test-team-2 should exist across paginated results")
	assertAllMetricsExact(t, team, "sdk-test-team-2")
}

func TestQueryTeamUsageStats_Pagination_SdkTestTeam3_ExactValues(t *testing.T) {
	c := initUsageStatsClient(t)

	team := findTeamAcrossPages(t, c, "sdk-test-team-3")
	require.NotNil(t, team, "sdk-test-team-3 should exist across paginated results")
	assertAllMetricsExact(t, team, "sdk-test-team-3")
}
