package transport

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestReadBuffer(t *testing.T) {
	tests := []struct {
		name           string
		chunk          []byte
		expectMsgCount int
		expected       []schema.JsonRpcMessage
		wantErr        bool
	}{
		{
			name:           "normal: cannot read anything from empty buffer",
			chunk:          []byte{},
			expectMsgCount: 0,
			expected:       nil,
		},
		{
			name:           "normal: JSON without newline cannot be read as a message",
			chunk:          []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}`),
			expectMsgCount: 0,
			expected:       nil,
		},
		{
			name:           "normal: can read one JSONRPC message with newline",
			chunk:          []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}` + "\n"),
			expectMsgCount: 1,
			expected: []schema.JsonRpcMessage{
				schema.JsonRpcRequest{
					BaseMessage: schema.BaseMessage{
						Jsonrpc: schema.JSON_RPC_VERSION,
						Id:      1,
					},
					Request: &schema.PingRequestSchema{
						MethodName: "ping",
					},
				},
			},
		},
		{
			name:           "normal: can read multiple JSONRPC messages separated by newlines",
			chunk:          []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}` + "\n" + `{"jsonrpc":"2.0","id":2,"method":"initialize","params":{"protocolVersion":"2025-01-01","capabilities":{"sampling":{}},"clientInfo":{"name":"test-client","version":"1.0.0"}}}` + "\n"),
			expectMsgCount: 2,
			expected: []schema.JsonRpcMessage{
				schema.JsonRpcRequest{
					BaseMessage: schema.BaseMessage{
						Jsonrpc: schema.JSON_RPC_VERSION,
						Id:      1,
					},
					Request: &schema.PingRequestSchema{
						MethodName: "ping",
					},
				},
				schema.JsonRpcRequest{
					BaseMessage: schema.BaseMessage{
						Jsonrpc: schema.JSON_RPC_VERSION,
						Id:      2,
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
		},
		{
			name:           "normal: can read JSONRPC message separated by CRLF",
			chunk:          []byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}` + "\r\n"),
			expectMsgCount: 1,
			expected: []schema.JsonRpcMessage{
				schema.JsonRpcRequest{
					BaseMessage: schema.BaseMessage{
						Jsonrpc: schema.JSON_RPC_VERSION,
						Id:      1,
					},
					Request: &schema.PingRequestSchema{
						MethodName: "ping",
					},
				},
			},
		},
		{
			name:           "seminormal: error occurs for invalid JSON",
			chunk:          []byte(`{"jsonrpc":"2.0","id":1,INVALID}` + "\n"),
			expectMsgCount: 0,
			expected:       nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewReadBuffer()
			err := sut.Append(tt.chunk)
			if err != nil {
				t.Errorf("Append() error = %v", err)
				return
			}

			var got []schema.JsonRpcMessage
			for i := 0; i < tt.expectMsgCount; i++ {
				msg, err := sut.ReadMessage()
				if (err != nil) != tt.wantErr {
					t.Errorf("ReadMessage() error = %v", err)
					return
				}
				got = append(got, msg)
			}
			if len(got) != tt.expectMsgCount {
				t.Errorf("ReadMessage() got = %v, want %v", len(got), tt.expectMsgCount)
				return
			}
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("ReadMessage() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
