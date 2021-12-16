package stream_chat // nolint: golint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissions_RoleEndpoints(t *testing.T) {
	c := initClient(t)
	p := c.Permissions()
	ctx := context.Background()
	roleName := randomString(12)

	_, err := p.CreateRole(ctx, roleName)
	require.NoError(t, err)
	_, _ = p.DeleteRole(ctx, roleName)
	// Unfortunately the API is too slow to create roles
	// and we don't want to wait > 10 seconds.
	// So we swallow potential errors here until that's fixed.
	// Plus we add a cleanup as well.

	roles, err := p.ListRoles(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, roles)

	t.Cleanup(func() {
		resp, _ := p.ListRoles(ctx)
		for _, role := range resp.Roles {
			if role.Custom {
				_, _ = p.DeleteRole(ctx, role.Name)
			}
		}
	})
}

func TestPermissions_PermissionEndpoints(t *testing.T) {
	c := initClient(t)
	p := c.Permissions()
	ctx := context.Background()
	permName := randomString(12)

	_, err := p.CreatePermission(ctx, &Permission{
		ID:          permName,
		Name:        permName,
		Action:      "DeleteChannel",
		Description: "integration test",
		Condition: map[string]interface{}{
			"$subject.magic_custom_field": map[string]string{"$eq": "true"},
		},
	})
	require.NoError(t, err)

	perms, err := p.ListPermissions(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, perms)

	resp, err := p.GetPermission(ctx, "create-channel")
	require.NoError(t, err)

	perm := resp.Permission
	assert.Equal(t, "create-channel", perm.ID)
	assert.False(t, perm.Custom)
	assert.NotEmpty(t, perm.Condition)

	t.Cleanup(func() {
		resp, _ := p.ListPermissions(ctx)
		for _, perm := range resp.Permissions {
			if perm.Description == "integration test" {
				_, _ = p.DeletePermission(ctx, perm.ID)
			}
		}
	})
}
