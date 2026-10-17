# go-morphik

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-morphik.svg)](https://pkg.go.dev/github.com/psyb0t/go-morphik)
[![CI](https://github.com/psyb0t/go-morphik/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-morphik/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-morphik/badges/version.svg)](https://github.com/psyb0t/go-morphik/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-morphik/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-morphik/badges/importers.svg)](https://github.com/psyb0t/go-morphik/blob/badges/importers.md)

Unofficial Go client for the [Morphik API](https://docs.morphik.ai), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Morphik publishes. I am not affiliated with Morphik.

The package is `morphik`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-morphik
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-morphik"
)

func main() {
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
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

The full API is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-morphik); from inside your module, `go doc github.com/psyb0t/go-morphik` prints the same.

Method names come from the spec's operation IDs, such as `GetAvailableModelsModelsGetWithResponse`, `ListDocsDocumentsPostWithResponse` and `RetrieveDocumentsRetrieveDocsPostWithResponse`. The client covers all 67 operations across 58 paths, including ingestion, documents, retrieval, queries, models, cloud apps and connectors.

## Auth

Morphik's [TypeScript SDK](https://github.com/morphik-org/morphik-ts/blob/main/src/client.ts) uses bearer auth and defaults to `https://api.morphik.ai`. Read your API key from `MORPHIK_API_KEY`, then set `Authorization: Bearer <key>` through `WithRequestEditorFn`, as above. Pass your own deployment URL to `NewClientWithResponses` when using another Morphik instance.

The spec exposes optional `authorization` header parameters rather than a security scheme and does not declare a server URL. The generated client does not load credentials from the environment or add auth automatically.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` reads `openapi_spec_url` from [morphik-org/morphik-ts/.stats.yml](https://raw.githubusercontent.com/morphik-org/morphik-ts/main/.stats.yml), then downloads that Stainless-hosted spec to `api/openapi.yaml`. The content-hashed URL changes when Morphik updates it. The file contains JSON despite its `.yaml` suffix.
- `scripts/patch-spec.sh` writes `api/openapi.fixed.yaml`, converting OpenAPI 3.1 nullable unions and numeric exclusive bounds into OpenAPI 3.0 form for oapi-codegen.
- oapi-codegen runs with `config.yaml`, preserving every operation and mapping `number` to `float64` instead of its default `float32`.

`.github/workflows/regen.yml` does the same every day in CI. When the spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release. Fixes cannot edit the protected paths listed in that workflow.

## Agent integrations

The [skill](.agents/skills/go-morphik) teaches a coding agent how to use go-morphik when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-morphik@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-morphik@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The release workflow publishes the skill to ClawHub:

```bash
openclaw skills install @psyb0t/go-morphik
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from Morphik's published API specification. The specification itself is not redistributed here. Its metadata does not state a license or terms of service; the vendor SDK separately uses Apache-2.0.
