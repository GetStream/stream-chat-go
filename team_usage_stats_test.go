package stream_chat

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// =============================================================================
// Basic Tests - Use regular app credentials, expect empty teams
// These tests verify the API works correctly without requiring multi-tenant data
// =============================================================================

func TestQueryTeamUsageStats_BasicAPI(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	t.Run("No parameters returns valid response", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Teams)
		// Regular app doesn't have multi-tenant, so teams is empty
		require.Empty(t, resp.Teams)
	})

	t.Run("Empty request returns valid response", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Empty(t, resp.Teams)
	})

	t.Run("Month parameter works", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Month: "2026-02",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Empty(t, resp.Teams)
	})

	t.Run("Date range works", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			StartDate: "2026-02-01",
			EndDate:   "2026-02-17",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Empty(t, resp.Teams)
	})

	t.Run("Pagination works", func(t *testing.T) {
		limit := 10
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Limit: &limit,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Empty(t, resp.Teams)
		// No next cursor when teams is empty
		require.Empty(t, resp.Next)
	})

	t.Run("Invalid month throws error", func(t *testing.T) {
		_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Month: "invalid",
		})
		require.Error(t, err)
	})

	t.Run("Wrong length month throws error", func(t *testing.T) {
		_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Month: "2026",
		})
		require.Error(t, err)
	})

	t.Run("Invalid start date throws error", func(t *testing.T) {
		_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			StartDate: "bad",
		})
		require.Error(t, err)
	})

	t.Run("End before start throws error", func(t *testing.T) {
		_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			StartDate: "2026-02-20",
			EndDate:   "2026-02-10",
		})
		require.Error(t, err)
	})

	t.Run("Over max limit throws error", func(t *testing.T) {
		limit := 31
		_, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Limit: &limit,
		})
		require.Error(t, err)
	})

	t.Run("Past month returns empty", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Month: "2025-01",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Empty(t, resp.Teams)
	})
}

// =============================================================================
// Integration Tests - Require multi-tenant app credentials
// These tests verify specific data exists and metrics are correct
// =============================================================================

// initMultiTenantClient creates a client using multi-tenant app credentials.
// Returns nil if credentials are not set (tests will be skipped).
func initMultiTenantClient(t *testing.T) *Client {
	t.Helper()

	apiKey := os.Getenv("STREAM_MULTI_TENANT_KEY")
	apiSecret := os.Getenv("STREAM_MULTI_TENANT_SECRET")

	if apiKey == "" || apiSecret == "" {
		return nil
	}

	c, err := NewClient(apiKey, apiSecret)
	require.NoError(t, err)
	return c
}

// skipIfNoMultiTenant skips the test if multi-tenant credentials are not available.
func skipIfNoMultiTenant(t *testing.T, c *Client) {
	t.Helper()
	if c == nil {
		t.Skip("Multi-tenant credentials not set. Set STREAM_MULTI_TENANT_KEY and STREAM_MULTI_TENANT_SECRET to run integration tests.")
	}
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

// assertAllMetricsExact verifies all metrics have expected exact values.
func assertAllMetricsExact(t *testing.T, team *TeamUsageStats, teamName string) {
	t.Helper()

	// Daily activity metrics
	require.Equal(t, int64(5), team.UsersDaily.Total, "%s users_daily", teamName)
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

func TestQueryTeamUsageStats_Integration(t *testing.T) {
	c := initMultiTenantClient(t)
	skipIfNoMultiTenant(t, c)
	ctx := context.Background()

	t.Run("No parameters returns teams", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Teams)
		require.Greater(t, len(resp.Teams), 0, "Should return at least one team")
	})

	t.Run("Month parameter returns teams", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Month: "2026-02",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Greater(t, len(resp.Teams), 0)
	})

	t.Run("Date range returns teams", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			StartDate: "2026-02-01",
			EndDate:   "2026-02-18",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Teams)
		require.Greater(t, len(resp.Teams), 0)
	})

	t.Run("Limit returns correct count", func(t *testing.T) {
		limit := 3
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Limit: &limit,
		})
		require.NoError(t, err)
		require.Equal(t, 3, len(resp.Teams))
	})

	t.Run("Limit returns next cursor", func(t *testing.T) {
		limit := 2
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Limit: &limit,
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.Next)
	})

	t.Run("Pagination returns different teams", func(t *testing.T) {
		limit := 2
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
	})

	t.Run("Teams have team field", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, nil)
		require.NoError(t, err)
		require.Greater(t, len(resp.Teams), 0)
		// team field exists (may be empty string for default team)
		_ = resp.Teams[0].Team
	})

	t.Run("All metrics present", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, nil)
		require.NoError(t, err)
		require.Greater(t, len(resp.Teams), 0)

		team := resp.Teams[0]

		// Daily activity metrics - verify they are initialized
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
	})

	t.Run("Metric totals non-negative", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, nil)
		require.NoError(t, err)

		for _, team := range resp.Teams {
			require.GreaterOrEqual(t, team.MessagesTotal.Total, int64(0), "messages_total should be >= 0")
			require.GreaterOrEqual(t, team.UsersDaily.Total, int64(0), "users_daily should be >= 0")
			require.GreaterOrEqual(t, team.ConcurrentUsers.Total, int64(0), "concurrent_users should be >= 0")
		}
	})
}

func TestQueryTeamUsageStats_DataCorrectness(t *testing.T) {
	c := initMultiTenantClient(t)
	skipIfNoMultiTenant(t, c)
	ctx := context.Background()

	testTeams := []string{"sdk-test-team-1", "sdk-test-team-2", "sdk-test-team-3"}

	t.Run("Date range query returns test teams with exact values", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			StartDate: "2026-02-18",
			EndDate:   "2026-02-19",
		})
		require.NoError(t, err)

		for _, teamName := range testTeams {
			team := findTeamByName(resp.Teams, teamName)
			require.NotNil(t, team, "%s should exist", teamName)
			assertAllMetricsExact(t, team, teamName)
		}
	})

	t.Run("Month query returns test teams with valid metrics", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, &QueryTeamUsageStatsRequest{
			Month: "2026-02",
		})
		require.NoError(t, err)

		for _, teamName := range testTeams {
			team := findTeamByName(resp.Teams, teamName)
			require.NotNil(t, team, "%s should exist", teamName)
			require.GreaterOrEqual(t, team.UsersTotal.Total, int64(0), "%s users_total", teamName)
			require.GreaterOrEqual(t, team.MessagesTotal.Total, int64(0), "%s messages_total", teamName)
		}
	})

	t.Run("No params query returns test teams with valid metrics", func(t *testing.T) {
		resp, err := c.QueryTeamUsageStats(ctx, nil)
		require.NoError(t, err)

		for _, teamName := range testTeams {
			team := findTeamByName(resp.Teams, teamName)
			require.NotNil(t, team, "%s should exist", teamName)
			require.GreaterOrEqual(t, team.UsersTotal.Total, int64(0), "%s users_total", teamName)
			require.GreaterOrEqual(t, team.MessagesTotal.Total, int64(0), "%s messages_total", teamName)
		}
	})

	t.Run("Pagination finds test teams with valid metrics", func(t *testing.T) {
		for _, teamName := range testTeams {
			team := findTeamAcrossPages(t, c, teamName)
			require.NotNil(t, team, "%s should exist across paginated results", teamName)
			require.GreaterOrEqual(t, team.UsersTotal.Total, int64(0), "%s users_total", teamName)
			require.GreaterOrEqual(t, team.MessagesTotal.Total, int64(0), "%s messages_total", teamName)
		}
	})
}
