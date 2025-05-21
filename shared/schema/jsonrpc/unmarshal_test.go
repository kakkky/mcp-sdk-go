package jsonrpc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		expected schema.JsonRpcMessage
	}{
		{
			name: "normal : able to unmarshal initialize request",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 1,
				"method": "initialize",
				"params": {
					"protocolVersion": "2025-01-01",
					"capabilities": {
						"sampling": {}
					},
					"clientInfo": {
						"name": "test-client",
						"version": "1.0.0"
					}
				}
			}`,
			expected: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      1,
				},
				Request: &schema.InitializeRequestSchema{
					MethodName: "initialize",
					ParamsData: schema.InitializeRequestParams{
						ProtocolVersion: "2025-01-01",
						Capabilities: schema.ClientCapabilities{
							Sampling: &schema.Sampling{},
						},
						ClientInfo: schema.Implementation{
							Name:    "test-client",
							Version: "1.0.0",
						},
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal ping request",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 2,
				"method": "ping"
			}`,
			expected: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      2,
				},
				Request: &schema.PingRequestSchema{
					MethodName: "ping",
				},
			},
		},
		{
			name: "normal : able to unmarshal read resource request",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 3,
				"method": "resources/read",
				"params": {
					"uri": "file:///example.txt"
				}
			}`,
			expected: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      3,
				},
				Request: &schema.ReadResourceRequestSchema{
					MethodName: "resources/read",
					ParamsData: schema.ReadResourceRequestParams{
						Uri: "file:///example.txt",
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal completion/complete request",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 4,
				"method": "completion/complete",
				"params": {
					"ref": {
						"type": "ref/resource",
						"uri": "file:///{path}"
					},
					"argument": {
						"name": "test",
						"value": "test"	
					}
				}
			}`,
			expected: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      4,
				},
				Request: &schema.CompleteRequestSchema{
					MethodName: "completion/complete",
					ParamsData: schema.CompleteRequestParams{
						Ref: &schema.ResourceReferenceSchema{
							TypeData: "ref/resource",
							UriData:  "file:///{path}",
						},
						Argument: schema.CompleteRequestParamsArg{
							Name:  "test",
							Value: "test",
						},
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal initialize response",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 5,
				"result": {
					"protocolVersion": "2025-01-01",
					"capabilities": {
						"resources": {
							"listChanged": true
						},
						"logging": {}
					},
					"serverInfo": {
						"name": "test-server",
						"version": "1.0.0"
					},
					"instructions": "Welcome to the test server"
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      5,
				},
				Result: &schema.InitializeResultSchema{
					ProtocolVersion: "2025-01-01",
					Capabilities: schema.ServerCapabilities{
						Resources: &schema.Resources{
							ListChanged: true,
						},
						Logging: &schema.Logging{},
					},
					ServerInfo: schema.Implementation{
						Name:    "test-server",
						Version: "1.0.0",
					},
					Instructions: "Welcome to the test server",
				},
			},
		},
		{
			name: "normal : able to unmarshal sampling/createMessage response",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 6,
				"result": {
					"model": "gpt-4",
					"role": "assistant",
					"stopReason": "endTurn", 
					"content": {
						"type": "text",
						"text": "Hello! How can I assist you today?"
					}
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      6,
				},
				Result: &schema.CreateMessageResultSchema[schema.TextContentSchema]{
					Model:      "gpt-4",
					Role:       "assistant",
					StopReason: "endTurn",
					Content: schema.TextContentSchema{
						Type: "text",
						Text: "Hello! How can I assist you today?",
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal resources/read response with blob & text contents",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 8,
				"result": {
					"contents": [
						{
							"uri": "file:///example.jpg",
							"mimeType": "image/jpeg",
							"blob": "base64encodeddata"
						},
						{
							"uri": "file:///example.txt",
							"mimeType": "text/plain",
							"text": "This is the content of example.txt"
						}
					]
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      8,
				},
				Result: &schema.ReadResourceResultSchema{
					Contents: []schema.ResourceContentSchema{
						&schema.BlobResourceContentsSchema{
							UriData:      "file:///example.jpg",
							MimeTypeData: "image/jpeg",
							ContentData:  "base64encodeddata",
						},
						&schema.TextResourceContentsSchema{
							UriData:      "file:///example.txt",
							MimeTypeData: "text/plain",
							ContentData:  "This is the content of example.txt",
						},
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal resources/list response",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 9,
				"result": {
					"resources": [
						{
							"uri": "file:///example.txt",
							"name": "Example Text",
							"description": "An example text file",
							"mimeType": "text/plain"
						},
						{
							"uri": "file:///image.jpg",
							"name": "Example Image",
							"mimeType": "image/jpeg"
						}
					]
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      9,
				},
				Result: &schema.ListResourcesResultSchema{
					Resources: []schema.ResourceSchema{
						{
							Uri:  "file:///example.txt",
							Name: "Example Text",
							ResourceMetadata: &schema.ResourceMetadata{
								Description: "An example text file",
								MimeType:    "text/plain",
							},
						},
						{
							Uri:  "file:///image.jpg",
							Name: "Example Image",
							ResourceMetadata: &schema.ResourceMetadata{
								MimeType: "image/jpeg",
							},
						},
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal resources/templates/list response",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 10,
				"result": {
					"resourceTemplates": [
						{
							"uriTemplate": "file:///{filename}.txt",
							"name": "Text Template",
							"description": "Template for text files",
							"mimeType": "text/plain"
						},
						{
							"uriTemplate": "file:///{filename}.md",
							"name": "Markdown Template",
							"mimeType": "text/markdown"
						}
					]
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      10,
				},
				Result: &schema.ListResourceTemplatesResultSchema{
					ResourceTemplates: []schema.ResourceTemplateSchema{
						{
							UriTemplate: "file:///{filename}.txt",
							Name:        "Text Template",
							ResourceMetadata: &schema.ResourceMetadata{
								Description: "Template for text files",
								MimeType:    "text/plain",
							},
						},
						{
							UriTemplate: "file:///{filename}.md",
							Name:        "Markdown Template",
							ResourceMetadata: &schema.ResourceMetadata{
								MimeType: "text/markdown",
							},
						},
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal completion/complete response",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 11,
				"result": {
					"completion": {
						"values": ["js", "python", "go"],
						"total": 3,
						"hasMore": false
					}
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      11,
				},
				Result: &schema.CompleteResultSchema{
					Completion: schema.CompletionSchema{
						Values:  []string{"js", "python", "go"},
						Total:   3,
						HasMore: func() *bool { b := false; return &b }(), // booleanをポインタで渡すため
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal empty result response",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 12,
				"result": {}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      12,
				},
				Result: &schema.EmptyResultSchema{},
			},
		},
		{
			name: "normal : able to unmarshal sampling/createMessage response with image content",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 13,
				"result": {
					"model": "dall-e-3",
					"role": "assistant",
					"stopReason": "complete", 
					"content": {
						"type": "image",
						"mimeType": "image/png",
						"data": "base64encodedimagedata"
					}
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      13,
				},
				Result: &schema.CreateMessageResultSchema[schema.ImageContentSchema]{
					Model:      "dall-e-3",
					Role:       "assistant",
					StopReason: "complete",
					Content: schema.ImageContentSchema{
						Type:     "image",
						MimeType: "image/png",
						Data:     "base64encodedimagedata",
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal sampling/createMessage response with audio content",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 14,
				"result": {
					"model": "whisper",
					"role": "assistant",
					"stopReason": "complete", 
					"content": {
						"type": "audio",
						"mimeType": "audio/mp3",
						"data": "base64encodedaudiodata"
					}
				}
			}`,
			expected: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      14,
				},
				Result: &schema.CreateMessageResultSchema[schema.AudioContentSchema]{
					Model:      "whisper",
					Role:       "assistant",
					StopReason: "complete",
					Content: schema.AudioContentSchema{
						Type:     "audio",
						MimeType: "audio/mp3",
						Data:     "base64encodedaudiodata",
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal initialized notification",
			jsonStr: `{
				"jsonrpc": "2.0",
				"method": "notifications/initialized"
			}`,
			expected: schema.JsonRpcNotification{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Notification: &schema.InitializeNotificationSchema{
					MethodName: "notifications/initialized",
				},
			},
		},
		{
			name: "normal : able to unmarshal logging message notification",
			jsonStr: `{
				"jsonrpc": "2.0",
				"method": "notifications/message",
				"params": {
					"level": "info",
					"logger": "",
					"data": "This is an informational message"
				}
			}`,
			expected: schema.JsonRpcNotification{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Notification: &schema.LoggingMessageNotificationSchema{
					MethodName: "notifications/message",
					ParamsData: schema.LoggingMessageNotificationParams{
						Level:  "info",
						Logger: "",
						Data:   "This is an informational message",
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal resources updated notification",
			jsonStr: `{
				"jsonrpc": "2.0",
				"method": "notifications/resources/updated",
				"params": {
					"uri": "file:///example.txt"
				}
			}`,
			expected: schema.JsonRpcNotification{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Notification: &schema.ResourceUpdatedNotificationSchema{
					MethodName: "notifications/resources/updated",
					ParamsData: schema.ResourceUpdatedNotificationParams{
						Uri: "file:///example.txt",
					},
				},
			},
		},
		{
			name: "normal : able to unmarshal resources list changed notification",
			jsonStr: `{
				"jsonrpc": "2.0", 
				"method": "notifications/resources/list_changed"
			}`,
			expected: schema.JsonRpcNotification{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Notification: &schema.ResourceListChangedNotificationSchema{
					MethodName: "notifications/resources/list_changed",
				},
			},
		},
		{
			name: "normal : able to unmarshal error response with simple error message",
			jsonStr: `{
				"jsonrpc": "2.0",
				"id": 100,
				"error": {
					"code": -32602,
					"message": "Invalid params"
				}
			}`,
			expected: schema.JsonRpcError{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      100,
				},
				Error: schema.Error{
					Code:    mcperr.INVALID_PARAMS,
					Message: "Invalid params",
				},
			},
		},
		{
			name: "normal : able to unmarshal error response with detailed error data",
			jsonStr: `{
				"jsonrpc": "2.0", 
				"id": 101,
				"error": {
					"code": -32601,
					"message": "Method not found",
					"data": {
						"details": "The requested method 'unknown_method' is not supported",
						"requestId": "req-123456"
					}
				}
			}`,
			expected: schema.JsonRpcError{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      101,
				},
				Error: schema.Error{
					Code:    mcperr.METHOD_NOT_FOUND,
					Message: "Method not found",
					Data: map[string]interface{}{
						"details":   "The requested method 'unknown_method' is not supported",
						"requestId": "req-123456",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal([]byte(tt.jsonStr))
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			if got == nil {
				t.Errorf("Unmarshal() got = nil")
				return
			}
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("Unmarshal() got(+) = %v, expected(-) %v, diff: %s", got, tt.expected, diff)
			}
		})
	}
}
