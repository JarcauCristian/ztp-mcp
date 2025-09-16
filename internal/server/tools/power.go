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

type Power struct{}

func (Power) RegisterTools(mcpServer *server.MCPServer) {
	mcpTools := []MCPTool{PowerState{}, ChangePowerState{}}

	for _, tool := range mcpTools {
		mcpServer.AddTool(tool.Create(), tool.Handle)
	}
}

type PowerState struct{}

func (PowerState) Create() mcp.Tool {
	return mcp.NewTool(
		"power_state",
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Pattern("^[0-9a-z]{6}$"),
			mcp.Description("The id of the machine to retrieve information for."),
		),
		mcp.WithDescription("Returns the power state of a particular machine."),
	)
}

func (PowerState) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	machineID, err := request.RequireString("id")
	if err != nil {
		zap.L().Error(fmt.Sprintf("Required parameter id not present err=%v", err))
	}

	path := fmt.Sprintf("/MAAS/api/2.0/machines/%s/op-query_power_state", machineID)

	client := maas_client.MustClient()

	zap.L().Info(fmt.Sprintf("Retrieving power state for machine with id %s...", machineID))
	resultData, err := client.Get(ctx, path)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to retrieve power state for machine with id %s err=%v", machineID, err))
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

type ChangePowerState struct{}

func (ChangePowerState) Create() mcp.Tool {
	return mcp.NewTool(
		"change_power_state",
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Pattern("^[0-9a-z]{6}$"),
			mcp.Description("The id of the machine to retrieve information for."),
		),
		mcp.WithBoolean(
			"state",
			mcp.Required(),
			mcp.Description("If true power on the machine else power off."),
		),
		mcp.WithDescription("Change the power state of a machine specified by id."),
	)
}

func (ChangePowerState) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	machineID, err := request.RequireString("id")
	if err != nil {
		zap.L().Error(fmt.Sprintf("Required parameter id not present err=%v", err))
	}

	state, err := request.RequireBool("state")
	if err != nil {
		zap.L().Error(fmt.Sprintf("Required parameter state not present err=%v", err))
	}

	var path string

	if state {
		path = fmt.Sprintf("/MAAS/api/2.0/machines/%s/op-power_on", machineID)
	} else {
		path = fmt.Sprintf("/MAAS/api/2.0/machines/%s/op-power_off", machineID)
	}

	client := maas_client.MustClient()

	powerName := "on"

	if !state {
		powerName = "off"
	}

	zap.L().Info(fmt.Sprintf("Power machine with id %s %s...", machineID, powerName))
	resultData, err := client.Get(ctx, path)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to power %s machine with id %s err=%v", powerName, machineID, err))
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
