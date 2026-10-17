# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.0, 2026-10-17

- Added `morphik`, a REST client generated with oapi-codegen from Morphik's OpenAPI document (title FastAPI, spec version 0.1.0, OpenAPI 3.1.0). It covers all 67 operations across 58 paths: ingestion, documents, retrieval, queries, models, cloud apps and connectors. Source: [Stainless-hosted spec](https://storage.googleapis.com/stainless-sdk-openapi-specs/morphik%2Fmorphik-2d9291c3821d4053d0cb0f8c0a132c4fccdbc2ce5eca47bd28eb1a4af7c1f123.yml), discovered through the vendor SDK's `.stats.yml`.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` resolves the current Stainless spec URL from `morphik-org/morphik-ts`, downloads the spec, converts nullable unions and exclusive bounds to OpenAPI 3.0 form, and regenerates the client. The spec is not committed.
- Added the `go-morphik` agent skill and Claude Code and Codex plugin manifests.
