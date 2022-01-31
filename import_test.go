package stream_chat

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportsEndToEnd(t *testing.T) {
	t.Skip("The backend isn't deployed yet.")
	filename := randomString(10)
	c := initClient(t)
	ctx := context.Background()

	createResp, err := c.CreateImport(ctx, filename)
	require.NoError(t, err)
	require.NotNil(t, createResp.ImportTask.ID)
	require.Equal(t, filename, createResp.ImportTask.Filename)
	require.NotEmpty(t, createResp.UploadURL)

	data := strings.NewReader("[]")
	err = c.makeRequest(ctx, http.MethodPut, createResp.UploadURL, nil, data, nil)
	require.NoError(t, err)

	getResp, err := c.GetImport(ctx, createResp.ImportTask.ID)
	require.NoError(t, err)
	require.Equal(t, createResp.ImportTask.ID, getResp.ImportTask.ID)

	listResp, err := c.ListImports(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, listResp.ImportTasks)
}
