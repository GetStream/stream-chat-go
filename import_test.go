package stream_chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportsEndToEnd(t *testing.T) {
	t.Skip("The backend isn't deployed yet.")
	filename := randomString(10) + ".json"
	content := "[]"
	c := initClient(t)
	ctx := context.Background()

	createURLResp, err := c.CreateImportURL(ctx, filename)
	require.NoError(t, err)
	require.NotEmpty(t, createURLResp.Path)
	require.NotEmpty(t, createURLResp.UploadURL)

	_, err = c.CreateImport(ctx, createURLResp.Path, "upsert")
	require.Error(t, err)

	data := strings.NewReader(content)
	r, err := http.NewRequestWithContext(ctx, http.MethodPut, createURLResp.UploadURL, data)
	require.NoError(t, err)

	r.Header.Set("Content-Type", "application/json")
	r.ContentLength = data.Size()
	uploadResp, err := c.HTTP.Do(r)
	require.NoError(t, err)
	uploadResp.Body.Close()

	createResp, err := c.CreateImport(ctx, createURLResp.Path, "upsert")
	require.NoError(t, err)
	require.NotNil(t, createResp.ImportTask.ID)
	require.True(t, strings.HasSuffix(createResp.ImportTask.Path, filename))

	getResp, err := c.GetImport(ctx, createResp.ImportTask.ID)
	require.NoError(t, err)
	require.Equal(t, createResp.ImportTask.ID, getResp.ImportTask.ID)

	listResp, err := c.ListImports(ctx, &ListImportsOptions{Limit: 1, Offset: 0})
	require.NoError(t, err)
	require.NotEmpty(t, listResp.ImportTasks)
}

func TestCreateImportRequest_MergeCustomJSON(t *testing.T) {
	t.Run("without merge_custom option", func(t *testing.T) {
		req := createImportRequest{
			Path: "some/path.json",
			Mode: string(UpsertMode),
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m))

		require.Equal(t, "some/path.json", m["path"])
		require.Equal(t, "upsert", m["mode"])
		require.NotContains(t, m, "merge_custom")
	})

	t.Run("with merge_custom true", func(t *testing.T) {
		req := createImportRequest{
			Path: "some/path.json",
			Mode: string(UpsertMode),
		}
		WithMergeCustom(true)(&req)

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m))

		require.Equal(t, "some/path.json", m["path"])
		require.Equal(t, "upsert", m["mode"])
		require.Equal(t, true, m["merge_custom"])
	})

	t.Run("with merge_custom false", func(t *testing.T) {
		req := createImportRequest{
			Path: "some/path.json",
			Mode: string(InsertMode),
		}
		WithMergeCustom(false)(&req)

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m))

		require.Equal(t, "some/path.json", m["path"])
		require.Equal(t, "insert", m["mode"])
		require.Equal(t, false, m["merge_custom"])
	})
}
