package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JarcauCristian/ztp-mcp/internal/server/templates"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type Templates struct{}

func (Templates) Register(mcpServer *server.MCPServer) {
	mcpTools := []MCPTool{RetrieveTemplates{}}

	for _, tool := range mcpTools {
		mcpServer.AddTool(tool.Create(), tool.Handle)
	}
}

type RetrieveTemplates struct{}

func (RetrieveTemplates) Create() mcp.Tool {
	return mcp.NewTool(
		"retrieve_templates",
		mcp.WithDescription("Retruns all deployment Cloud-Init templates that are available on the system."),
	)
}

func (RetrieveTemplates) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	descriptions, err := templates.Templates()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to retrieve the all the machines: %v", err))
		return nil, err
	}

	jsonData, err := json.Marshal(descriptions)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal result: %v", err)
		zap.L().Error(errMsg)
		return mcp.NewToolResultError(errMsg), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
