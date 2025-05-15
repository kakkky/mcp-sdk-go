package mcpserver

import (
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/shared/protocol"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

type McpServer struct {
	server                        *server.Server
	registeredResources           map[string]RegisteredResource
	registeredResourceTemplates   map[string]RegisteredResourceTemplate
	isResourceHandlersInitialized bool
}

func NewMcpServer(serverInfo schema.Implementation, options *server.ServerOptions) *McpServer {
	return &McpServer{
		server:                      server.NewServer(serverInfo, options),
		registeredResources:         make(map[string]RegisteredResource),
		registeredResourceTemplates: make(map[string]RegisteredResourceTemplate),
	}
}

func (m *McpServer) Connect(transport protocol.Transport) error {
	return m.server.Connect(transport)
}

func (m *McpServer) Close() error {
	return m.server.Close()
}

func (m *McpServer) isConnected() bool {
	return m.server.Transport() != nil
}
