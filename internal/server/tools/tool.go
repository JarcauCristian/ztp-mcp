package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type MCPTool interface {
	Create() mcp.Tool
	Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func CreateToolAnnotation(title string, readOnly, destructive, idempotency, openWorld bool) mcp.ToolAnnotation {
	return mcp.ToolAnnotation{
		Title:           title,
		ReadOnlyHint:    mcp.ToBoolPtr(readOnly),
		IdempotentHint:  mcp.ToBoolPtr(idempotency),
		DestructiveHint: mcp.ToBoolPtr(destructive),
		OpenWorldHint:   mcp.ToBoolPtr(openWorld),
	}
}
