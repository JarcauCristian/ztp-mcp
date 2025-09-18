package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/JarcauCristian/ztp-mcp/internal/server/maas_client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

const NUMBER_PATTERN = "^[0-9]+$"

type VMHosts struct{}

func (VMHosts) Register(mcpServer *server.MCPServer) {
	mcpTools := []MCPTool{ListVMHosts{}, ListVMHost{}, ComposeVM{}}

	for _, tool := range mcpTools {
		mcpServer.AddTool(tool.Create(), tool.Handle)
	}
}

type ListVMHosts struct{}

func (ListVMHosts) Create() mcp.Tool {
	return mcp.NewTool(
		"list_vm_hosts",
		mcp.WithDescription("Returns the available VM hosts from the ZTP agent conected."),
	)
}

func (ListVMHosts) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := "/MAAS/api/2.0/vm-hosts/"

	client := maas_client.MustClient()

	zap.L().Info("Retrieving all VM hosts...")
	resultData, err := client.Get(ctx, path)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to retrieve the VM hosts: %v", err))
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

type ListVMHost struct{}

func (ListVMHost) Create() mcp.Tool {
	return mcp.NewTool(
		"list_vm_host",
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Description("The ID of the VM host to query information for."),
			mcp.Pattern(NUMBER_PATTERN),
		),
		mcp.WithDescription("Returns information about a particular VM host specified by id on the ZTP agent conected."),
	)
}

func (ListVMHost) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	vmID, err := request.RequireString("id")
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/MAAS/api/2.0/vm-hosts/%s/", vmID)

	client := maas_client.MustClient()

	zap.L().Info(fmt.Sprintf("Retrieving VM host with ID %s...", vmID))
	resultData, err := client.Get(ctx, path)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to retreive VM host with ID %s, err=%v", vmID, err))
		return nil, err
	}

	jsonData, err := json.Marshal(resultData)
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to marshal result: %v", err))
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

type ComposeVM struct{}

func (ComposeVM) Create() mcp.Tool {
	return mcp.NewTool(
		"compose_vm_host",
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Description("ID of the VM host to compose the machine on."),
			mcp.Pattern(NUMBER_PATTERN),
		),
		mcp.WithString(
			"cores",
			mcp.Required(),
			mcp.Description("The number of cores the composed VM should have."),
			mcp.Pattern(NUMBER_PATTERN),
		),
		mcp.WithString(
			"memory",
			mcp.Required(),
			mcp.Description("How much RAM the composed VM should have (Should be in MiB)."),
			mcp.Pattern(NUMBER_PATTERN),
		),
		mcp.WithString(
			"storage",
			mcp.Required(),
			mcp.Description("How much storage the composed VM should have (Should be in GB)."),
			mcp.Pattern(NUMBER_PATTERN),
		),
		mcp.WithString(
			"hostname",
			mcp.Required(),
			mcp.Description("The name of the created VM (Give something random if not provided)."),
			mcp.Pattern("^[a-zA-Z0-9.-]+$"),
		),
		mcp.WithDescription("Compose a VM on a particular VM host specified by ID."),
	)
}

func (ComposeVM) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	vmHostID, err := request.RequireString("id")
	if err != nil {
		return nil, err
	}

	cores, err := request.RequireString("cores")
	if err != nil {
		return nil, err
	}

	memory, err := request.RequireString("memory")
	if err != nil {
		return nil, err
	}

	storage, err := request.RequireString("storage")
	if err != nil {
		return nil, err
	}

	hostname, err := request.RequireString("hostname")
	if err != nil {
		return nil, err
	}

	form := make(url.Values)
	form.Set("cores", cores)
	form.Set("memory", memory)
	form.Set("storage", storage)
	form.Set("hostname", hostname)

	path := fmt.Sprintf("/MAAS/api/2.0/vm-hosts/%s/op-compose", vmHostID)

	client := maas_client.MustClient()

	zap.L().Info(fmt.Sprintf("Composing VM on host %s with the following configuration:\nCores: %s\nMemory: %s\nStorage: %s\nHostname: %s", vmHostID, cores, memory, storage, hostname))
	resultData, err := client.Post(ctx, path, strings.NewReader(form.Encode()))
	if err != nil {
		zap.L().Error(fmt.Sprintf("Failed to compose VM err=%v", err))
		return nil, err
	}

	jsonData, err := json.Marshal(resultData)
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to marshal result: %v", err))
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}
