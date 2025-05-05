package protocol

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol/test"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestProtocol_Connect(t *testing.T) {
	// モックの Transport を作成
	transport := test.NewMockChannelTransport(make(chan schema.JsonRpcMessage, 1))

	// Protocol インスタンスを作成
	protocol := NewProtocol()

	// Connect メソッドを呼び出し
	protocol.Connect(transport)

	// Transport が正しく設定されているか確認
	if protocol.Transport() != transport {
		t.Errorf("Expected transport to be set, but it was not")
	}
}

func TestProtocolRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        schema.Request
		resultSchema   schema.Result
		expectedResult schema.Result
		expectedError  *mcp_err.McpErr
		isExpectedErr  bool
	}{
		{
			name: "nomal case : send request and receive response successfully",
			request: schema.Request{
				Method: "test",
			},
			resultSchema: schema.Result{
				Result: map[string]string{},
			},
			expectedResult: schema.Result{
				Result: map[string]string{
					"status": "success",
				},
			},
			isExpectedErr: false,
		},
		// {
		// 	name: "sminormal case : send unknown request and receive error response",
		// 	request: schema.Request{
		// 		Method: "unknown",
		// 	},
		// 	expectedError: &mcp_err.McpErr{
		// 		Code: mcp_err.METHOD_NOT_FOUND,
		// 	},
		// 	isExpectedErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// プロトコルを作成
			server := NewProtocol()
			client := NewProtocol()

			// サーバー側でリクエストハンドラを登録
			server.SetRequestHandler(schema.Request{Method: "test"}, func(request schema.JsonRpcRequest) (schema.Result, *mcp_err.McpErr) {
				return tt.expectedResult, nil
			})

			// トランスポートのモックを作成
			ch := make(chan schema.JsonRpcMessage, 1)
			defer close(ch)
			ServerTransport := test.NewMockChannelTransport(ch)
			ClientTransport := test.NewMockChannelTransport(ch)

			server.Connect(ServerTransport)
			client.Connect(ClientTransport)
			// リクエストを受け取ったら、レスポンスを返す
			got, err := client.Request(tt.request, tt.resultSchema)
			if err != nil {
				fmt.Println("Error:", err)
			}
			if diff := cmp.Diff(got, &tt.expectedResult); diff != "" {
				t.Errorf("Request() got = %v, want %v, diff %s", got, tt.expectedResult, diff)
			}
		})
	}
}
