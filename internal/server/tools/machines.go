package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JarcauCristian/ztp-mcp/internal/server/maas_client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type Machines struct{}

func (Machines) RegisterTools(mcpServer *server.MCPServer) {
	mcpServer.AddTool(ListMachines{}.Create(), ListMachines{}.Handle)
}

type ListMachines struct{}

func (ListMachines) Create() mcp.Tool {
	return mcp.NewTool(
		"list_machines",
		mcp.WithDescription("List all the available machines on the current ZTP agent conected."),
	)
}

func (ListMachines) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := "/MAAS/api/2.0/machines/"

	client := maas_client.MustClient()

	zap.L().Info("Retrieving all the machines...")
	resultData, err := client.Get(ctx, path)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to retrieve the all the machines: %v", err))
		return nil, err
	}

	jsonData, err := json.Marshal(resultData)
	if err != nil {
		errMsg := fmt.Sprintf("failed to marshal result: %v", err)
		zap.L().Error(errMsg)
		return mcp.NewToolResultError(errMsg), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
