package mcpserver

import (
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
)

type McpServer struct {
	*server.Server
	registeredResources         map[string]RegisteredResource
	registeredResourceTemplates map[string]RegisteredResourceTemplate
}
