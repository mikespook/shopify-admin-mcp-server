# Shopify Admin MCP Server

A lightweight [Model Context Protocol](https://modelcontextprotocol.io) server that connects an LLM (such as Claude) to your Shopify store via the Admin API. It lets AI assistants read and update your store's products, collections, pages, and navigation menus.

## Prerequisites

- A Shopify store with a **Custom App** that has an Admin API access token
- Go 1.22+ (to build from source), or use the pre-built binary / Docker image

### Required API scopes

When creating your Custom App in the Shopify Admin, enable these scopes:

| Scope | Used for |
|---|---|
| `read_products`, `write_products` | Products |
| `read_product_listings`, `write_product_listings` | Product listings |
| `read_content`, `write_content` | Pages |
| `read_online_store_navigation`, `write_online_store_navigation` | Menus |

## Installation

### Build from source

```sh
git clone https://github.com/mikespook/shopify-admin-mcp-server
cd shopify-admin-mcp-server
make build
# binary: ./build/shopify-admin-mcp-server
```

### Docker

```sh
docker pull mikespook/shopify-admin-mcp-server
```

## Usage

### Flags

| Flag | Env var | Default | Description |
|---|---|---|---|
| `-store` | `SHOPIFY_STORE` | — | Shopify store domain, e.g. `myshop.myshopify.com` |
| `-token` | `SHOPIFY_TOKEN` | — | Admin API access token (`shpca_...` or `shpat_...`) |
| `-addr` | — | `localhost:9093` | HTTP listen address |
| `-config` | — | `./config.yaml` | Path to YAML config file |
| `-stdio` | — | false | Use stdio transport instead of HTTP |
| `-version` | — | — | Print version and exit |

Priority order: **flag > env var > config file > interactive prompt**

### Credential resolution

On first run with an interactive terminal, if credentials are not supplied via flags or env vars, the server will prompt and save them to `config.yaml` (mode `0600`). Subsequent runs load from that file automatically.

In non-interactive environments (Docker, CI, stdio mode) the prompt is skipped — the server exits immediately with an error if credentials are missing.

### HTTP mode

```sh
./shopify-admin-mcp-server -store myshop.myshopify.com -token shpca_xxx
# Listening on localhost:9093
```

Or with env vars:

```sh
export SHOPIFY_STORE=myshop.myshopify.com
export SHOPIFY_TOKEN=shpca_xxx
./shopify-admin-mcp-server
```

Or with a config file:

```yaml
# config.yaml
store: myshop.myshopify.com
token: shpca_xxx
addr: localhost:9093
```

```sh
./shopify-admin-mcp-server -config ./config.yaml
```

### Stdio mode

Stdio mode is used by MCP clients (such as Claude Desktop or Claude Code) that launch the server as a subprocess. Logs are written to stderr so they do not interfere with the JSON-RPC stream on stdout.

```sh
./shopify-admin-mcp-server -stdio -store myshop.myshopify.com -token shpca_xxx
```

## Connecting to Claude

### Claude Code (`.mcp.json`)

Add to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "shopify-admin-stdio": {
      "type": "stdio",
      "command": "shopify-admin-mcp-server",
      "args": ["-stdio", "-store", "myshop.myshopify.com", "-token", "shpca_xxx"]
    }
  }
}
```

Or point at the running HTTP server:

```json
{
  "mcpServers": {
    "shopify-admin-http": {
      "type": "http",
      "url": "http://localhost:9093/mcp"
    }
  }
}
```

### Claude Desktop (`claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "shopify": {
      "command": "shopify-admin-mcp-server",
      "args": ["-stdio", "-store", "myshop.myshopify.com", "-token", "shpca_xxx"]
    }
  }
}
```

## Available tools

Once connected, the LLM can use these tools:

| Tool | Description |
|---|---|
| `shopify_list_products` | List all products |
| `shopify_get_product` | Get a product by ID |
| `shopify_create_product` | Create a new product |
| `shopify_update_product` | Update an existing product |
| `shopify_list_collections` | List all custom collections |
| `shopify_get_collection` | Get a collection by ID |
| `shopify_create_collection` | Create a new collection |
| `shopify_update_collection` | Update an existing collection |
| `shopify_list_pages` | List all pages |
| `shopify_get_page` | Get a page by ID |
| `shopify_create_page` | Create a new page |
| `shopify_update_page` | Update an existing page |
| `shopify_list_menus` | List all navigation menus |
| `shopify_get_menu` | Get a menu by handle |

## Docker

### docker run

```sh
docker run --rm \
  -e SHOPIFY_STORE=myshop.myshopify.com \
  -e SHOPIFY_TOKEN=shpca_xxx \
  -p 9093:9093 \
  mikespook/shopify-admin-mcp-server
```

The container listens on `0.0.0.0:9093` so the published port is reachable from the host.

### docker compose

```sh
cd docker
SHOPIFY_STORE=myshop.myshopify.com SHOPIFY_TOKEN=shpca_xxx docker compose up -d
```

Or create a `.env` file next to `docker-compose.yaml`:

```sh
SHOPIFY_STORE=myshop.myshopify.com
SHOPIFY_TOKEN=shpca_xxx
```

Then just:

```sh
docker compose -f docker/docker-compose.yaml up -d
```

A named volume (`config`) is mounted at `/home/appuser` inside the container, persisting `config.yaml` across restarts. On subsequent starts you can omit the env vars — the saved config will be used.

## Health check

```sh
curl http://localhost:9093/health
# ok
```

## License

MIT
