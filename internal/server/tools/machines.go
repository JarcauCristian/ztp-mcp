package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/JarcauCristian/ztp-mcp/internal/server/maas_client"
	"github.com/JarcauCristian/ztp-mcp/internal/server/parser"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

type Machines struct{}

func (Machines) Register(mcpServer *server.MCPServer) {
	mcpTools := []MCPTool{ListMachines{}, ListMachine{}, CommissionMachine{}}

	for _, tool := range mcpTools {
		mcpServer.AddTool(tool.Create(), tool.Handle)
	}
}

type ListMachines struct{}

func (ListMachines) Create() mcp.Tool {
	return mcp.NewTool(
		"list_machines",
		mcp.WithString(
			"status",
			mcp.Enum(
				"new",
				"commissioning",
				"failed_commissioning",
				"ready",
				"deploying",
				"deployed",
				"releasing",
				"failed_deployment",
				"allocated",
				"retired",
				"broken",
				"recommissioning",
				"testing",
				"failed_testing",
				"rescuing",
				"disk_erasing",
				"failed_disk_erasing",
			),
			mcp.Description("The status of the machine that will be retrieved. Returns all machines if not provided."),
		),
		mcp.WithDescription("List all the available machines on the current ZTP agent conected."),
	)
}

func (ListMachines) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var path, errMsg string

	status := request.GetString("status", "")

	if status == "" {
		path = "/MAAS/api/2.0/machines/"
	} else {
		path = fmt.Sprintf("/MAAS/api/2.0/machines/?status=%s", status)
	}

	client := maas_client.MustClient()

	zap.L().Info("[ListMachines] Retrieving all the machines...")
	resultData, err := client.Get(ctx, path)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to retrieve all the machines: %v", err)
		zap.L().Error(fmt.Sprintf("[ListMachines] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	var machines []map[string]any

	err = json.Unmarshal([]byte(resultData), &machines)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to unmarshal the result err=%v", err)
		zap.L().Error(fmt.Sprintf("[ListMachines] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	// Filter out machines with the "protected" tag
	var filteredMachines []map[string]any
	for _, machine := range machines {
		if !parser.CheckForProtectedTag(machine) {
			filteredMachines = append(filteredMachines, machine)
		}
	}

	response, err := json.Marshal(filteredMachines)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to marshal filtered machines: %v", err)
		zap.L().Error(fmt.Sprintf("[ListMachines] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	return mcp.NewToolResultText(string(response)), nil
}

type ListMachine struct{}

func (ListMachine) Create() mcp.Tool {
	return mcp.NewTool(
		"list_machine",
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Pattern("^[0-9a-z]{6}$"),
			mcp.Description("The id of the machine to retrieve information for."),
		),
		mcp.WithDescription("Return the information about a particular machine."),
	)
}

func (ListMachine) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var errMsg string
	jsonData := make([]byte, 0)

	machineID, err := request.RequireString("id")
	if err != nil {
		zap.L().Error(fmt.Sprintf("[ListMachine] Required parameter id not present err=%v", err))
		return mcp.NewToolResultError(err.Error()), nil
	}

	path := fmt.Sprintf("/MAAS/api/2.0/machines/%s/", machineID)

	client := maas_client.MustClient()

	zap.L().Info(fmt.Sprintf("[ListMachine] Retrieving machine with id %s...", machineID))
	resultData, err := client.Get(ctx, path)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to retrieve the machine with id %s err=%v", machineID, err)
		zap.L().Error(fmt.Sprintf("[ListMachine] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	var machine map[string]any
	err = json.Unmarshal([]byte(resultData), &machine)
	if err != nil {
		errMsg = fmt.Sprintf("Failed to unmarshal the result: %v", err)
		zap.L().Error(fmt.Sprintf("[ListMachine] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	if !parser.CheckForProtectedTag(machine) {
		jsonData, err = json.Marshal(resultData)
		if err != nil {
			errMsg = fmt.Sprintf("failed to marshal result: %v", err)
			zap.L().Error(fmt.Sprintf("[ListMachine] %s", errMsg))
			return mcp.NewToolResultError(errMsg), nil
		}
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

type CommissionMachine struct{}

func (CommissionMachine) Create() mcp.Tool {
	return mcp.NewTool(
		"commission_machine",
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Pattern("^[0-9a-z]{6}$"),
			mcp.Description("The id of the machine to commission."),
		),
		mcp.WithDescription("Start the commissioning process on a particular machine."),
	)
}

func (CommissionMachine) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var errMsg string

	machineID, err := request.RequireString("id")
	if err != nil {
		zap.L().Error(fmt.Sprintf("[CommissionMachine] Required parameter id not present err=%v", err))
		return mcp.NewToolResultError(err.Error()), nil
	}

	path := fmt.Sprintf("/MAAS/api/2.0/machines/%s/op-commission", machineID)

	client := maas_client.MustClient()

	form := make(url.Values)
	form.Add("enable_ssh", "1")

	zap.L().Info(fmt.Sprintf("[CommissionMachine] Commissioning machine with id %s...", machineID))
	resultData, err := client.Post(ctx, path, strings.NewReader(form.Encode()))
	if err != nil {
		errMsg = fmt.Sprintf("Failed to commission the machine with id %s err=%v", machineID, err)
		zap.L().Error(fmt.Sprintf("[CommissionMachine] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	jsonData, err := json.Marshal(resultData)
	if err != nil {
		errMsg = fmt.Sprintf("failed to marshal result: %v", err)
		zap.L().Error(fmt.Sprintf("[CommissionMachine] %s", errMsg))
		return mcp.NewToolResultError(errMsg), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
