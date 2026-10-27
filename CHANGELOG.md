# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.2, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- The issue-pull workflow runs at hours derived from the repository name instead of the same four hours as every other repository.

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.0, 2026-10-17

- Added `morphik`, a REST client generated with oapi-codegen from Morphik's OpenAPI document (title FastAPI, spec version 0.1.0, OpenAPI 3.1.0). It covers all 67 operations across 58 paths: ingestion, documents, retrieval, queries, models, cloud apps and connectors. Source: [Stainless-hosted spec](https://storage.googleapis.com/stainless-sdk-openapi-specs/morphik%2Fmorphik-2d9291c3821d4053d0cb0f8c0a132c4fccdbc2ce5eca47bd28eb1a4af7c1f123.yml), discovered through the vendor SDK's `.stats.yml`.
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` resolves the current Stainless spec URL from `morphik-org/morphik-ts`, downloads the spec, converts nullable unions and exclusive bounds to OpenAPI 3.0 form, and regenerates the client. The spec is not committed.
- Added the `go-morphik` agent skill and Claude Code and Codex plugin manifests.
