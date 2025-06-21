package main

import (
	"sync"

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
				Logging: &schema.Logging{},
			},
		})
	transport := transport.NewStdioServerTransport()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := mcpServer.Connect(transport)
		if err != nil {
			panic(err)
		}
	}()
	<-server.OperationPhaseStartNotify
	if err := mcpServer.Server.SendLoggingMessage(
		schema.LoggingMessageNotificationParams{
			Level: schema.NOTICE,
			Data:  "Server started successfully",
		},
	); err != nil {
		panic(err)
	}
	_, _ = mcpServer.Server.Ping()
	_, _ = mcpServer.Server.ListRoots()
	wg.Wait()
}
