package stream_chat

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestRateLimit asserts that rate limit headers are correctly decoded into the error object.
// We use DeleteUsers endpoint, it requires a very low number of requests (6/min).
func TestRateLimit(t *testing.T) {
	c := initClient(t)

	users := make([]*User, 0, 8)
	for i := 0; i < 8; i++ {
		users = append(users, randomUser(t, c))
	}

	for _, u := range users {
		_, err := c.DeleteUsers(context.Background(), []string{u.ID}, DeleteUserOptions{
			User:     SoftDelete,
			Messages: HardDelete,
		})
		if err != nil {
			var apiErr Error
			ok := errors.As(err, &apiErr)
			require.True(t, ok)
			require.Equal(t, http.StatusTooManyRequests, apiErr.StatusCode)
			require.NotZero(t, apiErr.RateLimit.Limit)
			require.NotZero(t, apiErr.RateLimit.Reset)
			require.EqualValues(t, 0, apiErr.RateLimit.Remaining)
			return
		}
	}
}

// TestContextExceeded asserts that the context error is correctly returned.
func TestContextExceeded(t *testing.T) {
	c := initClient(t)
	user := randomUser(t, c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := c.DeleteUsers(ctx, []string{user.ID}, DeleteUserOptions{
		User:     SoftDelete,
		Messages: HardDelete,
	})
	require.Error(t, err)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
