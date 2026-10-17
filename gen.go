// Package morphik is a Go client for the Morphik API,
// generated with oapi-codegen from the OpenAPI document Morphik publishes.
package morphik

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.yaml
