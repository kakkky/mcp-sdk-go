package protocol

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kakkky/mcp-sdk-go/shared/mcp_err"
	"github.com/kakkky/mcp-sdk-go/shared/protocol/test"
	"github.com/kakkky/mcp-sdk-go/shared/schema"
)

func TestProtocol_Connect(t *testing.T) {
	// モックの Transport を作成
	transport := test.NewMockChannelServerTransport(
		make(chan schema.JsonRpcMessage, 1),
		make(chan schema.JsonRpcMessage, 1),
	)

	// Protocol インスタンスを作成
	protocol := NewProtocol()

	// Connect メソッドを呼び出し
	protocol.Connect(transport)

	// Transport が正しく設定されているか確認
	if protocol.Transport() != transport {
		t.Errorf("Expected transport to be set, but it was not")
	}
}

func TestProtocol_Request(t *testing.T) {
	tests := []struct {
		name             string
		request          schema.Request
		resultSchema     schema.Result
		expectedResult   schema.Result
		expectedError    *mcp_err.McpErr
		isExpectedMcpErr bool
	}{
		{
			name: "nomal case :client send request and receive response successfully",
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
			isExpectedMcpErr: false,
		},
		{
			name: "sminormal case :client send unknown request and receive error response",
			request: schema.Request{
				Method: "unknown",
			},
			expectedError: &mcp_err.McpErr{
				Code: mcp_err.METHOD_NOT_FOUND,
			},
			isExpectedMcpErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// プロトコルのインスタンスを作成
			server := NewProtocol()
			client := NewProtocol()

			// トランスポートのモックを作成
			serverToClientCh := make(chan schema.JsonRpcMessage, 1)
			clientToServerCh := make(chan schema.JsonRpcMessage, 1)

			// トランスポートを初期化
			serverTransport := test.NewMockChannelServerTransport(clientToServerCh, serverToClientCh)
			clientTransport := test.NewMockChannelClientTransport(clientToServerCh, serverToClientCh)

			// Close時コールバックを登録
			server.SetOnClose(func() {
				close(serverToClientCh)
			})
			client.SetOnClose(func() {
				close(clientToServerCh)
			})

			// サーバー側でリクエストハンドラを登録
			server.SetRequestHandler(schema.Request{Method: "test"}, func(request schema.JsonRpcRequest) (schema.Result, *mcp_err.McpErr) {
				return tt.expectedResult, nil
			})

			// 通信を開始
			server.Connect(serverTransport)
			client.Connect(clientTransport)
			// クリーンアップ
			defer func() {
				server.Close()
				client.Close()
			}()
			// リクエストを受け取ったら、レスポンスを返すことを確認する
			got, err := client.Request(tt.request, tt.resultSchema)
			// テストケースがエラーを期待する場合、エラーが期待通りか確認
			if tt.isExpectedMcpErr {
				e, ok := err.(*mcp_err.McpErr)
				if !ok {
					t.Errorf("Request() got error = %v, want %v", err, tt.expectedError)
					return
				}
				if e.Code != tt.expectedError.Code {
					t.Errorf("Request() got error code = %v, want %v", e.Code, tt.expectedError.Code)
					return
				}
				return
			}
			if err != nil {
				t.Errorf("Request() got error = %v", err)
				return
			}
			if diff := cmp.Diff(got, &tt.expectedResult); diff != "" {
				t.Errorf("Request() got = %v, want %v, diff %s", got, tt.expectedResult, diff)
			}
		})
	}
}
