package stream_chat

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestChannelWithSharedLocations(t *testing.T, c *Client) (*Channel, *User) {
	t.Helper()

	ctx := context.Background()
	userID := "test-user-" + randomString(10)
	channelID := "test-channel-" + randomString(10)

	user := &User{ID: userID}
	resp, err := c.UpsertUser(ctx, user)
	require.NoError(t, err)
	user = resp.User

	channelResp, err := c.CreateChannelWithMembers(ctx, "messaging", channelID, userID)
	require.NoError(t, err)

	_, err = channelResp.Channel.PartialUpdate(ctx, PartialUpdate{
		Set: map[string]interface{}{
			"config_overrides": map[string]interface{}{
				"shared_locations": true,
			},
		},
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_, _ = c.DeleteUsers(ctx, []string{userID}, DeleteUserOptions{
			User:          HardDelete,
			Messages:      HardDelete,
			Conversations: HardDelete,
		})
		_, _ = c.DeleteChannels(ctx, []string{channelResp.Channel.CID}, true)
	})

	return channelResp.Channel, user
}

func TestClient_LiveLocation(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	settings := NewAppSettings().SetSharedLocationsEnabled(true)

	_, err := c.UpdateAppSettings(ctx, settings)
	require.NoError(t, err)

	// Create a user
	channel, user := createTestChannelWithSharedLocations(t, c)

	longitude := -122.4194
	latitude := 38.999

	// Create a shared location
	location := &SharedLocationRequest{
		Longitude:         &longitude,
		Latitude:          &latitude,
		EndAt:             timePtr(time.Now().Add(1 * time.Hour)),
		CreatedByDeviceID: "test-device",
	}

	messageResp, err := channel.SendMessage(ctx, &Message{
		SharedLocation: location,
		Text:           "Test message for shared location",
	}, user.ID)
	require.NoError(t, err)
	message := messageResp.Message
	fmt.Println(message.SharedLocation)

	newLocation := &SharedLocation{
		MessageID:         message.ID,
		Longitude:         -122.4194,
		Latitude:          38.999,
		EndAt:             timePtr(time.Now().Add(10 * time.Hour)),
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
		if loc.MessageID == messageResp.Message.ID {
			found = true
			assert.Equal(t, *messageResp.Message.SharedLocation.Latitude, loc.Latitude)
			assert.Equal(t, *messageResp.Message.SharedLocation.Longitude, loc.Longitude)
			break
		}
	}
	assert.True(t, found, "Should find the updated location")
}
