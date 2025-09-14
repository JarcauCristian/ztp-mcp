package main

import (
	"log"

	"github.com/JarcauCristian/ztp-mcp/internal/server/resources"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer(
		"Zero-Touch Provisioning MPC Server",
		"0.1.0",
		server.WithInstructions("This server is used to comunicate with the ZTP agent in order to deploy, interact and retrive the status of machines inside an Ubuntu MAAS instance."),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	mcpServer.AddResource(
		resources.CreateAvailableHosts(),
		resources.HandleAvailableHosts,
	)

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatal(err)
	}

	log.Println("Server started successfully!")
}
