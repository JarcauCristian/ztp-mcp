package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type MCPTool interface {
	Create() mcp.Tool
	Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
