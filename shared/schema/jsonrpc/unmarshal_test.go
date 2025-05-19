package json

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
