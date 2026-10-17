package morphik_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-morphik"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAuthorization = "Bearer EXAMPLE-DO-NOT-USE"

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *morphik.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := morphik.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", testAuthorization)

			return nil
		},
	)

	client, err := morphik.NewClientWithResponses(
		srv.URL,
		auth,
		morphik.WithHTTPClient(srv.Client()),
	)
	require.NoError(t, err)

	return client
}

func writeJSON(t *testing.T, w http.ResponseWriter, body string) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(body))
	assert.NoError(t, err)
}

func TestGetModels(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/models", r.URL.Path)
		assert.Equal(t, testAuthorization, r.Header.Get("Authorization"))

		writeJSON(t, w, `{
			"chat_models":[{"name":"model-A"}],
			"embedding_models":[],
			"default_models":{"chat":"model-A","embedding":null},
			"providers":["provider-A"]
		}`)
	})

	res, err := client.GetAvailableModelsModelsGetWithResponse(
		t.Context(),
		&morphik.GetAvailableModelsModelsGetParams{},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.ChatModels, 1)
	assert.Equal(t, "model-A", res.JSON200.ChatModels[0]["name"])
	assert.Equal(t, []string{"provider-A"}, res.JSON200.Providers)
	assert.Equal(t, new("model-A"), res.JSON200.DefaultModels["chat"])
	assert.Nil(t, res.JSON200.DefaultModels["embedding"])
}

func TestIngestText(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/ingest/text", r.URL.Path)
		assert.Equal(t, testAuthorization, r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body morphik.IngestTextRequest

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Example document", body.Content)
		assert.Nil(t, body.Filename)

		writeJSON(t, w, `{
			"external_id":"doc-A","content_type":"text/plain",
			"filename":null,"summary_version":null,
			"system_metadata":{"status":"processing"}
		}`)
	})

	res, err := client.IngestTextIngestTextPostWithResponse(
		t.Context(),
		&morphik.IngestTextIngestTextPostParams{},
		morphik.IngestTextRequest{Content: "Example document"},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.ExternalId)
	assert.Equal(t, "doc-A", *res.JSON200.ExternalId)
	assert.Equal(t, "text/plain", res.JSON200.ContentType)
	assert.Nil(t, res.JSON200.Filename)
	assert.Nil(t, res.JSON200.SummaryVersion)
}

func TestRetrieveChunks(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/retrieve/chunks", r.URL.Path)
		assert.Equal(t, testAuthorization, r.Header.Get("Authorization"))

		var body morphik.RetrieveRequest

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, new("Find a document"), body.Query)
		assert.Equal(t, new(morphik.OutputFormatText), body.OutputFormat)
		assert.Equal(t, new(2), body.K)
		assert.Equal(t, new(0.123456789012345), body.MinScore)

		writeJSON(t, w, `[{
			"content":"Example chunk","score":0.987654321012345,
			"document_id":"doc-A","chunk_number":0,"metadata":{},
			"content_type":"text/plain","filename":null
		}]`)
	})

	res, err := client.RetrieveChunksRetrieveChunksPostWithResponse(
		t.Context(),
		&morphik.RetrieveChunksRetrieveChunksPostParams{},
		morphik.RetrieveRequest{
			Query:        new("Find a document"),
			OutputFormat: new(morphik.OutputFormatText),
			K:            new(2),
			MinScore:     new(0.123456789012345),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.Len(t, *res.JSON200, 1)
	chunk := (*res.JSON200)[0]
	assert.Equal(t, "doc-A", chunk.DocumentId)
	assert.InDelta(t, 0.987654321012345, chunk.Score, 0)
	assert.Nil(t, chunk.Filename)
}

func TestUnauthorized(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, testAuthorization, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte(`{"detail":"invalid credentials"}`))
		assert.NoError(t, err)
	})

	res, err := client.GetAvailableModelsModelsGetWithResponse(
		t.Context(),
		&morphik.GetAvailableModelsModelsGetParams{},
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Nil(t, res.JSON200)
	assert.JSONEq(t, `{"detail":"invalid credentials"}`, string(res.Body))
}

func TestMalformedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(t, w, `{"chat_models":`)
	})

	_, err := client.GetAvailableModelsModelsGetWithResponse(
		t.Context(),
		&morphik.GetAvailableModelsModelsGetParams{},
	)

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
}

func TestCanceledRequest(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(_ http.ResponseWriter, _ *http.Request) {
		t.Error("a canceled request reached the server")
	})
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := client.GetAvailableModelsModelsGetWithResponse(
		ctx,
		&morphik.GetAvailableModelsModelsGetParams{},
	)
	require.ErrorIs(t, err, context.Canceled)
}

func TestGetDocument(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/documents/doc%20A", r.URL.EscapedPath())
		assert.Equal(t, testAuthorization, r.Header.Get("Authorization"))

		writeJSON(t, w, `{
			"external_id":"doc A","content_type":"text/plain",
			"filename":"example.txt","summary_version":2
		}`)
	})

	res, err := client.GetDocumentDocumentsDocumentIdGetWithResponse(
		t.Context(),
		"doc A",
		&morphik.GetDocumentDocumentsDocumentIdGetParams{},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	assert.Equal(t, new("example.txt"), res.JSON200.Filename)
	assert.Equal(t, new(2), res.JSON200.SummaryVersion)
}

func TestValidationError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err := w.Write([]byte(`{
			"detail":[{"loc":["body",0],"msg":"invalid","type":"value"}]
		}`))
		assert.NoError(t, err)
	})

	res, err := client.GetAvailableModelsModelsGetWithResponse(
		t.Context(),
		&morphik.GetAvailableModelsModelsGetParams{},
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode())
	assert.Nil(t, res.JSON200)
	require.NotNil(t, res.JSON422)
	require.NotNil(t, res.JSON422.Detail)
	require.Len(t, *res.JSON422.Detail, 1)
	assert.Equal(t, "invalid", (*res.JSON422.Detail)[0].Msg)
}
