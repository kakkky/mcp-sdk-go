package main

import (
	"fmt"
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
				Tools: &schema.Tools{
					ListChanged: true,
				},
				Prompts: &schema.Prompts{
					ListChanged: true,
				},
				Completion: &schema.Completion{},
			},
		})

	// Resourceを登録
	if _, err := mcpServer.Resource(
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
		}); err != nil {
		panic(err)
	}

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

	if _, err := mcpServer.ResourceTemplate(
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
	); err != nil {
		panic(err)
	}

	// Toolを登録
	if _, err := mcpServer.Tool(
		"calculate",
		"This tool calculates the sum of two numbers",
		schema.PropertySchema{
			"first": schema.PropertyInfoSchema{
				Type:        "number",
				Description: "This is the first parameter",
			},
			"second": schema.PropertyInfoSchema{
				Type:        "array",
				Description: "This is the second parameter, which is an array of numbers",
			},
		},
		nil,
		func(args map[string]any) (schema.CallToolResultSchema, error) {
			first, ok1 := args["first"].(float64)
			second, ok2 := args["second"].([]any)
			if !ok1 || !ok2 {
				return schema.CallToolResultSchema{
					Content: []schema.ToolContentSchema{},
					IsError: true,
				}, nil
			}
			var secondValue float64
			for _, v := range second {
				if num, ok := v.(float64); ok {
					secondValue += num
				} else {
					return schema.CallToolResultSchema{
						Content: []schema.ToolContentSchema{},
						IsError: true,
					}, nil
				}
			}
			return schema.CallToolResultSchema{
				Content: []schema.ToolContentSchema{
					&schema.TextContentSchema{
						Type: "text",
						Text: "The result of the addition is: " + fmt.Sprintf("%v", first+secondValue),
					},
				},
			}, nil
		},
	); err != nil {
		panic(err)
	}
	if _, err := mcpServer.Prompt(
		"example-prompt",
		"This is an example prompt",
		[]schema.PromptAugmentSchema{
			{
				Name:             "input",
				Description:      "This is an input parameter",
				Required:         true,
				CompletionValues: []string{"value1", "value2", "value3"},
			},
		},
		func(args []schema.PromptAugmentSchema) (schema.GetPromptResultSchema, error) {
			var promptMessages []schema.PromptMessageSchema
			for _, arg := range args {
				if arg.Name == "input" {
					promptMessages = append(promptMessages, schema.PromptMessageSchema{
						Role: "user",
						Content: &schema.TextContentSchema{
							Type: "text",
							Text: "You provided input: " + arg.CompletionValues[0], //
						},
					})
				}
			}
			return schema.GetPromptResultSchema{
				Description: "This is a response from the example prompt",
				Messages:    promptMessages,
			}, err
		},
	); err != nil {
		panic(err)
	}

	transport := transport.NewStdioServerTransport()
	if err := mcpServer.Connect(transport); err != nil {
		panic(err)
	}
}
