package stream_chat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_PathSegmentsAreEscaped(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c, err := NewClient("key", "secret")
	require.NoError(t, err)
	c.BaseURL = srv.URL

	const id = "a/b?c#d"
	const escaped = "a%2Fb%3Fc%23d"
	ctx := context.Background()

	tests := []struct {
		name string
		call func() error
		want string
	}{
		{"GetBlocklist", func() error { _, err := c.GetBlocklist(ctx, id); return err }, "/blocklists/" + escaped},
		{"UpdateBlocklist", func() error { _, err := c.UpdateBlocklist(ctx, id, nil); return err }, "/blocklists/" + escaped},
		{"DeleteBlocklist", func() error { _, err := c.DeleteBlocklist(ctx, id); return err }, "/blocklists/" + escaped},
		{"DeletePushProvider", func() error { _, err := c.DeletePushProvider(ctx, "fire/base", id); return err }, "/push_providers/fire%2Fbase/" + escaped},
		{"DeleteRole", func() error { _, err := c.Permissions().DeleteRole(ctx, id); return err }, "/roles/" + escaped},
		{"GetPermission", func() error { _, err := c.Permissions().GetPermission(ctx, id); return err }, "/permissions/" + escaped},
		{"UpdatePermission", func() error { _, err := c.Permissions().UpdatePermission(ctx, id, &Permission{}); return err }, "/permissions/" + escaped},
		{"DeletePermission", func() error { _, err := c.Permissions().DeletePermission(ctx, id); return err }, "/permissions/" + escaped},
		{"ReviewFlagReport", func() error {
			_, err := c.ReviewFlagReport(ctx, id, &ReviewFlagReportRequest{})
			return err
		}, "/moderation/reports/" + escaped},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath = ""
			require.NoError(t, tt.call())
			require.Equal(t, tt.want, gotPath)
		})
	}
}
