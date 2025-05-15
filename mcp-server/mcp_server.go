package mcpserver

import (
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type McpServer struct {
	*server.Server
	registeredResources           map[string]RegisteredResource
	registeredResourceTemplates   map[string]RegisteredResourceTemplate
	isResourceHandlersInitialized bool
}

func NewMcpServer(serverInfo schema.Implementation, options *server.ServerOptions) *McpServer {
	return &McpServer{
		Server:                      server.NewServer(serverInfo, options),
		registeredResources:         make(map[string]RegisteredResource),
		registeredResourceTemplates: make(map[string]RegisteredResourceTemplate),
	}
}

func (m *McpServer) Connect(transport protocol.Transport) error {
	return m.Server.Connect(transport)
}

func (m *McpServer) Close() error {
	return m.Server.Close()
}

func (m *McpServer) isConnected() bool {
	return m.Transport() != nil
}
