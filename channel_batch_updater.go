package stream_chat

import (
	"context"
)

// ChannelBatchUpdater provides convenience methods for batch channel operations.
type ChannelBatchUpdater struct {
	client *Client
}

// ChannelBatchMemberRequest represents a member in batch operations.
type ChannelBatchMemberRequest struct {
	UserID      string `json:"user_id"`
	ChannelRole string `json:"channel_role,omitempty"`
}

// AddMembers adds members to channels matching the filter.
func (u *ChannelBatchUpdater) AddMembers(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationAddMembers,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// RemoveMembers removes members from channels matching the filter.
func (u *ChannelBatchUpdater) RemoveMembers(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationRemoveMembers,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// InviteMembers invites members to channels matching the filter.
func (u *ChannelBatchUpdater) InviteMembers(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationInvites,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// AddModerators adds moderators to channels matching the filter.
func (u *ChannelBatchUpdater) AddModerators(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationAddModerators,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// DemoteModerators removes moderator role from members in channels matching the filter.
func (u *ChannelBatchUpdater) DemoteModerators(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationDemoteModerators,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// AssignRoles assigns roles to members in channels matching the filter.
func (u *ChannelBatchUpdater) AssignRoles(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationAssignRoles,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// Hide hides channels matching the filter for the specified members.
func (u *ChannelBatchUpdater) Hide(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationHide,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// Show shows channels matching the filter for the specified members.
func (u *ChannelBatchUpdater) Show(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationShow,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// Archive archives channels matching the filter for the specified members.
func (u *ChannelBatchUpdater) Archive(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationArchive,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// Unarchive unarchives channels matching the filter for the specified members.
func (u *ChannelBatchUpdater) Unarchive(ctx context.Context, filter ChannelsBatchFilters, members []ChannelBatchMemberRequest) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationUnarchive,
		Filter:    filter,
		Members:   members,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// UpdateData updates data on channels matching the filter.
func (u *ChannelBatchUpdater) UpdateData(ctx context.Context, filter ChannelsBatchFilters, data *ChannelDataUpdate) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation: BatchUpdateOperationUpdateData,
		Filter:    filter,
		Data:      data,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// AddFilterTags adds filter tags to channels matching the filter.
func (u *ChannelBatchUpdater) AddFilterTags(ctx context.Context, filter ChannelsBatchFilters, tags []string) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation:        BatchUpdateOperationAddFilterTags,
		Filter:           filter,
		FilterTagsUpdate: tags,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}

// RemoveFilterTags removes filter tags from channels matching the filter.
func (u *ChannelBatchUpdater) RemoveFilterTags(ctx context.Context, filter ChannelsBatchFilters, tags []string) (*AsyncTaskResponse, error) {
	options := &ChannelsBatchOptions{
		Operation:        BatchUpdateOperationRemoveFilterTags,
		Filter:           filter,
		FilterTagsUpdate: tags,
	}
	return u.client.UpdateChannelsBatch(ctx, options)
}
