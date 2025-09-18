package registry

import "github.com/mark3labs/mcp-go/server"

type Registry interface {
	Register(mcpServer *server.MCPServer)
}
