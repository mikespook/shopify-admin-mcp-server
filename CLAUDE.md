# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Purpose

Shopify Admin MCP Server — a pure Go MCP (Model Context Protocol) server that connects to a Shopify store via Admin API access token, enabling LLMs to read and update products, collections, menus, and pages.

## Commands

```makefile
make dev              # Run server in development mode (version = "dev")
make build            # Build binary to build/ (version injected from .version)
make test             # Run Go tests
make test-coverage    # Generate test coverage report
make lint             # Run golangci-lint
make fmt              # Run gofmt
make tidy             # go mod tidy
make docker-build-images  # Build Docker image (auto-versioned from .version)
make docker-push-images   # Push to Docker Hub
```

To run a single Go test:
```sh
go test -C backend ./internal/<package>/... -run TestFunctionName
```

## Architecture

### Structure
```
shopify-admin-mcp-server/
├── backend/
│   ├── cmd/shopify-admin-mcp-server/main.go   # Entry point, flag/env/config resolution
│   ├── internal/
│   │   ├── mcp/         # MCP JSON-RPC 2.0 protocol (server.go, tools.go, types.go)
│   │   └── shopify/     # Shopify Admin API client + tool implementations
│   ├── go.mod
│   └── go.sum
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yaml
├── utils/
│   ├── docker-build-images.sh
│   └── docker-push.sh
├── .version             # Semver string (e.g. v1.0.0) — read by Makefile at build time
├── build/
└── Makefile
```

### Key Design Decisions
- **No external MCP SDK** — custom JSON-RPC 2.0, protocol version `2025-11-25`
- **No web framework** — standard `net/http` with `github.com/mikespook/possum` middleware
- **Shopify SDK** — `github.com/bold-commerce/go-shopify/v4` (package `goshopify`)
- **No CGO** — `CGO_ENABLED=0`, no SQLite
- **Go module:** `github.com/mikespook/shopify-admin-mcp-server`

### MCP Protocol
- HTTP endpoint: `POST /mcp`
- Stdio: run with `-stdio` flag (JSON-RPC over stdin/stdout, logs to stderr)
- Methods: `initialize`, `tools/list`, `tools/call`, `notifications/initialized`
- Tools defined as `ToolDefinition` structs with JSON schema; returned by static `GetTools()`

### Shopify Tools
14 tools across 4 resources:
- **Products:** `shopify_list_products`, `shopify_get_product`, `shopify_create_product`, `shopify_update_product`
- **Collections:** `shopify_list_collections`, `shopify_get_collection`, `shopify_create_collection`, `shopify_update_collection`
- **Pages:** `shopify_list_pages`, `shopify_get_page`, `shopify_create_page`, `shopify_update_page`
- **Menus:** `shopify_list_menus`, `shopify_get_menu` (raw REST — not in go-shopify SDK)

### HTTP Routes
```
GET  /health   — liveness check
POST /mcp      — MCP JSON-RPC endpoint (logging + CORS via possum.Chain)
```

### Configuration (priority order: flag > env > config file > interactive prompt)
| Flag / Env var                   | Default           | Purpose                        |
| -------------------------------- | ----------------- | ------------------------------ |
| `-store` / `SHOPIFY_STORE`       | —                 | Shopify store domain           |
| `-token` / `SHOPIFY_TOKEN`       | —                 | Shopify Admin API access token |
| `-addr`                          | `localhost:9093`  | HTTP listen address            |
| `-config`                        | `./config.yaml`   | YAML config file path          |
| `-stdio`                         | false             | Use stdio transport (no HTTP)  |
| `-version`                       | —                 | Print version and exit         |

Config file (`config.yaml`) is auto-created on first run if it does not exist. `-stdio` is never written to the config file.

### Versioning
- Version string lives in [.version](.version) (e.g. `v1.0.0`)
- Injected at build time via `-ldflags "-X main.version=..."`
- Unbuilt / `go run` binaries report `dev`
- Reported in MCP `initialize` response (`serverInfo.version`)

### Docker
Single-stage Go build on `golang:1.26-alpine` (`CGO_ENABLED=0`) → `alpine:3` runtime.
Non-root user `appuser` (UID 1000), `WORKDIR /home/appuser`, exposes port `9093`, health check via `wget /health`.
Entrypoint overrides `-addr` to `0.0.0.0:9093` (binary default is `localhost:9093` for safe local use).
`docker-compose.yaml` mounts a named volume at `/home/appuser` to persist `config.yaml` across restarts.
In non-TTY environments (Docker, CI) the interactive credential prompt is skipped automatically via `term.IsTerminal`; the server exits with a clear error if credentials are missing.
