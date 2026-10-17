package morphik_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-morphik"
)

// This README example needs real credentials, so it only compiles.
//
//nolint:testableexamples // Requires a Morphik API key and a live server.
func Example() {
	auth := morphik.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization",
				"Bearer "+os.Getenv("MORPHIK_API_KEY"),
			)

			return nil
		},
	)

	client, err := morphik.NewClientWithResponses(
		"https://api.morphik.ai",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.GetAvailableModelsModelsGetWithResponse(
		context.Background(),
		&morphik.GetAvailableModelsModelsGetParams{},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
