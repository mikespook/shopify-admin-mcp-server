package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	goshopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/mikespook/shopify-admin-mcp-server/internal/shopify"
)

type Server struct {
	client  *goshopify.Client
	version string
}

func NewServer(client *goshopify.Client, version string) *Server {
	return &Server{client: client, version: version}
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, nil, -32700, "parse error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var result any
	var rpcErr *RPCError

	switch req.Method {
	case "initialize":
		result = InitializeResult{
			ProtocolVersion: "2025-11-25",
			Capabilities:    Capabilities{Tools: map[string]any{}},
			ServerInfo:      ServerInfo{Name: "shopify-admin-mcp-server", Version: s.version},
		}
	case "notifications/initialized":
		w.WriteHeader(http.StatusNoContent)
		return
	case "tools/list":
		result = ToolsListResult{Tools: GetTools()}
	case "tools/call":
		result, rpcErr = s.callTool(r.Context(), req.Params)
	default:
		rpcErr = &RPCError{Code: -32601, Message: fmt.Sprintf("method not found: %s", req.Method)}
	}

	resp := Response{JSONRPC: "2.0", ID: req.ID, Result: result, Error: rpcErr}
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) callTool(ctx context.Context, raw json.RawMessage) (any, *RPCError) {
	var params CallToolParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, &RPCError{Code: -32602, Message: "invalid params"}
	}

	text, err := s.dispatch(ctx, params.Name, params.Arguments)
	if err != nil {
		return CallToolResult{
			Content: []Content{{Type: "text", Text: err.Error()}},
			IsError: true,
		}, nil
	}
	return CallToolResult{
		Content: []Content{{Type: "text", Text: text}},
	}, nil
}

func (s *Server) dispatch(ctx context.Context, name string, args map[string]any) (string, error) {
	c := s.client
	switch name {
	case "shopify_list_products":
		return shopify.ListProducts(ctx, c, args)
	case "shopify_get_product":
		return shopify.GetProduct(ctx, c, args)
	case "shopify_create_product":
		return shopify.CreateProduct(ctx, c, args)
	case "shopify_update_product":
		return shopify.UpdateProduct(ctx, c, args)
	case "shopify_list_collections":
		return shopify.ListCollections(ctx, c, args)
	case "shopify_get_collection":
		return shopify.GetCollection(ctx, c, args)
	case "shopify_create_collection":
		return shopify.CreateCollection(ctx, c, args)
	case "shopify_update_collection":
		return shopify.UpdateCollection(ctx, c, args)
	case "shopify_list_pages":
		return shopify.ListPages(ctx, c, args)
	case "shopify_get_page":
		return shopify.GetPage(ctx, c, args)
	case "shopify_create_page":
		return shopify.CreatePage(ctx, c, args)
	case "shopify_update_page":
		return shopify.UpdatePage(ctx, c, args)
	case "shopify_list_menus":
		return shopify.ListMenus(ctx, c, args)
	case "shopify_get_menu":
		return shopify.GetMenu(ctx, c, args)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// ServeStdio runs the MCP server over stdin/stdout using newline-delimited JSON-RPC.
func (s *Server) ServeStdio() {
	scanner := bufio.NewScanner(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	ctx := context.Background()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			enc.Encode(Response{
				JSONRPC: "2.0",
				Error:   &RPCError{Code: -32700, Message: "parse error"},
			})
			continue
		}

		var result any
		var rpcErr *RPCError

		switch req.Method {
		case "initialize":
			result = InitializeResult{
				ProtocolVersion: "2025-11-25",
				Capabilities:    Capabilities{Tools: map[string]any{}},
				ServerInfo:      ServerInfo{Name: "shopify-admin-mcp-server", Version: s.version},
			}
		case "notifications/initialized":
			// no response for notifications
			continue
		case "tools/list":
			result = ToolsListResult{Tools: GetTools()}
		case "tools/call":
			result, rpcErr = s.callTool(ctx, req.Params)
		default:
			rpcErr = &RPCError{Code: -32601, Message: fmt.Sprintf("method not found: %s", req.Method)}
		}

		enc.Encode(Response{JSONRPC: "2.0", ID: req.ID, Result: result, Error: rpcErr})
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "stdio read error: %v\n", err)
	}
}

func writeError(w http.ResponseWriter, id json.RawMessage, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg},
	}
	json.NewEncoder(w).Encode(resp)
}
