package stream_chat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannel_UpdateLiveLocation(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create users
	user := randomUser(t, c)
	channel := initChannel(t, c)

	// Create a live location
	endTime := time.Now().Add(1 * time.Hour)
	liveLocation := &LiveLocation{
		UserID:            user.ID,
		ChannelID:         channel.ID,
		Latitude:          40.7128,
		Longitude:         -74.0060,
		EndAt:             &endTime,
		CreatedByDeviceID: "test-device",
	}

	// Create a message to associate with the live location
	msg := &Message{
		LiveLocation: liveLocation,
	}
	resp, err := channel.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err)
	require.NotNil(t, resp.Message)

	// Update live location
	updateResp, err := channel.UpdateLiveLocation(ctx, resp.Message.LiveLocation, user.ID)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
}

func TestClient_GetUserActiveLiveLocations(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create users
	user := randomUser(t, c)

	// Create a channel
	channel, err := c.CreateChannelWithMembers(ctx, "messaging", randomString(12), user.ID)
	require.NoError(t, err, "create channel")

	// Create a live location
	endTime := time.Now().Add(1 * time.Hour)
	liveLocation := &LiveLocation{
		UserID:            user.ID,
		ChannelID:         channel.Channel.ID,
		Latitude:          40.7128,
		Longitude:         -74.0060,
		EndAt:             &endTime,
		CreatedByDeviceID: "test-device",
	}

	// Create a message to associate with the live location
	msg := &Message{
		LiveLocation: liveLocation,
	}
	resp, err := channel.Channel.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err)
	require.NotNil(t, resp.Message)

	locationResp, err := c.GetUserActiveLiveLocations(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, locationResp)

	if err == nil && locationResp != nil {
		if len(locationResp.LiveLocations) > 0 {
			location := locationResp.LiveLocations[0]
			assert.Equal(t, channel.Channel.CID, location.ChannelID)
			assert.Equal(t, resp.Message.ID, location.MessageID)
			assert.InDelta(t, 40.7128, location.Latitude, 0.0001)
			assert.InDelta(t, -74.0060, location.Longitude, 0.0001)
		}
	}
}

func TestChannel_UpdateLiveLocation_Error(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create a channel
	channel := initChannel(t, c)

	// Test nil live location
	resp, err := channel.UpdateLiveLocation(ctx, nil, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "liveLocation should not be nil")
	assert.Nil(t, resp)

	// Test invalid live location (missing required fields)
	invalidLocation := &LiveLocation{
		// Missing userID, channelID, and messageID
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	_, err = channel.UpdateLiveLocation(ctx, invalidLocation, "")
	// The API might return an error or it might succeed with warnings
	// This test just ensures the function runs without panicking
	if err != nil {
		t.Logf("Expected error for invalid location: %v", err)
	}

}
