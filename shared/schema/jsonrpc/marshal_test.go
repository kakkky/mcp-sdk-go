package jsonrpc

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	mcperr "github.com/kakkky/mcp-sdk-go/shared/mcp-err"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestMarshal(t *testing.T) {
	tests := []struct {
		name        string
		message     schema.JsonRpcMessage
		expectedStr string
	}{
		{
			name: "normal : able to marshal initialize request",
			message: schema.JsonRpcRequest{
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
			expectedStr: `{
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
		},
		{
			name: "normal : able to marshal ping request",
			message: schema.JsonRpcRequest{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      2,
				},
				Request: &schema.PingRequestSchema{
					MethodName: "ping",
				},
			},
			expectedStr: `{
						"jsonrpc": "2.0",
						"id": 2,
						"method": "ping"
					}`,
		},
		{
			name: "normal : able to marshal resources/read response with mixed contents",
			message: schema.JsonRpcResponse{
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
			expectedStr: `{
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
		},
		{
			name: "normal : able to marshal completion/complete response",
			message: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      11,
				},
				Result: &schema.CompleteResultSchema{
					Completion: schema.CompletionSchema{
						Values:  []string{"js", "python", "go"},
						Total:   3,
						HasMore: func() *bool { b := false; return &b }(),
					},
				},
			},
			expectedStr: `{
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
		},
		{
			name: "normal : able to marshal sampling/createMessage response with text content",
			message: schema.JsonRpcResponse{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      6,
				},
				Result: &schema.CreateMessageResultSchema[schema.TextContentSchema]{
					Model:      "gpt-4",
					StopReason: "endTurn",
					Role:       "assistant",
					Content: schema.TextContentSchema{
						Type: "text",
						Text: "Hello! How can I assist you today?",
					},
				},
			},
			expectedStr: `{
						"jsonrpc": "2.0",
						"id": 6,
						"result": {
							"model": "gpt-4",
							"stopReason": "endTurn",
							"role": "assistant",
							"content": {
							"type": "text",
							"text": "Hello! How can I assist you today?"
							}
						}
					}`,
		},
		{
			name: "normal : able to marshal initialized notification",
			message: schema.JsonRpcNotification{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Notification: &schema.InitializeNotificationSchema{
					MethodName: "notifications/initialized",
				},
			},
			expectedStr: `{
						"jsonrpc": "2.0",
						"method": "notifications/initialized"
					}`,
		},
		{
			name: "normal : able to marshal logging message notification",
			message: schema.JsonRpcNotification{
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
			expectedStr: `{
						"jsonrpc": "2.0",
						"method": "notifications/message",
						"params": {
							"level": "info",
							"data": "This is an informational message"
						}
					}`,
		},
		{
			name: "normal : able to marshal resources updated notification",
			message: schema.JsonRpcNotification{
				Jsonrpc: schema.JSON_RPC_VERSION,
				Notification: &schema.ResourceUpdatedNotificationSchema{
					MethodName: "notifications/resources/updated",
					ParamsData: schema.ResourceUpdatedNotificationParams{
						Uri: "file:///example.txt",
					},
				},
			},
			expectedStr: `{
						"jsonrpc": "2.0",
						"method": "notifications/resources/updated",
						"params": {
							"uri": "file:///example.txt"
						}
					}`,
		},
		{
			name: "normal : able to marshal error response with simple error",
			message: schema.JsonRpcError{
				BaseMessage: schema.BaseMessage{
					Jsonrpc: schema.JSON_RPC_VERSION,
					Id:      100,
				},
				Error: schema.Error{
					Code:    mcperr.INVALID_PARAMS,
					Message: "Invalid params",
				},
			},
			expectedStr: `{
						"jsonrpc": "2.0",
						"id": 100,
						"error": {
							"code": -32602,
							"message": "Invalid params"
						}
					}`,
		},
		{
			name: "normal : able to marshal error response with detailed error data",
			message: schema.JsonRpcError{
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
			expectedStr: `{
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, err := Marshal(test.message)
			if err != nil {
				t.Fatalf("Failed to marshal message: %v", err)
			}

			// 出力されたJSONを圧縮
			var jsonDataBuffer bytes.Buffer
			if err := json.Compact(&jsonDataBuffer, jsonData); err != nil {
				t.Fatalf("Failed to compact JSON data: %v", err)
			}

			// 期待値のJSONも同様に圧縮
			var expectedBuffer bytes.Buffer
			if err := json.Compact(&expectedBuffer, []byte(test.expectedStr)); err != nil {
				t.Fatalf("Failed to compact expected JSON: %v", err)
			}

			// 圧縮されたJSON同士を比較
			if diff := cmp.Diff(expectedBuffer.String(), jsonDataBuffer.String()); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
