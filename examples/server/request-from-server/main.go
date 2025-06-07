package main

import (
	mcpserver "github.com/kakkky/mcp-sdk-go/mcp-server"
	"github.com/kakkky/mcp-sdk-go/mcp-server/server"
	"github.com/kakkky/mcp-sdk-go/mcp-server/transport"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func main() {
	mcpServer := mcpserver.NewMcpServer(
		schema.Implementation{
			Name:    "example-server",
			Version: "1.0.0",
		},
		&server.ServerOptions{
			Capabilities: schema.ServerCapabilities{
				Resources: &schema.Resources{
					ListChanged: true,
				},
			},
		})
	transport := transport.NewStdioServerTransport()
	go func() {
		err := mcpServer.Connect(transport)
		if err != nil {
			panic(err)
		}
	}()
	<-server.OperationPhaseStartNotify
	mcpServer.Server.Ping()
}
