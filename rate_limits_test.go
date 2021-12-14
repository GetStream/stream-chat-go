package stream_chat // nolint: golint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_GetRateLimits(t *testing.T) {
	c := initClient(t)

	t.Run("get all limits", func(t *testing.T) {
		limits, err := c.GetRateLimits(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, limits.Android)
		require.NotEmpty(t, limits.Web)
		require.NotEmpty(t, limits.IOS)
		require.NotEmpty(t, limits.ServerSide)
	})

	t.Run("get only a single platform", func(t *testing.T) {
		limits, err := c.GetRateLimits(context.Background(), WithServerSide())
		require.NoError(t, err)
		require.Empty(t, limits.Android)
		require.Empty(t, limits.Web)
		require.Empty(t, limits.IOS)
		require.NotEmpty(t, limits.ServerSide)
	})

	t.Run("get only a few endpoints", func(t *testing.T) {
		limits, err := c.GetRateLimits(context.Background(),
			WithServerSide(),
			WithAndroid(),
			WithEndpoints(
				"GetRateLimits",
				"SendMessage",
			),
		)
		require.NoError(t, err)
		require.Empty(t, limits.Web)
		require.Empty(t, limits.IOS)

		require.NotEmpty(t, limits.Android)
		require.Len(t, limits.Android, 2)
		require.Equal(t, limits.Android["GetRateLimits"].Limit, limits.Android["GetRateLimits"].Remaining)

		require.NotEmpty(t, limits.ServerSide)
		require.Len(t, limits.ServerSide, 2)
		require.Greater(t, limits.ServerSide["GetRateLimits"].Limit, limits.ServerSide["GetRateLimits"].Remaining)
	})
}
