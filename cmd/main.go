package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JarcauCristian/ztp-mcp/internal/server/tools"
	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/server"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	mcpServer := server.NewMCPServer(
		"Zero-Touch Provisioning MPC Server",
		"0.1.0",
		server.WithInstructions("This server is used to communicate with the ZTP agent in order to deploy, interact and retrieve the status of machines inside an Ubuntu MAAS instance."),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	tools.VMHosts{}.RegisterTools(mcpServer)
	tools.Machines{}.RegisterTools(mcpServer)

	transport := os.Getenv("MCP_TRANSPORT")
	addr := os.Getenv("MCP_ADDRESS")
	switch transport {
	case "SSE", "sse":
		zap.L().Info("Starting MCP server in SSE mode...")
		sseServer := server.NewSSEServer(mcpServer)
		if err := sseServer.Start(addr); err != nil {
			log.Fatal(err)
		}
	case "HTTP", "http":
		zap.L().Info(fmt.Sprintf("Starting MCP server in Streamable HTTP mode on %s...", addr))
		httpServer := server.NewStreamableHTTPServer(mcpServer)
		if err := httpServer.Start(addr); err != nil {
			log.Fatal(err)
		}
	default:
		zap.L().Info("Starting MCP server in stdio mode...")
		if err := server.ServeStdio(mcpServer); err != nil {
			zap.L().Fatal(err.Error())
		}
	}
}
