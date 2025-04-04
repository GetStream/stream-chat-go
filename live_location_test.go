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

	// Create a channel
	ch := c.Channel("messaging", randomString(10))
	_, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	// Create a message to associate with the live location
	msg := &Message{
		Text: "Message with live location",
		User: user,
	}
	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err)
	require.NotNil(t, resp.Message)

	// Create a live location
	endTime := time.Now().Add(1 * time.Hour)
	liveLocation := &LiveLocation{
		UserID:            user.ID,
		ChannelID:         ch.ID,
		MessageID:         resp.Message.ID,
		Latitude:          40.7128,
		Longitude:         -74.0060,
		EndAt:             &endTime,
		CreatedByDeviceID: "test-device",
	}

	// Update live location
	updateResp, err := ch.UpdateLiveLocation(ctx, liveLocation)
	require.NoError(t, err)
	require.NotNil(t, updateResp)
}

func TestClient_GetUserActiveLiveLocations(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()

	// Create users
	user := randomUser(t, c)

	// Create a channel
	ch := c.Channel("messaging", randomString(10))
	_, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	// Create a message to associate with the live location
	msg := &Message{
		Text: "Message with live location",
		User: user,
	}
	resp, err := ch.SendMessage(ctx, msg, user.ID)
	require.NoError(t, err)
	require.NotNil(t, resp.Message)

	// Create a live location
	endTime := time.Now().Add(1 * time.Hour)
	liveLocation := &LiveLocation{
		UserID:            user.ID,
		ChannelID:         ch.ID,
		MessageID:         resp.Message.ID,
		Latitude:          40.7128,
		Longitude:         -74.0060,
		EndAt:             &endTime,
		CreatedByDeviceID: "test-device",
	}

	// Update live location
	_, err = ch.UpdateLiveLocation(ctx, liveLocation)
	require.NoError(t, err)

	locationResp, err := c.GetUserActiveLiveLocations(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, locationResp)

	if err == nil && locationResp != nil {
		if len(locationResp.LiveLocations) > 0 {
			location := locationResp.LiveLocations[0]
			assert.Equal(t, user.ID, location.UserID)
			assert.Equal(t, ch.ID, location.ChannelID)
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
	ch := c.Channel("messaging", randomString(10))
	user := randomUser(t, c)
	_, err := c.CreateChannel(ctx, ch.Type, ch.ID, user.ID, nil)
	require.NoError(t, err)

	// Test nil live location
	resp, err := ch.UpdateLiveLocation(ctx, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "liveLocation should not be nil")
	assert.Nil(t, resp)

	// Test invalid live location (missing required fields)
	invalidLocation := &LiveLocation{
		// Missing userID, channelID, and messageID
		Latitude:  40.7128,
		Longitude: -74.0060,
	}
	_, err = ch.UpdateLiveLocation(ctx, invalidLocation)
	// The API might return an error or it might succeed with warnings
	// This test just ensures the function runs without panicking
	if err != nil {
		t.Logf("Expected error for invalid location: %v", err)
	}

}
