package stream_chat //nolint: golint

import (
	"context"
	"errors"
	"net/http"
	"path"
	"time"
)

type Permission struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Action      string                 `json:"action"`
	Owner       bool                   `json:"owner"`
	SameTeam    bool                   `json:"same_team"`
	Condition   map[string]interface{} `json:"condition"`
	Custom      bool                   `json:"custom"`
	Level       string                 `json:"level"`
}

type getPermissionResponse struct {
	Permission *Permission `json:"permission"`
}

type listPermissionsResponse struct {
	Permissions []*Permission `json:"permissions"`
}

type Role struct {
	Name      string    `json:"name"`
	Custom    bool      `json:"custom"`
	Scopes    []string  `json:"scoped"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type roleResponse struct {
	Roles []*Role `json:"roles"`
}

type PermissionClient struct {
	client *Client
}

// CreateRole creates a new role.
func (p *PermissionClient) CreateRole(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	return p.client.makeRequest(ctx, http.MethodPost, "roles", nil, map[string]interface{}{
		"name": name,
	}, nil)
}

// DeleteRole deletes an existing role by name.
func (p *PermissionClient) DeleteRole(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("name is required")
	}

	uri := path.Join("roles", name)
	return p.client.makeRequest(ctx, http.MethodDelete, uri, nil, nil, nil)
}

// ListRole returns all roles.
func (p *PermissionClient) ListRoles(ctx context.Context) ([]*Role, error) {
	var r roleResponse
	err := p.client.makeRequest(ctx, http.MethodGet, "roles", nil, nil, &r)
	return r.Roles, err
}

// CreatePermission creates a new permission.
func (p *PermissionClient) CreatePermission(ctx context.Context, perm *Permission) error {
	return p.client.makeRequest(ctx, http.MethodPost, "permissions", nil, perm, nil)
}

// GetPermission returns a permission by id.
func (p *PermissionClient) GetPermission(ctx context.Context, id string) (*Permission, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	var perm getPermissionResponse
	uri := path.Join("permissions", id)
	err := p.client.makeRequest(ctx, http.MethodGet, uri, nil, nil, &perm)
	return perm.Permission, err
}

// UpdatePermission updates an existing permission by id. Only custom permissions can be updated.
func (p *PermissionClient) UpdatePermission(ctx context.Context, id string, perm *Permission) error {
	if id == "" {
		return errors.New("id is required")
	}

	uri := path.Join("permissions", id)
	return p.client.makeRequest(ctx, http.MethodPut, uri, nil, perm, nil)
}

// ListPermissions returns all permissions of an app.
func (p *PermissionClient) ListPermissions(ctx context.Context) ([]*Permission, error) {
	var perm listPermissionsResponse
	err := p.client.makeRequest(ctx, http.MethodGet, "permissions", nil, nil, &perm)
	return perm.Permissions, err
}

// DeletePermission deletes a permission by id.
func (p *PermissionClient) DeletePermission(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	uri := path.Join("permissions", id)
	return p.client.makeRequest(ctx, http.MethodDelete, uri, nil, nil, nil)
}
