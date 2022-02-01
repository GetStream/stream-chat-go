package stream_chat

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_TestBlocklistsEndToEnd(t *testing.T) {
	c := initClient(t)
	ctx := context.Background()
	blocklistName := randomString(10)
	blocklistReq := &BlocklistCreateRequest{BlocklistBase{Name: blocklistName, Words: []string{"test"}}}

	_, err := c.CreateBlocklist(ctx, blocklistReq)
	require.NoError(t, err)

	getResp, err := c.GetBlocklist(ctx, blocklistName)
	require.NoError(t, err)
	require.Equal(t, blocklistName, getResp.Blocklist.Name)
	require.Equal(t, blocklistReq.Words, getResp.Blocklist.Words)

	listResp, err := c.ListBlocklists(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, listResp.Blocklists)

	_, err = c.UpdateBlocklist(ctx, blocklistName, []string{"test2"})
	require.NoError(t, err)

	_, err = c.DeleteBlocklist(ctx, blocklistName)
	require.NoError(t, err)
}
