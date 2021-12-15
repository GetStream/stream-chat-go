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

	require.NoError(t, p.CreateRole(ctx, roleName))
	_ = p.DeleteRole(ctx, roleName)
	// Unfortunately the API is too slow to create roles
	// and we don't want to wait > 10 seconds.
	// So we swallow potential errors here until that's fixed.
	// Plus we add a cleanup as well.

	roles, err := p.ListRoles(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, roles)

	t.Cleanup(func() {
		roles, _ = p.ListRoles(ctx)
		for _, role := range roles {
			if role.Custom {
				_ = p.DeleteRole(ctx, role.Name)
			}
		}
	})
}

func TestPermissions_PermissionEndpoints(t *testing.T) {
	c := initClient(t)
	p := c.Permissions()
	ctx := context.Background()
	permName := randomString(12)

	err := p.CreatePermission(ctx, &Permission{
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

	perm, err := p.GetPermission(ctx, "create-channel")
	require.NoError(t, err)
	assert.Equal(t, "create-channel", perm.ID)
	assert.False(t, perm.Custom)
	assert.NotEmpty(t, perm.Condition)

	t.Cleanup(func() {
		perms, _ = p.ListPermissions(ctx)
		for _, perm := range perms {
			if perm.Description == "integration test" {
				_ = p.DeletePermission(ctx, perm.ID)
			}
		}
	})
}
