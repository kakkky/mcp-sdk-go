package main

import (
	"net/url"

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
	// Resourceを登録
	mcpServer.Resource(
		"example",
		"file://sample/uri",
		nil,
		&schema.ResourceMetadata{
			Description: "This is an example resource",
			MimeType:    "text/plain",
		},
		func(url url.URL) (schema.ReadResourceResultSchema, error) {
			return schema.ReadResourceResultSchema{Contents: []schema.ResourceContentSchema{
				&schema.TextResourceContentsSchema{
					UriData:      url.String(),
					MimeTypeData: "text/plain",
					ContentData:  "This is the content of the example resource.",
				},
			}}, nil
		},
		nil)
	transport := transport.NewStdioServerTransport()
	err := mcpServer.Connect(transport)
	if err != nil {
		panic(err)
	}
}
