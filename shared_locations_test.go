package stream_chat_test

import (
	"context"
	"testing"
	"time"

	stream_chat "github.com/GetStream/stream-chat-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestChannelWithSharedLocations(t *testing.T, c *stream_chat.Client) (*stream_chat.Channel, *stream_chat.User) {
	t.Helper()

	ctx := context.Background()
	userID := "test-user-" + randomString(10)
	channelID := "test-channel-" + randomString(10)

	user := &stream_chat.User{ID: userID}
	resp, err := c.UpsertUser(ctx, user)
	require.NoError(t, err)
	user = resp.User

	channelResp, err := c.CreateChannelWithMembers(ctx, "messaging", channelID, userID)
	require.NoError(t, err)

	_, err = channelResp.Channel.PartialUpdate(ctx, stream_chat.PartialUpdate{
		Set: map[string]interface{}{
			"config_overrides": map[string]interface{}{
				"shared_locations": true,
			},
		},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = c.DeleteUsers(ctx, []string{userID}, stream_chat.DeleteUserOptions{
			User:          stream_chat.HardDelete,
			Messages:      stream_chat.HardDelete,
			Conversations: stream_chat.HardDelete,
		})
		_, _ = c.DeleteChannels(ctx, []string{channelResp.Channel.CID}, true)
	})

	return channelResp.Channel, user
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestClient_LiveLocation(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	settings := stream_chat.NewAppSettings().SetSharedLocationsEnabled(true)

	_, err := c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)

	// Create a user
	channel, user := createTestChannelWithSharedLocations(t, c)

	// Create a shared location
	location := &stream_chat.SharedLocation{
		UserID:            user.ID,
		MessageID:         randomString(10),
		Longitude:         -122.4194,
		Latitude:          37.7749,
		EndAt:             timePtr(time.Now().Add(1 * time.Hour)),
		CreatedByDeviceID: "test-device",
	}

	messageResp, err := channel.SendMessage(ctx, &stream_chat.Message{
		SharedLocation: location,
	}, user.ID)
	require.NoError(t, err)
	message := messageResp.Message

	newLocation := &stream_chat.SharedLocation{
		UserID:            user.ID,
		MessageID:         message.ID,
		Longitude:         -122.4194,
		Latitude:          38.999,
		CreatedByDeviceID: "test-device",
	}

	// Update the location
	updateResp1, err := c.UpdateUserActiveLocation(ctx, user.ID, newLocation)
	require.NoError(t, err, "UpdateUserActiveLocation should not return an error")
	require.NotNil(t, updateResp1)
	assert.Equal(t, newLocation.Latitude, updateResp1.Latitude)
	assert.Equal(t, newLocation.Longitude, updateResp1.Longitude)

	// Get active live locations
	getResp, err := c.GetUserActiveLocations(ctx, user.ID)
	require.NoError(t, err, "GetUserActiveLocations should not return an error")
	require.NotNil(t, getResp)
	require.NotEmpty(t, getResp.ActiveLiveLocations, "Should have active live locations")

	// Verify the location data
	found := false
	for _, loc := range getResp.ActiveLiveLocations {
		if loc.MessageID == location.MessageID {
			found = true
			assert.Equal(t, location.Latitude, loc.Latitude)
			assert.Equal(t, location.Longitude, loc.Longitude)
			assert.Equal(t, location.UserID, loc.UserID)
			assert.Equal(t, location.ChannelCID, loc.ChannelCID)
			break
		}
	}
	assert.True(t, found, "Should find the updated location")

	// Update the location with new coordinates
	location.Latitude = 37.7833
	location.Longitude = -122.4167

	updateResp2, err := c.UpdateUserActiveLocation(ctx, user.ID, location)
	require.NoError(t, err, "UpdateUserActiveLocation should not return an error")
	require.NotNil(t, updateResp2)
}
