package san_mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ctxKey int

const callToolRequestKey ctxKey = iota

// WithCallToolRequest returns a new context that carries the originating MCP CallToolRequest.
// Generated MCP tool adapters set this before invoking the underlying RPC client so downstream
// handlers can recover the MCP-level request (headers, tool name, request id, etc.).
func WithCallToolRequest(ctx context.Context, req *mcp.CallToolRequest) context.Context {
	return context.WithValue(ctx, callToolRequestKey, req)
}

// CallToolRequestFromContext returns the MCP CallToolRequest stored on ctx, or nil if absent.
func CallToolRequestFromContext(ctx context.Context) *mcp.CallToolRequest {
	req, _ := ctx.Value(callToolRequestKey).(*mcp.CallToolRequest)
	return req
}
