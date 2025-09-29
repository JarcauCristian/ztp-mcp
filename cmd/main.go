package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/JarcauCristian/ztp-mcp/internal/server/registry"
	"github.com/JarcauCristian/ztp-mcp/internal/server/tools"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mark3labs/mcp-go/server"
)

func init() {
	var logger *zap.Logger

	config := zap.NewDevelopmentConfig()

	config.DisableCaller = false
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.FunctionKey = "function"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	logger = zap.Must(config.Build())

	zap.ReplaceGlobals(logger)
}

func registerTools(mcpServer *server.MCPServer) {
	registries := []registry.Registry{tools.VMHosts{}, tools.Machines{}, tools.Power{}, tools.Templates{}}

	for _, reg := range registries {
		reg.Register(mcpServer)
	}
}

func main() {
	var version string
	info, ok := debug.ReadBuildInfo()

	if !ok {
		version = "0.1.0"
	} else {
		version = info.Main.Version
	}

	mcpServer := server.NewMCPServer(
		"Zero-Touch Provisioning MPC Server",
		version,
		server.WithInstructions("This server is used to communicate with the ZTP agent in order to deploy, interact and retrieve the status of machines inside an Ubuntu MAAS instance."),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	registerTools(mcpServer)

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
