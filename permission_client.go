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

type Role struct {
	Name      string    `json:"name"`
	Custom    bool      `json:"custom"`
	Scopes    []string  `json:"scoped"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionClient struct {
	client *Client
}

// CreateRole creates a new role.
func (p *PermissionClient) CreateRole(ctx context.Context, name string) (*Response, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	var resp Response
	err := p.client.makeRequest(ctx, http.MethodPost, "roles", nil, map[string]interface{}{
		"name": name,
	}, &resp)
	return &resp, err
}

// DeleteRole deletes an existing role by name.
func (p *PermissionClient) DeleteRole(ctx context.Context, name string) (*Response, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	uri := path.Join("roles", name)

	var resp Response
	err := p.client.makeRequest(ctx, http.MethodDelete, uri, nil, nil, &resp)
	return &resp, err
}

type RolesResponse struct {
	Roles []*Role `json:"roles"`
	Response
}

// ListRole returns all roles.
func (p *PermissionClient) ListRoles(ctx context.Context) (*RolesResponse, error) {
	var r RolesResponse
	err := p.client.makeRequest(ctx, http.MethodGet, "roles", nil, nil, &r)
	return &r, err
}

// CreatePermission creates a new permission.
func (p *PermissionClient) CreatePermission(ctx context.Context, perm *Permission) (*Response, error) {
	var resp Response
	err := p.client.makeRequest(ctx, http.MethodPost, "permissions", nil, perm, &resp)
	return &resp, err
}

type GetPermissionResponse struct {
	Permission *Permission `json:"permission"`
	Response
}

// GetPermission returns a permission by id.
func (p *PermissionClient) GetPermission(ctx context.Context, id string) (*GetPermissionResponse, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	uri := path.Join("permissions", id)

	var perm GetPermissionResponse
	err := p.client.makeRequest(ctx, http.MethodGet, uri, nil, nil, &perm)
	return &perm, err
}

// UpdatePermission updates an existing permission by id. Only custom permissions can be updated.
func (p *PermissionClient) UpdatePermission(ctx context.Context, id string, perm *Permission) (*Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	uri := path.Join("permissions", id)

	var resp Response
	err := p.client.makeRequest(ctx, http.MethodPut, uri, nil, perm, &resp)
	return &resp, err
}

type ListPermissionsResponse struct {
	Permissions []*Permission `json:"permissions"`
	Response
}

// ListPermissions returns all permissions of an app.
func (p *PermissionClient) ListPermissions(ctx context.Context) (*ListPermissionsResponse, error) {
	var perm ListPermissionsResponse
	err := p.client.makeRequest(ctx, http.MethodGet, "permissions", nil, nil, &perm)
	return &perm, err
}

// DeletePermission deletes a permission by id.
func (p *PermissionClient) DeletePermission(ctx context.Context, id string) (*Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	uri := path.Join("permissions", id)

	var resp Response
	err := p.client.makeRequest(ctx, http.MethodDelete, uri, nil, nil, &resp)
	return &resp, err
}
