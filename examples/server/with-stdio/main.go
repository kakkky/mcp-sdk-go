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
				Completion: &schema.Completion{},
			},
		})
	// Resourceを登録
	mcpServer.Resource(
		"example",
		"file://sample/uri",
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
		})

	// Resource Templateを登録
	template, err := mcpserver.NewResourceTemplate("file://sample/{variable}", &mcpserver.ResourceTemplateCallbacks{
		List: func() schema.ListResourcesResultSchema {
			return schema.ListResourcesResultSchema{
				Resources: []schema.ResourceSchema{
					{
						Name: "example-template",
						Uri:  "file://sample/example-value",
					},
					{
						Name: "another-template",
						Uri:  "file://sample/another-value",
					},
				},
			}
		},
		Complete: map[string]mcpserver.CompleteResourceCallback{
			"variable": func(value string) []string {
				// ここでは単純に固定の値を返すが、実際には何らかのロジックで候補を生成することができる
				return []string{"example-value", "another-value"}
			},
		},
	})
	if err != nil {
		panic(err)
	}
	mcpServer.ResourceTemplate(
		"example-template",
		template,
		&schema.ResourceMetadata{
			Description: "This is an example resource template",
			MimeType:    "text/plain",
		},
		func(url url.URL, variables map[string]any) (schema.ReadResourceResultSchema, error) {
			switch variables["variable"] {
			case "example-value":
				return schema.ReadResourceResultSchema{Contents: []schema.ResourceContentSchema{
					&schema.TextResourceContentsSchema{
						UriData:      "file://sample/example-value",
						MimeTypeData: "text/plain",
						ContentData:  "This is the content of the example resource template",
					},
				}}, nil
			case "another-value":
				return schema.ReadResourceResultSchema{Contents: []schema.ResourceContentSchema{
					&schema.TextResourceContentsSchema{
						UriData:      "file://sample/another-value",
						MimeTypeData: "text/plain",
						ContentData:  "This is the content of another example resource template",
					},
				}}, nil
			}
			return schema.ReadResourceResultSchema{}, nil
		},
	)
	transport := transport.NewStdioServerTransport()
	err = mcpServer.Connect(transport)
	if err != nil {
		panic(err)
	}
}
